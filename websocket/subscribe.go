package websocket

import (
	"context"
	"gitee.com/unitedrhino/share/utils"
	"github.com/mitchellh/mapstructure"
	"github.com/zeromicro/go-zero/core/logx"
)

func subscribeHandle(ctx context.Context, c *connection, body WsReq) {
	var info SubscribeInfo
	err := mapstructure.Decode(body.Body, &info)
	if err != nil {
		logx.WithContext(ctx).Errorf("websocket userSubscribe Decode userID:%v connectID:%v body:%v err:%v",
			c.userID, c.connectID, utils.Fmt(body), err)
		c.errorSend(err)
		return
	}
	if checkSubscribe != nil {
		err = checkSubscribe(ctx, &info)
		if err != nil {
			logx.WithContext(ctx).Errorf("websocket userSubscribe checkSubscribe userID:%v connectID:%v body:%v err:%v",
				c.userID, c.connectID, utils.Fmt(body), err)
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
	c.userSubscribe[key] = info.Params
	logx.WithContext(ctx).Infof("websocket userSubscribe userID:%v connectID:%v info:%v key:%v subList:%v",
		c.userID, c.connectID, utils.Fmt(info), key, utils.Fmt(c.userSubscribe))
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
