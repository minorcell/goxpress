// Package goxpress 是一个类似 Express.js 的 Go Web 框架
// 定位：提供请求和响应的上下文封装
// 作用：封装 HTTP 请求和响应，提供便捷的方法处理数据
// 使用方法：
//  1. 通过 c := goxpress.NewContext(w, req) 创建实例
//  2. 使用 c.Param() 获取路由参数
//  3. 使用 c.Query() 获取查询参数
//  4. 使用 c.JSON() 发送 JSON 响应
//  5. 使用 c.Next() 调用下一个中间件
package goxpress

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// 定义一个Context池，用于复用Context对象
var contextPool = sync.Pool{
	New: func() interface{} {
		return &Context{
			params: make(map[string]string),
			store:  make(map[string]interface{}),
			index:  -1,
		}
	},
}

// Context 封装了请求、响应和上下文数据
// 定位：HTTP 请求处理的上下文对象
// 作用：
//  1. 封装 *http.Request 和 http.ResponseWriter
//  2. 管理 URL 参数
//  3. 管理中间件链
//  4. 提供请求和响应的便捷方法
//  5. 管理错误和中止状态
//  6. 提供键值存储功能
//
// 使用方法：
//  1. 通过 NewContext() 创建实例
//  2. 在中间件和处理函数中使用其方法
type Context struct {
	// Embed the standard context
	// 嵌入标准 context，提供标准上下文功能
	context.Context

	// Request and response
	// HTTP 请求和响应对象
	Request  *http.Request       // 原始 HTTP 请求
	Response http.ResponseWriter // HTTP 响应写入器

	// URL parameters
	// URL 参数映射表，存储路由参数
	params map[string]string

	// middleware chain
	// 中间件链相关字段
	handlers []HandlerFunc // 中间件处理函数列表
	index    int           // 当前处理函数索引

	// error for error handling middleware
	// 错误处理相关字段
	err error // 中间件链中产生的错误

	// abort flag
	// 中止标志，用于中止中间件链执行
	aborted bool

	// store for Set/Get methods
	// 键值存储，用于在中间件间传递数据
	store map[string]interface{}

	// status code has been written flag
	// 状态码是否已写入的标志
	statusCodeWritten bool
}

// NewContext 创建一个新的 Context 实例
// 定位：Context 结构体的构造函数
// 作用：初始化 Context 实例及其所有字段
// 使用方法：
//
//	c := NewContext(w, req)
//	在框架内部使用，通常不需要用户直接调用
func NewContext(w http.ResponseWriter, req *http.Request) *Context {
	// 从池中获取Context或创建新的
	c := contextPool.Get().(*Context)

	// 初始化请求相关的字段
	c.Context = req.Context()
	c.Request = req
	c.Response = w

	// 重置其他字段
	c.index = -1
	c.aborted = false
	c.err = nil
	c.statusCodeWritten = false

	// 清空maps而不是重新分配，减少内存分配
	for k := range c.params {
		delete(c.params, k)
	}

	for k := range c.store {
		delete(c.store, k)
	}

	return c
}

// reset 重置Context状态，为复用做准备
// 定位：Context重置方法
// 作用：清理Context中的请求相关数据，为下次复用做准备
// 使用方法：框架内部使用
func (c *Context) reset() {
	// 清空maps
	for k := range c.params {
		delete(c.params, k)
	}

	for k := range c.store {
		delete(c.store, k)
	}

	// 重置其他字段
	c.index = -1
	c.aborted = false
	c.err = nil
	c.handlers = nil
	c.Context = nil
	c.Request = nil
	c.Response = nil
	c.statusCodeWritten = false
}

// Param 返回指定 URL 参数的值
// 定位：URL 参数获取方法
// 作用：获取路由中定义的参数值
// 使用方法：
//
//	// 对于路由 /users/:id
//	id := c.Param("id")
func (c *Context) Param(key string) string {
	return c.params[key]
}

// Query 返回指定查询参数的值
// 定位：查询参数获取方法
// 作用：获取 URL 查询字符串中的参数值
// 使用方法：
//
//	// 对于请求 /users?page=1&size=10
//	page := c.Query("page")
//	size := c.Query("size")
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// BindJSON 将请求体绑定到指定结构体
// 定位：请求体解析方法
// 作用：将 JSON 格式的请求体解析到 Go 结构体中
// 使用方法：
//
//	var user User
//	if err := c.BindJSON(&user); err != nil {
//	    // 处理错误
//	}
func (c *Context) BindJSON(obj interface{}) error {
	return json.NewDecoder(c.Request.Body).Decode(obj)
}

// Status 设置 HTTP 响应状态码
// 定位：HTTP 状态码设置方法
// 作用：设置响应的状态码
// 使用方法：
//
//	c.Status(200) // 设置状态码为 200
func (c *Context) Status(code int) {
	if !c.statusCodeWritten {
		c.Response.WriteHeader(code)
		c.statusCodeWritten = true
	}
}

// JSON 将指定结构体序列化为 JSON 并写入响应体
// 定位：JSON 响应发送方法
// 作用：将数据序列化为 JSON 格式并发送响应
// 使用方法：
//
//	c.JSON(200, map[string]interface{}{"message": "success"})
func (c *Context) JSON(code int, obj interface{}) error {
	if !c.statusCodeWritten {
		c.Response.Header().Set("Content-Type", "application/json")
		c.Response.WriteHeader(code)
		c.statusCodeWritten = true
	}
	return json.NewEncoder(c.Response).Encode(obj)
}

// String 将格式化字符串写入响应体
// 定位：文本响应发送方法
// 作用：将格式化字符串作为响应内容发送
// 使用方法：
//
//	c.String(200, "Hello %s", name)
func (c *Context) String(code int, format string, values ...interface{}) error {
	if !c.statusCodeWritten {
		c.Response.Header().Set("Content-Type", "text/plain; charset=utf-8")
		c.Response.WriteHeader(code)
		c.statusCodeWritten = true
	}
	_, err := c.Response.Write([]byte(fmt.Sprintf(format, values...)))
	return err
}

// Next 调用中间件链中的下一个处理函数
// 定位：中间件链控制方法
// 作用：执行中间件链中的下一个处理函数
// 使用方法：
//
//	func Middleware(c *Context) {
//	    // 前置处理
//	    c.Next()
//	    // 后置处理
//	}
func (c *Context) Next(err ...error) {
	// 如果提供了错误，则设置错误
	if len(err) > 0 && err[0] != nil {
		c.err = err[0]
	}

	// 增加索引并执行后续处理函数
	c.index++
	for c.index < len(c.handlers) {
		// 如果已中止，则返回
		if c.aborted {
			return
		}
		// 执行当前处理函数
		c.handlers[c.index](c)
		c.index++
	}
}

// Abort 中止后续中间件的执行
// 定位：中间件链中止方法
// 作用：设置中止标志，阻止后续中间件执行
// 使用方法：
//
//	func AuthMiddleware(c *Context) {
//	    if !authorized {
//	        c.Abort()
//	        c.JSON(401, map[string]interface{}{"error": "unauthorized"})
//	        return
//	    }
//	    c.Next()
//	}
func (c *Context) Abort() {
	c.aborted = true
}

// IsAborted 返回上下文是否已被中止
// 定位：中止状态检查方法
// 作用：检查中间件链是否已被中止
// 使用方法：
//
//	if c.IsAborted() {
//	    // 处理中止情况
//	}
func (c *Context) IsAborted() bool {
	return c.aborted
}

// Set 在上下文中存储键值对
// 定位：上下文数据存储方法
// 作用：在上下文中存储数据，供后续中间件使用
// 使用方法：
//
//	c.Set("user", user)
func (c *Context) Set(key string, value interface{}) {
	c.store[key] = value
}

// Get 从上下文中获取指定键的值
// 定位：上下文数据获取方法
// 作用：从上下文中获取之前存储的数据
// 使用方法：
//
//	if user, exists := c.Get("user"); exists {
//	    // 使用 user 数据
//	}
func (c *Context) Get(key string) (interface{}, bool) {
	value, exists := c.store[key]
	return value, exists
}

// MustGet 从上下文中获取指定键的值，如果不存在则 panic
// 定位：强制上下文数据获取方法
// 作用：获取必须存在的上下文数据，不存在时会 panic
// 使用方法：
//
//	user := c.MustGet("user").(User)
func (c *Context) MustGet(key string) interface{} {
	if value, exists := c.store[key]; exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

// GetString 从上下文中获取字符串类型的值
// 定位：字符串类型上下文数据获取方法
// 作用：方便地获取字符串类型的上下文数据
// 使用方法：
//
//	if name, ok := c.GetString("name"); ok {
//	    // 使用 name 字符串
//	}
func (c *Context) GetString(key string) (string, bool) {
	if val, ok := c.Get(key); ok && val != nil {
		if str, ok := val.(string); ok {
			return str, true
		}
	}
	return "", false
}
