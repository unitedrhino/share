package utils

//
//import (
//	"gitee.com/unitedrhino/share/utils/notify"
//	"sync/atomic"
//)
//
//var (
//	startCount = atomic.Int64{}
//)
//
//func Start() {
//	startCount.Add(1)
//}
//
//func Ready() {
//	if startCount.Add(-1) <= 0 { //全部服务就绪
//		notify.Ready()
//	}
//}
//
//func Stopping() {
//	notify.Stopping()
//}
