package events

import (
	"context"
	"encoding/json"
	"time"

	"gitee.com/unitedrhino/share/ctxs"
	"gitee.com/unitedrhino/share/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"go.opentelemetry.io/otel/trace"
)

type MySpanContextConfig struct {
	TraceID string
	SpanID  string
}

// HandleFunc 通用的消息处理函数，与具体消息队列无关
type HandleFunc func(ctx context.Context, ts time.Time, msg []byte) error

type (
	// MsgHead 消息队列的头
	MsgHead struct {
		Trace     string        `json:"trace"`             //追踪tid
		Timestamp int64         `json:"timestamp"`         //发送时毫秒级时间戳
		Data      string        `json:"data,omitempty"`    //传送的内容
		UserCtx   *ctxs.UserCtx `json:"userCtx,omitempty"` ////context中携带的上下文,如用户信息,租户信息等
	}

	EventHandle interface {
		GetCtx() context.Context
		GetTs() time.Time
		GetData() []byte
	}
)

func NewEventMsg(ctx context.Context, data []byte) []byte {
	//生成新的消息时，使用go-zero的链路追踪接口，从ctx中提取span信息，并放入MsgHead中的Trace字段
	span := trace.SpanFromContext(ctx)
	traceinfo, _ := span.SpanContext().MarshalJSON()

	msg := MsgHead{
		Trace:     string(traceinfo),
		Timestamp: time.Now().UnixMilli(),
		Data:      string(data),
		UserCtx:   ctxs.GetUserCtx(ctx).ClearInner(),
	}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return nil
	}
	return msgBytes
}

func GetEventMsg(data []byte) EventHandle {
	msg := MsgHead{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil
	}
	return &msg
}

func (m *MsgHead) GetCtx() context.Context {
	var msg MySpanContextConfig
	err := json.Unmarshal([]byte(m.Trace), &msg)
	if err != nil {
		logx.Errorf("[GetCtx]|json Unmarshal trace.SpanContextConfig MsgHead:%v  err:%v", utils.Fmt(m), err)
		return context.Background()
	}
	//将MsgHead 中的msg链路信息 重新注入ctx中并返回
	t, err := trace.TraceIDFromHex(msg.TraceID)
	if err != nil {
		logx.Infof("[GetCtx]|TraceIDFromHex msg:%v  err:%v", m.Data, err)
		return context.Background()
	}
	s, err := trace.SpanIDFromHex(msg.SpanID)
	if err != nil {
		logx.Errorf("[GetCtx]|SpanIDFromHex MsgHead:%v  err:%v", utils.Fmt(m), err)
		return context.Background()
	}
	parent := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    t,
		SpanID:     s,
		TraceFlags: 0x1,
	})
	ctx := trace.ContextWithRemoteSpanContext(context.Background(), parent)
	return ctxs.SetUserCtx(ctx, m.UserCtx)
}

func (m *MsgHead) GetTs() time.Time {
	return time.UnixMilli(m.Timestamp)
}

func (m *MsgHead) GetData() []byte {
	if m == nil {
		return []byte{}
	}
	return []byte(m.Data)
}
