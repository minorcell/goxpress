package goxpress

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"
	"time"
)

// BenchmarkEngine_Simple tests basic engine performance
func BenchmarkEngine_Simple(b *testing.B) {
	app := New()
	app.GET("/", func(c *Context) {
		c.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// BenchmarkEngine_JSON tests JSON response performance
func BenchmarkEngine_JSON(b *testing.B) {
	app := New()
	app.GET("/json", func(c *Context) {
		c.JSON(200, map[string]interface{}{
			"message": "Hello, World!",
			"status":  "success",
			"data":    []int{1, 2, 3, 4, 5},
		})
	})

	req := httptest.NewRequest("GET", "/json", nil)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// BenchmarkEngine_Params tests route parameter extraction performance
func BenchmarkEngine_Params(b *testing.B) {
	app := New()
	app.GET("/users/:id/posts/:postId", func(c *Context) {
		userID := c.Param("id")
		postID := c.Param("postId")
		c.JSON(200, map[string]string{
			"user_id": userID,
			"post_id": postID,
		})
	})

	req := httptest.NewRequest("GET", "/users/123/posts/456", nil)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// BenchmarkEngine_QueryParams tests query parameter parsing performance
func BenchmarkEngine_QueryParams(b *testing.B) {
	app := New()
	app.GET("/search", func(c *Context) {
		query := c.Query("q")
		page := c.Query("page")
		limit := c.Query("limit")
		c.JSON(200, map[string]string{
			"query": query,
			"page":  page,
			"limit": limit,
		})
	})

	req := httptest.NewRequest("GET", "/search?q=golang&page=1&limit=10", nil)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// BenchmarkEngine_PostJSON tests JSON request parsing performance
func BenchmarkEngine_PostJSON(b *testing.B) {
	app := New()
	app.POST("/users", func(c *Context) {
		var user struct {
			Name  string `json:"name"`
			Email string `json:"email"`
			Age   int    `json:"age"`
		}
		if err := c.BindJSON(&user); err != nil {
			c.JSON(400, map[string]string{"error": "invalid json"})
			return
		}
		c.JSON(201, user)
	})

	jsonData := `{"name":"John Doe","email":"john@example.com","age":30}`
	reqBody := bytes.NewBufferString(jsonData)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/users", bytes.NewBufferString(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
		reqBody.Reset()
		reqBody.WriteString(jsonData)
	}
}

// BenchmarkEngine_Middleware1 tests performance with 1 middleware
func BenchmarkEngine_Middleware1(b *testing.B) {
	app := New()
	app.Use(func(c *Context) {
		c.Next()
	})
	app.GET("/", func(c *Context) {
		c.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// BenchmarkEngine_Middleware5 tests performance with 5 middlewares
func BenchmarkEngine_Middleware5(b *testing.B) {
	app := New()
	for i := 0; i < 5; i++ {
		app.Use(func(c *Context) {
			c.Next()
		})
	}
	app.GET("/", func(c *Context) {
		c.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// BenchmarkEngine_Middleware10 tests performance with 10 middlewares
func BenchmarkEngine_Middleware10(b *testing.B) {
	app := New()
	for i := 0; i < 10; i++ {
		app.Use(func(c *Context) {
			c.Next()
		})
	}
	app.GET("/", func(c *Context) {
		c.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// BenchmarkEngine_BuiltinMiddleware tests performance with built-in middlewares
func BenchmarkEngine_BuiltinMiddleware(b *testing.B) {
	app := New()
	app.Use(Logger())
	app.Use(Recover())
	app.GET("/", func(c *Context) {
		c.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// BenchmarkRouter_StaticRoutes tests static route performance
func BenchmarkRouter_StaticRoutes(b *testing.B) {
	router := NewRouter()
	handler := func(c *Context) { c.String(200, "OK") }

	// Add many static routes
	routes := []string{
		"/",
		"/api",
		"/api/v1",
		"/api/v1/users",
		"/api/v1/posts",
		"/api/v1/comments",
		"/api/v2",
		"/api/v2/users",
		"/api/v2/posts",
		"/admin",
		"/admin/users",
		"/admin/settings",
		"/public",
		"/public/css",
		"/public/js",
		"/public/images",
	}

	for _, route := range routes {
		router.GET(route, handler)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		router.getRoute("GET", "/api/v1/users")
	}
}

// BenchmarkRouter_ParamRoutes tests parameterized route performance
func BenchmarkRouter_ParamRoutes(b *testing.B) {
	router := NewRouter()
	handler := func(c *Context) { c.String(200, "OK") }

	// Add parameterized routes
	routes := []string{
		"/users/:id",
		"/users/:id/posts/:postId",
		"/users/:id/posts/:postId/comments/:commentId",
		"/api/:version/users/:id",
		"/files/:category/:filename",
		"/shop/:category/:subcategory/:item",
	}

	for _, route := range routes {
		router.GET(route, handler)
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		router.getRoute("GET", "/users/123/posts/456/comments/789")
	}
}

// BenchmarkRouter_WildcardRoutes tests wildcard route performance
func BenchmarkRouter_WildcardRoutes(b *testing.B) {
	router := NewRouter()
	handler := func(c *Context) { c.String(200, "OK") }

	router.GET("/files/*filepath", handler)
	router.GET("/assets/*path", handler)
	router.GET("/static/*filename", handler)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		router.getRoute("GET", "/files/images/avatars/user123.png")
	}
}

// BenchmarkRouter_MixedRoutes tests performance with mixed route types
func BenchmarkRouter_MixedRoutes(b *testing.B) {
	router := NewRouter()
	handler := func(c *Context) { c.String(200, "OK") }

	// Static routes
	router.GET("/", handler)
	router.GET("/api", handler)
	router.GET("/health", handler)

	// Parameterized routes
	router.GET("/users/:id", handler)
	router.GET("/posts/:id", handler)

	// Wildcard routes
	router.GET("/files/*path", handler)

	testPaths := []string{
		"/",
		"/api",
		"/health",
		"/users/123",
		"/posts/456",
		"/files/css/style.css",
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		path := testPaths[i%len(testPaths)]
		router.getRoute("GET", path)
	}
}

// BenchmarkContext_Param tests parameter extraction performance
func BenchmarkContext_Param(b *testing.B) {
	req := httptest.NewRequest("GET", "/users/123", nil)
	w := httptest.NewRecorder()
	c := NewContext(w, req)
	c.params = map[string]string{
		"id":     "123",
		"name":   "john",
		"email":  "john@example.com",
		"status": "active",
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.Param("id")
	}
}

// BenchmarkContext_Query tests query parameter parsing performance
func BenchmarkContext_Query(b *testing.B) {
	req := httptest.NewRequest("GET", "/search?q=golang&page=1&limit=10&sort=name&order=asc", nil)
	w := httptest.NewRecorder()
	c := NewContext(w, req)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.Query("q")
	}
}

// BenchmarkContext_JSON tests JSON encoding performance
func BenchmarkContext_JSON(b *testing.B) {
	req := httptest.NewRequest("GET", "/", nil)
	data := map[string]interface{}{
		"message": "Hello, World!",
		"status":  "success",
		"data": []map[string]interface{}{
			{"id": 1, "name": "John", "email": "john@example.com"},
			{"id": 2, "name": "Jane", "email": "jane@example.com"},
			{"id": 3, "name": "Bob", "email": "bob@example.com"},
		},
		"meta": map[string]interface{}{
			"total": 3,
			"page":  1,
			"limit": 10,
		},
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		c := NewContext(w, req)
		c.JSON(200, data)
	}
}

// BenchmarkContext_BindJSON tests JSON decoding performance
func BenchmarkContext_BindJSON(b *testing.B) {
	jsonData := `{
		"name": "John Doe",
		"email": "john@example.com",
		"age": 30,
		"address": {
			"street": "123 Main St",
			"city": "New York",
			"country": "USA"
		},
		"hobbies": ["reading", "coding", "traveling"]
	}`

	var user struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Age     int    `json:"age"`
		Address struct {
			Street  string `json:"street"`
			City    string `json:"city"`
			Country string `json:"country"`
		} `json:"address"`
		Hobbies []string `json:"hobbies"`
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/", strings.NewReader(jsonData))
		w := httptest.NewRecorder()
		c := NewContext(w, req)
		c.BindJSON(&user)
	}
}

// BenchmarkContext_SetGet tests context storage performance
func BenchmarkContext_SetGet(b *testing.B) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	c := NewContext(w, req)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c.Set("key", "value")
		c.Get("key")
	}
}

// BenchmarkContext_Pool tests context pool performance
func BenchmarkContext_Pool(b *testing.B) {
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		c := NewContext(w, req)
		c.reset()
		contextPool.Put(c)
	}
}

// BenchmarkParsePattern_Performance tests URL pattern parsing performance
func BenchmarkParsePattern_Performance(b *testing.B) {
	patterns := []string{
		"/",
		"/api/v1/users",
		"/users/:id",
		"/users/:id/posts/:postId",
		"/files/*filepath",
		"/api/:version/users/:id/posts/:postId/comments/*path",
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		pattern := patterns[i%len(patterns)]
		parsePattern(pattern)
	}
}

// BenchmarkConcurrency tests concurrent request handling
func BenchmarkConcurrency(b *testing.B) {
	app := New()
	app.Use(func(c *Context) {
		// Simulate some work
		time.Sleep(time.Microsecond)
		c.Next()
	})
	app.GET("/users/:id", func(c *Context) {
		id := c.Param("id")
		c.JSON(200, map[string]string{"user_id": id})
	})

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		req := httptest.NewRequest("GET", "/users/123", nil)
		for pb.Next() {
			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)
		}
	})
}

// BenchmarkMemoryUsage tests memory allocation patterns
func BenchmarkMemoryUsage(b *testing.B) {
	app := New()
	app.GET("/api/data", func(c *Context) {
		// Create a moderately complex response
		data := make([]map[string]interface{}, 100)
		for i := 0; i < 100; i++ {
			data[i] = map[string]interface{}{
				"id":     i,
				"name":   fmt.Sprintf("Item %d", i),
				"active": i%2 == 0,
			}
		}
		c.JSON(200, map[string]interface{}{
			"items": data,
			"total": 100,
		})
	})

	req := httptest.NewRequest("GET", "/api/data", nil)

	// Force GC before benchmark
	runtime.GC()

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// BenchmarkRouteGroups tests route group performance
func BenchmarkRouteGroups(b *testing.B) {
	app := New()

	// Create nested route groups
	api := app.Route("/api")
	v1 := api.Group("/v1")
	v2 := api.Group("/v2")

	v1.GET("/users/:id", func(c *Context) {
		c.JSON(200, map[string]string{"version": "v1", "id": c.Param("id")})
	})

	v2.GET("/users/:id", func(c *Context) {
		c.JSON(200, map[string]string{"version": "v2", "id": c.Param("id")})
	})

	requests := []string{
		"/api/v1/users/123",
		"/api/v2/users/456",
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		path := requests[i%len(requests)]
		req := httptest.NewRequest("GET", path, nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// BenchmarkErrorHandling tests error handling performance
func BenchmarkErrorHandling(b *testing.B) {
	app := New()
	app.UseError(func(err error, c *Context) {
		c.JSON(500, map[string]string{"error": err.Error()})
	})
	app.Use(Recover())
	app.GET("/error", func(c *Context) {
		c.Next(fmt.Errorf("test error"))
	})

	req := httptest.NewRequest("GET", "/error", nil)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// BenchmarkLargePayload tests performance with large JSON payloads
func BenchmarkLargePayload(b *testing.B) {
	app := New()
	app.POST("/upload", func(c *Context) {
		var data map[string]interface{}
		if err := c.BindJSON(&data); err != nil {
			c.JSON(400, map[string]string{"error": "invalid json"})
			return
		}
		c.JSON(200, map[string]string{"status": "received"})
	})

	// Create a large JSON payload (about 1MB)
	largeData := make(map[string]interface{})
	for i := 0; i < 1000; i++ {
		largeData[fmt.Sprintf("key_%d", i)] = strings.Repeat("x", 1000)
	}
	jsonData, _ := json.Marshal(largeData)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/upload", bytes.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}
