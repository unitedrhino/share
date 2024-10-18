package wxClient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	ctx2 "github.com/silenceper/wechat/v2/work/context"
	"github.com/silenceper/wechat/v2/work/robot"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"
	"net/url"
)

const wx_token_url = "https://qyapi.weixin.qq.com/cgi-bin/gettoken"
const wx_msg_url = "https://qyapi.weixin.qq.com/cgi-bin/message/send" //需要拼接token
const wx_login_url = "https://login.work.weixin.qq.com/wwlogin/sso/login"
const wx_get_userInfo = "https://qyapi.weixin.qq.com/cgi-bin/auth/getuserinfo"
const wx_auth_url = "https://open.weixin.qq.com/connect/oauth2/authorize"

type WxCorAppClient struct {
	CorpID     string
	CorpSecret string
	AgentID    int64
}

type WxTokenRes struct {
	Errcode     int64  `json:"errcode"`      //出错返回码，为0表示成功，非0表示调用失败
	Errmsg      string `json:"errmsg"`       //返回码提示语
	Accesstoken string `json:"access_token"` //获取到的凭证，最长为512字节
	ExpiresIn   int64  `json:"expires_in"`   //凭证的有效时间（秒）
}

type WxMsgReq struct {
	ToUser                 string       `json:"touser,omitempty"`                   //指定接收消息的成员，成员ID列表（多个接收者用‘|’分隔，最多支持1000个）。特殊情况：指定为"@all"，则向该企业应用的全部成员发送
	ToParty                string       `json:"toparty,omitempty"`                  //指定接收消息的部门，部门ID列表，多个接收者用‘|’分隔，最多支持100个。当touser为"@all"时忽略本参数
	ToTag                  string       `json:"totag,omitempty"`                    //指定接收消息的标签，标签ID列表，多个接收者用‘|’分隔，最多支持100个。当touser为"@all"时忽略本参数
	MsgType                string       `json:"msgtype,omitempty"`                  //消息类型，此时固定为：text
	Agentid                int64        `json:"agentid,omitempty"`                  //企业应用的id，整型。企业内部开发，可在应用的设置页面查看；第三方服务商，可通过接口 获取企业授权信息 获取该参数值
	Text                   WxMsgReqText `json:"text"`                               //消息内容
	Safe                   int64        `json:"safe,omitempty"`                     //表示是否是保密消息，0表示可对外分享，1表示不能分享且内容显示水印，默认为0
	EnableIdTrans          int64        `json:"enable_id_trans,omitempty"`          //	表示是否开启id转译，0表示否，1表示是，默认0。仅第三方应用需要用到，企业自建应用可以忽略。
	EnableDuplicateCheck   int64        `json:"enable_duplicate_check,omitempty"`   //表示是否开启重复消息检查，0表示否，1表示是，默认0
	DuplicateCheckInterval int64        `json:"duplicate_check_interval,omitempty"` //表示是否重复消息检查的时间间隔，默认1800s，最大不超过4小时
}

type WxMsgResp struct {
	Errcode        int64  `json:"errcode"`
	Errmsg         string `json:"errmsg"`
	Invaliduser    string `json:"invaliduser"`
	Invalidparty   string `json:"invalidparty"`
	Invalidtag     string `json:"invalidtag"`
	Unlicenseduser string `json:"unlicenseduser"`
	Msgid          string `json:"msgid"`
	ResponseCode   string `json:"response_code"`
}

type WxMsgReqText struct {
	Content string `json:"content"`
}

type WxUserInfoRes struct {
	Errcode    int64  `json:"errcode"`
	Errmsg     string `json:"errmsg"`
	UserID     string `json:"userid"`
	UserTicket string `json:"user_ticket"`
}

func NewWxCorAppClient(corpid string, appSecret string, agentId int64) *WxCorAppClient {
	return &WxCorAppClient{
		CorpID:     corpid,
		CorpSecret: appSecret,
		AgentID:    agentId,
	}
}

func (c *WxCorAppClient) GetWxCopAccesstoken() (string, error) {
	token_url := wx_token_url + "?corpid=" + c.CorpID + "&corpsecret=" + c.CorpSecret
	resp, err := http.Get(token_url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var tokenData WxTokenRes
	err = json.Unmarshal(body, &tokenData)
	if err != nil {
		return "", err
	}

	if tokenData.Errcode != 0 {
		return "", errors.New(tokenData.Errmsg)
	}
	return tokenData.Accesstoken, nil
}

func (c *WxCorAppClient) SendWxCopMsg(msg string, toUser string) error {
	accessToken, err := c.GetWxCopAccesstoken()
	if err != nil {
		return err
	}

	req := WxMsgReq{
		ToUser:  toUser,
		MsgType: "text",
		Agentid: c.AgentID,
		Text:    WxMsgReqText{Content: msg},
	}
	reqJson, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := http.Post(wx_msg_url+"?access_token="+accessToken, "application/json", bytes.NewBuffer(reqJson))
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var msgData WxMsgResp
	err = json.Unmarshal(body, &msgData)
	if err != nil {
		return err
	}

	if msgData.Errcode != 0 {
		return errors.New(msgData.Errmsg)
	}
	return nil
}

func (c *WxCorAppClient) GetLoginUrl(redirectUri string, state string, agentID string) string {
	loginUrl := wx_login_url + "?login_type=CorpApp&appid=" + c.CorpID + "&agentid=" + agentID + "&redirect_uri=" + url.QueryEscape(redirectUri) + "&state=" + url.QueryEscape(state)
	return loginUrl
}

func (c *WxCorAppClient) GetUserInfo(code string) (string, error) {
	accessToken, err := c.GetWxCopAccesstoken()
	if err != nil {
		return "", err
	}

	url := wx_get_userInfo + "?access_token=" + accessToken + "&code=" + code
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var infoRes WxUserInfoRes
	err = json.Unmarshal(body, &infoRes)
	if err != nil {
		return "", err
	}
	fmt.Println(infoRes)

	if infoRes.Errcode != 0 {
		return "", errors.New(infoRes.Errmsg)
	}

	if infoRes.UserID == "" {
		return "", errors.New("无法获取您的信息，请使用企业微信扫码")
	}

	return infoRes.UserID, err
}

func (c *WxCorAppClient) GetAuthUrl(redirectUri string, state string, agentID string) string {
	auth_url := wx_auth_url + "?appid=" + c.CorpID + "&redirect_uri=" + url.QueryEscape(redirectUri) + "&response_type=code&scope=snsapi_base&state=" + url.QueryEscape(state) + "&agentid=" + agentID + "#wechat_redirect"
	return auth_url
}

func SendRobotMsg(ctx context.Context, token string, msg string) error {
	u, err := url.Parse(token)
	if err == nil {
		params, err := url.ParseQuery(u.RawQuery)
		if err != nil {
			return err
		}
		token = params.Get("key")
	}
	ret, err := robot.NewClient(&ctx2.Context{
		Config:            nil,
		AccessTokenHandle: nil,
	}).RobotBroadcast(token, robot.WebhookSendTextOption{
		MsgType: "text",
		Text: struct {
			Content             string   `json:"content"`
			MentionedList       []string `json:"mentioned_list"`
			MentionedMobileList []string `json:"mentioned_mobile_list"`
		}(struct {
			Content             string
			MentionedList       []string
			MentionedMobileList []string
		}{
			Content: msg,
		}),
	})
	if err != nil {
		logx.WithContext(ctx).Errorf("SendRobotMsg ret:%v err:%v", ret, err)
		return err
	}
	return nil
}
