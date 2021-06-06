package slim

import (
	"time"
)

// AccessLog access log handler
func AccessLog() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()

		ctx := c.Request.Context()

		reqID := c.Request.Header.Get(XRequestID.String())
		if reqID == "" {
			reqID = RndUUIDMd5()
		}

		// set x-request-id to ctx
		if ctxReqID := ctx.Value(XRequestID); ctxReqID != nil {
			c.Request = ContextSet(c.Request, XRequestID, ctxReqID)
		} else {
			c.Request = ContextSet(c.Request, XRequestID, reqID)
		}

		debugPrintf("x-request-id: %s", reqID)

		// Process request
		c.Next()

		// Calculate resolution time
		debugPrintf("status_code: [%d] request_uri: %s exec_seconds: %.4f\n",
			c.StatusCode, c.Request.RequestURI, time.Since(t).Seconds())
	}
}
