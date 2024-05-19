package websocket

import (
	"context"
	"crypto/md5"
	"gitee.com/i-Things/share/utils"
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
	//err = NewUserSubscribe(store).Add(ctx, c.userID, &info)
	//if err != nil {
	//	logx.Error(err)
	//	c.errorSend(err)
	//	return
	//}
	md := utils.Md5Map(info.Params)
	logx.Infof("userSubscribe info:%v md5sum:%v", info, md)
	c.userSubscribe[md] = info.Params
	func() {
		dp.userSubscribeMutex.Lock()
		defer dp.userSubscribeMutex.Unlock()
		if _, ok := dp.userSubscribe[md]; !ok {
			dp.userSubscribe[md] = map[int64]*connection{}
		}
		dp.userSubscribe[md][c.connectID] = c
	}()
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
	md := md5.Sum([]byte(utils.MarshalNoErr(info.Params)))
	delete(c.userSubscribe, md)
	func() {
		dp.userSubscribeMutex.Lock()
		defer dp.userSubscribeMutex.Unlock()
		if _, ok := dp.userSubscribe[md]; ok {
			delete(dp.userSubscribe[md], c.connectID)
		}
	}()
	var resp WsResp
	resp.WsBody.Type = UnSubRet
	c.sendMessage(resp)
}
