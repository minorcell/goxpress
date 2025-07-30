// Package relay_test 提供 relay 框架的测试
// 定位：框架核心功能测试
// 作用：测试 Engine 的基本功能，包括路由注册、中间件和请求处理
// 使用方法：使用 go test 命令运行测试
package relay_test

import (
	"net/http/httptest"
	"testing"

	"relay"
)

// TestEngineNew 测试 Engine 的创建
// 定位：Engine 构造函数测试
// 作用：验证 New() 函数是否正确创建 Engine 实例
// 使用方法：go test -run TestEngineNew
func TestEngineNew(t *testing.T) {
	engine := relay.New()

	if engine == nil {
		t.Error("Expected engine to be created, but got nil")
	}
}

// TestEngineUse 测试中间件注册功能
// 定位：中间件注册测试
// 作用：验证 Use() 方法是否正确注册中间件
// 使用方法：go test -run TestEngineUse
func TestEngineUse(t *testing.T) {
	engine := relay.New()

	// 创建一个简单的测试中间件
	middleware := func(c *relay.Context) {}

	// 注册中间件
	result := engine.Use(middleware)

	// 验证返回值是否正确（支持链式调用）
	if result != engine {
		t.Error("Expected Use() to return engine for chaining")
	}
}

// TestEngineUseError 测试错误处理中间件注册功能
// 定位：错误处理中间件注册测试
// 作用：验证 UseError() 方法是否正确注册错误处理中间件
// 使用方法：go test -run TestEngineUseError
func TestEngineUseError(t *testing.T) {
	engine := relay.New()

	// 创建一个简单的错误处理中间件
	errorHandler := func(err error, c *relay.Context) {}

	// 注册错误处理中间件
	result := engine.UseError(errorHandler)

	// 验证返回值是否正确（支持链式调用）
	if result != engine {
		t.Error("Expected UseError() to return engine for chaining")
	}
}

// TestEngineHTTPMethods 测试 HTTP 方法路由注册
// 定位：HTTP 方法路由注册测试
// 作用：验证各种 HTTP 方法的路由注册功能
// 使用方法：go test -run TestEngineHTTPMethods
func TestEngineHTTPMethods(t *testing.T) {
	engine := relay.New()

	// 测试处理函数
	handler := func(c *relay.Context) {}

	// 测试各种 HTTP 方法
	methods := []func(string, ...relay.HandlerFunc) *relay.Engine{
		engine.GET,
		engine.POST,
		engine.PUT,
		engine.PATCH,
		engine.DELETE,
		engine.HEAD,
		engine.OPTIONS,
	}

	for _, method := range methods {
		result := method("/test", handler)
		if result != engine {
			t.Error("Expected HTTP method functions to return engine for chaining")
		}
	}
}

// TestEngineRoute 测试路由组功能
// 定位：路由组测试
// 作用：验证 Route() 方法是否正确创建路由组
// 使用方法：go test -run TestEngineRoute
func TestEngineRoute(t *testing.T) {
	engine := relay.New()

	router := engine.Route("/api")

	if router == nil {
		t.Error("Expected Route() to return a router")
	}
}

// TestEngineServeHTTP 测试 HTTP 请求处理
// 定位：HTTP 请求处理测试
// 作用：验证 ServeHTTP 方法是否正确处理 HTTP 请求
// 使用方法：go test -run TestEngineServeHTTP
func TestEngineServeHTTP(t *testing.T) {
	engine := relay.New()

	// 注册一个简单的路由
	engine.GET("/", func(c *relay.Context) {
		c.String(200, "Hello World")
	})

	// 创建一个 HTTP 请求
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// 处理请求
	engine.ServeHTTP(w, req)

	// 验证响应
	if w.Code != 200 {
		t.Errorf("Expected status code 200, but got %d", w.Code)
	}

	if w.Body.String() != "Hello World" {
		t.Errorf("Expected body 'Hello World', but got '%s'", w.Body.String())
	}
}
