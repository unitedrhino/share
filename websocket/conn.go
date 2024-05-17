package websocket

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"gitee.com/asktop_golib/util/aslice"
	"gitee.com/i-Things/share/ctxs"
	"gitee.com/i-Things/share/errors"
	"gitee.com/i-Things/share/eventBus"
	"gitee.com/i-Things/share/utils"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/kv"
	"github.com/zeromicro/go-zero/core/trace"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

var (
	dp             *dispatcher //ws调度器
	once           sync.Once
	store          kv.Store
	nodeID         int64
	checkSubscribe func(ctx context.Context, in *SubscribeInfo) error
)

const (
	errorCount    = 5                     //错误次数
	interval      = 5 * time.Second       //心跳间隔
	keepAliveType = websocket.PingMessage //心跳类型
)

type connection struct {
	r        *http.Request
	uc       *ctxs.UserCtx
	server   *Server
	ws       *websocket.Conn   //ws连接实例
	userID   int64             //ws连接实例唯一标识
	closed   bool              //ws连接已关闭
	send     chan []byte       //发送信息管道
	topics   map[string]string //订阅信息
	pingErrs []int64           //发送的心跳失败次数
	pongErrs []int64           //收到的心跳失败次数
}

// ws调度器
type dispatcher struct {
	s2cGzip  bool                  //发送的信息是否gzip压缩
	connPool map[int64]*connection //ws连接池 map[clientId]*connection
	mu       sync.RWMutex          // 互斥锁
}
type WsPublish struct {
	UserID int64
	Code   string
	Data   any
}
type WsPublishes []WsPublish

// 创建ws调度器
func StartWsDp(s2cGzip bool, NodeID int64, event *eventBus.FastEvent, c cache.ClusterConf) {
	nodeID = NodeID
	once.Do(func() {
		store = kv.NewStore(c)
		dp = newDp(s2cGzip)
		event.Subscribe(fmt.Sprintf(eventBus.CoreApiUserPublish, NodeID), func(ctx context.Context, t time.Time, body []byte) error {
			var pbs = WsPublishes{}
			logx.Infof("StartWsDp.sendMessage nodeID:%v publishs:%v", nodeID, string(body))
			err := json.Unmarshal(body, &pbs)
			if err != nil {
				return err
			}
			for _, pb := range pbs {
				err := func() error {
					dp.mu.RLock()
					defer dp.mu.RUnlock()
					conn := dp.connPool[pb.UserID]
					if conn == nil {
						return err
					}
					conn.sendMessage(WsResp{
						WsBody: WsBody{
							Type: Pub,
							Path: pb.Code,
							Body: pb.Data,
						},
					})
					return nil
				}()
				if err != nil {
					logx.WithContext(ctx).Error(err)
					continue
				}
			}
			return nil
		})
	})
}

func RegisterSubscribeCheck(f func(ctx context.Context, in *SubscribeInfo) error) {
	checkSubscribe = f
}

// 创建ws调度器
func newDp(s2cGzip bool) *dispatcher {
	d := &dispatcher{
		s2cGzip:  s2cGzip,
		connPool: make(map[int64]*connection),
	}
	return d
}

// 读ping心跳
func (c *connection) pingRead(message []byte) error {
	logx.Infof("%s.[ws] message:%s userID:%v", utils.FuncName(), string(message), c.userID)
	if aslice.ContainInt64(c.pingErrs, int64(binary.BigEndian.Uint64(message))) {
		c.pingErrs = []int64{}
	} else {
		c.writeMessage(websocket.PingMessage, []byte("ping error message :"+string(message)))
	}
	return nil
}

// 读pong心跳
func (c *connection) pongRead(message []byte) error {
	//logx.Infof("%s.[ws] message:%s userID:%v", utils.FuncName(), string(message), c.userID)
	if aslice.ContainInt64(c.pongErrs, int64(binary.BigEndian.Uint64(message))) {
		c.pongErrs = []int64{}
	} else {
		c.writeMessage(websocket.PongMessage, []byte("pong error message :"+string(message)))
	}
	return nil
}

// 发送ping心跳
func (c *connection) pingSend() error {
	if len(c.pingErrs) >= errorCount || len(c.pongErrs) >= errorCount {
		//连续5次没有收到ping心跳 关闭连接
		return errors.System.AddMsg("connection timeout")
	}
	nowTime := []byte(strconv.FormatInt(time.Now().Unix(), 10))
	if err := c.writeMessage(websocket.PingMessage, nowTime); err != nil {
		c.pingErrs = append(c.pingErrs, int64(binary.BigEndian.Uint64(nowTime)))
	} else {
		c.pingErrs = []int64{}
		c.pongErrs = append(c.pongErrs, int64(binary.BigEndian.Uint64(nowTime)))
	}
	return nil
}

// 发送pong心跳
func (c *connection) pongSend() error {
	if len(c.pingErrs) >= errorCount || len(c.pongErrs) >= errorCount {
		//连续5次没有收到pong心跳 关闭连接
		return errors.System.AddMsg("connection timeout")
	}
	nowTime := []byte(strconv.FormatInt(time.Now().Unix(), 10))
	if err := c.writeMessage(websocket.PongMessage, nowTime); err != nil {
		c.pongErrs = append(c.pongErrs, int64(binary.BigEndian.Uint64(nowTime)))
	} else {
		c.pongErrs = []int64{}
		c.pingErrs = append(c.pingErrs, int64(binary.BigEndian.Uint64(nowTime)))
	}
	return nil
}

// 发送订阅信息
func SendSub(ctx context.Context, body WsResp) {
	clientToken := trace.TraceIDFromContext(ctx)
	body.Handler = map[string]string{"Traceparent": clientToken}
}

// 创建ws连接
func NewConn(ctx context.Context, userID int64, server *Server, r *http.Request, wsConn *websocket.Conn) *connection {
	conn := &connection{
		server: server,
		ws:     wsConn,
		uc:     ctxs.GetUserCtx(ctx),
		r:      r,
		userID: userID,
		send:   make(chan []byte, 10000),
		topics: make(map[string]string),
	}
	dp.connPool[userID] = conn
	logx.Infof("%s.[ws]创建连接成功 RemoteAddr::%s userID:%v", utils.FuncName(), wsConn.RemoteAddr().String(), userID)
	resp := WsResp{}
	clientToken := trace.TraceIDFromContext(ctx)
	resp.Handler = map[string]string{"Traceparent": clientToken}
	conn.sendMessage(resp)
	return conn
}

// 开启读取进程
func (c *connection) StartRead() {
	defer func() {
		c.Close("read message error")
	}()
	c.ws.SetPongHandler(func(message string) error {
		c.pongRead([]byte(message))
		return nil
	})
	c.ws.SetPingHandler(func(message string) error {
		c.pingRead([]byte(message))
		return nil
	})
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		logx.Infof("%s.[ws] message:%s userID:%v", utils.FuncName(), string(message), c.userID)
		var data map[string]interface{}
		err = json.Unmarshal(message, &data)
		if err != nil {
			c.errorSend(errors.Type.AddDetail("error reading message"))
			continue
		}
		c.handleRequest(message)
	}
}
func (c *connection) errorSend(data error) {
	e := errors.Fmt(data)
	resp := WsResp{
		Code: e.GetCode(),
		Msg:  e.GetI18nMsg(""),
	}
	c.sendMessage(resp)
}

func (c *connection) handleRequest(message []byte) {
	var body WsReq
	err := json.Unmarshal(message, &body)
	if err != nil {
		c.errorSend(errors.Parameter)
		return
	}
	if err := isDataComplete(body.Type, body); err != nil {
		c.errorSend(err)
		return
	}
	if len(body.Handler) > 0 {
		for k, v := range body.Handler {
			c.r.Header.Set(k, v)
		}
	}
	ctx := ctxs.SetUserCtx(context.Background(), c.uc)
	switch body.Type {
	case Control:
		downControl(c, body)
	case Sub:
		subscribeHandle(ctx, c, body)
	case UnSub:
		unSubscribeHandle(ctx, c, body)
	default:
	}
}

func isDataComplete(wsType WsType, body WsReq) error {
	if wsType == "" {
		return errors.Parameter.AddDetail("type is  null")
	}
	switch wsType {
	case Control:
		if body.Path == "" || body.Method == "" || body.Body == "" {
			return errors.Parameter.AddDetail("path|method|body is  null")
		}
	case Sub, UnSub:
		if _, ok := body.Body.(map[string]interface{}); !ok {
			return errors.Parameter.AddDetail("body is  null")
		}
	case Pub:
		return errors.NotRealize
	default:
	}
	return nil
}

func downControl(c *connection, body WsReq) {
	reqBody, err := getRequestBody(body.Body)
	if err != nil {
		// 处理编码错误
	}
	bodyBytes, err := json.Marshal(body.Body)
	length := len(bodyBytes)
	header := c.r.Header
	header.Set("Content-Type", "application/json")
	header.Set("Content-Length", strconv.Itoa(length))
	r := &http.Request{
		Method: body.Method,
		Host:   c.r.Host,
		URL: &url.URL{
			Path: body.Path,
		},
		Header:        header,
		Body:          reqBody,
		ContentLength: int64(length),
	}
	w := response{req: &body, resp: WsResp{WsBody: WsBody{Handler: map[string]string{}, Type: ControlRet}}}
	c.server.ServeHTTP(&w, r)
	if token := w.Header().Get(ctxs.UserSetTokenKey); token != "" { //登录态保持更新
		c.r.Header.Set(ctxs.UserSetTokenKey, token)
	}
	c.sendMessage(w.resp)
}

// 将请求体转换为io.ReadCloser类型
func getRequestBody(body interface{}) (io.ReadCloser, error) {
	var reqBody io.ReadCloser
	if body != nil {
		switch body.(type) {
		case string:
			reqBody = io.NopCloser(bytes.NewBufferString(body.(string)))
		case []byte:
			reqBody = io.NopCloser(bytes.NewBuffer(body.([]byte)))
		case map[string]interface{}:
			bodyBytes, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			reqBody = io.NopCloser(bytes.NewReader(bodyBytes))
		default:
			// 处理其他类型
		}
	}
	return reqBody, nil
}

// 开启发送进程
func (c *connection) StartWrite() {
	ticker := time.NewTicker(interval)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		//发送心跳
		case <-ticker.C:
			if c.closed {
				return
			}
			var err error
			switch keepAliveType {
			case websocket.PingMessage:
				err = c.pingSend()
			case websocket.PongMessage:
				err = c.pongSend()
			}
			if err != nil {
				c.Close("connection timeout")
				return
			}
		//发送信息
		case message := <-c.send:
			if c.closed {
				return
			}
			if err := c.writeMessage(websocket.TextMessage, message); err != nil {
				c.Close("send message error")
				return
			}
		}
	}
}

// 关闭ws连接
func (c *connection) Close(msg string) {
	dp.mu.Lock()
	defer dp.mu.Unlock()
	_, ok := dp.connPool[c.userID]
	if ok || !c.closed {
		c.closed = true
		close(c.send)
		delete(dp.connPool, c.userID)
		NewUserSubscribe(store).Clear(context.Background(), c.userID)
		c.ws.Close()
		logx.Infof("%s.[ws]关闭连接  userID:%v", utils.FuncName(), c.userID)
	}
}

// 发送信息
func (c *connection) sendMessage(body WsResp) {
	if body.Code == 0 {
		body.Code = errors.OK.Code
		body.Msg = errors.OK.GetMsg()
	}
	message, _ := json.Marshal(body)
	if !c.closed {
		c.send <- message
	}
}

// 写消息
func (c *connection) writeMessage(messageType int, message []byte) error {
	if message == nil {
		logx.Infof("%s.[ws]error message: is  null ")
	}
	switch messageType {
	case websocket.PingMessage, websocket.PongMessage:
		err := c.ws.WriteControl(messageType, message, time.Time{})
		if err != nil {
			logx.Infof("%s.[ws]error message::%s userID:%v", utils.FuncName(), string(message), c.userID)
		}
	case websocket.TextMessage:
		err := c.ws.WriteMessage(messageType, message)
		if err != nil {
			logx.Infof("%s.[ws]error message::%s userID:%v", utils.FuncName(), string(message), c.userID)
		}
	}
	logx.Debugf("%s.[ws] message:%s userID:%v", utils.FuncName(), string(message), c.userID)
	return nil
}
