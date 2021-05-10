package slim

import (
	"net/http"
)

// NoRoute no route HandlersChain
func NoRoute() HandlerFunc {
	return func(c *Context) {
		c.String(http.StatusNotFound, "404 not found: %s\n", c.Path)
	}
}
