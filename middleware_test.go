// Package relay_test 提供 relay 框架的测试
// 定位：内置中间件功能测试
// 作用：测试框架提供的内置中间件，如 Logger 和 Recover
// 使用方法：使用 go test 命令运行测试
package relay_test

import (
	"bytes"
	"log"
	"net/http/httptest"
	"strings"
	"testing"

	"relay"
)

// TestLoggerMiddleware 测试日志中间件功能
// 定位：Logger 中间件测试
// 作用：验证 Logger() 中间件是否正确记录请求信息
// 使用方法：go test -run TestLoggerMiddleware
func TestLoggerMiddleware(t *testing.T) {
	// 捕获日志输出
	var buf bytes.Buffer
	log.SetOutput(&buf)

	engine := relay.New()

	// 使用 Logger 中间件
	engine.Use(relay.Logger())

	// 注册测试路由
	engine.GET("/test", func(c *relay.Context) {
		c.String(200, "test")
	})

	// 创建请求
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// 处理请求
	engine.ServeHTTP(w, req)

	// 验证响应
	if w.Code != 200 {
		t.Errorf("Expected status code 200, but got %d", w.Code)
	}

	// 验证日志是否输出
	logOutput := buf.String()
	if !strings.Contains(logOutput, "GET") || !strings.Contains(logOutput, "/test") {
		t.Errorf("Expected log to contain request info, but got: %s", logOutput)
	}
}

// TestRecoverMiddleware 测试恢复中间件功能
// 定位：Recover 中间件测试
// 作用：验证 Recover() 中间件是否正确捕获和处理 panic
// 使用方法：go test -run TestRecoverMiddleware
func TestRecoverMiddleware(t *testing.T) {
	// 捕获日志输出
	var buf bytes.Buffer
	log.SetOutput(&buf)

	engine := relay.New()

	// 使用 Recover 中间件
	engine.Use(relay.Recover())

	// 注册会 panic 的路由
	engine.GET("/panic", func(c *relay.Context) {
		panic("test panic")
	})

	// 注册错误处理中间件来捕获 Recover 传递的错误
	var errorHandlerCalled = false
	var handledError error
	engine.UseError(func(err error, c *relay.Context) {
		errorHandlerCalled = true
		handledError = err
		c.String(500, "error: %v", err)
	})

	// 创建请求
	req := httptest.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()

	// 处理请求
	engine.ServeHTTP(w, req)

	// 验证响应
	if w.Code != 500 {
		t.Errorf("Expected status code 500, but got %d", w.Code)
	}

	// 验证错误处理中间件是否被调用
	if !errorHandlerCalled {
		t.Error("Expected error handler to be called")
	}

	// 验证错误是否正确传递
	if handledError == nil {
		t.Error("Expected error to be passed to error handler")
	} else if handledError.Error() != "test panic" {
		t.Errorf("Expected error message to be 'test panic', but got '%s'", handledError.Error())
	}

	// 验证日志是否输出
	logOutput := buf.String()
	if !strings.Contains(logOutput, "Panic recovered") || !strings.Contains(logOutput, "test panic") {
		t.Errorf("Expected log to contain panic info, but got: %s", logOutput)
	}
}

// TestMiddlewareChaining 测试中间件链功能
// 定位：中间件链测试
// 作用：验证多个中间件是否按正确顺序执行（洋葱模型）
// 使用方法：go test -run TestMiddlewareChaining
func TestMiddlewareChaining(t *testing.T) {
	engine := relay.New()

	// 记录执行顺序
	executionOrder := make([]string, 0)

	// 创建测试中间件
	middleware1 := func(c *relay.Context) {
		executionOrder = append(executionOrder, "middleware1-before")
		c.Next()
		executionOrder = append(executionOrder, "middleware1-after")
	}

	middleware2 := func(c *relay.Context) {
		executionOrder = append(executionOrder, "middleware2-before")
		c.Next()
		executionOrder = append(executionOrder, "middleware2-after")
	}

	// 注册中间件
	engine.Use(middleware1, middleware2)

	// 注册路由处理器
	engine.GET("/test", func(c *relay.Context) {
		executionOrder = append(executionOrder, "handler")
		c.String(200, "test")
	})

	// 创建请求
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// 处理请求
	engine.ServeHTTP(w, req)

	// 验证执行顺序（洋葱模型）
	expectedOrder := []string{
		"middleware1-before",
		"middleware2-before",
		"handler",
		"middleware2-after",
		"middleware1-after",
	}

	if len(executionOrder) != len(expectedOrder) {
		t.Errorf("Expected execution order length %d, but got %d", len(expectedOrder), len(executionOrder))
	}

	for i, expected := range expectedOrder {
		if executionOrder[i] != expected {
			t.Errorf("Expected execution order %v, but got %v", expectedOrder, executionOrder)
			break
		}
	}
}

// TestMiddlewareAbort 测试中间件中止功能
// 定位：中间件中止测试
// 作用：验证中间件是否能正确中止后续处理
// 使用方法：go test -run TestMiddlewareAbort
func TestMiddlewareAbort(t *testing.T) {
	engine := relay.New()

	// 记录是否执行
	middleware2Executed := false
	handlerExecuted := false

	// 创建会中止的中间件
	abortMiddleware := func(c *relay.Context) {
		c.Abort()
		c.String(401, "unauthorized")
	}

	// 创建第二个中间件
	secondMiddleware := func(c *relay.Context) {
		middleware2Executed = true
		c.Next()
	}

	// 创建处理函数
	handler := func(c *relay.Context) {
		handlerExecuted = true
		c.String(200, "ok")
	}

	// 注册中间件和路由
	engine.Use(abortMiddleware, secondMiddleware)
	engine.GET("/test", handler)

	// 创建请求
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// 处理请求
	engine.ServeHTTP(w, req)

	// 验证响应
	if w.Code != 401 {
		t.Errorf("Expected status code 401, but got %d", w.Code)
	}

	if w.Body.String() != "unauthorized" {
		t.Errorf("Expected body 'unauthorized', but got '%s'", w.Body.String())
	}

	// 验证后续中间件和处理函数未执行
	if middleware2Executed {
		t.Error("Expected second middleware not to be executed after abort")
	}

	if handlerExecuted {
		t.Error("Expected handler not to be executed after abort")
	}
}
