# Relay - A Lightweight Express-style Web Framework for Go

![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![License](https://img.shields.io/badge/license-MIT-blue)
![Version](https://img.shields.io/badge/version-1.0.0-orange)

## 项目简介

Relay 是一个轻量级的 Go Web 框架，旨在为开发者提供类似 Node.js 中 Express 框架的开发体验。它采用高性能的 Radix Tree 路由算法，支持 Express.js 风格的 API 设计，提供灵活的中间件机制和高效的请求处理能力。

## 快速入门

### 安装

```bash
# 安装 relay 框架
go get github.com/minorcell/relay
```

### 简单示例

```go
package main

import (
	"log"
	"relay"
)

func main() {
	// 创建一个新的 relay 引擎
	app := relay.New().
		Use(relay.Logger()).
		Use(relay.Recover())

	// 定义一个简单的路由
	app.GET("/", func(c *relay.Context) {
		c.String(200, "Hello, Relay!")
	})

	// 定义一个带参数的路由
	app.GET("/users/:id", func(c *relay.Context) {
		id := c.Param("id")
		c.JSON(200, map[string]interface{}{
			"id":   id,
			"name": "User " + id,
		})
	})

	// 定义一个路由组
	api := app.Route("/api").
		Use(func(c *relay.Context) {
			c.Set("version", "v1")
			c.Next()
	})

	api.GET("/status", func(c *relay.Context) {
		version, _ := c.GetString("version")
		c.JSON(200, map[string]interface{}{
			"status":  "ok",
			"version": version,
		})
	})

	// 错误处理中间件
	app.UseError(func(err error, c *relay.Context) {
		c.JSON(500, map[string]interface{}{
			"error": err.Error(),
		})
	})

	// 启动服务器
	log.Println("Server starting on :8080")
	err := app.Listen(":8080", nil)
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
```

### 运行示例

```bash
# 运行示例程序
go run main.go

# 测试 API
curl http://localhost:8080/
curl http://localhost:8080/users/123
curl http://localhost:8080/api/status
curl http://localhost:8080/not-exit  # 404 测试
```

## 许可证

Relay 框架采用 MIT 许可证，详情请参见 [LICENSE](LICENSE) 文件。
