package slim

import (
	"net/http"
	"time"
)

// Option Server option
type Option func(s *Server)

// Apply apply option for Server
func (s *Server) Apply(opts ...Option) {
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		opt(s)
	}
}

// WithHandler 设置server handler
func WithHandler(handler http.Handler) Option {
	return func(s *Server) {
		if handler == nil {
			panic("please set handler")
		}

		s.server.Handler = handler
	}
}

// WithHTTPServer 设置http server
func WithHTTPServer(server *http.Server) Option {
	return func(s *Server) {
		s.server = server
	}
}

// WithAddress 设置运行地址addr
func WithAddress(addr string) Option {
	return func(s *Server) {
		s.address = addr
	}
}

// WithRecovery 设置recovery
func WithRecovery(fn func()) Option {
	return func(s *Server) {
		s.recovery = fn
	}
}

// WithGracefulWait 设置平滑退出时间
func WithGracefulWait(gracefulWait time.Duration) Option {
	return func(s *Server) {
		s.gracefulWait = gracefulWait
	}
}

// WithShutdownFunc 设置shutdown func
func WithShutdownFunc(fn func()) Option {
	return func(s *Server) {
		s.shutdownFunc = fn
	}
}

// WithLogger 设置logger
func WithLogger(l Logger) Option {
	return func(s *Server) {
		s.logger = l
	}
}
