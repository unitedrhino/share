package caches

import (
	"context"
	"encoding/json"
	"gitee.com/i-Things/share/domain/tenant"
)

// 生产用户数据权限缓存key
func genTenantKey() string {
	return "tenant"
}

func InitTenant(ctx context.Context, tenants ...*tenant.Info) error {
	if len(tenants) == 0 {
		return nil
	}
	return store.HmsetCtx(ctx, genTenantKey(), DoToTenantMap(tenants...))
}

func SetTenant(ctx context.Context, t *tenant.Info) error {
	val, _ := json.Marshal(t)
	return store.HsetCtx(ctx, genTenantKey(), t.Code, string(val))
}

func DelTenant(ctx context.Context, code string) error {
	_, err := store.HdelCtx(ctx, genTenantKey(), code)
	return err
}

func GetTenant(ctx context.Context, code string) (*tenant.Info, error) {
	val, err := store.HgetCtx(ctx, genTenantKey(), code)
	if err != nil {
		return nil, err
	}
	var ret tenant.Info
	err = json.Unmarshal([]byte(val), &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func GetTenantCodes(ctx context.Context) ([]string, error) {
	return store.HkeysCtx(ctx, genTenantKey())
}

func DoToTenantMap(tenants ...*tenant.Info) map[string]string {
	var ret = map[string]string{}
	for _, v := range tenants {
		val, _ := json.Marshal(v)
		ret[v.Code] = string(val)
	}
	return ret
}
