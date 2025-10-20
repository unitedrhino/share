package i18ns

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

func init() {
	l, ok := os.LookupEnv("SYS_OS_LANG")
	if ok {
		tags, _, err := language.ParseAcceptLanguage(l)
		logx.Must(err)
		bundle = i18n.NewBundle(tags[0])
		return
	}
	bundle = i18n.NewBundle(language.English)
}

var bundle *i18n.Bundle = i18n.NewBundle(language.Chinese)

func init() {
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
}

// InitWithEmbedFS 从嵌入文件系统初始化i18n bundle
func InitWithEmbedFS(LocaleFS embed.FS, dir string) error {
	dirs, err := LocaleFS.ReadDir(dir)
	if err != nil {
		logx.Errorf("读取语言文件目录失败: %v", err)
		return err
	}

	for _, v := range dirs {
		path := dir + "/" + v.Name()
		_, err = bundle.LoadMessageFileFS(LocaleFS, path)
		if err != nil {
			logx.Errorf("加载语言文件失败 %s: %v", path, err)
			continue
		}
		logx.Infof("成功加载语言文件: %s", path)
	}
	return nil
}

// InitWithFS 从文件系统初始化i18n bundle
func InitWithFS(dir string) error {
	dirs, err := os.ReadDir(dir)
	if err != nil {
		logx.Errorf("读取语言文件目录失败: %v", err)
		return err
	}

	for _, v := range dirs {
		if !(strings.HasSuffix(v.Name(), "json") || strings.HasSuffix(v.Name(), "toml") || strings.HasSuffix(v.Name(), "yaml")) {
			continue
		}
		path := filepath.Join(dir, v.Name())
		_, err = bundle.LoadMessageFile(path)
		if err != nil {
			logx.Errorf("加载语言文件失败 %s: %v", path, err)
			continue
		}
		logx.Infof("成功加载语言文件: %s", path)
	}
	return nil
}

// 示例:  	msg := i18ns.LocalizeMsgWithLang("en_US", "nodered.protocol.unsupported", "vewwrfw3")
func LocalizeMsgWithLang(lang string, format string, args ...interface{}) string {
	localizer := i18n.NewLocalizer(bundle, lang)
	msg, e := localizer.LocalizeMessage(&i18n.Message{ID: format})
	if e != nil {
		msg = format
	}
	if len(args) == 0 {
		return msg
	}
	return fmt.Sprintf(msg, args...)
}
func LocalizeMsg(format string, args ...interface{}) string {
	return LocalizeMsgWithLang("", format, args...)
}
