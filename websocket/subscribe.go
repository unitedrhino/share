package websocket

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/zeromicro/go-zero/core/logx"
)

func subscribeHandle(ctx context.Context, c *connection, body WsReq) {
	var info SubscribeInfo
	err := mapstructure.Decode(body.Body, &info)
	if err != nil {
		logx.Error(err)
		c.errorSend(err)
		return
	}
	if checkSubscribe != nil {
		err = checkSubscribe(ctx, &info)
		if err != nil {
			logx.Error(err)
			c.errorSend(err)
			return
		}
	}
	err = NewUserSubscribe(store).Add(ctx, c.userID, &info)
	if err != nil {
		logx.Error(err)
		c.errorSend(err)
		return
	}
	var resp WsResp
	resp.WsBody.Type = SubRet
	c.sendMessage(resp)
}

func unSubscribeHandle(ctx context.Context, c *connection, body WsReq) {
	var info SubscribeInfo
	err := mapstructure.Decode(body.Body, &info)
	if err != nil {
		logx.Error(err)
		c.errorSend(err)
		return
	}
	err = NewUserSubscribe(store).Del(ctx, c.userID, &info)
	if err != nil {
		logx.Error(err)
		c.errorSend(err)
		return
	}
	var resp WsResp
	resp.WsBody.Type = UnSubRet
	c.sendMessage(resp)

}
