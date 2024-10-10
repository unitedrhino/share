package stores

import (
	"context"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"math/rand"
	"sync"
	"time"
)

type AsyncInsert[t any] struct {
	once       sync.Once
	insertChan chan *t
}

const (
	asyncExecMax = 200 //异步执行sql最大数量
	asyncRunMax  = 40
)

func NewAsyncInsert[t any]() (a *AsyncInsert[t]) {
	a = &AsyncInsert[t]{
		insertChan: make(chan *t, asyncExecMax*10),
	}
	for i := 0; i < asyncRunMax; i++ {
		utils.Go(context.Background(), func() {
			a.asyncInsertRuntime()
		})
	}
	return a
}

func (a *AsyncInsert[t]) AsyncInsert(stu *t) {
	a.insertChan <- stu
}

func (a *AsyncInsert[t]) asyncInsertRuntime() {
	r := rand.Intn(1000)
	tick := time.Tick(time.Second/2 + time.Millisecond*time.Duration(r))
	execCache := make([]*t, 0, asyncExecMax*2)
	exec := func() {
		if len(execCache) == 0 {
			return
		}
		err := GetTenantConn(ctxs.WithRoot(context.Background())).CreateInBatches(execCache, 100).Error
		if err != nil {
			logx.Error(err)
		}
		execCache = execCache[0:0] //清空切片
	}
	for {
		select {
		case _ = <-tick:
			exec()
		case e := <-a.insertChan:
			execCache = append(execCache, e)
			if len(execCache) > asyncExecMax {
				exec()
			}
		}
	}
}
