package slim

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
)

// CatchPanic catch panic
func CatchPanic() {
	if err := recover(); err != nil {
		LogEntry.Printf("panic error: %v\n", err)
		LogEntry.Printf("full stack: %s\n", string(CatchStack()))
	}
}

// CatchStack 捕获指定stack信息,一般在处理panic/recover中处理
// 返回完整的堆栈信息和函数调用信息
func CatchStack() []byte {
	return debug.Stack()
}

// Recovery recovery handler
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				errMsg := fmt.Sprintf("%v", err)
				debugPrintf("exec panic: %s\n", errMsg)
				debugPrintf("full stack: %s\n", string(CatchStack()))

				// 是否是 brokenPipe类型的错误
				// 如果是该类型的错误，就不需要返回任何数据给客户端
				// 代码参考gin recovery.go RecoveryWithWriter方法实现
				// If the connection is dead, we can't write a status to it.
				// if broken pipe,return nothing.
				if isBroken(err) {
					// ctx.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}

// isBroken check for a broken connection, as it is not really a
// condition that warrants a panic stack trace.
func isBroken(err interface{}) bool {
	if ne, ok := err.(*net.OpError); ok {
		if se, ok := ne.Err.(*os.SyscallError); ok {
			errMsg := strings.ToLower(se.Error())
			debugPrintf("os syscall error:%s", errMsg)

			if strings.Contains(errMsg, "broken pipe") ||
				strings.Contains(errMsg, "reset by peer") ||
				strings.Contains(errMsg, "request headers: small read buffer") ||
				strings.Contains(errMsg, "unexpected EOF") ||
				strings.Contains(errMsg, "i/o timeout") {
				return true
			}
		}
	}

	return false
}
