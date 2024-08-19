package systems

import (
	"context"
	"fmt"
	"gitee.com/i-Things/share/clients/dingClient"
	"gitee.com/i-Things/share/utils"
	"github.com/zeromicro/go-zero/core/proc"
	"os"
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
