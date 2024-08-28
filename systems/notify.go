package systems

import (
	"context"
	"fmt"
	"gitee.com/i-Things/share/clients/dingClient"
	"gitee.com/i-Things/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
	"os"
	"runtime"
	"runtime/debug"
)

func init() {
	token := os.Getenv("dingRobotToken")
	if token != "" {
		c := dingClient.NewDingRobotClient(token)
		utils.SetPanicNotify(func(s string) {
			c.SendRobotMsg(dingClient.NewTextMessage("抓到panic:" + s))
		})
		proc.AddShutdownListener(func() {
			e, _ := os.Executable()
			c.SendRobotMsg(dingClient.NewTextMessage(fmt.Sprintf("iThings程序退出:%v", e)))
		})
	} else {
		proc.AddShutdownListener(func() {
			pc := make([]uintptr, 1)
			runtime.Callers(3, pc)
			f := runtime.FuncForPC(pc[0])
			msg := fmt.Sprintf("Shutdown|func=%s|stack=%s\n", f, string(debug.Stack()))
			logx.Error(msg)
		})
	}
}

func SysNotify(in string) {
	token := os.Getenv("dingRobotToken")
	if token != "" {
		ctx := context.Background()
		utils.Go(ctx, func() {
			c := dingClient.NewDingRobotClient(token)
			c.SendRobotMsg(dingClient.NewTextMessage("iThings系统通知:" + in))
		})
	}
}
