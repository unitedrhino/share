package ctxs

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"go.opentelemetry.io/otel/trace"
	"time"
)

var ContextKeys = []string{UserToken2Key, UserTokenKey, UserSetTokenKey, UserRoleKey, MetadataKey, UserAppCodeKey, UserAppCodeKey2, UserTenantCodeKey}

func CopyCtx(ctx context.Context) context.Context {
	newCtx := NewUserCtx(ctx)
	newCtx = trace.ContextWithSpanContext(newCtx, trace.SpanContextFromContext(ctx))
	for _, k := range ContextKeys {
		if v := ctx.Value(k); v != nil {
			newCtx = context.WithValue(newCtx, k, v)
		}
	}
	return newCtx
}

func GoNewCtx(ctx context.Context, f func(ctx2 context.Context)) {
	ctx = CopyCtx(ctx)
	go func() {
		defer utils.Recover(ctx)
		f(ctx)
	}()
}

func GetDeadLine(ctx context.Context, defaultDeadLine time.Time) time.Time {
	dead, ok := ctx.Deadline()
	if !ok {
		return defaultDeadLine
	}
	return dead
}
