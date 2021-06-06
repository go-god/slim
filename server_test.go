package slim

import (
	"log"
	"net/http"
	"runtime/debug"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	server := NewServer(":8080")
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Welcome here"))
	})

	handler.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello world"))
	})

	server.server.Handler = handler

	server.Run()
}

func TestServerWithOptions(t *testing.T) {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Welcome here"))
	})

	handler.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("hello world"))
	})

	server := NewServer(":8090",
		WithGracefulWait(5*time.Second),
		WithHandler(handler),
		WithRecovery(func() {
			if err := recover(); err != nil {
				log.Printf("exec panic: %v", err)
				log.Printf("full stack: %s", string(debug.Stack()))
			}
		}),
	)

	server.Run()
}
