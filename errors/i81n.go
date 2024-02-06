package errors

import (
	"embed"
	"fmt"
	"gitee.com/i-Things/share/i18ns"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"strings"
)

//go:embed locale/*.json
var LocaleFS embed.FS
var bundle *i18n.Bundle

type I18nImpl interface {
	I81n(localizer *i18n.Localizer) string
	String() string
}

type Msgf struct {
	Format string
	A      []any
}
type String string

func init() {
	bundle = i18ns.InitWithLocaleFS(LocaleFS, "locale")
}

func (c CodeError) GetI18nMsg(accept string) string {
	localizer := i18n.NewLocalizer(bundle, accept)
	var msgs []string
	for _, v := range c.Msg {
		msg := v.I81n(localizer)
		msgs = append(msgs, msg)
	}
	return strings.Join(msgs, ":")
}

func (m Msgf) I81n(localizer *i18n.Localizer) string {
	msg, e := localizer.LocalizeMessage(&i18n.Message{ID: m.Format})
	if e != nil {
		msg = m.Format
	}
	return fmt.Sprintf(msg, m.A...)
}
func (m Msgf) String() string {
	return fmt.Sprintf(m.Format, m.A...)
}

func (s String) I81n(localizer *i18n.Localizer) string {
	msg, e := localizer.LocalizeMessage(&i18n.Message{ID: string(s)})
	if e != nil {
		msg = string(s)
	}
	return msg
}
func (m String) String() string {
	return string(m)
}
func stringMsgs(in []I18nImpl) string {
	var msgs []string
	for _, v := range in {
		msgs = append(msgs, v.String())
	}
	return strings.Join(msgs, ":")
}
