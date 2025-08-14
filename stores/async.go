package stores

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm/clause"
)

type AsyncInsert[t any] struct {
	db         *DB
	once       sync.Once
	insertChan chan *t
	tableName  string
}

const (
	asyncExecMax = 200 //异步执行sql最大数量
	asyncRunMax  = 40
)

func NewAsyncInsert[t any](db *DB, tableName string) (a *AsyncInsert[t]) {
	a = &AsyncInsert[t]{
		insertChan: make(chan *t, asyncExecMax*10),
		db:         db,
		tableName:  tableName,
	}
	for i := 0; i < asyncRunMax; i++ {
		utils.Go(context.Background(), func() {
			a.asyncInsertRuntime()
		})
	}
	return a
}

func (a *AsyncInsert[t]) AsyncInsert(stu *t) {
	if a == nil {
		return
	}
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
		db := a.db.WithContext(ctxs.WithRoot(context.Background()))
		if a.tableName != "" {
			db = db.Table(a.tableName)
		}
		err := db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(execCache, 100).Error
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
