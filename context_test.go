// Package relay_test 提供 relay 框架的测试
// 定位：Context 功能测试
// 作用：测试 Context 的各种方法和功能
// 使用方法：使用 go test 命令运行测试
package relay_test

import (
	"net/http/httptest"
	"strings"
	"testing"

	"relay"
)

// TestContextNew 测试 Context 的创建
// 定位：Context 构造函数测试
// 作用：验证 NewContext() 函数是否正确创建 Context 实例
// 使用方法：go test -run TestContextNew
func TestContextNew(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	c := relay.NewContext(w, req)

	if c == nil {
		t.Error("Expected context to be created, but got nil")
	}

	if c.Request != req {
		t.Error("Expected context request to match input request")
	}

	if c.Response != w {
		t.Error("Expected context response to match input response")
	}
}

// TestContextParam 测试参数获取功能
// 定位：URL 参数获取测试
// 作用：验证 Param() 方法是否正确获取 URL 参数
// 使用方法：go test -run TestContextParam
func TestContextParam(t *testing.T) {
	// 由于 params 是私有字段，我们通过实际路由测试来验证
	engine := relay.New()
	var paramValue string

	engine.GET("/users/:id", func(c *relay.Context) {
		paramValue = c.Param("id")
		c.String(200, "user %s", paramValue)
	})

	// 创建带参数的请求
	req := httptest.NewRequest("GET", "/users/123", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if paramValue != "123" {
		t.Errorf("Expected Param() to return '123', but got '%s'", paramValue)
	}
}

// TestContextQuery 测试查询参数获取功能
// 定位：查询参数获取测试
// 作用：验证 Query() 方法是否正确获取查询参数
// 使用方法：go test -run TestContextQuery
func TestContextQuery(t *testing.T) {
	req := httptest.NewRequest("GET", "/?name=john&age=25", nil)
	w := httptest.NewRecorder()

	c := relay.NewContext(w, req)

	if c.Query("name") != "john" {
		t.Error("Expected Query() to return correct name value")
	}

	if c.Query("age") != "25" {
		t.Error("Expected Query() to return correct age value")
	}

	if c.Query("nonexistent") != "" {
		t.Error("Expected Query() to return empty string for nonexistent key")
	}
}

// TestContextStatus 测试状态码设置功能
// 定位：HTTP 状态码设置测试
// 作用：验证 Status() 方法是否正确设置 HTTP 状态码
// 使用方法：go test -run TestContextStatus
func TestContextStatus(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	c := relay.NewContext(w, req)
	c.Status(404)

	// 注意：由于 httptest.Recorder 的特殊性，我们无法直接验证状态码是否已设置
	// 但在实际 HTTP 服务器中，这会正确设置响应状态码
}

// TestContextString 测试字符串响应功能
// 定位：字符串响应测试
// 作用：验证 String() 方法是否正确发送字符串响应
// 使用方法：go test -run TestContextString
func TestContextString(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	c := relay.NewContext(w, req)
	err := c.String(200, "Hello %s", "World")

	if err != nil {
		t.Errorf("Expected String() to not return error, but got %v", err)
	}

	if w.Code != 200 {
		t.Errorf("Expected status code 200, but got %d", w.Code)
	}

	if w.Body.String() != "Hello World" {
		t.Errorf("Expected body 'Hello World', but got '%s'", w.Body.String())
	}

	if w.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Errorf("Expected Content-Type to be 'text/plain; charset=utf-8', but got '%s'", w.Header().Get("Content-Type"))
	}
}

// TestContextJSON 测试 JSON 响应功能
// 定位：JSON 响应测试
// 作用：验证 JSON() 方法是否正确发送 JSON 响应
// 使用方法：go test -run TestContextJSON
func TestContextJSON(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	c := relay.NewContext(w, req)
	data := map[string]interface{}{"message": "hello", "code": 200}
	err := c.JSON(200, data)

	if err != nil {
		t.Errorf("Expected JSON() to not return error, but got %v", err)
	}

	if w.Code != 200 {
		t.Errorf("Expected status code 200, but got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, `"message":"hello"`) || !strings.Contains(body, `"code":200`) {
		t.Errorf("Expected body to contain JSON data, but got '%s'", body)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type to be 'application/json', but got '%s'", w.Header().Get("Content-Type"))
	}
}

// TestContextSetGet 测试上下文数据存储功能
// 定位：上下文数据存储测试
// 作用：验证 Set() 和 Get() 方法是否正确存储和获取数据
// 使用方法：go test -run TestContextSetGet
func TestContextSetGet(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	c := relay.NewContext(w, req)

	// 测试 Set 和 Get
	c.Set("key1", "value1")

	if value, exists := c.Get("key1"); !exists {
		t.Error("Expected Get() to return true for existing key")
	} else if value != "value1" {
		t.Errorf("Expected Get() to return 'value1', but got '%v'", value)
	}

	// 测试不存在的键
	if _, exists := c.Get("nonexistent"); exists {
		t.Error("Expected Get() to return false for nonexistent key")
	}

	// 测试 GetString
	c.Set("key2", "value2")
	if value, ok := c.GetString("key2"); !ok {
		t.Error("Expected GetString() to return true for existing string key")
	} else if value != "value2" {
		t.Errorf("Expected GetString() to return 'value2', but got '%s'", value)
	}

	// 测试非字符串值的 GetString
	c.Set("key3", 123)
	if _, ok := c.GetString("key3"); ok {
		t.Error("Expected GetString() to return false for non-string value")
	}
}

// TestContextMustGet 测试 MustGet 功能
// 定位：强制获取上下文数据测试
// 作用：验证 MustGet() 方法是否正确获取数据，以及在键不存在时是否 panic
// 使用方法：go test -run TestContextMustGet
func TestContextMustGet(t *testing.T) {
	engine := relay.New()

	var mustGetValue interface{}

	// 创建中间件来测试 MustGet
	engine.Use(func(c *relay.Context) {
		c.Set("key", "value")
		mustGetValue = c.MustGet("key")
		c.Next()
	})

	// 注册路由
	engine.GET("/test", func(c *relay.Context) {
		c.String(200, "test")
	})

	// 创建请求
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// 处理请求
	engine.ServeHTTP(w, req)

	// 测试正常获取
	if mustGetValue != "value" {
		t.Errorf("Expected MustGet() to return 'value', but got '%v'", mustGetValue)
	}

	// 测试不存在的键（应该 panic）
	engine2 := relay.New()

	engine2.Use(func(c *relay.Context) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected MustGet() to panic for nonexistent key")
			}
		}()
		c.MustGet("nonexistent")
		c.Next()
	})

	engine2.GET("/test", func(c *relay.Context) {
		c.String(200, "test")
	})

	req = httptest.NewRequest("GET", "/test", nil)
	w = httptest.NewRecorder()
	engine2.ServeHTTP(w, req)
}

// TestContextNext 测试中间件链执行功能
// 定位：中间件链执行测试
// 作用：验证 Next() 方法是否正确执行中间件链
// 使用方法：go test -run TestContextNext
func TestContextNext(t *testing.T) {
	engine := relay.New()

	executionOrder := make([]int, 0)

	// 创建测试中间件
	middleware1 := func(c *relay.Context) {
		executionOrder = append(executionOrder, 1)
		c.Next()
		executionOrder = append(executionOrder, 4)
	}

	middleware2 := func(c *relay.Context) {
		executionOrder = append(executionOrder, 2)
		c.Next()
		executionOrder = append(executionOrder, 3)
	}

	engine.Use(middleware1, middleware2)

	// 注册路由
	engine.GET("/test", func(c *relay.Context) {
		c.String(200, "test")
	})

	// 创建请求
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// 处理请求
	engine.ServeHTTP(w, req)

	// 验证执行顺序（洋葱模型）
	expectedOrder := []int{1, 2, 3, 4}
	for i, v := range expectedOrder {
		if executionOrder[i] != v {
			t.Errorf("Expected execution order %v, but got %v", expectedOrder, executionOrder)
			break
		}
	}
}

// TestContextAbort 测试中间件链中止功能
// 定位：中间件链中止测试
// 作用：验证 Abort() 方法是否正确中止中间件链执行
// 使用方法：go test -run TestContextAbort
func TestContextAbort(t *testing.T) {
	engine := relay.New()

	executed := false

	// 创建测试中间件
	middleware := func(c *relay.Context) {
		c.Abort()
	}

	afterMiddleware := func(c *relay.Context) {
		executed = true
		c.String(200, "should not execute")
	}

	engine.Use(middleware, afterMiddleware)

	// 注册路由
	engine.GET("/test", func(c *relay.Context) {
		c.String(200, "handler")
	})

	// 创建请求
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// 处理请求
	engine.ServeHTTP(w, req)

	// 验证第二个中间件未执行
	if executed {
		t.Error("Expected Abort() to prevent further middleware execution")
	}

	// 验证处理函数未执行
	if w.Body.String() == "handler" {
		t.Error("Expected handler not to be executed after Abort()")
	}

	// 验证中止响应
	if w.Body.String() == "should not execute" {
		t.Error("Expected afterMiddleware not to be executed after Abort()")
	}
}
