package utils

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"time"
)

func SingletonRun(ctx context.Context, store kv.Store, singletonKey string, f func(ctx2 context.Context)) {
	key := "singleton:" + singletonKey
	for { //定时任务为单例执行模式,有效期15秒,如果服务挂了,其他服务每隔6秒检测到就抢到执行
		ok, err := store.SetnxExCtx(ctx, key, time.Now().Format("2006-01-02 15:04:05.999"), 15)
		if err != nil {
			logx.WithContext(ctx).Errorf("%s.Store.SetnxExCtx singletonKey:%v err:%v", FuncName(), key, err)
			time.Sleep(time.Second * 6)
			continue
		}
		if ok { //抢到锁了
			break
		}
		//没抢到锁,6秒钟后继续
		time.Sleep(time.Second * 6)
	}
	logx.WithContext(ctx).Infof("SingletonRun start running singletonKey:%v", key)
	//抢到锁需要维系锁
	//每隔6秒刷新锁,如果服务挂了,锁才能退出
	Go(ctx, func() {
		defer Recover(ctx)
		ticker := time.NewTicker(time.Second * 6)
		for range ticker.C {
			for i := 0; i < 3; i++ { //如果超过三次都没有设置成功,则应该是redis有问题了,其他服务也没法注册成功
				err := store.SetexCtx(ctx, key, time.Now().Format("2006-01-02 15:04:05.999"), 15)
				if err == nil {
					break
				}
			}
		}
	})
	f(ctx)
}
