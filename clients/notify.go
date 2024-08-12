package clients

import (
	"gitee.com/i-Things/share/utils"
	"os"
)

func init() {
	token := os.Getenv("dingRobotToken")
	if token != "" {
		c := NewDingRobotClient(token)
		utils.SetPanicNotify(func(s string) {
			c.SendRobotMsg(NewTextMessage("抓到panic:" + s))
		})
	}
}

func SysNotify(in string) {
	token := os.Getenv("dingRobotToken")
	if token != "" {
		c := NewDingRobotClient(token)
		c.SendRobotMsg(NewTextMessage("系统通知:" + in))

	}
}
