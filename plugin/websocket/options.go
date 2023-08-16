package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

type Options struct {
	HandshakeTimeout time.Duration

	ReadBufferSize, WriteBufferSize int

	WriteBufferPool websocket.BufferPool

	Subprotocols []string

	Error func(w http.ResponseWriter, r *http.Request, status int, reason error)

	CheckOrigin func(r *http.Request) bool

	EnableCompression bool
}

func WithHandshakeTimeout(handshakeTimeout time.Duration) Option {
	return func(opt *Options) {
		opt.HandshakeTimeout = handshakeTimeout
	}
}

func WithReadBufferSize(readBuffSize int) Option {
	return func(opt *Options) {
		opt.ReadBufferSize = readBuffSize
	}
}

func WithWriteBufferSize(writeBuffSize int) Option {
	return func(opt *Options) {
		opt.WriteBufferSize = writeBuffSize
	}
}

func WithWriteBufferPool(writeBufferPool websocket.BufferPool) Option {
	return func(opt *Options) {
		opt.WriteBufferPool = writeBufferPool
	}
}

func WithSubprotocols(subProtocols []string) Option {
	return func(opt *Options) {
		opt.Subprotocols = subProtocols
	}
}

func WithError(err func(w http.ResponseWriter, r *http.Request, status int, reason error)) Option {
	return func(opt *Options) {
		opt.Error = err
	}
}

func WithCheckOrigin(checkOrigin func(r *http.Request) bool) Option {
	return func(opt *Options) {
		opt.CheckOrigin = checkOrigin
	}
}

func WithEnableCompression(enableCompression bool) Option {
	return func(opt *Options) {
		opt.EnableCompression = enableCompression
	}
}
