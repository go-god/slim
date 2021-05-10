package slim

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestNestedGroup(t *testing.T) {
	r := New()
	v1 := r.Group("/v1")
	v2 := v1.Group("/v2")
	v3 := v2.Group("/v3")
	if v2.prefix != "/v1/v2" {
		t.Fatal("v2 prefix should be /v1/v2")
	}
	if v3.prefix != "/v1/v2/v3" {
		t.Fatal("v2 prefix should be /v1/v2")
	}
}

func TestEngineRun(t *testing.T) {
	// engine := New()
	// engine.Use(Recovery(), AccessLog())

	engine := Default()
	engine.GET("/", func(ctx *Context) {
		ctx.ApiSuccess(0, "ok", H{
			"a": 1,
			"b": 2,
		})
	})

	v1 := engine.Group("/v1", func(c *Context) {
		log.Println("abc")
		log.Println("handler name: ", c.HandlerName())
	}, func(c *Context) {
		log.Println("abc2")
		log.Println("handler name2: ", c.HandlerName())
	})

	v1.GET("/index", func(ctx *Context) {
		ctx.JSON(200, H{
			"code":    0,
			"message": "ok",
			"data": H{
				"lang": []string{
					"go", "js", "rust",
				},
			},
		})
	})

	engine.GET("/test", func(ctx *Context) {
		ctx.String(200, "hello")
	})

	engine.Run(":8080")
}

func TestRunWithServer(t *testing.T) {
	var (
		// 运行地址
		address = fmt.Sprintf("0.0.0.0:1337")

		// 平滑退出等待时间
		gracefulWait = 5 * time.Second
	)

	// 创建slim engine router引擎
	engine := New()
	engine.Use(Recovery(), AccessLog())

	engine.GET("/", func(ctx *Context) {
		ctx.ApiSuccess(0, "ok", H{
			"a": 1,
			"b": 2,
		})
	})

	engine.GET("/test", func(ctx *Context) {
		ctx.String(200, "hello")
	})

	// engine.Run(address)

	server := NewServer(address,
		WithGracefulWait(gracefulWait),
		WithHandler(engine),
	)

	server.Run()
}
