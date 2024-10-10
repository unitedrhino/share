package barrier

import (
	"context"
	"database/sql"
	"gitee.com/unitedrhino/share/stores"
	"github.com/dtm-labs/client/dtmgrpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 屏障分布式事务
func BarrierTransaction(ctx context.Context, fc func(tx *gorm.DB) error) error {
	conn := stores.GetCommonConn(ctx)
	barrier, _ := dtmgrpc.BarrierFromGrpc(ctx)
	if barrier == nil { //如果没有开启分布式事务,则直接走普通事务即可
		return conn.Transaction(fc)
	}
	db, _ := conn.DB()
	return barrier.CallWithDB(db, func(tx *sql.Tx) error {
		gdb, _ := gorm.Open(mysql.New(mysql.Config{
			Conn: tx,
		}), &gorm.Config{})
		return fc(gdb)
	})
}
