// Package goxpress 是一个类似 Express.js 的 Go Web 框架
// 定位：框架的核心引擎，负责协调路由、中间件和请求处理
// 作用：实现 http.Handler 接口，管理路由和中间件，提供 HTTP 服务启动功能
// 使用方法：
//  1. 创建引擎实例：app := goxpress.New()
//  2. 注册中间件：app.Use(middleware)
//  3. 定义路由：app.GET("/path", handler)
//  4. 启动服务：app.Listen(":8080", nil)
package goxpress

import (
	"net/http"
)

// HandlerFunc 定义处理 HTTP 请求的函数类型
// 定位：中间件和路由处理器的统一函数签名
// 作用：为中间件和路由处理函数提供统一的接口
// 使用方法：实现此函数类型，接收 *Context 参数，处理请求逻辑
type HandlerFunc func(*Context)

// ErrorHandlerFunc 定义处理错误的函数类型
// 定位：专门用于处理错误的中间件函数签名
// 作用：捕获和处理在请求处理过程中发生的错误
// 使用方法：实现此函数类型，接收 error 和 *Context 参数，处理错误逻辑
type ErrorHandlerFunc func(error, *Context)

// Engine 是 goxpress 框架的主要结构体
// 定位：框架的核心引擎，协调所有组件
// 作用：
//  1. 实现 http.Handler 接口
//  2. 管理路由系统
//  3. 管理全局中间件
//  4. 管理错误处理中间件
//
// 使用方法：
//  1. 通过 goxpress.New() 创建实例
//  2. 使用 Use() 方法注册中间件
//  3. 使用 UseError() 方法注册错误处理中间件
//  4. 使用 HTTP 方法函数（GET/POST等）定义路由
//  5. 使用 Listen() 启动 HTTP 服务
type Engine struct {
	router        *Router            // 路由器，管理所有路由
	middlewares   []HandlerFunc      // 全局中间件列表
	errorHandlers []ErrorHandlerFunc // 错误处理中间件列表
}

// New 创建一个新的 goxpress 引擎实例
// 定位：Engine 结构体的构造函数
// 作用：初始化 Engine 实例及其依赖组件
// 使用方法：app := goxpress.New()
func New() *Engine {
	engine := &Engine{
		router:        NewRouter(),
		middlewares:   make([]HandlerFunc, 0),
		errorHandlers: make([]ErrorHandlerFunc, 0),
	}
	return engine
}

// Use 注册全局中间件函数
// 定位：中间件注册方法
// 作用：将中间件添加到全局中间件链中，对所有请求生效
// 使用方法：
//
//	app.Use(Logger(), Recover())
//	返回 *Engine 实例，支持链式调用
func (e *Engine) Use(middleware ...HandlerFunc) *Engine {
	e.middlewares = append(e.middlewares, middleware...)
	return e
}

// UseError 注册错误处理中间件函数
// 定位：错误处理中间件注册方法
// 作用：将错误处理中间件添加到错误处理链中
// 使用方法：
//
//	app.UseError(func(err error, c *Context) {
//	    // 错误处理逻辑
//	})
//	返回 *Engine 实例，支持链式调用
func (e *Engine) UseError(handler ...ErrorHandlerFunc) *Engine {
	e.errorHandlers = append(e.errorHandlers, handler...)
	return e
}

// GET 注册一个新的 GET 路由
// 定位：HTTP GET 方法路由注册方法
// 作用：为指定路径注册 GET 请求处理函数
// 使用方法：
//
//	app.GET("/users/:id", getUserHandler)
//	返回 *Engine 实例，支持链式调用
func (e *Engine) GET(pattern string, handlers ...HandlerFunc) *Engine {
	e.router.GET(pattern, handlers...)
	return e
}

// POST 注册一个新的 POST 路由
// 定位：HTTP POST 方法路由注册方法
// 作用：为指定路径注册 POST 请求处理函数
// 使用方法：
//
//	app.POST("/users", createUserHandler)
//	返回 *Engine 实例，支持链式调用
func (e *Engine) POST(pattern string, handlers ...HandlerFunc) *Engine {
	e.router.POST(pattern, handlers...)
	return e
}

// PUT 注册一个新的 PUT 路由
// 定位：HTTP PUT 方法路由注册方法
// 作用：为指定路径注册 PUT 请求处理函数
// 使用方法：
//
//	app.PUT("/users/:id", updateUserHandler)
//	返回 *Engine 实例，支持链式调用
func (e *Engine) PUT(pattern string, handlers ...HandlerFunc) *Engine {
	e.router.PUT(pattern, handlers...)
	return e
}

// DELETE 注册一个新的 DELETE 路由
// 定位：HTTP DELETE 方法路由注册方法
// 作用：为指定路径注册 DELETE 请求处理函数
// 使用方法：
//
//	app.DELETE("/users/:id", deleteUserHandler)
//	返回 *Engine 实例，支持链式调用
func (e *Engine) DELETE(pattern string, handlers ...HandlerFunc) *Engine {
	e.router.DELETE(pattern, handlers...)
	return e
}

// PATCH 注册一个新的 PATCH 路由
// 定位：HTTP PATCH 方法路由注册方法
// 作用：为指定路径注册 PATCH 请求处理函数
// 使用方法：
//
//	app.PATCH("/users/:id", patchUserHandler)
//	返回 *Engine 实例，支持链式调用
func (e *Engine) PATCH(pattern string, handlers ...HandlerFunc) *Engine {
	e.router.PATCH(pattern, handlers...)
	return e
}

// HEAD 注册一个新的 HEAD 路由
// 定位：HTTP HEAD 方法路由注册方法
// 作用：为指定路径注册 HEAD 请求处理函数
// 使用方法：
//
//	app.HEAD("/users/:id", headUserHandler)
//	返回 *Engine 实例，支持链式调用
func (e *Engine) HEAD(pattern string, handlers ...HandlerFunc) *Engine {
	e.router.HEAD(pattern, handlers...)
	return e
}

// OPTIONS 注册一个新的 OPTIONS 路由
// 定位：HTTP OPTIONS 方法路由注册方法
// 作用：为指定路径注册 OPTIONS 请求处理函数
// 使用方法：
//
//	app.OPTIONS("/users", optionsUserHandler)
//	返回 *Engine 实例，支持链式调用
func (e *Engine) OPTIONS(pattern string, handlers ...HandlerFunc) *Engine {
	e.router.OPTIONS(pattern, handlers...)
	return e
}

// Route 创建一个新的路由组
// 定位：路由组创建方法
// 作用：创建具有共同前缀的路由组，便于组织相关路由
// 使用方法：
//
//	api := app.Route("/api")
//	api.GET("/users", getUsersHandler)
func (e *Engine) Route(prefix string) *Router {
	return e.router.Group(prefix)
}

// ServeHTTP 实现 http.Handler 接口
// 定位：HTTP 请求处理入口
// 作用：处理所有进入的 HTTP 请求，协调路由匹配和中间件执行
// 使用方法：由 HTTP 服务器自动调用，无需手动调用
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 从池中获取Context
	c := NewContext(w, req)

	// 确保请求处理完成后将Context放回池中
	defer func() {
		c.reset()
		contextPool.Put(c)
	}()

	// 查找匹配的路由
	node, params := e.router.getRoute(req.Method, req.URL.Path)

	// 设置 URL 参数
	if params != nil {
		c.params = params
	}

	// 准备中间件链
	handlers := make([]HandlerFunc, 0)
	handlers = append(handlers, e.middlewares...)

	// 添加路由处理器（如果路由存在）
	if node != nil {
		handlers = append(handlers, node.handlers...)
	} else {
		// 如果没有找到路由，设置404处理器
		handlers = append(handlers, func(c *Context) {
			c.Status(http.StatusNotFound)
			c.String(http.StatusNotFound, "404 page not found")
		})
	}

	c.handlers = handlers

	// 执行中间件链
	c.Next()

	// 处理错误（如果有）
	if c.err != nil && len(e.errorHandlers) > 0 {
		for _, handler := range e.errorHandlers {
			handler(c.err, c)
		}
	}
}

// Listen 启动 HTTP 服务器
// 定位：HTTP 服务器启动方法
// 作用：在指定地址启动 HTTP 服务器
// 使用方法：
//
//	app.Listen(":8080", func() { fmt.Println("Server started") })
func (e *Engine) Listen(addr string, cb func()) error {
	server := &http.Server{
		Addr:    addr,
		Handler: e,
	}

	if cb != nil {
		cb()
	}

	return server.ListenAndServe()
}

// ListenTLS 启动 HTTPS 服务器
// 定位：HTTPS 服务器启动方法
// 作用：在指定地址启动 HTTPS 服务器
// 使用方法：
//
//	app.ListenTLS(":443", "cert.pem", "key.pem", func() { fmt.Println("HTTPS Server started") })
func (e *Engine) ListenTLS(addr, certFile, keyFile string, cb func()) error {
	server := &http.Server{
		Addr:    addr,
		Handler: e,
	}

	if cb != nil {
		cb()
	}

	return server.ListenAndServeTLS(certFile, keyFile)
}
