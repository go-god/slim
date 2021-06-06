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
	server       *http.Server  // http server object
	address      string        // http server host eg: ip:port
	recovery     func()        // goroutine exec recover catch stack
	gracefulWait time.Duration // when server exit graceful wait time
	shutdownFunc func()        // shutdown callback func
	logger       Logger        // server logger
}

// NewServer create Server entry through Functional Options
func NewServer(address string, opts ...Option) *Server {
	if address == "" {
		panic("server run address is empty")
	}

	s := &Server{
		server:       NewHTTPServer(),
		address:      address,
		recovery:     CatchPanic,
		gracefulWait: 5 * time.Second,
		// 采用接口函数模式，设置logger,默认采用log.Printf
		logger: LoggerFunc(log.Printf),
	}

	s.shutdownFunc = func() {
		s.logger.Printf("server shutdown...\n")
	}

	if len(opts) > 0 {
		s.Apply(opts...)
	}

	if s.server == nil {
		panic("http server is nil")
	}

	if s.server.Handler == nil {
		s.server.Handler = Default() // 默认engine
	}

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

	// graceful exit
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
