package eventBus

import (
	"context"
	"sync"
	"time"

	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/events"
	"gitee.com/unitedrhino/share/utils"
	"github.com/nats-io/nats.go"
	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/atomic"
)

/*
服务消息,不需要模糊匹配的,发送给所有订阅者的可以用这个来简化实现
*/

type FastEvent struct {
	natsCli       eventCli
	handlers      map[string]*handleInfo
	queueHandlers map[string]*handleInfo
	serverName    string
	queueMutex    sync.RWMutex
	handlerMutex  sync.RWMutex
	isStart       bool
}

type handleInfo struct {
	Handle map[int64]FastFunc
	Sub    subscription
}

type FastFunc func(ctx context.Context, t time.Time, body []byte) error

var (
	fastEvent *FastEvent
	fastOnce  sync.Once
	idGen     atomic.Int64
)

type (
	subscription interface {
		Unsubscribe() error
	}
	eventCli interface {
		Subscribe(subj string, cb events.HandleFunc) (*natsSubscription, error)
		Publish(ctx context.Context, topic string, arg []byte) error
		QueueSubscribe(subj, queue string, cb events.HandleFunc) (*natsSubscription, error)
	}
)

func NewFastEvent(c conf.EventConf, serverName string, nodeID int64) (s *FastEvent, err error) {
	fastOnce.Do(func() {
		fastEvent = &FastEvent{handlers: map[string]*handleInfo{}, queueHandlers: map[string]*handleInfo{}, serverName: serverName}
		switch c.Mode {
		case conf.EventModeNats, conf.EventModeNatsJs:
			fastEvent.natsCli, err = NewNatsEvent(c, serverName, nodeID)
		default:
			err = errors.Parameter.AddMsgf("mode:%v not support", c.Mode)
		}
	})
	return fastEvent, err
}

func (bus *FastEvent) subscribe(topic string) (subscription, error) {
	sub, err := bus.natsCli.Subscribe(topic, func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
		natsMsg.Ack()
		ctx = ctxs.CopyCtx(ctx)
		bus.handlerMutex.RLock()
		defer bus.handlerMutex.RUnlock()
		if _, ok := bus.handlers[topic]; !ok {
			return nil
		}
		for _, f := range bus.handlers[topic].Handle {
			ff := f
			ctxs.GoNewCtx(ctx, func(ctx context.Context) {
				err := ff(ctx, events.GetEventMsg(natsMsg.Data).GetTs(), msg)
				if err != nil {
					logx.WithContext(ctx).Error(err)
				}
			})
		}
		return nil
	})
	return sub, err
}

func (bus *FastEvent) queueSubscribe(topic string) (subscription, error) {
	sub, err := bus.natsCli.QueueSubscribe(topic, bus.serverName, func(ctx context.Context, msg []byte, natsMsg *nats.Msg) error {
		ctx = ctxs.CopyCtx(ctx)
		bus.queueMutex.RLock()
		defer bus.queueMutex.RUnlock()
		if _, ok := bus.queueHandlers[topic]; !ok {
			return nil
		}
		for _, f := range bus.queueHandlers[topic].Handle {
			run := f
			ctxs.GoNewCtx(ctx, func(ctx context.Context) {
				err := run(ctx, events.GetEventMsg(natsMsg.Data).GetTs(), msg)
				if err != nil {
					logx.WithContext(ctx).Error(err)
				}
			})
		}
		return nil
	})
	return sub, err
}

func (bus *FastEvent) Start() error {
	if bus.isStart == true {
		return nil
	}
	bus.isStart = true
	for topic, h := range bus.handlers {
		sub, err := bus.subscribe(topic)
		if err != nil {
			return err
		}
		h.Sub = sub
	}
	for topic, h := range bus.queueHandlers {
		sub, err := bus.queueSubscribe(topic)
		if err != nil {
			return err
		}
		h.Sub = sub
	}
	return nil
}

// Subscribe 订阅
func (bus *FastEvent) Subscribe(topic string, f FastFunc) error {
	_, err := bus.SubscribeWithID(topic, f)
	return err
}
func (bus *FastEvent) SubscribeWithID(topic string, f FastFunc) (int64, error) {
	bus.handlerMutex.Lock()
	defer bus.handlerMutex.Unlock()
	handler, ok := bus.handlers[topic]
	if !ok {
		handler = &handleInfo{
			Handle: map[int64]FastFunc{},
			Sub:    nil,
		}
	}
	id := idGen.Add(1)
	handler.Handle[id] = f
	bus.handlers[topic] = handler
	if !ok && bus.isStart { //如果已经启动,且没有监听这个topic,则需要加入
		sub, err := bus.subscribe(topic)
		if err != nil {
			return 0, err
		}
		handler.Sub = sub
		return id, err
	}
	return id, nil
}
func (bus *FastEvent) UnSubscribeWithID(topic string, id int64) error {
	bus.handlerMutex.Lock()
	defer bus.handlerMutex.Unlock()
	if _, ok := bus.handlers[topic]; !ok {
		return nil
	}
	delete(bus.handlers[topic].Handle, id)
	if len(bus.handlers[topic].Handle) == 0 {
		err := bus.handlers[topic].Sub.Unsubscribe()
		if err != nil {
			logx.Error(err)
			return nil
		}
		delete(bus.handlers, topic)
	}
	return nil
}

func (bus *FastEvent) QueueSubscribeWithID(topic string, f FastFunc) (int64, error) {
	bus.queueMutex.Lock()
	defer bus.queueMutex.Unlock()
	handler, ok := bus.queueHandlers[topic]
	if !ok {
		handler = &handleInfo{
			Handle: map[int64]FastFunc{},
			Sub:    nil,
		}
	}
	id := idGen.Add(1)
	handler.Handle[id] = f
	bus.queueHandlers[topic] = handler
	if !ok && bus.isStart { //如果已经启动,且没有监听这个topic,则需要加入
		sub, err := bus.queueSubscribe(topic)
		if err != nil {
			return 0, err
		}
		handler.Sub = sub
		return id, err
	}
	return id, nil
}

func (bus *FastEvent) QueueSubscribe(topic string, f FastFunc) error {
	_, err := bus.QueueSubscribeWithID(topic, f)
	return err
}

func (bus *FastEvent) UnQueueSubscribeWithID(topic string, id int64) error {
	bus.queueMutex.Lock()
	defer bus.queueMutex.Unlock()
	if _, ok := bus.queueHandlers[topic]; !ok {
		return nil
	}
	delete(bus.queueHandlers[topic].Handle, id)
	if len(bus.queueHandlers[topic].Handle) == 0 {
		err := bus.queueHandlers[topic].Sub.Unsubscribe()
		if err != nil {
			logx.Error(err)
			return nil
		}
		delete(bus.queueHandlers, topic)
	}
	return nil
}

// Publish 发布
// 这里异步执行，并且不会等待返回结果
func (bus *FastEvent) Publish(ctx context.Context, topic string, arg any) error {
	err := bus.natsCli.Publish(ctx, topic, []byte(utils.ToString(arg)))
	return err
}
