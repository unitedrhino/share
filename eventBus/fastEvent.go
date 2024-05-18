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
	"sync"
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
	queueMutex    sync.RWMutex
	handlerMutex  sync.RWMutex
	isStart       bool
}

type FastFunc func(ctx context.Context, t time.Time, body []byte) error

var (
	fastEvent *FastEvent
	fastOnce  sync.Once
)

func NewFastEvent(c conf.EventConf, serverName string, nodeID int64) (s *FastEvent, err error) {
	fastOnce.Do(func() {
		fastEvent = &FastEvent{handlers: map[string][]FastFunc{}, queueHandlers: map[string][]FastFunc{}, serverName: serverName}
		switch c.Mode {
		case conf.EventModeNats, conf.EventModeNatsJs:
			fastEvent.natsCli, err = clients.NewNatsClient2(c.Mode, serverName, c.Nats, nodeID)
		default:
			err = errors.Parameter.AddMsgf("mode:%v not support", c.Mode)
		}
	})
	return fastEvent, err
}

func (bus *FastEvent) subscribe(topic string) error {
	_, err := bus.natsCli.Subscribe(topic, func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
		natsMsg.Ack()
		ctx = ctxs.CopyCtx(ctx)
		bus.handlerMutex.RLock()
		defer bus.handlerMutex.RUnlock()
		for _, f := range bus.handlers[topic] {
			ff := f
			utils.Go(ctx, func() {
				err := ff(ctx, events.GetEventMsg(natsMsg.Data).GetTs(), msg)
				if err != nil {
					logx.WithContext(ctx).Error(err)
				}
			})
		}
		return nil
	})
	return err
}

func (bus *FastEvent) queueSubscribe(topic string) error {
	_, err := bus.natsCli.QueueSubscribe(topic, bus.serverName, func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
		ctx = ctxs.CopyCtx(ctx)
		bus.queueMutex.RLock()
		defer bus.queueMutex.RUnlock()
		for _, f := range bus.queueHandlers[topic] {
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
	return err
}

func (bus *FastEvent) Start() error {
	if bus.isStart == true {
		return nil
	}
	bus.isStart = true
	for topic := range bus.handlers {
		err := bus.subscribe(topic)
		if err != nil {
			return err
		}
	}
	for topic := range bus.queueHandlers {
		err := bus.queueSubscribe(topic)
		if err != nil {
			return err
		}
	}
	return nil
}

// Subscribe 订阅
func (bus *FastEvent) Subscribe(topic string, f FastFunc) error {
	bus.handlerMutex.Lock()
	defer bus.handlerMutex.Unlock()
	handler, ok := bus.handlers[topic]
	if !ok {
		handler = []FastFunc{}
	}
	handler = append(handler, f)
	bus.handlers[topic] = handler
	if !ok && bus.isStart { //如果已经启动,且没有监听这个topic,则需要加入
		err := bus.subscribe(topic)
		return err
	}
	return nil
}

func (bus *FastEvent) QueueSubscribe(topic string, f FastFunc) error {
	bus.queueMutex.Lock()
	defer bus.queueMutex.Unlock()
	handler, ok := bus.queueHandlers[topic]
	if !ok {
		handler = []FastFunc{}
	}
	handler = append(handler, f)
	bus.queueHandlers[topic] = handler
	if !ok && bus.isStart { //如果已经启动,且没有监听这个topic,则需要加入
		err := bus.queueSubscribe(topic)
		return err
	}
	return nil
}

// Publish 发布
// 这里异步执行，并且不会等待返回结果
func (bus *FastEvent) Publish(ctx context.Context, topic string, arg any) error {
	err := bus.natsCli.Publish(ctx, topic, []byte(utils.MarshalNoErr(arg)))
	return err
}
