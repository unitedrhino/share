package tools

import (
	"context"
	"fmt"
	"gitee.com/unitedrhino/share/caches"
	"gitee.com/unitedrhino/share/conf"
	"github.com/spf13/cast"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
)

type Limit struct {
	Conf  []conf.Limit
	Type  string //验证的类型 如账号 ip
	Scene string //使用场景,如登录 验证码
	store kv.Store
}

func NewLimit(c []conf.Limit, scene string, t string, defaultC []conf.Limit) *Limit {
	if len(c) == 0 {
		c = defaultC
	}
	return &Limit{
		Conf:  c,
		Type:  t,
		Scene: scene,
		store: caches.GetStore(),
	}
}

func (l Limit) genCountKey(i int, key string) string {
	return fmt.Sprintf("limit:%s:%s:count:%v:%s", l.Scene, l.Type, i, key)
}

func (l Limit) genForbiddenKey(key string) string {
	return fmt.Sprintf("limit:%s:%s:forbidden:%s", l.Scene, l.Type, key)
}

func (l Limit) CheckLimit(ctx context.Context, key string) bool {
	cacheKey := l.genForbiddenKey(key)
	ret, err := l.store.GetCtx(ctx, cacheKey)
	if err != nil {
		return false
	}
	if ret != "" {
		return true
	}
	return false
}

func (l Limit) CleanLimit(ctx context.Context, key string) error {
	cacheKey := l.genForbiddenKey(key)
	_, err := l.store.DelCtx(ctx, cacheKey)
	if err != nil {
		logx.WithContext(ctx).Error(err)
	}
	return err
}

// 错误了之后限制这个操作
func (l Limit) LimitIt(ctx context.Context, key string) error {
	for i, v := range l.Conf {
		cacheKey := l.genCountKey(i, key)
		ret, err := l.store.GetCtx(ctx, cacheKey)
		if ret == "" {
			err = l.store.SetexCtx(ctx, cacheKey, "1", v.Timeout)
			if err != nil {
				logx.WithContext(ctx).Error(err)
				return err
			}
			continue
		}
		//错误加一次
		_, err = l.store.Incr(cacheKey)
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
		//重置key的过期时间
		err = l.store.ExpireCtx(ctx, cacheKey, v.Timeout)
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
		//如果达到封禁的次数,则添加封禁
		if cast.ToInt(ret)+1 >= v.TriggerTime {
			cacheKey := l.genForbiddenKey(key)
			err = l.store.SetexCtx(ctx, cacheKey, cast.ToString(v.ForbiddenTime), v.ForbiddenTime)
			if err != nil {
				logx.WithContext(ctx).Error(err)
			}
		}
		return nil
	}

	return nil
}
