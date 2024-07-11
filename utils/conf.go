package utils

import (
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"path/filepath"
)

func ConfMustLoad(path string, v any, opts ...conf.Option) {
	if path[0] == '/' {
		conf.MustLoad(path, v, opts...)
		return
	}
	_, err := os.Stat(path)
	if err == nil {
		conf.MustLoad(path, v, opts...)
		return
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	logx.Must(err)
	newPath := filepath.Join(dir, path)
	conf.MustLoad(newPath, v, opts...)
}
