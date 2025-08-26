package eventBus

import (
	"context"

	"gitee.com/unitedrhino/share/clients"
	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/events"
	"github.com/nats-io/nats.go"
)

type natsEvent struct {
	natsCli *clients.NatsClient
}
type natsSubscription struct {
	*nats.Subscription
}

func NewNatsEvent(c conf.EventConf, serverName string, nodeID int64) (*natsEvent, error) {
	cli, err := clients.NewNatsClient2(c.Mode, serverName, c.Nats, nodeID)
	if err != nil {
		return nil, err
	}
	return &natsEvent{cli}, nil
}

func (n *natsEvent) Subscribe(subj string, cb events.HandleFunc) (subscription, error) {
	// 将通用 HandleFunc 转换为 NATS 特定的处理函数
	natsHandler := events.GenericToNatsHandleFunc(cb)
	sub, err := n.natsCli.Subscribe(subj, natsHandler)
	if err != nil {
		return nil, err
	}
	return &natsSubscription{sub}, nil
}

func (n *natsEvent) Publish(ctx context.Context, topic string, payload []byte) error {
	return n.natsCli.Publish(ctx, topic, payload)
}

func (n *natsEvent) QueueSubscribe(subj, queue string, cb events.HandleFunc) (subscription, error) {
	// 将通用 HandleFunc 转换为 NATS 特定的处理函数
	natsHandler := events.GenericToNatsHandleFunc(cb)
	sub, err := n.natsCli.QueueSubscribe(subj, queue, natsHandler)
	if err != nil {
		return nil, err
	}
	return &natsSubscription{sub}, nil
}

func (n *natsSubscription) Unsubscribe() error {
	return n.Subscription.Unsubscribe()
}
