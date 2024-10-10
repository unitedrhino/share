package wxClient

import (
	"context"
	"gitee.com/unitedrhino/share/conf"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/credential"
	"github.com/silenceper/wechat/v2/miniprogram"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	zeroCache "github.com/zeromicro/go-zero/core/stores/cache"
)

type MiniProgram = miniprogram.MiniProgram

func NewWxMiniProgram(ctx context.Context, conf *conf.ThirdConf, redisConf zeroCache.ClusterConf) (*MiniProgram, error) {
	if conf == nil {
		return nil, nil
	}
	wc := wechat.NewWechat()
	memory := cache.NewRedis(ctx, &cache.RedisOpts{
		Host:     redisConf[0].Host,
		Password: redisConf[0].Pass,
	})
	cfg := &miniConfig.Config{
		AppID:     conf.AppID,
		AppSecret: conf.AppSecret,
		Cache:     memory,
	}
	program := wc.GetMiniProgram(cfg)
	program.SetAccessTokenHandle(credential.NewStableAccessToken(cfg.AppID, cfg.AppSecret, credential.CacheKeyMiniProgramPrefix, cfg.Cache))
	return program, nil
}

type WxOfficialAccount = officialaccount.OfficialAccount

func NewWxOfficialAccount(ctx context.Context, conf *conf.ThirdConf, redisConf zeroCache.ClusterConf) (*WxOfficialAccount, error) {
	if conf == nil {
		return nil, nil
	}
	wc := wechat.NewWechat()
	memory := cache.NewRedis(ctx, &cache.RedisOpts{
		Host:     redisConf[0].Host,
		Password: redisConf[0].Pass,
	})
	cfg := &offConfig.Config{
		AppID:     conf.AppID,
		AppSecret: conf.AppSecret,
		Cache:     memory,
	}
	program := wc.GetOfficialAccount(cfg)
	program.SetAccessTokenHandle(credential.NewStableAccessToken(cfg.AppID, cfg.AppSecret, credential.CacheKeyMiniProgramPrefix, cfg.Cache))
	return program, nil
}
