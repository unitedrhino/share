package mqttClient

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/utils"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"net/url"
	"time"
)

type Action = string

const (
	ActionConnected    Action = "connected"
	ActionDisconnected Action = "disconnected"
)

type DevConn struct {
	UserName  string `json:"username"`
	Timestamp int64  `json:"timestamp"` //毫秒时间戳
	Address   string `json:"addr"`
	ClientID  string `json:"clientID"`
	/*
		https://www.emqx.com/en/blog/emqx-mqtt-broker-connection-troubleshooting
			normal：客户端主动断开连接；
			kicked：被服务器踢出，通过 REST API；
			keepalive_timeout：保持活动超时；
			not_authorized：认证失败，或者acl_nomatch=disconnect时，没有权限的Pub/Sub会主动断开客户端连接；
			tcp_closed：对端关闭了网络连接；
			discard：因为相同ClientID的客户端上线了，并且设置了clean_start=true；
			takeovered：因为相同ClientID的客户端上线，并且设置了clean_start=false；
			internal_error：格式错误的消息或其他未知错误。
	*/
	Reason     string `json:"reason"`
	Action     Action `json:"action"` //登录 onLogin 登出 onLogout
	ProductID  string `json:"productID"`
	DeviceName string `json:"deviceName"`
}

type EmqResp struct {
	Code    string `json:"code"` //如果不在线返回: CLIENTID_NOT_FOUND
	Message string `json:"message"`
}

type MutSubReq struct {
	Topic string `json:"topic"`
	Qos   int    `json:"qos"`
	Nl    int    `json:"nl"`
	Rap   int    `json:"rap"`
	Rh    int    `json:"rh"`
}

func (m MqttClient) SetClientMutSub(ctx context.Context, clientID string, topics []string, qos int) error {
	if len(topics) == 0 {
		return nil
	}
	logx.WithContext(ctx).Infof("SetClientMut clientID:%v,topics:%v", clientID, topics)
	ts, err := m.GetClientSub(ctx, clientID)
	if err != nil {
		return err
	}
	var topicSet = map[string]struct{}{}
	for _, topic := range topics {
		topicSet[topic] = struct{}{}
	}
	for _, t := range ts {
		if _, ok := topicSet[t]; ok {
			delete(topicSet, t)
		}
	}
	topics = utils.SetToSlice(topicSet)
	if len(topics) == 0 {
		return nil
	}
	if m.cfg.OpenApi == nil {
		return errors.System.AddMsg("没有配置秘钥")
	}
	oa := m.cfg.OpenApi
	greq := gorequest.New().Retry(1, time.Second*2)
	greq.SetBasicAuth(oa.ApiKey, oa.SecretKey)
	var ret []*MutSubReq
	var req []*MutSubReq
	for _, v := range topics {
		req = append(req, &MutSubReq{
			Topic: v,
			Qos:   qos,
		})
	}
	var errs []error
	var body []byte
	var tryTime = 5
	for i := tryTime; i > 0; i-- {
		_, body, errs = greq.Post(fmt.Sprintf("%s/api/v5/clients/%s/subscribe/bulk", oa.Host,
			url.QueryEscape(clientID))).Send(&req).EndStruct(&ret)
		if errs != nil {
			time.Sleep(time.Second / 2 * time.Duration(tryTime) / time.Duration(i))
			continue
		}
		break
	}
	if errs != nil {
		return errors.System.AddDetail(topics, qos, errs, string(body))
	}

	return nil
}
func (m MqttClient) GetClientSub(ctx context.Context, clientID string) ([]string, error) {
	logx.WithContext(ctx).Infof("GetClientSub clientID:%v", clientID)
	if m.cfg.OpenApi == nil {
		return nil, errors.System.AddMsg("没有配置秘钥")
	}
	oa := m.cfg.OpenApi
	greq := gorequest.New().Retry(1, time.Second*2)
	greq.SetBasicAuth(oa.ApiKey, oa.SecretKey)
	var ret []*MutSubReq
	var errs []error
	var body []byte
	var tryTime = 5
	for i := tryTime; i > 0; i-- {
		_, body, errs = greq.Get(fmt.Sprintf("%s/api/v5/clients/%s/subscriptions", oa.Host,
			url.QueryEscape(clientID))).EndStruct(&ret)
		if errs != nil {
			time.Sleep(time.Second / 2 * time.Duration(tryTime) / time.Duration(i))
			continue
		}
		break
	}
	if errs != nil {
		return nil, errors.System.AddDetail(errs, string(body))
	}
	var topics []string
	for _, v := range ret {
		topics = append(topics, v.Topic)
	}
	return topics, nil
}

func (m MqttClient) SetClientMutUnSub(ctx context.Context, clientID string, topics []string) error {
	if len(topics) == 0 {
		return nil
	}
	logx.WithContext(ctx).Infof("SetClientMut clientID:%v,topics:%v", clientID, topics)
	if m.cfg.OpenApi == nil {
		return errors.System.AddMsg("没有配置秘钥")
	}
	oa := m.cfg.OpenApi
	greq := gorequest.New().Retry(1, time.Second*2)
	greq.SetBasicAuth(oa.ApiKey, oa.SecretKey)
	var ret []*MutSubReq
	var req []*MutSubReq
	for _, v := range topics {
		req = append(req, &MutSubReq{
			Topic: v,
			Qos:   0,
		})
	}
	var errs []error
	var body []byte
	var tryTime = 5
	for i := tryTime; i > 0; i-- {
		_, body, errs = greq.Post(fmt.Sprintf("%s/api/v5/clients/%s/unsubscribe/bulk", oa.Host,
			url.QueryEscape(clientID))).Send(&req).EndStruct(&ret)
		if errs != nil {
			time.Sleep(time.Second / 2 * time.Duration(tryTime) / time.Duration(i))
			continue
		}
		break
	}
	if errs != nil {
		return errors.System.AddDetail(errs, string(body))
	}
	return nil
}

// // https://www.emqx.io/docs/zh/v5.5/admin/api-docs.html#tag/Clients/paths/~1clients~1%7Bclientid%7D/get
//
//	func (m MqttClient) CheckIsOnline(ctx context.Context, clientID string) (bool, error) {
//		if m.cfg.OpenApi == nil {
//			return false, errors.System.AddMsg("没有配置秘钥")
//		}
//		oa := m.cfg.OpenApi
//		greq := gorequest.New().Retry(1, time.Second*2)
//		greq.SetBasicAuth(oa.ApiKey, oa.SecretKey)
//		var ret EmqResp
//		resp, rets, errs := greq.Get(fmt.Sprintf("%s/api/v5/clients/%s", oa.Host, url.QueryEscape(clientID))).EndStruct(&ret)
//		if errs != nil {
//			return false, errors.System.AddDetail(errs)
//		}
//		if resp.StatusCode != http.StatusOK {
//			return false, errors.System.AddDetail(string(rets))
//		}
//		if ret.Code == "" {
//			return true, nil
//		}
//		return false, nil
//	}
type EmqGetClientsData struct {
	CreatedAt      time.Time `json:"created_at"`
	Connected      bool      `json:"connected"`
	IpAddress      string    `json:"ip_address"`
	ProtoVer       int       `json:"proto_ver"`
	Mountpoint     string    `json:"mountpoint"`
	ProtoName      string    `json:"proto_name"`
	Port           int       `json:"port"`
	ConnectedAt    time.Time `json:"connected_at"`
	ExpiryInterval int       `json:"expiry_interval"`
	Username       *string   `json:"username"`
	Keepalive      int       `json:"keepalive"`
	Clientid       string    `json:"clientid"`
}

type EmqGetClientsResp struct {
	Data []EmqGetClientsData `json:"data"`
	Meta EmqGetClientsMeta   `json:"meta"`
}
type EmqGetClientsMeta struct {
	Count int64 `json:"count"`
	Limit int   `json:"limit"`
	Page  int   `json:"page"`
}

type OnlineClientsInfo struct {
	ClientID    string
	UserName    string
	ConnectedAt time.Time
}
type GetOnlineClientsFilter struct {
	UserName string
}

type PageInfo struct {
	Page int64 `json:"page" form:"page"`         // 页码
	Size int64 `json:"pageSize" form:"pageSize"` // 每页大小
}

func (m MqttClient) GetOnlineClients(ctx context.Context, f GetOnlineClientsFilter, page *PageInfo) ([]*DevConn, int64, error) {
	if m.cfg.OpenApi == nil {
		return nil, 0, errors.System.AddMsg("没有配置秘钥")
	}
	oa := m.cfg.OpenApi
	greq := gorequest.New().Retry(1, time.Second*2)
	greq.SetBasicAuth(oa.ApiKey, oa.SecretKey)
	greq.Get(fmt.Sprintf("%s/api/v5/clients", oa.Host))
	if f.UserName != "" {
		greq.Query("username=" + f.UserName)
	}
	if page != nil {
		greq.Query(fmt.Sprintf("page=%v", page.Page))
		greq.Query(fmt.Sprintf("limit=%v", page.Size))
	}
	var ret EmqGetClientsResp
	resp, rets, errs := greq.EndStruct(&ret)
	if errs != nil {
		return nil, 0, errors.System.AddDetail(errs)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, 0, errors.System.AddDetail(string(rets))
	}
	var infos []*DevConn
	for _, v := range ret.Data {
		infos = append(infos, &DevConn{
			ClientID:  v.Clientid,
			UserName:  cast.ToString(v.Username),
			Timestamp: v.ConnectedAt.UnixMilli(),
		})
	}
	return infos, ret.Meta.Count, nil
}
