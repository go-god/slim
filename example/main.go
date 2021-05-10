package main

import (
	"net/http"

	"github.com/go-god/slim"
)

func main() {
	r := slim.Default()
	r.GET("/", func(c *slim.Context) {
		c.String(http.StatusOK, "hello word")
	})

	// index out of range for testing Recovery()
	r.GET("/panic", func(c *slim.Context) {
		names := []string{"abc"}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":1337")
}
