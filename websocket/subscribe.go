package websocket

import (
	"context"
	"gitee.com/i-Things/share/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/zeromicro/go-zero/core/logx"
)

func subscribeHandle(ctx context.Context, c *connection, body WsReq) {
	var info SubscribeInfo
	err := mapstructure.Decode(body.Body, &info)
	if err != nil {
		logx.Errorf("userSubscribe Decode body:%v err:%v", utils.Fmt(body), err)
		c.errorSend(err)
		return
	}
	if checkSubscribe != nil {
		err = checkSubscribe(ctx, &info)
		if err != nil {
			logx.Errorf("userSubscribe checkSubscribe body:%v err:%v", utils.Fmt(body), err)
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
	key := info.Code + ":" + md
	logx.Infof("userSubscribe info:%v key:%v", info, key)
	c.userSubscribe[key] = info.Params
	func() {
		dp.userSubscribeMutex.Lock()
		defer dp.userSubscribeMutex.Unlock()
		if _, ok := dp.userSubscribe[key]; !ok {
			dp.userSubscribe[key] = map[int64]*connection{}
		}
		dp.userSubscribe[key][c.connectID] = c
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
	//err = NewUserSubscribe(store).Del(ctx, c.userID, &info)
	if err != nil {
		logx.Error(err)
		c.errorSend(err)
		return
	}
	md := utils.Md5Map(info.Params)
	key := info.Code + ":" + md
	delete(c.userSubscribe, key)
	func() {
		dp.userSubscribeMutex.Lock()
		defer dp.userSubscribeMutex.Unlock()
		if _, ok := dp.userSubscribe[key]; ok {
			delete(dp.userSubscribe[key], c.connectID)
		}
	}()
	var resp WsResp
	resp.WsBody.Type = UnSubRet
	c.sendMessage(resp)
}
