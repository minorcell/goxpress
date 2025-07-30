// Package relay 是一个类似 Express.js 的 Go Web 框架
// 定位：内置中间件模块
// 作用：提供常用的内置中间件，如日志记录和错误恢复
// 使用方法：
//   1. 使用 relay.Logger() 获取日志中间件
//   2. 使用 relay.Recover() 获取错误恢复中间件
//   3. 通过 app.Use() 注册中间件
package relay

import (
	"fmt"
	"log"
	"time"
)

// Logger 返回记录 HTTP 请求的日志中间件
// 定位：日志记录中间件
// 作用：记录每个 HTTP 请求的基本信息，包括方法、路径、客户端地址和处理时间
// 使用方法：
//   app := relay.New()
//   app.Use(relay.Logger())
//   该中间件会自动记录请求信息到标准日志输出
func Logger() HandlerFunc {
	return func(c *Context) {
		// 启动计时器
		t := time.Now()
		
		// 处理请求
		c.Next()
		
		// 记录请求日志
		log.Printf("[%s] %s %s %v", c.Request.Method, c.Request.URL.Path, c.Request.RemoteAddr, time.Since(t))
	}
}

// Recover 返回从 panic 中恢复的中间件
// 定位：错误恢复中间件
// 作用：捕获处理过程中发生的 panic，防止服务崩溃，并将 panic 转换为错误传递给错误处理中间件
// 使用方法：
//   app := relay.New()
//   app.Use(relay.Recover())
//   该中间件会自动捕获并处理 panic，防止服务中断
func Recover() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.Abort()
				// 将 panic 作为错误传递给错误处理中间件
				if e, ok := err.(error); ok {
					c.Next(e)
				} else {
					c.Next(fmt.Errorf("%v", err))
				}
			}
		}()
		c.Next()
	}
}