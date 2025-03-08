package services

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

var (
	apisvrs []func(server *rest.Server) error
)

func RegisterApisvr(run func(server *rest.Server) error) {
	apisvrs = append(apisvrs, run)
}

func InitApisvrs(svr *rest.Server) *rest.Server {
	for _, r := range apisvrs {
		err := r(svr)
		logx.Must(err)
	}
	return svr
}
