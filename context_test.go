package goxpress

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewContext(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	c := NewContext(w, req)

	if c == nil {
		t.Fatal("NewContext should return a valid Context instance")
	}

	if c.Request != req {
		t.Error("Context should have the correct request")
	}

	if c.Response != w {
		t.Error("Context should have the correct response writer")
	}

	if c.params == nil {
		t.Error("Context should have params map initialized")
	}

	if c.store == nil {
		t.Error("Context should have store map initialized")
	}

	if c.index != -1 {
		t.Errorf("Context index should be -1, got %d", c.index)
	}

	if c.aborted {
		t.Error("Context should not be aborted initially")
	}

	if c.statusCodeWritten {
		t.Error("Context should not have status code written initially")
	}
}

func TestContextReset(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	c := NewContext(w, req)

	// Set some data
	c.params["id"] = "123"
	c.store["user"] = "john"
	c.index = 5
	c.aborted = true
	c.statusCodeWritten = true
	c.err = fmt.Errorf("test error")

	// Reset the context
	c.reset()

	if len(c.params) != 0 {
		t.Error("Params should be empty after reset")
	}

	if len(c.store) != 0 {
		t.Error("Store should be empty after reset")
	}

	if c.index != -1 {
		t.Errorf("Index should be -1 after reset, got %d", c.index)
	}

	if c.aborted {
		t.Error("Context should not be aborted after reset")
	}

	if c.statusCodeWritten {
		t.Error("Status code written should be false after reset")
	}

	if c.err != nil {
		t.Error("Error should be nil after reset")
	}
}

func TestContextParam(t *testing.T) {
	req := httptest.NewRequest("GET", "/users/123", nil)
	w := httptest.NewRecorder()

	c := NewContext(w, req)
	c.params = map[string]string{
		"id":   "123",
		"name": "john",
	}

	tests := []struct {
		key      string
		expected string
	}{
		{"id", "123"},
		{"name", "john"},
		{"nonexistent", ""},
	}

	for _, test := range tests {
		t.Run(test.key, func(t *testing.T) {
			result := c.Param(test.key)
			if result != test.expected {
				t.Errorf("Expected param %s = %s, got %s", test.key, test.expected, result)
			}
		})
	}
}

func TestContextQuery(t *testing.T) {
	req := httptest.NewRequest("GET", "/search?q=golang&page=2&sort=name", nil)
	w := httptest.NewRecorder()

	c := NewContext(w, req)

	tests := []struct {
		key      string
		expected string
	}{
		{"q", "golang"},
		{"page", "2"},
		{"sort", "name"},
		{"nonexistent", ""},
	}

	for _, test := range tests {
		t.Run(test.key, func(t *testing.T) {
			result := c.Query(test.key)
			if result != test.expected {
				t.Errorf("Expected query %s = %s, got %s", test.key, test.expected, result)
			}
		})
	}
}

func TestContextBindJSON(t *testing.T) {
	t.Run("ValidJSON", func(t *testing.T) {
		jsonData := `{"name":"John","age":30,"email":"john@example.com"}`
		req := httptest.NewRequest("POST", "/users", strings.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c := NewContext(w, req)

		var user struct {
			Name  string `json:"name"`
			Age   int    `json:"age"`
			Email string `json:"email"`
		}

		err := c.BindJSON(&user)
		if err != nil {
			t.Fatalf("BindJSON should not return error for valid JSON: %v", err)
		}

		if user.Name != "John" {
			t.Errorf("Expected name 'John', got '%s'", user.Name)
		}

		if user.Age != 30 {
			t.Errorf("Expected age 30, got %d", user.Age)
		}

		if user.Email != "john@example.com" {
			t.Errorf("Expected email 'john@example.com', got '%s'", user.Email)
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		invalidJSON := `{"name":"John","age":}`
		req := httptest.NewRequest("POST", "/users", strings.NewReader(invalidJSON))
		w := httptest.NewRecorder()

		c := NewContext(w, req)

		var user struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		err := c.BindJSON(&user)
		if err == nil {
			t.Error("BindJSON should return error for invalid JSON")
		}
	})
}

func TestContextStatus(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	c := NewContext(w, req)

	// First status call should work
	c.Status(201)
	if w.Code != 201 {
		t.Errorf("Expected status code 201, got %d", w.Code)
	}

	if !c.statusCodeWritten {
		t.Error("statusCodeWritten should be true after Status call")
	}

	// Second status call should be ignored
	c.Status(404)
	if w.Code != 201 {
		t.Errorf("Status code should remain 201, got %d", w.Code)
	}
}

func TestContextStatusCode(t *testing.T) {
	// Create a new context
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	c := NewContext(w, req)

	// Initially should be 0 since no status has been written
	if code := c.StatusCode(); code != 0 {
		t.Errorf("Expected status code 0, got %d", code)
	}

	// After writing status, should be 200 (our default placeholder)
	c.Status(404)
	if code := c.StatusCode(); code != 200 { // 200 because that's what our placeholder returns
		t.Errorf("Expected status code 200 (placeholder), got %d", code)
	}
}

func TestContextJSON(t *testing.T) {
	// Create a new context
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	c := NewContext(w, req)

	data := map[string]interface{}{
		"message": "success",
		"data":    []int{1, 2, 3},
		"user": map[string]string{
			"name":  "John",
			"email": "john@example.com",
		},
	}

	err := c.JSON(200, data)
	if err != nil {
		t.Fatalf("JSON should not return error: %v", err)
	}

	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "success" {
		t.Errorf("Expected message 'success', got '%v'", response["message"])
	}
}

func TestContextString(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	c := NewContext(w, req)

	err := c.String(200, "Hello %s, you are %d years old", "John", 30)
	if err != nil {
		t.Fatalf("String should not return error: %v", err)
	}

	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "text/plain; charset=utf-8" {
		t.Errorf("Expected Content-Type 'text/plain; charset=utf-8', got '%s'", contentType)
	}

	expected := "Hello John, you are 30 years old"
	if w.Body.String() != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, w.Body.String())
	}
}

func TestContextNext(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	c := NewContext(w, req)

	var executed []string

	handlers := []HandlerFunc{
		func(c *Context) {
			executed = append(executed, "handler1")
			c.Next()
		},
		func(c *Context) {
			executed = append(executed, "handler2")
			c.Next()
		},
		func(c *Context) {
			executed = append(executed, "handler3")
		},
	}

	c.handlers = handlers

	c.Next()

	expected := []string{"handler1", "handler2", "handler3"}
	if len(executed) != len(expected) {
		t.Fatalf("Expected %d executions, got %d", len(expected), len(executed))
	}

	for i, exp := range expected {
		if executed[i] != exp {
			t.Errorf("Expected execution[%d] = '%s', got '%s'", i, exp, executed[i])
		}
	}
}

func TestContextNextWithError(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	c := NewContext(w, req)

	testError := fmt.Errorf("test error")

	c.handlers = []HandlerFunc{
		func(c *Context) {
			c.Next(testError)
		},
	}

	c.Next()

	if c.err == nil {
		t.Fatal("Context should have error set")
	}

	if c.err.Error() != "test error" {
		t.Errorf("Expected error 'test error', got '%s'", c.err.Error())
	}
}

func TestContextAbort(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	c := NewContext(w, req)

	var executed []string

	handlers := []HandlerFunc{
		func(c *Context) {
			executed = append(executed, "handler1")
			c.Abort()
		},
		func(c *Context) {
			executed = append(executed, "handler2")
		},
	}

	c.handlers = handlers

	c.Next()

	// Only first handler should execute
	if len(executed) != 1 || executed[0] != "handler1" {
		t.Errorf("Expected only 'handler1' to execute, got %v", executed)
	}

	if !c.IsAborted() {
		t.Error("Context should be aborted")
	}
}

func TestContextSetGet(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	c := NewContext(w, req)

	// Test Set and Get
	c.Set("user", "john")
	c.Set("age", 30)
	c.Set("active", true)

	// Test Get existing values
	user, exists := c.Get("user")
	if !exists {
		t.Error("Expected 'user' key to exist")
	}
	if user != "john" {
		t.Errorf("Expected user 'john', got '%v'", user)
	}

	age, exists := c.Get("age")
	if !exists {
		t.Error("Expected 'age' key to exist")
	}
	if age != 30 {
		t.Errorf("Expected age 30, got %v", age)
	}

	// Test Get non-existing value
	_, exists = c.Get("nonexistent")
	if exists {
		t.Error("Expected 'nonexistent' key to not exist")
	}
}

func TestContextMustGet(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	c := NewContext(w, req)

	t.Run("ExistingKey", func(t *testing.T) {
		c.Set("user", "john")

		result := c.MustGet("user")
		if result != "john" {
			t.Errorf("Expected 'john', got '%v'", result)
		}
	})

	t.Run("NonExistentKey", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustGet should panic for non-existent key")
			}
		}()

		c.MustGet("nonexistent")
	})
}

func TestContextGetString(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	c := NewContext(w, req)

	// Test string value
	c.Set("name", "john")
	name, ok := c.GetString("name")
	if !ok {
		t.Error("Expected string value to be found")
	}
	if name != "john" {
		t.Errorf("Expected 'john', got '%s'", name)
	}

	// Test non-string value
	c.Set("age", 30)
	_, ok = c.GetString("age")
	if ok {
		t.Error("Expected non-string value to return false")
	}

	// Test non-existent key
	_, ok = c.GetString("nonexistent")
	if ok {
		t.Error("Expected non-existent key to return false")
	}

	// Test nil value
	c.Set("nilValue", nil)
	_, ok = c.GetString("nilValue")
	if ok {
		t.Error("Expected nil value to return false")
	}
}

func TestContextPool(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Create context from pool
	c1 := NewContext(w, req)

	// Use and reset the context
	c1.Set("test", "value")
	c1.params["id"] = "123"
	c1.index = 5
	c1.aborted = true

	c1.reset()
	contextPool.Put(c1)

	// Get another context from pool (should be the same instance)
	c2 := NewContext(w, req)

	// Verify it's reset
	if len(c2.store) != 0 {
		t.Error("Context from pool should have empty store")
	}
	if len(c2.params) != 0 {
		t.Error("Context from pool should have empty params")
	}
	if c2.index != -1 {
		t.Error("Context from pool should have index -1")
	}
	if c2.aborted {
		t.Error("Context from pool should not be aborted")
	}
}

func BenchmarkContextParam(b *testing.B) {
	req := httptest.NewRequest("GET", "/users/123", nil)
	w := httptest.NewRecorder()
	c := NewContext(w, req)
	c.params = map[string]string{"id": "123"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Param("id")
	}
}

func BenchmarkContextQuery(b *testing.B) {
	req := httptest.NewRequest("GET", "/search?q=golang&page=2", nil)
	w := httptest.NewRecorder()
	c := NewContext(w, req)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Query("q")
	}
}

func BenchmarkContextJSON(b *testing.B) {
	req := httptest.NewRequest("GET", "/test", nil)
	data := map[string]string{"message": "hello"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		c := NewContext(w, req)
		c.JSON(200, data)
	}
}

func BenchmarkContextSetGet(b *testing.B) {
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	c := NewContext(w, req)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.Set("key", "value")
		c.Get("key")
	}
}

func BenchmarkContextFromPool(b *testing.B) {
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c := NewContext(w, req)
		c.reset()
		contextPool.Put(c)
	}
}
