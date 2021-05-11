package slim

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Server http server struct
type Server struct {
	handler      http.Handler  // http server handler
	server       *http.Server  // http server object
	address      string        // http server host eg: ip:port
	recovery     func()        // goroutine exec recover catch stack
	gracefulWait time.Duration // when server exit graceful wait time
	shutdownFunc func()        // shutdown callback func
	logger       Logger        // server logger
}

// Option Server option
type Option func(s *Server)

// NewServer create http server
func NewServer(address string, opts ...Option) *Server {
	if address == "" {
		panic("server run address is empty")
	}

	s := &Server{
		address:      address,
		recovery:     CatchPanic,
		gracefulWait: 5 * time.Second,
		logger:       LoggerFunc(log.Printf), // 默认采用log.Printf
	}

	s.shutdownFunc = func() {
		s.logger.Printf("server shutdown...\n")
	}

	for _, opt := range opts {
		opt(s)
	}

	if s.server == nil {
		s.server = InitHTTPServer() // 初始化默认http server
	}

	if s.handler == nil {
		s.handler = Default() // 默认engine
	}

	s.server.Handler = s.handler

	return s
}

// Run run server
func (s *Server) Run() {
	s.server.Addr = s.address
	if s.server.Handler == nil {
		panic("please set handler")
	}

	// 在独立协程中运行服务
	s.logger.Printf("server run on: %s\n", s.address)
	s.logger.Printf("server pid: %d\n", os.Getppid())

	// 注册平滑退出时候shutdown callback func
	s.server.RegisterOnShutdown(s.shutdownFunc)

	go func() {
		defer s.recovery()

		if err := s.server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				s.logger.Printf("server close error: %s\n", err.Error())
				return
			}

			s.logger.Printf("server will exit...\n")
		}
	}()

	// 平滑重启
	ch := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// recv signal to exit main goroutine
	// window signal
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)

	// linux signal if you use linux on production,please use this code.
	signal.Notify(ch, InterruptSignals...)

	// Block until we receive our signal.
	sig := <-ch

	s.logger.Printf("exit signal: %s\n", sig.String())
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), s.gracefulWait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// if your application should wait for other services
	// to finalize based on context cancellation.
	go s.server.Shutdown(ctx)
	<-ctx.Done()

	s.logger.Printf("server shutting down\n")
}

// WithHandler 设置server handler
func WithHandler(handler http.Handler) Option {
	return func(s *Server) {
		if handler == nil {
			panic("please set handler")
		}

		s.handler = handler
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
