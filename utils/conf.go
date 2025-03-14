package utils

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	confPrefix     = os.Getenv("confPrefix")
	confSuffix     = os.Getenv("confSuffix")
	commonFileName = "common.yaml"
)

func getFileName(path string) string {
	if confPrefix != "" || confSuffix != "" {
		dir, file := filepath.Split(path)
		files := strings.Split(file, ".")
		if len(files) > 1 && confSuffix != "" {
			file = fmt.Sprintf("%s%s.%s", files[0], confSuffix, files[len(files)-1])
		}
		newPath := filepath.Join(dir, fmt.Sprintf("%s%s", confPrefix, file))
		_, err := os.Stat(newPath)
		if err == nil {
			return newPath
		}
	}
	return ""
}

func ConfMustLoad(path string, v any) {
	if path[0] == '/' { //绝对目录
		newPath := getFileName(path)
		if newPath != "" {
			MustLoad(newPath, v)
			return
		}
		MustLoad(path, v)
		return
	}
	//相对目录
	newPath := getFileName(path)
	if newPath != "" {
		MustLoad(newPath, v)
		return
	}
	_, err := os.Stat(path)
	if err == nil {
		MustLoad(path, v)
		return
	}
	{
		if strings.HasPrefix(path, "./") {
			path = path[2:]
		}
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		logx.Must(err)
		newPath := filepath.Join(dir, path)
		newPath2 := getFileName(newPath)
		if newPath2 != "" {
			MustLoad(newPath2, v)
			return
		}
		MustLoad(newPath, v)
	}
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
	return newPath
}

func LoadContent(content []byte, v any) error {
	content = []byte(os.ExpandEnv(string(content)))
	return conf.LoadFromYamlBytes(content, v)
}

func MustLoadContent(content []byte, v any) {
	if err := LoadContent(content, v); err != nil {
		log.Fatalf("error: config file  %s", err.Error())
	}
}

func MustLoad(path string, v any) {
	if err := Load(path, v); err != nil {
		log.Fatalf("error: config file  %s", err.Error())
	}
}

func Load(path string, v any) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	dir, _ := filepath.Split(abs)
	commonFile := filepath.Join(dir, commonFileName)
	newCommonFile := getFileName(commonFile)
	if newCommonFile != "" {
		cfg, err := os.ReadFile(newCommonFile)
		if err != nil {
			return err
		}
		cfg = append(cfg, '\n')
		content = append(cfg, content...)
	} else {
		cfg, _ := os.ReadFile(commonFile)
		if cfg != nil {
			cfg = append(cfg, '\n')
			content = append(cfg, content...)
		}
	}

	content = []byte(os.ExpandEnv(string(content)))
	err = conf.LoadFromYamlBytes(content, v)
	return err
}
