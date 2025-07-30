// Package relay_test 提供 relay 框架的基准测试
// 定位：框架性能基准测试
// 作用：测试框架的路由匹配性能和中间件执行性能
// 使用方法：使用 go test -bench=. 命令运行基准测试
package relay_test

import (
	"net/http/httptest"
	"testing"

	"relay"
)

// BenchmarkStaticRoute 测试静态路由匹配性能
// 定位：静态路由匹配性能基准测试
// 作用：测量框架匹配静态路由的性能
// 使用方法：go test -bench=BenchmarkStaticRoute
func BenchmarkStaticRoute(b *testing.B) {
	engine := relay.New()

	// 注册静态路由
	engine.GET("/users", func(c *relay.Context) {
		c.String(200, "users")
	})

	// 创建请求
	req := httptest.NewRequest("GET", "/users", nil)

	b.ResetTimer()

	// 执行基准测试
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
	}
}

// BenchmarkDynamicRoute 测试动态路由匹配性能
// 定位：动态路由匹配性能基准测试
// 作用：测量框架匹配动态路由的性能
// 使用方法：go test -bench=BenchmarkDynamicRoute
func BenchmarkDynamicRoute(b *testing.B) {
	engine := relay.New()

	// 注册动态路由
	engine.GET("/users/:id", func(c *relay.Context) {
		id := c.Param("id")
		c.String(200, "user %s", id)
	})

	// 创建请求
	req := httptest.NewRequest("GET", "/users/123", nil)

	b.ResetTimer()

	// 执行基准测试
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
	}
}

// BenchmarkMiddleware 测试中间件执行性能
// 定位：中间件执行性能基准测试
// 作用：测量框架执行中间件的性能
// 使用方法：go test -bench=BenchmarkMiddleware
func BenchmarkMiddleware(b *testing.B) {
	engine := relay.New()

	// 注册中间件
	engine.Use(func(c *relay.Context) {
		// 空中间件，仅用于测试执行开销
		c.Next()
	})

	// 注册路由
	engine.GET("/test", func(c *relay.Context) {
		c.String(200, "test")
	})

	// 创建请求
	req := httptest.NewRequest("GET", "/test", nil)

	b.ResetTimer()

	// 执行基准测试
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
	}
}

// BenchmarkMultipleMiddleware 测试多个中间件执行性能
// 定位：多个中间件执行性能基准测试
// 作用：测量框架执行多个中间件的性能
// 使用方法：go test -bench=BenchmarkMultipleMiddleware
func BenchmarkMultipleMiddleware(b *testing.B) {
	engine := relay.New()

	// 注册多个中间件
	for i := 0; i < 5; i++ {
		engine.Use(func(c *relay.Context) {
			// 空中间件，仅用于测试执行开销
			c.Next()
		})
	}

	// 注册路由
	engine.GET("/test", func(c *relay.Context) {
		c.String(200, "test")
	})

	// 创建请求
	req := httptest.NewRequest("GET", "/test", nil)

	b.ResetTimer()

	// 执行基准测试
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
	}
}

// BenchmarkRouterMatchStatic 测试路由匹配静态路径性能
// 定位：路由匹配静态路径性能基准测试
// 作用：测量路由系统匹配静态路径的性能
// 使用方法：go test -bench=BenchmarkRouterMatchStatic
func BenchmarkRouterMatchStatic(b *testing.B) {
	// 由于 getRoute 是私有方法，我们通过 ServeHTTP 测试路由匹配性能
	engine := relay.New()

	// 注册路由
	engine.GET("/users", func(c *relay.Context) {
		c.String(200, "users")
	})

	// 创建请求
	req := httptest.NewRequest("GET", "/users", nil)

	b.ResetTimer()

	// 执行基准测试
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
	}
}

// BenchmarkRouterMatchDynamic 测试路由匹配动态路径性能
// 定位：路由匹配动态路径性能基准测试
// 作用：测量路由系统匹配动态路径的性能
// 使用方法：go test -bench=BenchmarkRouterMatchDynamic
func BenchmarkRouterMatchDynamic(b *testing.B) {
	// 由于 getRoute 是私有方法，我们通过 ServeHTTP 测试路由匹配性能
	engine := relay.New()

	// 注册路由
	engine.GET("/users/:id", func(c *relay.Context) {
		c.String(200, "user %s", c.Param("id"))
	})

	// 创建请求
	req := httptest.NewRequest("GET", "/users/123", nil)

	b.ResetTimer()

	// 执行基准测试
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)
	}
}

// BenchmarkContextCreation 测试上下文创建性能
// 定位：上下文创建性能基准测试
// 作用：测量创建 Context 实例的性能
// 使用方法：go test -bench=BenchmarkContextCreation
func BenchmarkContextCreation(b *testing.B) {
	// 创建请求和响应
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()

	// 执行基准测试
	for i := 0; i < b.N; i++ {
		_ = relay.NewContext(w, req)
	}
}
