package clients

import (
	"context"
	"database/sql"
	"fmt"
	_ "gitee.com/unitedrhino/driver-go/v3/taosRestful"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/stores"
	"gitee.com/unitedrhino/share/utils"
	"github.com/spf13/cast"
	"math/rand"
	"strings"
	"sync/atomic"
	"time"

	//tdengine 的cgo模式，这个模式是最快的，需要可以打开
	//_ "gitee.com/unitedrhino/driver-go/v3/taosSql"
	_ "gitee.com/unitedrhino/driver-go/v3/taosWS"
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
	td                = Td{}
	once              = sync.Once{}
	insertChan        = make(chan ExecArgs, 1000)
	insertNoDebugChan = make(chan ExecArgs, 1000)
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
		utils.Go(context.Background(), func() {
			td.countSql()
		})
		for i := 0; i < asyncRunMax; i++ {
			id := int64(i)
			utils.Go(context.Background(), func() {
				td.asyncInsertRuntime(id)
			})
			utils.Go(context.Background(), func() {
				td.asyncInsertNoDebugRuntime(id)
			})
		}
	})
	if err != nil {
		logx.Errorf("tdengine 初始化失败,err:%v", err)
	}
	return &td, err
}

var sendCount atomic.Int64

func (t *Td) countSql() {
	tick := time.Tick(time.Minute)
	for {
		select {
		case <-tick:
			e := sendCount.Swap(0)
			if e != 0 {
				logx.Infof("tdengineRuntimeCountSql %v/mim ", sendCount.Swap(0))
			}
		}
	}
}

func (t *Td) asyncInsertRuntime(id int64) {
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
			logx.Error(sql, args, err)
		}
		sendCount.Add(int64(len(execCache)))
		execCache = execCache[0:0] //清空切片
	}
	for {
		select {
		case _ = <-tick:
			exec()
		case e := <-insertChan:
			execCache = append(execCache, e)
			if len(execCache) > asyncExecMax {
				logx.Infof("tdengineRuntime id:%v, exec to much now num:%v", id, len(execCache))
				exec()
			}
		}
	}

}

func (t *Td) asyncInsertNoDebugRuntime(id int64) {
	r := rand.Intn(1000)
	tick := time.Tick(time.Second + time.Millisecond*time.Duration(r))
	execCache := make([]ExecArgs, 0, asyncExecMax*2)
	exec := func() {
		if len(execCache) == 0 {
			return
		}
		sql, args := t.genInsertSql(execCache...)
		var err error
		for i := 3; i > 0; i-- { //三次重试
			_, err = t.ExecContext(stores.SetIsDebug(context.Background(), false), sql, args...)
			if err == nil {
				break
			}
		}
		if err != nil {
			logx.Error(sql, args, err)
		}
		sendCount.Add(int64(len(execCache)))
		execCache = execCache[0:0] //清空切片
	}
	for {
		select {
		case _ = <-tick:
			exec()
		case e := <-insertNoDebugChan:
			execCache = append(execCache, e)
			if len(execCache) > asyncExecMax {
				logx.Infof("tdengineRuntime id:%v, exec to much now num:%v", id, len(execCache))
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

func (t *Td) AsyncInsertNoDebug(query string, args ...any) {
	insertNoDebugChan <- ExecArgs{
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

type TableDesc struct {
	Field    string
	Type     string
	Length   int64
	Note     string
	Encode   string
	Compress string
	Level    string
}

func (t *Td) STables(ctx context.Context, like string) ([]string, error) {
	rows, err := t.QueryContext(ctx, fmt.Sprintf("show stables like '%%%s%%';", like))
	if err != nil {
		return nil, errors.Fmt(err)
	}
	var datas []map[string]any
	err = stores.Scan(rows, &datas)
	if err != nil {
		return nil, errors.Fmt(err)
	}
	var ret []string
	for _, v := range datas {
		ret = append(ret, v["stable_name"].(string))
	}
	return ret, nil
}

func (t *Td) Tables(ctx context.Context, like string) ([]string, error) {
	rows, err := t.QueryContext(ctx, fmt.Sprintf("show tables like '%%%s%%';", like))
	if err != nil {
		return nil, errors.Fmt(err)
	}
	var datas []map[string]any
	err = stores.Scan(rows, &datas)
	if err != nil {
		return nil, errors.Fmt(err)
	}
	var ret []string
	for _, v := range datas {
		ret = append(ret, v["table_name"].(string))
	}
	return ret, nil
}

func (t *Td) TableTags(ctx context.Context, tbName string) (map[string]string, error) {

	rows, err := t.QueryContext(ctx, fmt.Sprintf("SHOW  TAGS FROM  `%s`;", tbName))
	if err != nil {
		return nil, errors.Fmt(err)
	}
	var datas []map[string]any
	err = stores.Scan(rows, &datas)
	if err != nil {
		return nil, errors.Fmt(err)
	}
	var ret = map[string]string{}
	for _, v := range datas {
		ret[v["tag_name"].(string)] = cast.ToString(v["tag_value"])
	}
	return ret, nil
}

func (t *Td) DescTable(ctx context.Context, tbName string) (map[string]*TableDesc, error) {
	rows, err := t.QueryContext(ctx, fmt.Sprintf("desc `%s`", tbName))
	if err != nil {
		return nil, errors.Fmt(err)
	}
	var datas []map[string]any
	err = stores.Scan(rows, &datas)
	if err != nil {
		return nil, errors.Fmt(err)
	}
	var ret = map[string]*TableDesc{}
	for _, d := range datas {
		ret[cast.ToString(d["field"])] = &TableDesc{
			Field:    cast.ToString(d["field"]),
			Type:     cast.ToString(d["type"]),
			Length:   cast.ToInt64(d["length"]),
			Note:     cast.ToString(d["note"]),
			Encode:   cast.ToString(d["encode"]),
			Compress: cast.ToString(d["compress"]),
			Level:    cast.ToString(d["level"]),
		}
	}
	return ret, nil
}
