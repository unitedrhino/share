package tools

import (
	"context"
	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

func RunAllTenants(ctx context.Context, f func(ctx context.Context) error) error {
	tenantCodes, err := caches.GetTenantCodes(ctx)
	if err != nil {
		return err
	}
	for _, v := range tenantCodes {
		ctx := ctxs.BindTenantCode(ctx, v, 0)
		utils.Go(ctx, func() {
			err := f(ctx)
			if err != nil {
				logx.Error(err)
			}
		})
	}
	return nil
}
