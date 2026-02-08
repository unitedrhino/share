package stores

import (
	"context"
	"sync"
	"time"

	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/errors"
	"github.com/glebarez/sqlite"
	"github.com/spf13/cast"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

var (
	instTenantConn        sync.Map
	instDBGetter          func(ctx context.Context, tenantCode string) (conf.Database, error)
	instTenantCodesGetter func(ctx context.Context) ([]string, error)
	instDataInitFunc      func(ctx context.Context, tenantCode string, db *gorm.DB) error
	instGetterMu          sync.RWMutex
)

func RegisterInstDsn(getter func(ctx context.Context, tenantCode string) (dbType string, dsn string, readDSN []string, err error)) {
	instGetterMu.Lock()
	defer instGetterMu.Unlock()
	if getter == nil {
		instDBGetter = nil
		return
	}
	instDBGetter = func(ctx context.Context, tenantCode string) (conf.Database, error) {
		dbType, dsn, readDSN, err := getter(ctx, tenantCode)
		if err != nil {
			return conf.Database{}, err
		}
		db := normalizeDatabase(commonDB)
		if dbType != "" {
			db.DBType = dbType
		}
		db.DSN = dsn
		db.ReadDSN = readDSN
		return db, nil
	}
}

func RegisterInstTenantCodes(getter func(ctx context.Context) ([]string, error)) {
	instGetterMu.Lock()
	defer instGetterMu.Unlock()
	instTenantCodesGetter = getter
}

// RegisterInstDataInit 注册数据初始化函数
// 该函数会在以下场景被调用：
// 1. InstDsnInit 启动时初始化所有租户
// 2. GetInstTenantConn 新租户首次连接时
// 3. InstMigrate 批量迁移时
func RegisterInstDataInit(initFunc func(ctx context.Context, tenantCode string, db *gorm.DB) error) {
	instGetterMu.Lock()
	defer instGetterMu.Unlock()
	instDataInitFunc = initFunc
}

func GetInstDsn(tenantCode string) string {
	db, err := getInstDatabase(context.Background(), tenantCode)
	if err != nil {
		return ""
	}
	return db.DSN
}

func getInstDatabase(ctx context.Context, tenantCode string) (conf.Database, error) {
	instGetterMu.RLock()
	getter := instDBGetter
	instGetterMu.RUnlock()
	if getter != nil {
		return getter(ctx, tenantCode)
	}

	db := normalizeDatabase(commonDB)
	db.DSN = tenantDSN(db.DSN, tenantCode)
	if len(db.ReadDSN) > 0 {
		readDSN := make([]string, 0, len(db.ReadDSN))
		for _, tpl := range db.ReadDSN {
			readDSN = append(readDSN, tenantDSN(tpl, tenantCode))
		}
		db.ReadDSN = readDSN
	}
	return db, nil
}

func InstDsnInit(ctx context.Context) error {
	instGetterMu.RLock()
	getter := instTenantCodesGetter
	instGetterMu.RUnlock()
	if getter == nil {
		return errors.Parameter.AddDetail("未注册实例租户列表获取函数")
	}
	tenantCodes, err := getter(ctx)
	if err != nil {
		return errors.System.AddDetail("获取实例租户列表失败").AddDetail(err)
	}
	if len(tenantCodes) == 0 {
		return nil
	}

	ctx = SetIsDebug(ctx, false)
	taskPool := make(chan struct{}, 20)
	var wg sync.WaitGroup
	errCh := make(chan error, len(tenantCodes))
	for _, tenantCode := range tenantCodes {
		tc := tenantCode
		wg.Add(1)
		taskPool <- struct{}{}
		go func() {
			defer func() {
				<-taskPool
				wg.Done()
			}()
			tenantCtx := SetTenantCode(ctx, tc)
			db := GetInstTenantConn(tenantCtx)
			if db.Error != nil {
				errCh <- errors.System.AddDetail("初始化租户连接失败").AddDetail(tc).AddDetail(db.Error)
				return
			}
			sqlDB, err := db.DB()
			if err != nil {
				errCh <- errors.System.AddDetail("获取数据库连接失败").AddDetail(tc).AddDetail(err)
				return
			}
			if err := sqlDB.Ping(); err != nil {
				errCh <- errors.System.AddDetail("数据库连接ping失败").AddDetail(tc).AddDetail(err)
				return
			}
			// 调用数据初始化函数
			instGetterMu.RLock()
			initFunc := instDataInitFunc
			instGetterMu.RUnlock()
			if initFunc != nil {
				if err := initFunc(tenantCtx, tc, db); err != nil {
					errCh <- errors.System.AddDetail("租户数据初始化失败").AddDetail(tc).AddDetail(err)
					return
				}
			}
		}()
	}
	wg.Wait()
	close(errCh)
	for e := range errCh {
		if e != nil {
			return e
		}
	}
	return nil
}

func GetInstTenantConn(in any) *gorm.DB {
	ctx, conn := validateAndGetContext(in, commonConn)
	if conn != nil {
		return conn
	}
	conn = commonConn.WithContext(ctx)
	if commonConn == nil {
		conn.Error = errors.System.AddDetail("数据库连接未初始化")
		return conn
	}
	val := ctx.Value(dbCtxTenantCodeKey)
	if val == nil {
		conn.Error = errors.Permissions.AddDetail("没有传入租户号")
		return conn
	}
	tk := cast.ToString(val)
	if tk == "" {
		conn.Error = errors.Permissions.AddDetail("没有传入租户号")
		return conn
	}

	if cached, ok := instTenantConn.Load(tk); ok {
		cc := cached.(*gorm.DB)
		// 健康检查：验证连接是否有效
		if sqlDB, err := cc.DB(); err == nil {
			if err := sqlDB.Ping(); err == nil {
				// 连接正常，直接返回
				cc = cc.WithContext(ctx)
				if val := ctx.Value(dbCtxDebugKey); val != nil && cast.ToBool(val) == false {
					return cc
				}
				return cc.Debug()
			}
		}
		// 连接失效，从缓存中删除，重新建立连接
		instTenantConn.Delete(tk)
	}

	database, err := getInstDatabase(ctx, tk)
	if err != nil {
		conn.Error = errors.System.AddDetail("获取实例连接信息失败").AddDetail(err)
		return conn
	}
	writeDSN := database.DSN
	cfg := gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		DisableAutomaticPing:                     true,
		PrepareStmt:                              false,
		Logger:                                   commonConn.Logger,
		NamingStrategy:                           schema.NamingStrategy{SingularTable: true},
	}

	var cc *gorm.DB
	switch database.DBType {
	case conf.Pgsql:
		cc, err = gorm.Open(postgres.New(postgres.Config{DSN: writeDSN}), &cfg)
	case conf.Sqlite:
		cc, err = gorm.Open(sqlite.Open(writeDSN), &cfg)
	default:
		cc, err = gorm.Open(mysql.Open(writeDSN), &cfg)
	}
	if err != nil {
		conn.Error = errors.System.AddDetail("创建数据库连接失败").AddDetail(err)
		return conn
	}
	if err := applyConnPool(cc, database); err != nil {
		conn.Error = errors.System.AddDetail("设置数据库连接池失败").AddDetail(err)
		return conn
	}

	if len(database.ReadDSN) > 0 {
		if database.DBType == conf.Sqlite {
			conn.Error = errors.Parameter.AddDetail("sqlite不支持读写分离")
			return conn
		}
		replicas := make([]gorm.Dialector, 0, len(database.ReadDSN))
		for _, dsn := range database.ReadDSN {
			switch database.DBType {
			case conf.Pgsql:
				replicas = append(replicas, postgres.New(postgres.Config{DSN: dsn}))
			default:
				replicas = append(replicas, mysql.Open(dsn))
			}
		}
		err := cc.Use(dbresolver.Register(dbresolver.Config{
			Replicas: replicas,
			Policy:   dbresolver.RandomPolicy{},
		}).
			SetMaxIdleConns(database.MaxIdleConns).
			SetMaxOpenConns(database.MaxOpenConns).
			SetConnMaxIdleTime(time.Minute * 10).
			SetConnMaxLifetime(time.Minute * 30))
		if err != nil {
			conn.Error = errors.System.AddDetail("配置读写分离失败").AddDetail(err)
			return conn
		}
	}

	instTenantConn.Store(tk, cc)

	// 新租户首次连接时调用数据初始化
	instGetterMu.RLock()
	initFunc := instDataInitFunc
	instGetterMu.RUnlock()
	if initFunc != nil {
		if err := initFunc(ctx, tk, cc); err != nil {
			conn.Error = errors.System.AddDetail("租户数据初始化失败").AddDetail(err)
			return conn
		}
	}

	cc = cc.WithContext(ctx)
	if val := ctx.Value(dbCtxDebugKey); val != nil && cast.ToBool(val) == false {
		return cc
	}
	return cc.Debug()
}

// InstMigrate 批量更新所有租户表结构
// 用于版本更新时迁移所有租户的数据库表结构
func InstMigrate(ctx context.Context) error {
	instGetterMu.RLock()
	getter := instTenantCodesGetter
	initFunc := instDataInitFunc
	instGetterMu.RUnlock()

	if getter == nil {
		return errors.Parameter.AddDetail("未注册实例租户列表获取函数")
	}
	if initFunc == nil {
		return errors.Parameter.AddDetail("未注册数据初始化函数")
	}

	tenantCodes, err := getter(ctx)
	if err != nil {
		return errors.System.AddDetail("获取实例租户列表失败").AddDetail(err)
	}
	if len(tenantCodes) == 0 {
		return nil
	}

	ctx = SetIsDebug(ctx, false)
	taskPool := make(chan struct{}, 20)
	var wg sync.WaitGroup
	errCh := make(chan error, len(tenantCodes))

	for _, tenantCode := range tenantCodes {
		tc := tenantCode
		wg.Add(1)
		taskPool <- struct{}{}
		go func() {
			defer func() {
				<-taskPool
				wg.Done()
			}()
			tenantCtx := SetTenantCode(ctx, tc)
			db := GetInstTenantConn(tenantCtx)
			if db.Error != nil {
				errCh <- errors.System.AddDetail("获取租户连接失败").AddDetail(tc).AddDetail(db.Error)
				return
			}
			if err := initFunc(tenantCtx, tc, db); err != nil {
				errCh <- errors.System.AddDetail("租户数据迁移失败").AddDetail(tc).AddDetail(err)
				return
			}
		}()
	}

	wg.Wait()
	close(errCh)
	for e := range errCh {
		if e != nil {
			return e
		}
	}
	return nil
}

// RemoveInstTenantConn 手动剔除指定租户的连接缓存
// 用于租户删除或需要强制重建连接时
func RemoveInstTenantConn(tenantCode string) {
	if tenantCode == "" {
		return
	}
	instTenantConn.Delete(tenantCode)
}

// InstHealthCheck 批量检查所有租户连接的健康状态
// 返回失效的租户列表，并自动从缓存中剔除失效连接
func InstHealthCheck() []string {
	var invalidTenants []string
	instTenantConn.Range(func(key, value any) bool {
		tenantCode := key.(string)
		db := value.(*gorm.DB)
		sqlDB, err := db.DB()
		if err != nil {
			invalidTenants = append(invalidTenants, tenantCode)
			instTenantConn.Delete(tenantCode)
			return true
		}
		if err := sqlDB.Ping(); err != nil {
			invalidTenants = append(invalidTenants, tenantCode)
			instTenantConn.Delete(tenantCode)
		}
		return true
	})
	return invalidTenants
}
