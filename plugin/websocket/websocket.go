package websocket

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/lsjhtang/kid/injector"

	"github.com/gorilla/websocket"
	"github.com/lsjhtang/kid/config"
	"github.com/lsjhtang/kid/logger"
	"github.com/lsjhtang/kid/utils"
)

type Websocket struct {
	upgrade *websocket.Upgrader
	conns   sync.Map
	options *Options
}

type ReadMessage struct {
	MsgType int
	Msg     []byte
}

type Connection struct {
	Socket   *websocket.Conn
	ChRead   chan *ReadMessage
	ChWrite  chan interface{}
	IsClosed bool
	ChClose  chan struct{}
}

func (conn *Connection) ReadLoop() {
	for {
		msgType, msg, err := conn.Socket.ReadMessage()
		if err != nil {
			logger.Errorf("read message error: %s", err.Error())
			break
		}
		readMsg := &ReadMessage{
			MsgType: msgType,
			Msg:     msg,
		}
		select {
		case conn.ChRead <- readMsg:
		case <-conn.ChClose:
			logger.Infof("Connection is close")
			break
		}
	}
}

func (conn *Connection) WriteLoop() {
	for {
		select {
		case msg := <-conn.ChWrite:
			if err := conn.Socket.WriteJSON(msg); err != nil {
				logger.Errorf("write message error: %s", err.Error())
				break
			}
		case <-conn.ChClose:
			logger.Info("connection is close")
			break
		}
	}
}

func (conn *Connection) Write(data interface{}) {
	conn.ChWrite <- data
}

func (conn *Connection) Read() *ReadMessage {
	return <-conn.ChRead
}

func (conn *Connection) Close() error {
	conn.ChClose <- struct{}{}
	return conn.Socket.Close()
}

type (
	BufferPoolKey  struct{}
	ErrorKey       struct{}
	CheckOriginKey struct{}
)

func (ws *Websocket) Provide(ctx context.Context) interface{} {
	var confName string
	if name, ok := ctx.Value(injector.NameKey{}).(string); ok && len(name) > 0 {
		confName = "." + name
	}
	confPrefix := fmt.Sprintf("websocket%s", confName)
	if !config.Exist(confPrefix) {
		panic(fmt.Sprintf("config.yaml file not found configuration item [%s]", confPrefix))
	}
	bufferPool, _ := ctx.Value(BufferPoolKey{}).(websocket.BufferPool)
	errorFunc, _ := ctx.Value(ErrorKey{}).(func(w http.ResponseWriter, r *http.Request, status int, reason error))
	return New(
		WithHandshakeTimeout(config.GetDuration(utils.GetConfigurationItem(confPrefix, "handshakeTimeout"))*time.Second),
		WithReadBufferSize(config.GetInt(utils.GetConfigurationItem(confPrefix, "readBufferSize"))),
		WithWriteBufferSize(config.GetInt(utils.GetConfigurationItem(confPrefix, "writeBufferSize"))),
		WithSubprotocols(config.GetStringSlice(utils.GetConfigurationItem(confPrefix, "subProtocols"))),
		WithWriteBufferPool(bufferPool),
		WithError(errorFunc),
		WithCheckOrigin(func(r *http.Request) bool {
			return true
		}),
		WithEnableCompression(config.GetBool(utils.GetConfigurationItem(confPrefix, "enableCompression"))),
	)
}

func (ws *Websocket) Upgrade(key interface{}, w http.ResponseWriter, r *http.Request, respHeader http.Header) (*Connection, error) {
	socket, err := ws.upgrade.Upgrade(w, r, respHeader)
	if err != nil {
		return nil, err
	}
	conn := &Connection{
		Socket:   socket,
		ChRead:   make(chan *ReadMessage, 20),
		ChWrite:  make(chan interface{}, 20),
		IsClosed: false,
		ChClose:  make(chan struct{}),
	}
	go conn.ReadLoop()
	go conn.WriteLoop()
	ws.AddConnection(key, conn)
	return conn, nil
}

func (ws *Websocket) AddConnection(key interface{}, conn *Connection) {
	ws.conns.Store(key, conn)
}

func (ws *Websocket) Exist(key interface{}) bool {
	if _, exist := ws.conns.Load(key); exist {
		return true
	}
	return false
}

func (ws *Websocket) RemoveConnection(key interface{}) {
	conn, ok := ws.GetConnection(key)
	if !ok {
		return
	}
	ws.conns.Delete(key)
	conn.Close()
}

func (ws *Websocket) GetConnection(key interface{}) (*Connection, bool) {
	v, exist := ws.conns.Load(key)
	if !exist {
		return nil, false
	}
	conn, ok := v.(*Connection)
	if !ok {
		return nil, false
	}
	return conn, true
}

type Option func(*Options)

func New(opts ...Option) *Websocket {
	var opt = new(Options)
	for _, o := range opts {
		o(opt)
	}
	return &Websocket{
		upgrade: &websocket.Upgrader{
			HandshakeTimeout:  opt.HandshakeTimeout,
			ReadBufferSize:    opt.ReadBufferSize,
			WriteBufferSize:   opt.WriteBufferSize,
			WriteBufferPool:   opt.WriteBufferPool,
			Subprotocols:      opt.Subprotocols,
			Error:             opt.Error,
			CheckOrigin:       opt.CheckOrigin,
			EnableCompression: opt.EnableCompression,
		},
		conns:   sync.Map{},
		options: opt,
	}
}
