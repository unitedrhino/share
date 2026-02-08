# 实例级别隔离使用说明

## 概述

实例级别隔离（Instance Isolation）允许不同租户使用独立的数据库实例，实现数据的物理隔离。每个租户可以配置独立的主库DSN和读库DSN（可选）。

## 隔离模式对比

| 模式 | 函数 | 说明 |
|------|------|------|
| 共享数据库 | `GetTenantConn` | 所有租户共享同一数据库 |
| Schema隔离 | `GetSchemaTenantConn` | PostgreSQL schema级别隔离 |
| 实例隔离 | `GetInstTenantConn` | 每个租户独立数据库实例 |

## 表结构设计

租户DSN配置存储在公共数据库的 `sys_tenant_dsn` 表中：

```sql
CREATE TABLE sys_tenant_dsn (
    id            BIGINT PRIMARY KEY AUTO_INCREMENT,
    tenant_code   VARCHAR(64) NOT NULL UNIQUE COMMENT '租户代码',
    db_type       VARCHAR(16) NOT NULL DEFAULT 'mysql' COMMENT '数据库类型: mysql/pgsql/sqlite',
    dsn           TEXT NOT NULL COMMENT '主库DSN(加密存储)',
    read_dsn      TEXT COMMENT '读库DSN列表(加密存储,JSON数组)',
    status        SMALLINT NOT NULL DEFAULT 1 COMMENT '状态: 1启用 2禁用',
    created_time  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_time  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## 使用方式

### 1. 注册DSN获取器

在服务启动时，注册DSN获取函数和租户列表获取函数：

```go
import "gitee.com/unitedrhino/share/stores"

func main() {
    // 初始化公共数据库连接
    stores.InitConn(database)

    // 注册DSN获取器
    stores.RegisterInstDsn(func(ctx context.Context, tenantCode string) (dbType, dsn string, readDSN []string, err error) {
        // 从 sys_tenant_dsn 表查询租户DSN配置
        // 返回解密后的DSN
        return queryTenantDsn(ctx, tenantCode)
    })

    // 注册租户列表获取器（用于启动时预热）
    stores.RegisterInstTenantCodes(func(ctx context.Context) ([]string, error) {
        // 返回所有启用状态的租户代码列表
        return queryAllTenantCodes(ctx)
    })

    // 注册数据初始化函数（表结构迁移和基础数据初始化）
    stores.RegisterInstDataInit(func(ctx context.Context, tenantCode string, db *gorm.DB) error {
        // 自动迁移表结构
        if err := db.AutoMigrate(&User{}, &Order{}); err != nil {
            return err
        }
        // 初始化基础数据（如有需要）
        return initBaseData(ctx, db)
    })

    // 初始化所有租户连接（预热+数据初始化）
    if err := stores.InstDsnInit(ctx); err != nil {
        log.Fatal(err)
    }
}
```

### 2. 获取租户连接

在业务代码中使用 `GetInstTenantConn` 获取租户数据库连接：

```go
func GetUserList(ctx context.Context) ([]User, error) {
    // ctx 中需要包含租户代码
    ctx = stores.SetTenantCode(ctx, "tenant001")

    db := stores.GetInstTenantConn(ctx)
    if db.Error != nil {
        return nil, db.Error
    }

    var users []User
    err := db.Find(&users).Error
    return users, err
}
```

### 3. DSN加密存储

DSN中包含数据库密码，建议使用AES加密存储：

```go
import "gitee.com/unitedrhino/share/utils"

// 加密DSN
func encryptDsn(dsn, secret string) (string, error) {
    return utils.AesCbcBase64(dsn, secret)
}

// 解密DSN（需要自行实现解密函数）
func decryptDsn(encrypted, secret string) (string, error) {
    // 实现AES-CBC解密
}
```

## 连接缓存机制

- 启动时通过 `InstDsnInit` 预热所有租户连接
- 运行时首次访问新租户会自动查询DSN并建立连接
- 连接会被缓存到 `instTenantConn` 中，后续请求复用

## 读写分离

如果配置了 `read_dsn`，系统会自动配置读写分离：
- 写操作路由到主库
- 读操作随机路由到读库

```go
// 强制使用主库
db.Clauses(dbresolver.Write).Find(&users)

// 强制使用读库
db.Clauses(dbresolver.Read).Find(&users)
```

## 连接池配置

所有租户使用统一的连接池配置，继承自 `commonDB` 的配置：
- `MaxOpenConns`: 最大打开连接数（默认20）
- `MaxIdleConns`: 最大空闲连接数（默认10）
- `ConnMaxIdleTime`: 连接最大空闲时间（10分钟）
- `ConnMaxLifetime`: 连接最大生命周期（30分钟）

## 数据初始化

### 初始化时机

数据初始化函数会在以下场景被调用：

1. **服务启动时**: `InstDsnInit` 初始化所有已有租户
2. **新租户首次连接**: `GetInstTenantConn` 检测到新租户时自动初始化
3. **版本更新迁移**: `InstMigrate` 批量更新所有租户表结构

### 版本更新迁移

当服务版本更新需要修改表结构时，调用 `InstMigrate`：

```go
// 版本更新时批量迁移所有租户
if err := stores.InstMigrate(ctx); err != nil {
    log.Fatalf("租户数据迁移失败: %v", err)
}
```

### 初始化函数示例

```go
func initTenantData(ctx context.Context, tenantCode string, db *gorm.DB) error {
    // 1. 自动迁移表结构
    if err := db.AutoMigrate(
        &model.User{},
        &model.Order{},
        &model.Product{},
    ); err != nil {
        return fmt.Errorf("表结构迁移失败: %w", err)
    }

    // 2. 初始化基础数据（幂等操作）
    if err := initDefaultRoles(db); err != nil {
        return fmt.Errorf("初始化角色失败: %w", err)
    }

    return nil
}
```

## 连接健康检查

系统内置了连接健康检查机制，自动处理失效连接。

### 自动检查

`GetInstTenantConn` 在从缓存获取连接时会自动进行 ping 检查：
- 如果连接有效，直接返回
- 如果连接失效，自动从缓存删除并重新建立连接

### 手动剔除连接

当租户被删除时，可以手动剔除连接缓存：

```go
// 剔除指定租户的连接
stores.RemoveInstTenantConn("tenant001")
```

### 批量健康检查

定期检查所有租户连接的健康状态：

```go
// 返回失效的租户列表，并自动剔除
invalidTenants := stores.InstHealthCheck()
if len(invalidTenants) > 0 {
    log.Warnf("以下租户连接已失效: %v", invalidTenants)
}
```

## 注意事项

1. **安全性**: DSN必须加密存储，密钥应妥善保管
2. **预热**: 建议在服务启动时预热所有租户连接，避免首次请求延迟
3. **监控**: 建议监控各租户的连接池使用情况
4. **故障隔离**: 单个租户数据库故障不会影响其他租户
