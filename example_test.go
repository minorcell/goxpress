package goxpress

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
)

// Example demonstrates the basic usage of goxpress
func Example() {
	app := New()

	// Simple GET route
	app.GET("/", func(c *Context) {
		c.String(200, "Hello, World!")
	})

	// JSON response
	app.GET("/api/hello", func(c *Context) {
		c.JSON(200, map[string]string{
			"message": "Hello from goxpress!",
		})
	})

	// Parameter route
	app.GET("/users/:id", func(c *Context) {
		id := c.Param("id")
		c.JSON(200, map[string]string{
			"user_id": id,
		})
	})

	fmt.Println("Basic goxpress server setup complete")
	// Output: Basic goxpress server setup complete
}

// Example_middleware demonstrates middleware usage
func Example_middleware() {
	app := New()

	// Add built-in middleware
	app.Use(Logger())
	app.Use(Recover())

	// Custom middleware
	app.Use(func(c *Context) {
		c.Set("start_time", "middleware_executed")
		c.Next()
	})

	app.GET("/", func(c *Context) {
		startTime, _ := c.GetString("start_time")
		c.String(200, "Middleware data: "+startTime)
	})

	fmt.Println("Middleware setup complete")
	// Output: Middleware setup complete
}

// Example_routeGroups demonstrates route grouping
func Example_routeGroups() {
	app := New()

	// API v1 group
	v1 := app.Route("/api/v1")
	v1.GET("/users", func(c *Context) {
		c.JSON(200, []string{"user1", "user2"})
	})
	v1.POST("/users", func(c *Context) {
		c.JSON(201, map[string]string{"status": "created"})
	})

	// API v2 group
	v2 := app.Route("/api/v2")
	v2.GET("/users", func(c *Context) {
		c.JSON(200, map[string]interface{}{
			"users":   []string{"user1", "user2"},
			"version": "v2",
		})
	})

	fmt.Println("Route groups setup complete")
	// Output: Route groups setup complete
}

// Example_errorHandling demonstrates error handling
func Example_errorHandling() {
	app := New()

	// Error handler
	app.UseError(func(err error, c *Context) {
		c.JSON(500, map[string]string{
			"error": err.Error(),
		})
	})

	// Route that triggers an error
	app.GET("/error", func(c *Context) {
		c.Next(fmt.Errorf("something went wrong"))
	})

	// Route that panics
	app.GET("/panic", func(c *Context) {
		panic("panic occurred")
	})

	fmt.Println("Error handling setup complete")
	// Output: Error handling setup complete
}

// Example_restAPI demonstrates a complete REST API
func Example_restAPI() {
	app := New()

	// User data store (in memory for example)
	users := map[string]map[string]interface{}{
		"1": {"id": "1", "name": "John", "email": "john@example.com"},
		"2": {"id": "2", "name": "Jane", "email": "jane@example.com"},
	}

	// Middleware to add CORS headers
	app.Use(func(c *Context) {
		c.Response.Header().Set("Access-Control-Allow-Origin", "*")
		c.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Response.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Next()
	})

	// Users API group
	api := app.Route("/api")

	// GET /api/users - List all users
	api.GET("/users", func(c *Context) {
		userList := make([]map[string]interface{}, 0)
		for _, user := range users {
			userList = append(userList, user)
		}
		c.JSON(200, map[string]interface{}{
			"users": userList,
			"count": len(userList),
		})
	})

	// GET /api/users/:id - Get user by ID
	api.GET("/users/:id", func(c *Context) {
		id := c.Param("id")
		if user, exists := users[id]; exists {
			c.JSON(200, user)
		} else {
			c.JSON(404, map[string]string{"error": "User not found"})
		}
	})

	// POST /api/users - Create new user
	api.POST("/users", func(c *Context) {
		var newUser map[string]interface{}
		if err := c.BindJSON(&newUser); err != nil {
			c.JSON(400, map[string]string{"error": "Invalid JSON"})
			return
		}

		// Simple validation
		if newUser["name"] == nil || newUser["email"] == nil {
			c.JSON(400, map[string]string{"error": "Name and email are required"})
			return
		}

		// Generate new ID
		newID := fmt.Sprintf("%d", len(users)+1)
		newUser["id"] = newID
		users[newID] = newUser

		c.JSON(201, newUser)
	})

	// PUT /api/users/:id - Update user
	api.PUT("/users/:id", func(c *Context) {
		id := c.Param("id")
		if _, exists := users[id]; !exists {
			c.JSON(404, map[string]string{"error": "User not found"})
			return
		}

		var updatedUser map[string]interface{}
		if err := c.BindJSON(&updatedUser); err != nil {
			c.JSON(400, map[string]string{"error": "Invalid JSON"})
			return
		}

		updatedUser["id"] = id
		users[id] = updatedUser
		c.JSON(200, updatedUser)
	})

	// DELETE /api/users/:id - Delete user
	api.DELETE("/users/:id", func(c *Context) {
		id := c.Param("id")
		if _, exists := users[id]; exists {
			delete(users, id)
			c.JSON(200, map[string]string{"message": "User deleted"})
		} else {
			c.JSON(404, map[string]string{"error": "User not found"})
		}
	})

	fmt.Println("REST API setup complete")
	// Output: REST API setup complete
}

// Test the examples to ensure they work
func TestExamples(t *testing.T) {
	t.Run("BasicUsage", func(t *testing.T) {
		app := New()
		app.GET("/", func(c *Context) {
			c.String(200, "Hello, World!")
		})

		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
		if w.Body.String() != "Hello, World!" {
			t.Errorf("Expected 'Hello, World!', got '%s'", w.Body.String())
		}
	})

	t.Run("JSONResponse", func(t *testing.T) {
		app := New()
		app.GET("/api/hello", func(c *Context) {
			c.JSON(200, map[string]string{"message": "Hello from goxpress!"})
		})

		req := httptest.NewRequest("GET", "/api/hello", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		if response["message"] != "Hello from goxpress!" {
			t.Errorf("Expected correct message, got '%s'", response["message"])
		}
	})

	t.Run("ParameterRoute", func(t *testing.T) {
		app := New()
		app.GET("/users/:id", func(c *Context) {
			id := c.Param("id")
			c.JSON(200, map[string]string{"user_id": id})
		})

		req := httptest.NewRequest("GET", "/users/123", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		if response["user_id"] != "123" {
			t.Errorf("Expected user_id '123', got '%s'", response["user_id"])
		}
	})

	t.Run("Middleware", func(t *testing.T) {
		app := New()
		app.Use(func(c *Context) {
			c.Set("middleware", "executed")
			c.Next()
		})
		app.GET("/", func(c *Context) {
			value, _ := c.GetString("middleware")
			c.String(200, value)
		})

		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Body.String() != "executed" {
			t.Errorf("Expected 'executed', got '%s'", w.Body.String())
		}
	})

	t.Run("RouteGroups", func(t *testing.T) {
		app := New()
		api := app.Route("/api")
		api.GET("/users", func(c *Context) {
			c.JSON(200, []string{"user1", "user2"})
		})

		req := httptest.NewRequest("GET", "/api/users", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		app := New()
		var handledError error

		app.UseError(func(err error, c *Context) {
			handledError = err
			c.JSON(500, map[string]string{"error": err.Error()})
		})

		app.GET("/error", func(c *Context) {
			c.Next(fmt.Errorf("test error"))
		})

		req := httptest.NewRequest("GET", "/error", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if handledError == nil {
			t.Error("Error should have been handled")
		}
		if handledError.Error() != "test error" {
			t.Errorf("Expected 'test error', got '%s'", handledError.Error())
		}
	})

	t.Run("POST with JSON", func(t *testing.T) {
		app := New()
		app.POST("/users", func(c *Context) {
			var user map[string]interface{}
			if err := c.BindJSON(&user); err != nil {
				c.JSON(400, map[string]string{"error": "Invalid JSON"})
				return
			}
			c.JSON(201, user)
		})

		jsonData := `{"name":"John","email":"john@example.com"}`
		req := httptest.NewRequest("POST", "/users", strings.NewReader(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != 201 {
			t.Errorf("Expected status 201, got %d", w.Code)
		}
	})

	t.Run("QueryParameters", func(t *testing.T) {
		app := New()
		app.GET("/search", func(c *Context) {
			query := c.Query("q")
			page := c.Query("page")
			c.JSON(200, map[string]string{
				"query": query,
				"page":  page,
			})
		})

		req := httptest.NewRequest("GET", "/search?q=golang&page=1", nil)
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)

		if w.Code != 200 {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)
		if response["query"] != "golang" {
			t.Errorf("Expected query 'golang', got '%s'", response["query"])
		}
		if response["page"] != "1" {
			t.Errorf("Expected page '1', got '%s'", response["page"])
		}
	})
}

// BenchmarkFullStack benchmarks a complete request flow
func BenchmarkFullStack(b *testing.B) {
	app := New()
	app.Use(func(c *Context) {
		c.Set("middleware", "executed")
		c.Next()
	})
	app.GET("/api/users/:id", func(c *Context) {
		id := c.Param("id")
		middleware, _ := c.GetString("middleware")
		c.JSON(200, map[string]interface{}{
			"user_id":    id,
			"middleware": middleware,
			"status":     "success",
		})
	})

	req := httptest.NewRequest("GET", "/api/users/123", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}
