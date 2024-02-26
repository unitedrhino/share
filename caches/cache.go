package caches

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/i-Things/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

type Cache[dataT any] struct {
	kepPrefix string
}

func NewCache[dataT any](keyPrefix string) Cache[dataT] {
	return Cache[dataT]{kepPrefix: keyPrefix}
}

func (c Cache[dataT]) genCacheKey(key string) string {
	return fmt.Sprintf("cache:%s:%s", c.kepPrefix, key)
}

func (c Cache[dataT]) SetData(ctx context.Context, data *dataT, keys ...any) error {
	dataStr, err := json.Marshal(data)
	if err != nil {
		return err
	}
	for _, key := range keys {
		err := store.SetexCtx(ctx, c.genCacheKey(utils.Fmt(key)), string(dataStr), 10*60)
		if err != nil {
			logx.WithContext(ctx).Error(err)
		}
	}
	return nil
}

func (c Cache[dataT]) GetData(ctx context.Context, key any) (*dataT, error) {
	val, err := store.GetCtx(ctx, c.genCacheKey(utils.Fmt(key)))
	if err != nil {
		return nil, err
	}
	var ret dataT
	err = json.Unmarshal([]byte(val), &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
