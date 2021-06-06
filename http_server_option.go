package slim

import (
	"net/http"
	"time"
)

var (
	// DefaultMaxHeaderBytes default MaxHeaderBytes
	DefaultMaxHeaderBytes = 2 << 20 // header max 2MB

	// DefaultReadTimeout default ReadTimeout
	DefaultReadTimeout = 10 * time.Second

	// DefaultWriteTimeout default WriteTimeout
	DefaultWriteTimeout = 10 * time.Second

	// DefaultIdleTimeout default IdleTimeout
	DefaultIdleTimeout = 60 * time.Second
)

// HTTPServerOption http server option
type HTTPServerOption func(server *http.Server)

// NewHTTPServer 创建一个http server
func NewHTTPServer(opts ...HTTPServerOption) *http.Server {
	server := &http.Server{
		// Good practice to set timeouts to avoid Slowloris attacks.
		ReadTimeout:    DefaultReadTimeout,
		WriteTimeout:   DefaultWriteTimeout,
		IdleTimeout:    DefaultIdleTimeout,
		MaxHeaderBytes: DefaultMaxHeaderBytes,
	}

	for _, opt := range opts {
		opt(server)
	}

	return server
}

// WithHTTPReadTimeout 设置http read timeout
func WithHTTPReadTimeout(timeout time.Duration) HTTPServerOption {
	return func(server *http.Server) {
		server.ReadHeaderTimeout = timeout
	}
}

// WithHTTPAddress 设置http server address
func WithHTTPAddress(address string) HTTPServerOption {
	return func(server *http.Server) {
		server.Addr = address
	}
}

// WithHTTPWriteTimeout 设置http read timeout
func WithHTTPWriteTimeout(timeout time.Duration) HTTPServerOption {
	return func(server *http.Server) {
		server.WriteTimeout = timeout
	}
}

// WithHTTPIdleTimeout 设置 IdleTimeout
func WithHTTPIdleTimeout(timeout time.Duration) HTTPServerOption {
	return func(server *http.Server) {
		server.IdleTimeout = timeout
	}
}

// WithHTTPMaxHeaderBytes 设置MaxHeaderBytes
func WithHTTPMaxHeaderBytes(b int) HTTPServerOption {
	return func(server *http.Server) {
		server.MaxHeaderBytes = b
	}
}
