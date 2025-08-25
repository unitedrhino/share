package events

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/netx"
)

// NatsHandleFunc NATS 专用的处理函数类型，包含 NATS 消息对象
type NatsHandleFunc func(ctx context.Context, ts time.Time, msg []byte, natsMsg *nats.Msg) error

func NatsSubWithType[msgType any](handle func(ctx context.Context, msgIn msgType, natsMsg *nats.Msg) error) func(msg *nats.Msg) {
	return NatsSubscription(func(ctx context.Context, ts time.Time, msg []byte, natsMsg *nats.Msg) error {
		var tempInfo msgType
		err := json.Unmarshal(msg, &tempInfo)
		if err != nil {
			return err
		}
		return handle(ctx, tempInfo, natsMsg)
	})
}

func NatsSubscription(handle NatsHandleFunc) func(msg *nats.Msg) {
	return func(msg *nats.Msg) {
		msg.Ack()
		utils.Go(context.Background(), func() {
			var ctx context.Context
			utils.Recover(ctx)
			startTime := time.Now()
			emsg := GetEventMsg(msg.Data)
			if emsg == nil {
				logx.Error(msg.Subject, string(msg.Data))
				return
			}
			ctx = emsg.GetCtx()
			ctx, span := ctxs.StartSpan(ctx, msg.Subject, "")
			defer span.End()

			err := handle(ctx, emsg.GetTs(), emsg.GetData(), msg)
			duration := time.Now().Sub(startTime)
			if err != nil {
				logx.WithContext(ctx).WithDuration(duration).Errorf("nats subscription|startTime:%v,subject:%v,body:%v,err:%v",
					startTime, msg.Subject, string(emsg.GetData()), err)
			} else {
				logx.WithContext(ctx).WithDuration(duration).Debugf("nats subscription|startTime:%v,subject:%v,body:%v",
					startTime, msg.Subject, string(emsg.GetData()))
			}
		})
	}
}

// GenericToNatsHandleFunc 将通用处理函数转换为 NATS 特定的处理函数
func GenericToNatsHandleFunc(handle HandleFunc) NatsHandleFunc {
	return func(ctx context.Context, ts time.Time, msg []byte, natsMsg *nats.Msg) error {
		return handle(ctx, ts, msg)
	}
}

func GenNatsJsDurable(serverName string, topic string) string {
	ip := netx.InternalIp()
	ret := fmt.Sprintf("%s_%s_%s", serverName, ip, topic)
	ret = strings.ReplaceAll(ret, ".", "-")
	ret = strings.ReplaceAll(ret, "*", "+")
	ret = strings.ReplaceAll(ret, ">", "~")
	return ret
}
