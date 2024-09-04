package clients

import (
	"context"
	"database/sql"
	"fmt"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/utils"
	_ "github.com/i-Things/driver-go/v3/taosRestful"
	"math/rand"
	"strings"
	"time"

	//tdengine 的cgo模式，这个模式是最快的，需要可以打开
	//_ "github.com/i-Things/driver-go/v3/taosSql"
	_ "github.com/i-Things/driver-go/v3/taosWS"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
)

type Td struct {
	*sql.DB
}

type ExecArgs struct {
	Query string
	Args  []any
}

var (
	td         = Td{}
	once       = sync.Once{}
	insertChan = make(chan ExecArgs, 1000)
)

const (
	asyncExecMax = 200 //异步执行sql最大数量
	asyncRunMax  = 40
)

func NewTDengine(DataSource conf.TSDB) (TD *Td, err error) {
	once.Do(func() {
		if DataSource.Driver == "" {
			DataSource.Driver = "taosWS"
		}
		td.DB, err = sql.Open(DataSource.Driver, DataSource.DSN)
		if err != nil {
			return
		}
		td.DB.SetMaxIdleConns(50)
		td.DB.SetMaxOpenConns(50)
		td.DB.SetConnMaxIdleTime(time.Hour)
		td.DB.SetConnMaxLifetime(time.Hour)
		_, err = td.Exec("create database if not exists ithings;")
		if err != nil {
			return
		}
		for i := 0; i < asyncRunMax; i++ {
			utils.Go(context.Background(), func() {
				td.asyncInsertRuntime()
			})
		}
	})
	if err != nil {
		logx.Errorf("tdengine 初始化失败,err:%v", err)
	}
	return &td, err
}

func (t *Td) asyncInsertRuntime() {
	r := rand.Intn(1000)
	tick := time.Tick(time.Second/2 + time.Millisecond*time.Duration(r))
	execCache := make([]ExecArgs, 0, asyncExecMax*2)
	exec := func() {
		if len(execCache) == 0 {
			return
		}
		sql, args := t.genInsertSql(execCache...)
		var err error
		for i := 3; i > 0; i-- { //三次重试
			_, err = t.Exec(sql, args...)
			if err == nil {
				break
			}
		}
		if err != nil {
			logx.Error(err)
		}
		execCache = execCache[0:0] //清空切片
	}
	for {
		select {
		case _ = <-tick:
			exec()
		case e := <-insertChan:
			execCache = append(execCache, e)
			if len(execCache) > asyncExecMax {
				exec()
			}
		}
	}

}

func (t *Td) AsyncInsert(query string, args ...any) {
	insertChan <- ExecArgs{
		Query: query,
		Args:  args,
	}
}
func (t *Td) genInsertSql(eas ...ExecArgs) (query string, args []any) {
	qs := make([]string, 0, len(eas))
	as := make([]any, 0, len(eas))
	for _, e := range eas {
		qs = append(qs, e.Query)
		as = append(as, e.Args...)
	}
	return fmt.Sprintf("insert into %s;", strings.Join(qs, " ")), as
}
