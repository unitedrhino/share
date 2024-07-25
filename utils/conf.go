package utils

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"path/filepath"
	"strings"
)

func ConfMustLoad(path string, v any, opts ...conf.Option) {
	confPrefix := os.Getenv("confPrefix")
	if path[0] == '/' {
		if confPrefix != "" {
			dir, file := filepath.Split(path)
			newPath := filepath.Join(dir, fmt.Sprintf("%s/%s", confPrefix, file))
			_, err := os.Stat(newPath)
			if err == nil {
				conf.MustLoad(newPath, v, opts...)
				return
			}
		}
		conf.MustLoad(path, v, opts...)
		return
	}
	if confPrefix != "" {
		dir, file := filepath.Split(path)
		newPath := filepath.Join(dir, fmt.Sprintf("%s_%s", confPrefix, file))
		_, err := os.Stat(newPath)
		if err == nil {
			conf.MustLoad(newPath, v, opts...)
			return
		}
	}
	_, err := os.Stat(path)
	if err == nil {
		conf.MustLoad(path, v, opts...)
		return
	}
	if strings.HasPrefix(path, "./") {
		path = path[2:]
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	logx.Must(err)
	newPath := filepath.Join(dir, path)
	if confPrefix != "" {
		dir, file := filepath.Split(newPath)
		newPath2 := filepath.Join(dir, fmt.Sprintf("%s_%s", confPrefix, file))
		_, err := os.Stat(newPath2)
		if err == nil {
			conf.MustLoad(newPath2, v, opts...)
			return
		}
	}
	conf.MustLoad(newPath, v, opts...)
}

func GerRealPwd(path string) string {
	if path[0] == '/' {
		return path
	}
	_, err := os.Stat(path)
	if err == nil {
		return path
	}
	if strings.HasPrefix(path, "./") {
		path = path[2:]
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	logx.Must(err)
	newPath := filepath.Join(dir, path)
	fmt.Printf("\n\n\npath:%s\n\n", newPath)
	return newPath
}