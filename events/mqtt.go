package events

import (
	"context"
	"encoding/json"
	"time"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

// MqttSubscription MQTT 订阅适配器
func MqttSubscription(handle HandleFunc) func(topic string, payload []byte) {
	return func(topic string, payload []byte) {
		utils.Go(context.Background(), func() {
			var ctx context.Context
			utils.Recover(ctx)
			startTime := time.Now()
			emsg := GetEventMsg(payload)
			if emsg == nil {
				logx.Error(topic, string(payload))
				return
			}
			ctx = emsg.GetCtx()
			ctx, span := ctxs.StartSpan(ctx, topic, "")
			defer span.End()

			err := handle(ctx, emsg.GetTs(), emsg.GetData())
			duration := time.Now().Sub(startTime)
			if err != nil {
				logx.WithContext(ctx).WithDuration(duration).Errorf("mqtt subscription|startTime:%v,topic:%v,body:%v,err:%v",
					startTime, topic, string(emsg.GetData()), err)
			} else {
				logx.WithContext(ctx).WithDuration(duration).Debugf("mqtt subscription|startTime:%v,topic:%v,body:%v",
					startTime, topic, string(emsg.GetData()))
			}
		})
	}
}

// MqttSubWithType 支持类型化的 MQTT 订阅
func MqttSubWithType[msgType any](handle func(ctx context.Context, msgIn msgType) error) func(topic string, payload []byte) {
	return MqttSubscription(func(ctx context.Context, ts time.Time, msg []byte) error {
		var tempInfo msgType
		err := json.Unmarshal(msg, &tempInfo)
		if err != nil {
			return err
		}
		return handle(ctx, tempInfo)
	})
}
