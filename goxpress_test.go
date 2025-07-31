package goxpress

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	app := New()
	if app == nil {
		t.Fatal("New() should return a valid Engine instance")
	}
	if app.router == nil {
		t.Fatal("Engine should have a router instance")
	}
	if app.middlewares == nil {
		t.Fatal("Engine should have middlewares slice initialized")
	}
	if app.errorHandlers == nil {
		t.Fatal("Engine should have errorHandlers slice initialized")
	}
}

func TestEngineUse(t *testing.T) {
	app := New()

	middleware1 := func(c *Context) { c.Next() }
	middleware2 := func(c *Context) { c.Next() }

	// Test single middleware
	app.Use(middleware1)
	if len(app.middlewares) != 1 {
		t.Errorf("Expected 1 middleware, got %d", len(app.middlewares))
	}

	// Test multiple middlewares
	app.Use(middleware2)
	if len(app.middlewares) != 2 {
		t.Errorf("Expected 2 middlewares, got %d", len(app.middlewares))
	}

	// Test chaining
	result := app.Use(middleware1)
	if result != app {
		t.Error("Use() should return the same Engine instance for chaining")
	}
}

func TestEngineUseError(t *testing.T) {
	app := New()

	errorHandler1 := func(err error, c *Context) {}
	errorHandler2 := func(err error, c *Context) {}

	// Test single error handler
	app.UseError(errorHandler1)
	if len(app.errorHandlers) != 1 {
		t.Errorf("Expected 1 error handler, got %d", len(app.errorHandlers))
	}

	// Test multiple error handlers
	app.UseError(errorHandler2)
	if len(app.errorHandlers) != 2 {
		t.Errorf("Expected 2 error handlers, got %d", len(app.errorHandlers))
	}

	// Test chaining
	result := app.UseError(errorHandler1)
	if result != app {
		t.Error("UseError() should return the same Engine instance for chaining")
	}
}

func TestHTTPMethods(t *testing.T) {
	app := New()
	handler := func(c *Context) {
		c.String(200, "OK")
	}

	// Test all HTTP methods
	methods := []struct {
		method string
		fn     func(string, ...HandlerFunc) *Engine
	}{
		{"GET", app.GET},
		{"POST", app.POST},
		{"PUT", app.PUT},
		{"DELETE", app.DELETE},
		{"PATCH", app.PATCH},
		{"HEAD", app.HEAD},
		{"OPTIONS", app.OPTIONS},
	}

	for _, m := range methods {
		t.Run(m.method, func(t *testing.T) {
			result := m.fn("/test", handler)
			if result != app {
				t.Errorf("%s() should return the same Engine instance for chaining", m.method)
			}
		})
	}
}

func TestRoute(t *testing.T) {
	app := New()

	router := app.Route("/api")
	if router == nil {
		t.Fatal("Route() should return a valid Router instance")
	}
}

func TestServeHTTP(t *testing.T) {
	app := New()

	// Test basic route
	app.GET("/hello", func(c *Context) {
		c.String(200, "Hello World")
	})

	req := httptest.NewRequest("GET", "/hello", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", w.Body.String())
	}
}

func TestServeHTTP404(t *testing.T) {
	app := New()

	req := httptest.NewRequest("GET", "/notfound", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "404 page not found") {
		t.Errorf("Expected '404 page not found' in response, got '%s'", w.Body.String())
	}
}

func TestMiddlewareExecution(t *testing.T) {
	app := New()
	var executed []string

	// Add middleware
	app.Use(func(c *Context) {
		executed = append(executed, "middleware1")
		c.Next()
	})

	app.Use(func(c *Context) {
		executed = append(executed, "middleware2")
		c.Next()
	})

	// Add route handler
	app.GET("/test", func(c *Context) {
		executed = append(executed, "handler")
		c.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	expected := []string{"middleware1", "middleware2", "handler"}
	if len(executed) != len(expected) {
		t.Fatalf("Expected %d executions, got %d", len(expected), len(executed))
	}

	for i, exp := range expected {
		if executed[i] != exp {
			t.Errorf("Expected execution[%d] = '%s', got '%s'", i, exp, executed[i])
		}
	}
}

func TestMiddlewareAbort(t *testing.T) {
	app := New()
	var executed []string

	// Add middleware that aborts
	app.Use(func(c *Context) {
		executed = append(executed, "middleware1")
		c.Abort()
		c.JSON(401, map[string]string{"error": "unauthorized"})
	})

	app.Use(func(c *Context) {
		executed = append(executed, "middleware2")
		c.Next()
	})

	// Add route handler
	app.GET("/test", func(c *Context) {
		executed = append(executed, "handler")
		c.String(200, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	// Only middleware1 should execute
	if len(executed) != 1 || executed[0] != "middleware1" {
		t.Errorf("Expected only 'middleware1' to execute, got %v", executed)
	}

	if w.Code != 401 {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestErrorHandling(t *testing.T) {
	app := New()
	var handledError error

	// Add error handler
	app.UseError(func(err error, c *Context) {
		handledError = err
		c.JSON(500, map[string]string{"error": err.Error()})
	})

	// Add route that generates error
	app.GET("/error", func(c *Context) {
		c.Next(fmt.Errorf("test error"))
	})

	req := httptest.NewRequest("GET", "/error", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if handledError == nil {
		t.Fatal("Error handler should have been called")
	}

	if handledError.Error() != "test error" {
		t.Errorf("Expected error 'test error', got '%s'", handledError.Error())
	}
}

func TestJSONResponse(t *testing.T) {
	app := New()

	app.GET("/json", func(c *Context) {
		c.JSON(200, map[string]interface{}{
			"message": "success",
			"data":    []int{1, 2, 3},
		})
	})

	req := httptest.NewRequest("GET", "/json", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal JSON response: %v", err)
	}

	if response["message"] != "success" {
		t.Errorf("Expected message 'success', got '%v'", response["message"])
	}
}

func TestRouteParams(t *testing.T) {
	app := New()

	app.GET("/users/:id/posts/:postId", func(c *Context) {
		userId := c.Param("id")
		postId := c.Param("postId")
		c.JSON(200, map[string]string{
			"userId": userId,
			"postId": postId,
		})
	})

	req := httptest.NewRequest("GET", "/users/123/posts/456", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal JSON response: %v", err)
	}

	if response["userId"] != "123" {
		t.Errorf("Expected userId '123', got '%s'", response["userId"])
	}

	if response["postId"] != "456" {
		t.Errorf("Expected postId '456', got '%s'", response["postId"])
	}
}

func TestQueryParams(t *testing.T) {
	app := New()

	app.GET("/search", func(c *Context) {
		query := c.Query("q")
		page := c.Query("page")
		c.JSON(200, map[string]string{
			"query": query,
			"page":  page,
		})
	})

	req := httptest.NewRequest("GET", "/search?q=golang&page=2", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal JSON response: %v", err)
	}

	if response["query"] != "golang" {
		t.Errorf("Expected query 'golang', got '%s'", response["query"])
	}

	if response["page"] != "2" {
		t.Errorf("Expected page '2', got '%s'", response["page"])
	}
}

func TestPostJSON(t *testing.T) {
	app := New()

	app.POST("/users", func(c *Context) {
		var user struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}

		if err := c.BindJSON(&user); err != nil {
			c.JSON(400, map[string]string{"error": "Invalid JSON"})
			return
		}

		c.JSON(201, map[string]interface{}{
			"id":    1,
			"name":  user.Name,
			"email": user.Email,
		})
	})

	reqBody := map[string]string{
		"name":  "John Doe",
		"email": "john@example.com",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != 201 {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal JSON response: %v", err)
	}

	if response["name"] != "John Doe" {
		t.Errorf("Expected name 'John Doe', got '%v'", response["name"])
	}

	if response["email"] != "john@example.com" {
		t.Errorf("Expected email 'john@example.com', got '%v'", response["email"])
	}
}

func TestRouteGroups(t *testing.T) {
	app := New()

	api := app.Route("/api")
	api.GET("/users", func(c *Context) {
		c.JSON(200, map[string]string{"endpoint": "users"})
	})

	req := httptest.NewRequest("GET", "/api/users", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal JSON response: %v", err)
	}

	if response["endpoint"] != "users" {
		t.Errorf("Expected endpoint 'users', got '%s'", response["endpoint"])
	}
}

func BenchmarkEngine(b *testing.B) {
	app := New()
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

func BenchmarkEngineWithMiddleware(b *testing.B) {
	app := New()
	app.Use(func(c *Context) { c.Next() })
	app.Use(func(c *Context) { c.Next() })
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
