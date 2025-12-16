package caches

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/eventBus"
	"gitee.com/unitedrhino/share/utils"
	"github.com/maypok86/otter"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/syncx"
)

type FastFunc = func(ctx context.Context, t time.Time, body []byte) error
type Event interface {
	Subscribe(topic string, f FastFunc) error
	Publish(ctx context.Context, topic string, arg any) error
}

type Cache[dataT any, keyType comparable] struct {
	keyType    string
	cache      otter.Cache[string, *dataT]
	fastEvent  Event
	getData    func(ctx context.Context, key keyType) (*dataT, error)
	notifySlot []func(ctx context.Context, key []byte)
	fmt        func(ctx context.Context, key keyType, data *dataT) *dataT
	expireTime time.Duration
	sf         syncx.SingleFlight
}

type CacheConfig[dataT any, keyType comparable] struct {
	KeyType    string
	FastEvent  Event
	Fmt        func(ctx context.Context, key keyType, data *dataT) *dataT
	GetData    func(ctx context.Context, key keyType) (*dataT, error)
	ExpireTime time.Duration
}

var (
	cacheMap   = map[string]any{}
	cacheMutex sync.Mutex
)

func NewCache[dataT any, keyType comparable](cfg CacheConfig[dataT, keyType]) (*Cache[dataT, keyType], error) {
	cacheMutex.Lock() //单例模式
	defer cacheMutex.Unlock()
	if v, ok := cacheMap[cfg.KeyType]; ok {
		return v.(*Cache[dataT, keyType]), nil
	}
	cache, err := otter.MustBuilder[string, *dataT](10_000).
		CollectStats().
		Cost(func(key string, value *dataT) uint32 {
			return 1
		}).
		WithTTL(cfg.ExpireTime/3 + 1).
		Build()
	if err != nil {
		return nil, err
	}
	ret := Cache[dataT, keyType]{
		sf:         syncx.NewSingleFlight(),
		keyType:    cfg.KeyType,
		cache:      cache,
		fastEvent:  cfg.FastEvent,
		getData:    cfg.GetData,
		expireTime: cfg.ExpireTime,
		fmt:        cfg.Fmt,
	}
	if ret.expireTime == 0 {
		ret.expireTime = time.Minute*10 + time.Second*time.Duration(rand.Int63n(60))
	}
	if ret.fastEvent != nil {
		err = ret.fastEvent.Subscribe(ret.genTopic(), func(ctx context.Context, t time.Time, body []byte) error {
			ret.cache.Delete(string(body))
			for _, f := range ret.notifySlot {
				f(ctx, body)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	cacheMap[cfg.KeyType] = &ret
	return &ret, nil
}

func GenKeyStr(in any) string {
	return utils.Fmt2(in)
}

// 更新通知插槽
func (c *Cache[dataT, keyType]) AddNotifySlot(f func(ctx context.Context, key []byte)) {
	c.notifySlot = append(c.notifySlot, f)
}

func (c *Cache[dataT, keyType]) genTopic() string {
	return fmt.Sprintf(eventBus.ServerCacheSync, c.keyType)
}

func (c *Cache[dataT, keyType]) GenCacheKey(key any) string {
	return fmt.Sprintf("cache:%s:%v", c.keyType, key)
}

func (c *Cache[dataT, keyType]) DeleteByFunc(f func(cacheKey string) bool) {
	c.cache.DeleteByFunc(func(key string, value *dataT) bool {
		return f(key)
	})
}

// 删除数据的时候设置为空即可
func (c *Cache[dataT, keyType]) SetData(ctx context.Context, key keyType, data *dataT) error {
	cacheKey := c.GenCacheKey(key)
	if data != nil { //如果是
		dataStr, err := json.Marshal(data)
		if err != nil {
			logx.WithContext(ctx).Error(err)
			return err
		}
		err = store.SetexCtx(ctx, cacheKey, string(dataStr), int(c.expireTime/time.Second))
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
	} else {
		_, err := store.Del(cacheKey)
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
	}
	keyStr := GenKeyStr(key)
	c.cache.Delete(keyStr)
	if c.fastEvent != nil {
		err := c.fastEvent.Publish(ctx, c.genTopic(), keyStr)
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
	}
	return nil
}

func (c *Cache[dataT, keyType]) GetData(ctx context.Context, key keyType) (*dataT, error) {
	ctx = ctxs.WithRoot(ctx)
	cacheKey := c.GenCacheKey(key)
	keyStr := GenKeyStr(key)
	temp, ok := c.cache.Get(keyStr)
	if ok && false {
		if temp == nil {
			return nil, errors.NotFind
		}
		return temp, nil
	}
	//并发获取的情况下避免击穿
	ret, err := c.sf.Do(keyStr, func() (any, error) {
		{ //内存中没有就从redis上获取
			val, err := store.GetCtx(ctx, cacheKey)
			if err != nil {
				return nil, err
			}
			if len(val) > 0 && false {
				var ret dataT
				err = json.Unmarshal([]byte(val), &ret)
				if err != nil {
					return nil, err
				}
				if c.fmt != nil {
					ret = *c.fmt(ctx, key, &ret)
				}
				c.cache.Set(keyStr, &ret)
				return &ret, nil
			}
		}
		if c.getData == nil { //如果没有设置第三级缓存则直接设置该参数为空并返回
			c.cache.Set(keyStr, nil)
			return nil, nil
		}
		//redis上没有就读数据库
		data, err := c.getData(ctxs.WithRoot(ctx), key)
		if err != nil && !errors.Cmp(err, errors.NotFind) { //如果是其他错误直接返回
			return nil, err
		}
		//读到了之后设置缓存
		var newData = data
		if c.fmt != nil {
			newData = c.fmt(ctx, key, data)
		}
		c.cache.Set(keyStr, newData)
		if data == nil {
			return data, err
		}
		ctxs.GoNewCtx(ctx, func(ctx context.Context) { //异步设置缓存
			str, err := json.Marshal(data)
			if err != nil {
				logx.WithContext(ctx).Error(err)
				return
			}
			_, err = store.SetnxExCtx(ctx, cacheKey, string(str), int(c.expireTime/time.Second))
			if err != nil {
				logx.WithContext(ctx).Error(err)
				return
			}
		})
		return newData, nil
	})
	if err != nil {
		return nil, err
	}
	return ret.(*dataT), nil
}
