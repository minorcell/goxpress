package goxpress

import (
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	// Capture log output
	var logOutput strings.Builder
	log.SetOutput(&logOutput)
	defer log.SetOutput(os.Stderr) // Restore default output

	app := New()
	app.Use(Logger())
	app.GET("/test", func(c *Context) {
		time.Sleep(10 * time.Millisecond) // Simulate some processing time
		c.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	w := httptest.NewRecorder()

	start := time.Now()
	app.ServeHTTP(w, req)
	elapsed := time.Since(start)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "OK" {
		t.Errorf("Expected body 'OK', got '%s'", w.Body.String())
	}

	// Check log output
	logStr := logOutput.String()
	if !strings.Contains(logStr, "[GET]") {
		t.Error("Log should contain request method")
	}
	if !strings.Contains(logStr, "/test") {
		t.Error("Log should contain request path")
	}
	if !strings.Contains(logStr, "127.0.0.1:12345") {
		t.Error("Log should contain remote address")
	}

	// Check that duration is reasonable (should be at least 10ms due to sleep)
	if elapsed < 10*time.Millisecond {
		t.Errorf("Expected request to take at least 10ms, took %v", elapsed)
	}
}

func TestLoggerWithMultipleRequests(t *testing.T) {
	var logOutput strings.Builder
	log.SetOutput(&logOutput)
	defer log.SetOutput(os.Stderr)

	app := New()
	app.Use(Logger())
	app.GET("/users/:id", func(c *Context) {
		id := c.Param("id")
		c.String(200, "User "+id)
	})

	requests := []struct {
		method string
		path   string
		addr   string
	}{
		{"GET", "/users/123", "192.168.1.1:8080"},
		{"GET", "/users/456", "10.0.0.1:9090"},
		{"GET", "/users/789", "172.16.0.1:3000"},
	}

	for _, req := range requests {
		t.Run(req.path, func(t *testing.T) {
			httpReq := httptest.NewRequest(req.method, req.path, nil)
			httpReq.RemoteAddr = req.addr
			w := httptest.NewRecorder()

			app.ServeHTTP(w, httpReq)

			if w.Code != 200 {
				t.Errorf("Expected status 200, got %d", w.Code)
			}
		})
	}

	// Verify all requests were logged
	logStr := logOutput.String()
	for _, req := range requests {
		if !strings.Contains(logStr, req.path) {
			t.Errorf("Log should contain path %s", req.path)
		}
		if !strings.Contains(logStr, req.addr) {
			t.Errorf("Log should contain address %s", req.addr)
		}
	}
}

func TestRecover(t *testing.T) {
	var logOutput strings.Builder
	log.SetOutput(&logOutput)
	defer log.SetOutput(os.Stderr)

	app := New()
	var handledError error

	// Add error handler to capture the error
	app.UseError(func(err error, c *Context) {
		handledError = err
		c.JSON(500, map[string]string{"error": "Internal Server Error"})
	})

	app.Use(Recover())
	app.GET("/panic", func(c *Context) {
		panic("test panic")
	})

	req := httptest.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	// Check that panic was recovered and converted to error
	if handledError == nil {
		t.Fatal("Error handler should have been called")
	}

	if handledError.Error() != "test panic" {
		t.Errorf("Expected error 'test panic', got '%s'", handledError.Error())
	}

	// Check log output
	logStr := logOutput.String()
	if !strings.Contains(logStr, "Panic recovered") {
		t.Error("Log should contain 'Panic recovered'")
	}
	if !strings.Contains(logStr, "test panic") {
		t.Error("Log should contain the panic message")
	}
}

func TestRecoverWithErrorType(t *testing.T) {
	var logOutput strings.Builder
	log.SetOutput(&logOutput)
	defer log.SetOutput(os.Stderr)

	app := New()
	var handledError error

	app.UseError(func(err error, c *Context) {
		handledError = err
		c.JSON(500, map[string]string{"error": err.Error()})
	})

	app.Use(Recover())
	app.GET("/panic-error", func(c *Context) {
		panic(fmt.Errorf("custom error"))
	})

	req := httptest.NewRequest("GET", "/panic-error", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if handledError == nil {
		t.Fatal("Error handler should have been called")
	}

	if handledError.Error() != "custom error" {
		t.Errorf("Expected error 'custom error', got '%s'", handledError.Error())
	}
}

func TestRecoverWithNonErrorPanic(t *testing.T) {
	var logOutput strings.Builder
	log.SetOutput(&logOutput)
	defer log.SetOutput(os.Stderr)

	app := New()
	var handledError error

	app.UseError(func(err error, c *Context) {
		handledError = err
		c.JSON(500, map[string]string{"error": "Internal Server Error"})
	})

	app.Use(Recover())
	app.GET("/panic-int", func(c *Context) {
		panic(42) // Panic with non-error type
	})

	req := httptest.NewRequest("GET", "/panic-int", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if handledError == nil {
		t.Fatal("Error handler should have been called")
	}

	if handledError.Error() != "42" {
		t.Errorf("Expected error '42', got '%s'", handledError.Error())
	}
}

func TestRecoverDoesNotAffectNormalRequests(t *testing.T) {
	app := New()
	app.Use(Recover())
	app.GET("/normal", func(c *Context) {
		c.String(200, "Normal response")
	})

	req := httptest.NewRequest("GET", "/normal", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "Normal response" {
		t.Errorf("Expected 'Normal response', got '%s'", w.Body.String())
	}
}

func TestMiddlewareOrder(t *testing.T) {
	var logOutput strings.Builder
	log.SetOutput(&logOutput)
	defer log.SetOutput(os.Stderr)

	var executed []string

	app := New()

	// Add custom middleware that tracks execution
	app.Use(func(c *Context) {
		executed = append(executed, "custom1")
		c.Next()
	})

	app.Use(Logger())

	app.Use(func(c *Context) {
		executed = append(executed, "custom2")
		c.Next()
	})

	app.Use(Recover())

	app.GET("/test", func(c *Context) {
		executed = append(executed, "handler")
		c.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "127.0.0.1:8080"
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	expected := []string{"custom1", "custom2", "handler"}
	if len(executed) != len(expected) {
		t.Fatalf("Expected %d executions, got %d", len(expected), len(executed))
	}

	for i, exp := range expected {
		if executed[i] != exp {
			t.Errorf("Expected execution[%d] = '%s', got '%s'", i, exp, executed[i])
		}
	}

	// Logger should have logged the request
	logStr := logOutput.String()
	if !strings.Contains(logStr, "[GET]") {
		t.Error("Logger should have logged the request")
	}
}

func TestLoggerMiddlewareChaining(t *testing.T) {
	app := New()

	// Test that Logger returns a HandlerFunc that can be chained
	logger := Logger()
	if logger == nil {
		t.Fatal("Logger() should return a valid HandlerFunc")
	}

	// Test chaining multiple middlewares
	result := app.Use(Logger()).Use(Recover())
	if result != app {
		t.Error("Middleware chaining should return the same app instance")
	}

	if len(app.middlewares) != 2 {
		t.Errorf("Expected 2 middlewares, got %d", len(app.middlewares))
	}
}

func TestRecoverMiddlewareChaining(t *testing.T) {
	app := New()

	// Test that Recover returns a HandlerFunc that can be chained
	recover := Recover()
	if recover == nil {
		t.Fatal("Recover() should return a valid HandlerFunc")
	}

	// Test chaining
	result := app.Use(Recover()).Use(Logger())
	if result != app {
		t.Error("Middleware chaining should return the same app instance")
	}

	if len(app.middlewares) != 2 {
		t.Errorf("Expected 2 middlewares, got %d", len(app.middlewares))
	}
}

func TestBuiltinMiddlewareIntegration(t *testing.T) {
	var logOutput strings.Builder
	log.SetOutput(&logOutput)
	defer log.SetOutput(os.Stderr)

	app := New()
	var recoveredError error

	app.UseError(func(err error, c *Context) {
		recoveredError = err
		c.JSON(500, map[string]string{"error": "Server Error"})
	})

	// Use both built-in middlewares
	app.Use(Logger())
	app.Use(Recover())

	app.GET("/test-panic", func(c *Context) {
		panic("integration test panic")
	})

	app.GET("/test-normal", func(c *Context) {
		c.String(200, "Normal")
	})

	// Test panic route
	req1 := httptest.NewRequest("GET", "/test-panic", nil)
	req1.RemoteAddr = "127.0.0.1:8080"
	w1 := httptest.NewRecorder()

	app.ServeHTTP(w1, req1)

	if recoveredError == nil {
		t.Error("Panic should have been recovered")
	}

	// Test normal route
	req2 := httptest.NewRequest("GET", "/test-normal", nil)
	req2.RemoteAddr = "127.0.0.1:8080"
	w2 := httptest.NewRecorder()

	app.ServeHTTP(w2, req2)

	if w2.Code != 200 {
		t.Errorf("Expected status 200, got %d", w2.Code)
	}

	// Both requests should be logged
	logStr := logOutput.String()
	panicLogCount := strings.Count(logStr, "/test-panic")
	normalLogCount := strings.Count(logStr, "/test-normal")

	if panicLogCount < 1 {
		t.Error("Panic request should be logged")
	}
	if normalLogCount < 1 {
		t.Error("Normal request should be logged")
	}
}

func BenchmarkLogger(b *testing.B) {
	app := New()
	app.Use(Logger())
	app.GET("/bench", func(c *Context) {
		c.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/bench", nil)

	// Suppress log output during benchmark
	log.SetOutput(&strings.Builder{})
	defer log.SetOutput(os.Stderr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

func BenchmarkRecover(b *testing.B) {
	app := New()
	app.Use(Recover())
	app.GET("/bench", func(c *Context) {
		c.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/bench", nil)

	// Suppress log output during benchmark
	log.SetOutput(&strings.Builder{})
	// Restore log output
	defer log.SetOutput(os.Stderr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

func TestLoggerWithConfig_SkipPaths(t *testing.T) {
	var logOutput strings.Builder

	config := LoggerConfig{
		SkipPaths: []string{"/health", "/metrics", "/api/*/internal"},
		Output:    &logOutput,
	}

	app := New()
	app.Use(LoggerWithConfig(config))

	// Add routes
	app.GET("/health", func(c *Context) {
		c.String(200, "OK")
	})
	app.GET("/api/v1/internal", func(c *Context) {
		c.String(200, "Internal")
	})
	app.GET("/api/users", func(c *Context) {
		c.String(200, "Users")
	})

	// Test skipped paths
	tests := []struct {
		path           string
		shouldBeLogged bool
	}{
		{"/health", false},          // exact match skip
		{"/metrics", false},         // exact match skip
		{"/api/v1/internal", false}, // wildcard match skip
		{"/api/users", true},        // should be logged
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			logOutput.Reset()

			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)

			logStr := logOutput.String()
			if tt.shouldBeLogged && logStr == "" {
				t.Errorf("Path %s should be logged but wasn't", tt.path)
			}
			if !tt.shouldBeLogged && logStr != "" {
				t.Errorf("Path %s should be skipped but was logged: %s", tt.path, logStr)
			}
		})
	}
}

func TestLoggerWithConfig_CustomOutput(t *testing.T) {
	var buffer1, buffer2 strings.Builder

	config1 := LoggerConfig{Output: &buffer1}
	config2 := LoggerConfig{Output: &buffer2}

	app1 := New()
	app1.Use(LoggerWithConfig(config1))
	app1.GET("/test1", func(c *Context) {
		c.String(200, "OK1")
	})

	app2 := New()
	app2.Use(LoggerWithConfig(config2))
	app2.GET("/test2", func(c *Context) {
		c.String(200, "OK2")
	})

	// Test first app
	req1 := httptest.NewRequest("GET", "/test1", nil)
	w1 := httptest.NewRecorder()
	app1.ServeHTTP(w1, req1)

	// Test second app
	req2 := httptest.NewRequest("GET", "/test2", nil)
	w2 := httptest.NewRecorder()
	app2.ServeHTTP(w2, req2)

	// Check outputs are separate
	log1 := buffer1.String()
	log2 := buffer2.String()

	if !strings.Contains(log1, "/test1") {
		t.Error("Buffer1 should contain /test1 log")
	}
	if strings.Contains(log1, "/test2") {
		t.Error("Buffer1 should not contain /test2 log")
	}
	if !strings.Contains(log2, "/test2") {
		t.Error("Buffer2 should contain /test2 log")
	}
	if strings.Contains(log2, "/test1") {
		t.Error("Buffer2 should not contain /test1 log")
	}
}

func TestLoggerWithConfig_WildcardMatching(t *testing.T) {
	var logOutput strings.Builder

	config := LoggerConfig{
		SkipPaths: []string{"/api/*/health", "/admin/*/debug/*"},
		Output:    &logOutput,
	}

	app := New()
	app.Use(LoggerWithConfig(config))

	tests := []struct {
		path           string
		shouldBeLogged bool
	}{
		{"/api/v1/health", false},          // matches /api/*/health
		{"/api/v2/health", false},          // matches /api/*/health
		{"/api/users", true},               // doesn't match
		{"/admin/panel/debug/info", false}, // matches /admin/*/debug/*
		{"/admin/panel/users", true},       // doesn't match
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			logOutput.Reset()

			app.GET(tt.path, func(c *Context) {
				c.String(200, "OK")
			})

			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()
			app.ServeHTTP(w, req)

			logStr := logOutput.String()
			if tt.shouldBeLogged && logStr == "" {
				t.Errorf("Path %s should be logged but wasn't", tt.path)
			}
			if !tt.shouldBeLogged && logStr != "" {
				t.Errorf("Path %s should be skipped but was logged: %s", tt.path, logStr)
			}
		})
	}
}

func TestLoggerWithConfig_DefaultValues(t *testing.T) {
	// Test that default values are applied when config is empty
	config := LoggerConfig{} // empty config

	app := New()
	app.Use(LoggerWithConfig(config))
	app.GET("/test", func(c *Context) {
		c.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Should not panic and should work with defaults
	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func BenchmarkLoggerWithConfig(b *testing.B) {
	var buffer strings.Builder
	config := LoggerConfig{
		Output:    &buffer,
		SkipPaths: []string{"/health"},
	}

	app := New()
	app.Use(LoggerWithConfig(config))
	app.GET("/bench", func(c *Context) {
		c.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/bench", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

func BenchmarkLoggerAndRecover(b *testing.B) {
	app := New()
	app.Use(Logger())
	app.Use(Recover())
	app.GET("/bench", func(c *Context) {
		c.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/bench", nil)

	// Suppress log output during benchmark
	log.SetOutput(&strings.Builder{})
	defer log.SetOutput(os.Stderr)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}
