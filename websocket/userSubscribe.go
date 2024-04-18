package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/utils"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
)

type UserSubscribe struct {
	store kv.Store
}

func NewUserSubscribe(store kv.Store) *UserSubscribe {
	return &UserSubscribe{
		store: store,
	}
}

func (u *UserSubscribe) genUserKey(userID int64) string {
	return fmt.Sprintf("cache:subcribe:user:%v", userID)
}

func (u *UserSubscribe) genInfoKey(info *SubscribeInfo) string {
	params, _ := json.Marshal(info.Params)
	return fmt.Sprintf("cache:subcribe:info:%v:%v", info.Code, utils.MD5V(params))
}

func (u *UserSubscribe) Add(ctx context.Context, userID int64, info *SubscribeInfo) error {
	if dp.connPool[userID] == nil {
		return errors.Parameter.AddMsg("ws未连接")
	}
	nodeStr := cast.ToString(nodeID)
	err := u.store.HsetCtx(ctx, u.genInfoKey(info), cast.ToString(userID), nodeStr)
	if err != nil {
		return err
	}
	_, err = u.store.HsetnxCtx(ctx, u.genUserKey(userID), utils.MarshalNoErr(info), nodeStr)
	return err
}
func (u *UserSubscribe) Del(ctx context.Context, userID int64, info *SubscribeInfo) error {
	u.store.HdelCtx(ctx, u.genInfoKey(info), cast.ToString(userID))
	_, err := u.store.HdelCtx(ctx, u.genUserKey(userID), utils.MarshalNoErr(info))
	return err
}

func (u *UserSubscribe) Clear(ctx context.Context, userID int64) error {
	infos, err := u.store.HkeysCtx(ctx, u.genUserKey(userID))
	if err != nil {
		return err
	}
	for _, v := range infos {
		var info SubscribeInfo
		json.Unmarshal([]byte(v), &info)
		_, err := u.store.HdelCtx(ctx, u.genInfoKey(&info), cast.ToString(userID))
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
	}
	_, err = u.store.DelCtx(ctx, u.genUserKey(userID))
	return err
}

// 返回的key是nodeid,value是userIDs
func (u *UserSubscribe) IndexInfo(ctx context.Context, info *SubscribeInfo) (map[int64][]int64, error) {
	val, err := u.store.HgetallCtx(ctx, u.genInfoKey(info))
	if err != nil {
		return nil, err
	}
	var ret = map[int64][]int64{}
	for kStr, vStr := range val {
		k := cast.ToInt64(kStr)
		v := cast.ToInt64(vStr)
		if ret[v] == nil {
			ret[v] = []int64{k}
			continue
		}
		ret[v] = append(ret[v], k)
	}
	return ret, nil
}

func (u *UserSubscribe) Index(ctx context.Context, userID int64) ([]*SubscribeInfo, error) {
	members, err := u.store.HkeysCtx(ctx, u.genUserKey(userID))
	if err != nil {
		return nil, err
	}
	var ret []*SubscribeInfo
	for _, v := range members {
		var val SubscribeInfo
		err := json.Unmarshal([]byte(v), &val)
		if err != nil {
			continue
		}
		ret = append(ret, &val)
	}
	return ret, nil
}
