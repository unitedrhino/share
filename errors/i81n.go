package errors

import (
	"embed"
	"fmt"
	"strings"

	"gitee.com/unitedrhino/share/i18ns"
	"github.com/zeromicro/go-zero/core/logx"
)

//go:embed locale/*.json
var LocaleFS embed.FS

// Msg 格式化消息结构
type Msg struct {
	Format string
	Args   []any
}

func init() {
	var err error
	err = i18ns.InitWithEmbedFS(LocaleFS, "locale")
	logx.Must(err)
}

// GetI18nMsg 获取多语言消息
func (c CodeError) GetI18nMsg(accept string) string {
	var msgs []string
	for _, v := range c.Msg {
		msg := i18ns.LocalizeMsgWithLang(accept, v.Format, v.Args...)
		msgs = append(msgs, msg)
	}
	return strings.Join(msgs, ":")
}

// String 格式化消息的字符串实现
func (m Msg) String() string {
	return fmt.Sprintf(m.Format, m.Args...)
}

// stringMsgs 将消息数组转换为字符串
func stringMsgs(in []Msg) string {
	var msgs []string
	for _, v := range in {
		msgs = append(msgs, v.String())
	}
	return strings.Join(msgs, ":")
}
