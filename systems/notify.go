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

var token = os.Getenv("dingRobotToken")

func init() {
	if token != "" {
		c := dingClient.NewDingRobotClient(token)
		utils.SetPanicNotify(func(s string) {
			_, err := c.SendRobotMsg(dingClient.NewTextMessage("程序抓到panic:" + s))
			if err != nil {
				logx.Error(err)
			}

		})
		e, _ := os.Executable()
		proc.AddShutdownListener(func() {
			_, err := c.SendRobotMsg(dingClient.NewTextMessage(fmt.Sprintf("程序退出:%v", e)))
			if err != nil {
				logx.Error(err)
			}
		})
		_, err := c.SendRobotMsg(dingClient.NewTextMessage(fmt.Sprintf("程序启动:%v", e)))
		if err != nil {
			logx.Error(err)
		}
	} else {
		proc.AddShutdownListener(func() {
			pc := make([]uintptr, 1)
			runtime.Callers(3, pc)
			f := runtime.FuncForPC(pc[0])
			msg := fmt.Sprintf("程序Shutdown|func=%s|stack=%s\n", f, string(debug.Stack()))
			logx.Error(msg)
		})
	}
}

func SysNotify(in string) {
	if token != "" {
		ctx := context.Background()
		utils.Go(ctx, func() {
			c := dingClient.NewDingRobotClient(token)
			_, err := c.SendRobotMsg(dingClient.NewTextMessage("程序系统通知:" + in))
			if err != nil {
				logx.Error(err)
			}
		})
	}
}
