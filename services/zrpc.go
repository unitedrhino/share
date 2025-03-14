package services

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	RegisterFn = func(*grpc.Server)
)

func MustNewServer(c zrpc.RpcServerConf, register RegisterFn) *zrpc.RpcServer {
	if c.Etcd.Key == "" && c.Name != "" {
		c.Etcd.Key = c.Name
	}
	server, err := zrpc.NewServer(c, register)
	logx.Must(err)
	return server
}
