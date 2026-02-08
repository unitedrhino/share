package stores

import (
	"context"
	"sync"
	"testing"

	"gitee.com/unitedrhino/share/conf"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func TestTenantDSN(t *testing.T) {
	assert.Equal(t, "a.t1.b", tenantDSN("a.{.tenantCode}.b", "t1"))
	assert.Equal(t, "a.t1.b", tenantDSN("a.{{.tenantCode}}.b", "t1"))
	assert.Equal(t, "a.b", tenantDSN("a.b", "t1"))
}

func TestGetInstDsn_RegisterAndFallback(t *testing.T) {
	commonDB = conf.Database{
		DBType:  conf.Mysql,
		DSN:     "root:password{.tenantCode}@tcp(mysql.{.tenantCode}:3306)/iThings",
		ReadDSN: []string{"root:password{.tenantCode}@tcp(mysql-ro.{.tenantCode}:3306)/iThings"},
	}
	RegisterInstDsn(nil)
	assert.Equal(t, "root:passwordt1@tcp(mysql.t1:3306)/iThings", GetInstDsn("t1"))

	RegisterInstDsn(func(ctx context.Context, tenantCode string) (string, string, []string, error) {
		return conf.Mysql, "override", []string{"r1"}, nil
	})
	t.Cleanup(func() { RegisterInstDsn(nil) })
	assert.Equal(t, "override", GetInstDsn("t1"))
}

func TestGetInstTenantConn_CacheAndDialectorDSN(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	commonConn = db
	commonDB = conf.Database{
		DBType:       conf.Sqlite,
		DSN:          "file::memory:?cache=shared",
		MaxIdleConns: 1,
		MaxOpenConns: 2,
	}
	instTenantConn = sync.Map{}
	RegisterInstDsn(func(ctx context.Context, tenantCode string) (string, string, []string, error) {
		return conf.Sqlite, "file::memory:?cache=shared", nil, nil
	})
	RegisterInstTenantCodes(nil)
	t.Cleanup(func() {
		RegisterInstDsn(nil)
		RegisterInstTenantCodes(nil)
	})

	ctx := context.Background()
	ctx = SetTenantCode(ctx, "t1")
	ctx = SetIsDebug(ctx, false)

	got := GetInstTenantConn(ctx)
	assert.NotNil(t, got)
	assert.NoError(t, got.Error)

	cachedAny, ok := instTenantConn.Load("t1")
	assert.True(t, ok)

	got2 := GetInstTenantConn(ctx)
	assert.NoError(t, got2.Error)

	cachedAny2, ok := instTenantConn.Load("t1")
	assert.True(t, ok)
	assert.Same(t, cachedAny, cachedAny2)

	assert.NotPanics(t, func() {
		_ = got2.Clauses(dbresolver.Read)
	})
}

func TestInstDsnInit(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	commonConn = db
	commonDB = conf.Database{
		DBType:       conf.Sqlite,
		DSN:          "file::memory:?cache=shared",
		MaxIdleConns: 1,
		MaxOpenConns: 2,
	}
	instTenantConn = sync.Map{}

	RegisterInstDsn(func(ctx context.Context, tenantCode string) (string, string, []string, error) {
		return conf.Sqlite, "file::memory:?cache=shared", nil, nil
	})
	RegisterInstTenantCodes(func(ctx context.Context) ([]string, error) {
		return []string{"t1", "t2"}, nil
	})
	t.Cleanup(func() {
		RegisterInstDsn(nil)
		RegisterInstTenantCodes(nil)
	})

	err = InstDsnInit(context.Background())
	assert.NoError(t, err)

	_, ok := instTenantConn.Load("t1")
	assert.True(t, ok)
	_, ok = instTenantConn.Load("t2")
	assert.True(t, ok)
}
