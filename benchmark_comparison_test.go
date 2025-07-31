package goxpress

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// standardHTTPServer creates a standard library HTTP server for comparison
func standardHTTPServer() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Hello, World!"))
	})

	mux.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]string{"message": "Hello, World!"})
	})

	mux.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/user/")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]string{"user_id": id})
	})

	return mux
}

// BenchmarkStandardLibrary_Simple tests standard library performance
func BenchmarkStandardLibrary_Simple(b *testing.B) {
	handler := standardHTTPServer()
	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}
}

// BenchmarkStandardLibrary_JSON tests standard library JSON performance
func BenchmarkStandardLibrary_JSON(b *testing.B) {
	handler := standardHTTPServer()
	req := httptest.NewRequest("GET", "/json", nil)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}
}

// BenchmarkStandardLibrary_Params tests standard library parameter handling
func BenchmarkStandardLibrary_Params(b *testing.B) {
	handler := standardHTTPServer()
	req := httptest.NewRequest("GET", "/user/123", nil)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}
}

// BenchmarkGoxpress_Simple tests goxpress basic performance
func BenchmarkGoxpress_Simple(b *testing.B) {
	app := New()
	app.GET("/", func(c *Context) {
		c.String(200, "Hello, World!")
	})

	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// BenchmarkGoxpress_JSON tests goxpress JSON performance
func BenchmarkGoxpress_JSON(b *testing.B) {
	app := New()
	app.GET("/json", func(c *Context) {
		c.JSON(200, map[string]string{"message": "Hello, World!"})
	})

	req := httptest.NewRequest("GET", "/json", nil)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// BenchmarkGoxpress_Params tests goxpress parameter handling
func BenchmarkGoxpress_Params(b *testing.B) {
	app := New()
	app.GET("/user/:id", func(c *Context) {
		id := c.Param("id")
		c.JSON(200, map[string]string{"user_id": id})
	})

	req := httptest.NewRequest("GET", "/user/123", nil)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// BenchmarkComparison_MiddlewareChain compares middleware performance
func BenchmarkGoxpress_MiddlewareChain(b *testing.B) {
	app := New()

	// Add 5 middlewares
	for i := 0; i < 5; i++ {
		app.Use(func(c *Context) {
			c.Set(fmt.Sprintf("middleware_%d", i), "executed")
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

// BenchmarkStandardLibrary_MiddlewareChain simulates middleware with standard library
func BenchmarkStandardLibrary_MiddlewareChain(b *testing.B) {
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	// Simulate 5 middlewares
	var handler http.Handler = finalHandler
	for i := 0; i < 5; i++ {
		middleware := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Simulate middleware work
				next.ServeHTTP(w, r)
			})
		}
		handler = middleware(handler)
	}

	req := httptest.NewRequest("GET", "/", nil)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}
}

// BenchmarkRouting_LargeRouteSet tests performance with many routes
func BenchmarkGoxpress_LargeRouteSet(b *testing.B) {
	app := New()

	// Add 100 routes
	for i := 0; i < 100; i++ {
		path := fmt.Sprintf("/api/v1/resource%d", i)
		app.GET(path, func(c *Context) {
			c.String(200, "OK")
		})
	}

	// Test against the 50th route
	req := httptest.NewRequest("GET", "/api/v1/resource50", nil)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// BenchmarkStandardLibrary_LargeRouteSet tests standard library with many routes
func BenchmarkStandardLibrary_LargeRouteSet(b *testing.B) {
	mux := http.NewServeMux()

	// Add 100 routes
	for i := 0; i < 100; i++ {
		path := fmt.Sprintf("/api/v1/resource%d", i)
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write([]byte("OK"))
		})
	}

	// Test against the 50th route
	req := httptest.NewRequest("GET", "/api/v1/resource50", nil)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
	}
}

// BenchmarkComplexRouting tests complex route patterns
func BenchmarkGoxpress_ComplexRouting(b *testing.B) {
	app := New()

	// Mix of static, parameter, and wildcard routes
	app.GET("/", func(c *Context) { c.String(200, "home") })
	app.GET("/api", func(c *Context) { c.String(200, "api") })
	app.GET("/api/v1", func(c *Context) { c.String(200, "v1") })
	app.GET("/users/:id", func(c *Context) { c.String(200, "user") })
	app.GET("/users/:id/posts/:postId", func(c *Context) { c.String(200, "post") })
	app.GET("/files/*path", func(c *Context) { c.String(200, "file") })
	app.GET("/admin/users", func(c *Context) { c.String(200, "admin") })
	app.GET("/public/css/style.css", func(c *Context) { c.String(200, "css") })

	routes := []string{
		"/",
		"/api",
		"/api/v1",
		"/users/123",
		"/users/123/posts/456",
		"/files/images/avatar.png",
		"/admin/users",
		"/public/css/style.css",
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		route := routes[i%len(routes)]
		req := httptest.NewRequest("GET", route, nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// BenchmarkMemoryEfficiency tests memory usage patterns
func BenchmarkGoxpress_MemoryEfficiency(b *testing.B) {
	app := New()
	app.GET("/api/data/:id", func(c *Context) {
		id := c.Param("id")
		data := map[string]interface{}{
			"id":     id,
			"name":   "Item " + id,
			"status": "active",
			"meta": map[string]interface{}{
				"created": "2023-01-01",
				"updated": "2023-01-02",
			},
		}
		c.JSON(200, data)
	})

	req := httptest.NewRequest("GET", "/api/data/123", nil)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// BenchmarkStandardLibrary_MemoryEfficiency tests standard library memory usage
func BenchmarkStandardLibrary_MemoryEfficiency(b *testing.B) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/data/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/api/data/")
		data := map[string]interface{}{
			"id":     id,
			"name":   "Item " + id,
			"status": "active",
			"meta": map[string]interface{}{
				"created": "2023-01-01",
				"updated": "2023-01-02",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(data)
	})

	req := httptest.NewRequest("GET", "/api/data/123", nil)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		mux.ServeHTTP(w, req)
	}
}

// BenchmarkConcurrentRequests tests concurrent request handling
func BenchmarkGoxpress_ConcurrentRequests(b *testing.B) {
	app := New()
	app.GET("/api/endpoint", func(c *Context) {
		c.JSON(200, map[string]string{"status": "success"})
	})

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		req := httptest.NewRequest("GET", "/api/endpoint", nil)
		for pb.Next() {
			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)
		}
	})
}

// BenchmarkStandardLibrary_ConcurrentRequests tests standard library concurrent handling
func BenchmarkStandardLibrary_ConcurrentRequests(b *testing.B) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/endpoint", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	})

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		req := httptest.NewRequest("GET", "/api/endpoint", nil)
		for pb.Next() {
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
		}
	})
}

// BenchmarkContextOperations tests context-specific operations
func BenchmarkGoxpress_ContextOperations(b *testing.B) {
	app := New()
	app.GET("/api/user/:id", func(c *Context) {
		// Simulate typical context operations
		id := c.Param("id")
		page := c.Query("page")
		c.Set("processed", true)
		processed, _ := c.Get("processed")

		c.JSON(200, map[string]interface{}{
			"user_id":   id,
			"page":      page,
			"processed": processed,
		})
	})

	req := httptest.NewRequest("GET", "/api/user/123?page=1", nil)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// BenchmarkErrorHandling tests error handling performance
func BenchmarkGoxpress_ErrorHandling(b *testing.B) {
	app := New()
	app.UseError(func(err error, c *Context) {
		c.JSON(500, map[string]string{"error": err.Error()})
	})

	app.GET("/error", func(c *Context) {
		c.Next(fmt.Errorf("simulated error"))
	})

	req := httptest.NewRequest("GET", "/error", nil)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

// BenchmarkRouteGrouping tests route group performance
func BenchmarkGoxpress_RouteGrouping(b *testing.B) {
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

	routes := []string{
		"/api/v1/users/123",
		"/api/v2/users/456",
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		route := routes[i%len(routes)]
		req := httptest.NewRequest("GET", route, nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}
