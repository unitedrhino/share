package dingClient

import (
	"encoding/json"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/errors"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/zhaoyunxing92/dingtalk/v2"
	"net/url"
)

type DingTalk = dingtalk.DingTalk

func NewDingTalkClient(c *conf.ThirdConf) (*DingTalk, error) {
	if c == nil {
		return nil, nil
	}
	cli, err := dingtalk.NewClient(c.AppKey, c.AppSecret)
	if err != nil {
		return nil, errors.System.AddDetail(err)
	}
	return cli, nil
}

type DingRobot = dingtalk.Robot

func NewDingRobotClient(token string) DingRobot {
	u, err := url.Parse(token)
	if err != nil {
		return dingtalk.NewRobot(token)
	}
	params, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return dingtalk.NewRobot(token)
	}
	at := params.Get("access_token")
	if at != "" {
		return dingtalk.NewRobot(at)
	}
	return dingtalk.NewRobot(token)
}

// 钉钉消息结构体
type text struct {
	Content string `json:"content" validate:"required"`
}

// 文本消息
type textMessage struct {
	MsgType string `json:"msgtype" validate:"required,oneof=text image voice file link oa markdown action_card feedCard"`
	text    `json:"text" validate:"required"`
}

func (t *textMessage) Validate(valid *validator.Validate, trans ut.Translator) error {
	return nil
}

func (t *textMessage) String() string {
	str, _ := json.Marshal(t)
	return string(str)
}

func (t *textMessage) MessageType() string {
	return "text"
}

// NewTextMessage 文本对象
func NewTextMessage(context string) *textMessage {
	msg := &textMessage{}
	msg.MsgType = msg.MessageType()
	msg.text = text{Content: context}
	return msg
}
