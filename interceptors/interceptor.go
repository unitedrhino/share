package interceptors

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func Error(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		logx.WithContext(ctx).Errorf("err=%s", errors.Fmt(err).Error())
	} else {
		body := utils.MarshalNoErr(resp)
		if len(body) > 1024 {
			body = body[:1024]
		}
		logx.WithContext(ctx).Infof("resp=%v", body)
	}
	var acceptLanguage string
	if uc := ctxs.GetUserCtx(ctx); uc != nil {
		acceptLanguage = uc.AcceptLanguage
	}
	err = errors.ToRpc(err, acceptLanguage)
	return resp, err
}

func Ctxs(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	ctx2 := func() context.Context {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return ctx
		}
		info := md[ctxs.UserInfoKey]
		if len(info) == 0 {
			return ctx
		}
		var val ctxs.UserCtx
		str, err := base64.StdEncoding.DecodeString(info[0])
		if err != nil {
			return ctx
		}
		if err := json.Unmarshal(str, &val); err != nil {
			return ctx
		}
		return ctxs.SetUserCtx(ctx, &val)
	}()
	resp, err := handler(ctx2, req)
	return resp, err
}
