package caches

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/eventBus"
	"github.com/dgraph-io/ristretto"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
	"time"
)

type CacheSyncStu struct {
	KeyType string `json:"keyType"`
}

type Cache[dataT any] struct {
	keyType    string
	cache      *ristretto.Cache
	fastEvent  *eventBus.FastEvent
	getData    func(ctx context.Context, key string) (*dataT, error)
	expireTime time.Duration
}

type CacheConfig[dataT any] struct {
	KeyType    string
	FastEvent  *eventBus.FastEvent
	GetData    func(ctx context.Context, key string) (*dataT, error)
	ExpireTime time.Duration
}

var (
	cacheMap   = map[string]any{}
	cacheMutex sync.Mutex
)

func NewCache[dataT any](cfg CacheConfig[dataT]) (*Cache[dataT], error) {
	cacheMutex.Lock() //单例模式
	defer cacheMutex.Unlock()
	if v, ok := cacheMap[cfg.KeyType]; ok {
		return v.(*Cache[dataT]), nil
	}
	cache, _ := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	ret := Cache[dataT]{
		keyType:    cfg.KeyType,
		cache:      cache,
		fastEvent:  cfg.FastEvent,
		getData:    cfg.GetData,
		expireTime: cfg.ExpireTime,
	}
	err := ret.fastEvent.Subscribe(ret.genTopic(), func(ctx context.Context, t time.Time, body []byte) error {
		cacheKey := string(body)
		ret.cache.Del(cacheKey)
		return nil
	})
	if err != nil {
		return nil, err
	}
	cacheMap[cfg.KeyType] = &ret
	return &ret, nil
}

func (c Cache[dataT]) genTopic() string {
	return fmt.Sprintf(eventBus.ServerCacheSync, c.keyType)
}

func (c Cache[dataT]) genCacheKey(key string) string {
	return fmt.Sprintf("cache:%s:%s", c.keyType, key)
}

// 删除数据的时候设置为空即可
func (c Cache[dataT]) SetData(ctx context.Context, key string, data *dataT) error {
	cacheKey := c.genCacheKey(key)
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
	c.cache.Del(cacheKey)
	err := c.fastEvent.Publish(ctx, c.genTopic(), cacheKey)
	if err != nil {
		logx.WithContext(ctx).Error(err)
	}
	return nil
}

func (c Cache[dataT]) GetData(ctx context.Context, key string) (*dataT, error) {
	cacheKey := c.genCacheKey(key)
	temp, ok := c.cache.Get(cacheKey)
	if ok {
		if temp == nil {
			return nil, errors.NotFind
		}
		return temp.(*dataT), nil
	}
	{ //内存中没有就从redis上获取
		val, err := store.GetCtx(ctx, cacheKey)
		if err != nil {
			return nil, err
		}
		if len(val) > 0 {
			var ret dataT
			err = json.Unmarshal([]byte(val), &ret)
			if err != nil {
				return nil, err
			}
			c.cache.SetWithTTL(cacheKey, &ret, 1, c.expireTime)
			return &ret, nil
		}
	}
	if c.getData == nil { //如果没有设置第三级缓存则直接设置该参数为空并返回
		c.cache.SetWithTTL(cacheKey, nil, 1, c.expireTime)
		return nil, nil
	}
	//redis上没有就读数据库
	data, err := c.getData(ctx, key)
	if err != nil && !errors.Cmp(err, errors.NotFind) { //如果是其他错误直接返回
		return nil, err
	}
	//读到了之后设置缓存
	c.cache.SetWithTTL(cacheKey, data, 1, c.expireTime)
	if data == nil {
		return data, nil
	}
	ctxs.GoNewCtx(ctx, func(ctx2 context.Context) { //异步设置缓存
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
	return data, nil
}
