package i18ns

import (
	"embed"
	"encoding/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/text/language"
	"io/fs"
)

func InitWithLocaleFS(LocaleFS embed.FS, dir string) *i18n.Bundle {
	dirs, err := LocaleFS.ReadDir(dir)
	logx.Must(err)
	bundle := i18n.NewBundle(language.Chinese)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	var dataOne = map[string]string{}
	for _, v := range dirs {
		path := dir + "/" + v.Name()
		_, err = bundle.LoadMessageFileFS(LocaleFS, path)
		logx.Must(err)
	}
	path := dir + "/" + dirs[0].Name()
	data, err := fs.ReadFile(LocaleFS, path)
	logx.Must(err)
	json.Unmarshal(data, &dataOne)
	for k := range dataOne {
		dataOne[k] = k
	}
	data, err = json.Marshal(dataOne)
	logx.Must(err)
	_, err = bundle.ParseMessageFileBytes(data, "zh.json")
	logx.Must(err)
	return bundle
}
