package eventBus

import (
	"context"
	"gitee.com/i-Things/share/clients"
	"gitee.com/i-Things/share/conf"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/events"
	"gitee.com/i-Things/share/utils"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
)

/*
服务消息,不需要模糊匹配的,发送给所有订阅者的可以用这个来简化实现
*/

type FastEvent struct {
	natsCli       *clients.NatsClient
	handlers      map[string][]FastFunc
	queueHandlers map[string][]FastFunc
	serverName    string
}

type FastFunc func(ctx context.Context, t time.Time, body []byte) error

func NewFastEvent(c conf.EventConf, serverName string) (s *FastEvent, err error) {
	serverMsg := FastEvent{handlers: map[string][]FastFunc{}, queueHandlers: map[string][]FastFunc{}, serverName: serverName}
	switch c.Mode {
	case conf.EventModeNats, conf.EventModeNatsJs:
		serverMsg.natsCli, err = clients.NewNatsClient2(c.Mode, serverName, c.Nats)
	default:
		err = errors.Parameter.AddMsgf("mode:%v not support", c.Mode)
	}
	return &serverMsg, err
}

func (bus *FastEvent) Start() error {
	for topic, handles := range bus.handlers {
		hs := handles
		_, err := bus.natsCli.Subscribe(topic, func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			ctx = ctxs.CopyCtx(ctx)
			for _, f := range hs {
				utils.Go(ctx, func() {
					err := f(ctx, events.GetEventMsg(natsMsg.Data).GetTs(), msg)
					if err != nil {
						logx.WithContext(ctx).Error(err)
					}
				})
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	for topic, handles := range bus.queueHandlers {
		hs := handles
		_, err := bus.natsCli.QueueSubscribe(topic, bus.serverName, func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
			ctx = ctxs.CopyCtx(ctx)
			for _, f := range hs {
				run := f
				utils.Go(ctx, func() {
					err := run(ctx, events.GetEventMsg(natsMsg.Data).GetTs(), msg)
					if err != nil {
						logx.WithContext(ctx).Error(err)
					}
				})
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Subscribe 订阅
func (bus *FastEvent) Subscribe(topic string, f FastFunc) {
	handler, ok := bus.handlers[topic]
	if !ok {
		handler = []FastFunc{}
	}
	handler = append(handler, f)
	bus.handlers[topic] = handler
	return
}

func (bus *FastEvent) QueueSubscribe(topic string, f FastFunc) {
	handler, ok := bus.queueHandlers[topic]
	if !ok {
		handler = []FastFunc{}
	}
	handler = append(handler, f)
	bus.queueHandlers[topic] = handler
	return
}

// Publish 发布
// 这里异步执行，并且不会等待返回结果
func (bus *FastEvent) Publish(ctx context.Context, topic string, arg any) error {
	err := bus.natsCli.Publish(ctx, topic, []byte(utils.Fmt(arg)))
	return err
}
