// Package relay_test 提供 relay 框架的测试
// 定位：Router 和路由功能测试
// 作用：测试 Router 的路由注册、匹配和路由组功能
// 使用方法：使用 go test 命令运行测试
package relay_test

import (
	"net/http/httptest"
	"testing"

	"relay"
)

// TestRouterNew 测试 Router 的创建
// 定位：Router 构造函数测试
// 作用：验证 NewRouter() 函数是否正确创建 Router 实例
// 使用方法：go test -run TestRouterNew
func TestRouterNew(t *testing.T) {
	router := relay.NewRouter()

	if router == nil {
		t.Error("Expected router to be created, but got nil")
	}
}

// TestRouterUse 测试路由器中间件注册功能
// 定位：路由器中间件注册测试
// 作用：验证 Use() 方法是否正确注册中间件
// 使用方法：go test -run TestRouterUse
func TestRouterUse(t *testing.T) {
	router := relay.NewRouter()

	// 创建一个简单的测试中间件
	middleware := func(c *relay.Context) {}

	// 注册中间件
	result := router.Use(middleware)

	// 验证返回值是否正确（支持链式调用）
	if result != router {
		t.Error("Expected Use() to return router for chaining")
	}
}

// TestRouterGroup 测试路由组功能
// 定位：路由组测试
// 作用：验证 Group() 方法是否正确创建路由组
// 使用方法：go test -run TestRouterGroup
func TestRouterGroup(t *testing.T) {
	router := relay.NewRouter()

	group := router.Group("/api")

	if group == nil {
		t.Error("Expected Group() to return a router")
	}

	// 验证前缀是否正确设置
	// 注意：由于 prefix 字段是私有的，我们无法直接访问它
	// 但在实际使用中，这个前缀会被用于路由匹配
}

// TestRouterHTTPMethods 测试路由器的 HTTP 方法路由注册
// 定位：路由器 HTTP 方法路由注册测试
// 作用：验证路由器各种 HTTP 方法的路由注册功能
// 使用方法：go test -run TestRouterHTTPMethods
func TestRouterHTTPMethods(t *testing.T) {
	router := relay.NewRouter()

	// 测试处理函数
	handler := func(c *relay.Context) {}

	// 测试各种 HTTP 方法
	methods := []func(string, ...relay.HandlerFunc) *relay.Router{
		router.GET,
		router.POST,
		router.PUT,
		router.PATCH,
		router.DELETE,
		router.HEAD,
		router.OPTIONS,
	}

	for _, method := range methods {
		result := method("/test", handler)
		if result != router {
			t.Error("Expected HTTP method functions to return router for chaining")
		}
	}
}

// TestRouterStaticRoutes 测试静态路由匹配
// 定位：静态路由匹配测试
// 作用：验证路由器是否正确匹配静态路由
// 使用方法：go test -run TestRouterStaticRoutes
func TestRouterStaticRoutes(t *testing.T) {
	engine := relay.New()

	// 注册静态路由
	called := false
	engine.GET("/users", func(c *relay.Context) {
		called = true
		c.String(200, "users")
	})

	// 创建请求
	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()

	// 处理请求
	engine.ServeHTTP(w, req)

	// 验证路由是否匹配
	if !called {
		t.Error("Expected handler to be called for static route")
	}

	if w.Code != 200 {
		t.Errorf("Expected status code 200, but got %d", w.Code)
	}
}

// TestRouterDynamicRoutes 测试动态路由匹配
// 定位：动态路由匹配测试
// 作用：验证路由器是否正确匹配包含参数的动态路由
// 使用方法：go test -run TestRouterDynamicRoutes
func TestRouterDynamicRoutes(t *testing.T) {
	engine := relay.New()

	// 注册动态路由
	var id string
	engine.GET("/users/:id", func(c *relay.Context) {
		id = c.Param("id")
		c.String(200, "user %s", id)
	})

	// 创建请求
	req := httptest.NewRequest("GET", "/users/123", nil)
	w := httptest.NewRecorder()

	// 处理请求
	engine.ServeHTTP(w, req)

	// 验证路由是否匹配且参数正确
	if id != "123" {
		t.Errorf("Expected id to be '123', but got '%s'", id)
	}

	if w.Code != 200 {
		t.Errorf("Expected status code 200, but got %d", w.Code)
	}
}

// TestRouterWildcardRoutes 测试通配符路由匹配
// 定位：通配符路由匹配测试
// 作用：验证路由器是否正确匹配包含通配符的路由
// 使用方法：go test -run TestRouterWildcardRoutes
func TestRouterWildcardRoutes(t *testing.T) {
	engine := relay.New()

	// 注册通配符路由
	var path string
	engine.GET("/static/*filepath", func(c *relay.Context) {
		path = c.Param("filepath")
		c.String(200, "file %s", path)
	})

	// 创建请求
	req := httptest.NewRequest("GET", "/static/css/style.css", nil)
	w := httptest.NewRecorder()

	// 处理请求
	engine.ServeHTTP(w, req)

	// 验证路由是否匹配且参数正确
	if path != "css/style.css" {
		t.Errorf("Expected path to be 'css/style.css', but got '%s'", path)
	}

	if w.Code != 200 {
		t.Errorf("Expected status code 200, but got %d", w.Code)
	}
}

// TestRouterComplexRoutes 测试复杂路由匹配
// 定位：复杂路由匹配测试
// 作用：验证路由器是否正确匹配包含多个参数和静态段的复杂路由
// 使用方法：go test -run TestRouterComplexRoutes
func TestRouterComplexRoutes(t *testing.T) {
	engine := relay.New()

	// 注册复杂路由
	var lang, section, id string
	engine.GET("/docs/:lang/:section/:id", func(c *relay.Context) {
		lang = c.Param("lang")
		section = c.Param("section")
		id = c.Param("id")
		c.String(200, "doc %s/%s/%s", lang, section, id)
	})

	// 创建请求
	req := httptest.NewRequest("GET", "/docs/en/getting-started/installation", nil)
	w := httptest.NewRecorder()

	// 处理请求
	engine.ServeHTTP(w, req)

	// 验证路由是否匹配且参数正确
	if lang != "en" {
		t.Errorf("Expected lang to be 'en', but got '%s'", lang)
	}

	if section != "getting-started" {
		t.Errorf("Expected section to be 'getting-started', but got '%s'", section)
	}

	if id != "installation" {
		t.Errorf("Expected id to be 'installation', but got '%s'", id)
	}

	if w.Code != 200 {
		t.Errorf("Expected status code 200, but got %d", w.Code)
	}
}

// TestRouterMethodSpecificity 测试路由方法特异性
// 定位：路由方法特异性测试
// 作用：验证路由器是否正确区分不同 HTTP 方法的相同路径
// 使用方法：go test -run TestRouterMethodSpecificity
func TestRouterMethodSpecificity(t *testing.T) {
	engine := relay.New()

	// 注册相同路径的不同方法
	getCalled := false
	postCalled := false

	engine.GET("/users", func(c *relay.Context) {
		getCalled = true
		c.String(200, "get users")
	})

	engine.POST("/users", func(c *relay.Context) {
		postCalled = true
		c.String(200, "post users")
	})

	// 测试 GET 请求
	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if !getCalled {
		t.Error("Expected GET handler to be called")
	}

	if postCalled {
		t.Error("Expected POST handler not to be called for GET request")
	}

	// 重置标志
	getCalled = false
	postCalled = false

	// 测试 POST 请求
	req = httptest.NewRequest("POST", "/users", nil)
	w = httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	if !postCalled {
		t.Error("Expected POST handler to be called")
	}

	if getCalled {
		t.Error("Expected GET handler not to be called for POST request")
	}
}
