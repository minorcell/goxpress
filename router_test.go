package goxpress

import (
	"net/http/httptest"
	"testing"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter()
	if router == nil {
		t.Fatal("NewRouter() should return a valid Router instance")
	}
	if router.subRouters == nil {
		t.Fatal("Router should have subRouters map initialized")
	}
	if router.routes == nil {
		t.Fatal("Router should have routes map initialized")
	}
}

func TestRouterHTTPMethods(t *testing.T) {
	router := NewRouter()
	handler := func(c *Context) { c.String(200, "OK") }

	methods := []struct {
		name   string
		method func(string, ...HandlerFunc) *Router
		verb   string
	}{
		{"GET", router.GET, "GET"},
		{"POST", router.POST, "POST"},
		{"PUT", router.PUT, "PUT"},
		{"DELETE", router.DELETE, "DELETE"},
		{"PATCH", router.PATCH, "PATCH"},
		{"HEAD", router.HEAD, "HEAD"},
		{"OPTIONS", router.OPTIONS, "OPTIONS"},
	}

	for _, m := range methods {
		t.Run(m.name, func(t *testing.T) {
			result := m.method("/test", handler)
			if result != router {
				t.Errorf("%s() should return the same Router instance for chaining", m.name)
			}

			// Verify the route was registered
			_, exists := router.routes[m.verb]
			if !exists {
				t.Errorf("Route tree for %s method should exist", m.verb)
				return
			}

			node, _ := router.getRoute(m.verb, "/test")
			if node == nil {
				t.Errorf("Route /test should be registered for %s method", m.verb)
			}
		})
	}
}

func TestRouterStaticRoutes(t *testing.T) {
	router := NewRouter()

	router.GET("/", func(c *Context) { c.String(200, "root") })
	router.GET("/hello", func(c *Context) { c.String(200, "hello") })
	router.GET("/hello/world", func(c *Context) { c.String(200, "hello world") })

	tests := []struct {
		path     string
		expected bool
	}{
		{"/", true},
		{"/hello", true},
		{"/hello/world", true},
		{"/notfound", false},
		{"/hello/notfound", false},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			node, _ := router.getRoute("GET", test.path)
			if test.expected && node == nil {
				t.Errorf("Expected route %s to be found", test.path)
			}
			if !test.expected && node != nil {
				t.Errorf("Expected route %s to not be found", test.path)
			}
		})
	}
}

func TestRouterParameterRoutes(t *testing.T) {
	router := NewRouter()

	router.GET("/users/:id", func(c *Context) { c.String(200, "user") })
	router.GET("/users/:id/posts/:postId", func(c *Context) { c.String(200, "post") })
	router.GET("/files/*filepath", func(c *Context) { c.String(200, "file") })

	tests := []struct {
		path           string
		expectedFound  bool
		expectedParams map[string]string
	}{
		{
			path:           "/users/123",
			expectedFound:  true,
			expectedParams: map[string]string{"id": "123"},
		},
		{
			path:           "/users/abc",
			expectedFound:  true,
			expectedParams: map[string]string{"id": "abc"},
		},
		{
			path:           "/users/123/posts/456",
			expectedFound:  true,
			expectedParams: map[string]string{"id": "123", "postId": "456"},
		},
		{
			path:           "/files/images/avatar.png",
			expectedFound:  true,
			expectedParams: map[string]string{"filepath": "images/avatar.png"},
		},
		{
			path:           "/users",
			expectedFound:  false,
			expectedParams: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			node, params := router.getRoute("GET", test.path)

			if test.expectedFound && node == nil {
				t.Errorf("Expected route %s to be found", test.path)
				return
			}
			if !test.expectedFound && node != nil {
				t.Errorf("Expected route %s to not be found", test.path)
				return
			}

			if test.expectedFound {
				if len(params) != len(test.expectedParams) {
					t.Errorf("Expected %d parameters, got %d", len(test.expectedParams), len(params))
					return
				}

				for key, expectedValue := range test.expectedParams {
					if actualValue, exists := params[key]; !exists {
						t.Errorf("Expected parameter %s to exist", key)
					} else if actualValue != expectedValue {
						t.Errorf("Expected parameter %s = %s, got %s", key, expectedValue, actualValue)
					}
				}
			}
		})
	}
}

func TestRouterGroups(t *testing.T) {
	router := NewRouter()

	// Create a group
	api := router.Group("/api")
	if api == nil {
		t.Fatal("Group() should return a valid Router instance")
	}

	// Add routes to the group
	api.GET("/users", func(c *Context) { c.String(200, "users") })
	api.POST("/users", func(c *Context) { c.String(201, "create user") })

	// Test the grouped routes
	tests := []struct {
		method string
		path   string
		found  bool
	}{
		{"GET", "/api/users", true},
		{"POST", "/api/users", true},
		{"GET", "/users", false},
		{"POST", "/users", false},
	}

	for _, test := range tests {
		t.Run(test.method+" "+test.path, func(t *testing.T) {
			node, _ := router.getRoute(test.method, test.path)
			if test.found && node == nil {
				t.Errorf("Expected route %s %s to be found", test.method, test.path)
			}
			if !test.found && node != nil {
				t.Errorf("Expected route %s %s to not be found", test.method, test.path)
			}
		})
	}
}

func TestRouterNestedGroups(t *testing.T) {
	router := NewRouter()

	// Create nested groups
	api := router.Group("/api")
	v1 := api.Group("/v1")
	v1.GET("/users", func(c *Context) { c.String(200, "v1 users") })

	// Test nested group route
	node, _ := router.getRoute("GET", "/api/v1/users")
	if node == nil {
		t.Error("Expected nested group route /api/v1/users to be found")
	}
}

func TestRouterMultipleHandlers(t *testing.T) {
	router := NewRouter()

	handler1 := func(c *Context) { c.Next() }
	handler2 := func(c *Context) { c.String(200, "OK") }

	router.GET("/test", handler1, handler2)

	node, _ := router.getRoute("GET", "/test")
	if node == nil {
		t.Fatal("Expected route /test to be found")
	}

	if len(node.handlers) != 2 {
		t.Errorf("Expected 2 handlers, got %d", len(node.handlers))
	}
}

func TestRouterConflictingRoutes(t *testing.T) {
	router := NewRouter()

	// Add static route first
	router.GET("/users/new", func(c *Context) { c.String(200, "new user form") })

	// Add parameter route
	router.GET("/users/:id", func(c *Context) { c.String(200, "user detail") })

	// Static route should take precedence
	node, params := router.getRoute("GET", "/users/new")
	if node == nil {
		t.Error("Expected static route /users/new to be found")
	}
	if len(params) != 0 {
		t.Error("Static route should not have parameters")
	}

	// Parameter route should still work for other paths
	node, params = router.getRoute("GET", "/users/123")
	if node == nil {
		t.Error("Expected parameter route /users/:id to be found")
	}
	if params["id"] != "123" {
		t.Errorf("Expected parameter id = 123, got %s", params["id"])
	}
}

func TestParsePattern(t *testing.T) {
	tests := []struct {
		pattern  string
		expected []string
	}{
		{"/", []string{}},
		{"/hello", []string{"hello"}},
		{"/users/:id", []string{"users", ":id"}},
		{"/api/v1/users/:id/posts/:postId", []string{"api", "v1", "users", ":id", "posts", ":postId"}},
		{"/files/*filepath", []string{"files", "*filepath"}},
		{"//double//slashes//", []string{"double", "slashes"}},
		{"/trailing/slash/", []string{"trailing", "slash"}},
	}

	for _, test := range tests {
		t.Run(test.pattern, func(t *testing.T) {
			result := parsePattern(test.pattern)
			if len(result) != len(test.expected) {
				t.Errorf("Expected %d parts, got %d", len(test.expected), len(result))
				return
			}

			for i, expected := range test.expected {
				if result[i] != expected {
					t.Errorf("Expected part[%d] = %s, got %s", i, expected, result[i])
				}
			}
		})
	}
}

func TestRouterTreeBasicOperations(t *testing.T) {
	router := NewRouter()
	handler := func(c *Context) { c.String(200, "OK") }

	// Test adding routes
	router.GET("/users", handler)
	router.GET("/users/:id", handler)
	router.GET("/posts/*filepath", handler)

	// Test searching routes
	tests := []struct {
		path           string
		expectedFound  bool
		expectedParams map[string]string
	}{
		{"/users", true, map[string]string{}},
		{"/users/123", true, map[string]string{"id": "123"}},
		{"/posts/2023/article.html", true, map[string]string{"filepath": "2023/article.html"}},
		{"/notfound", false, nil},
	}

	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			node, params := router.getRoute("GET", test.path)

			if test.expectedFound && node == nil {
				t.Errorf("Expected route %s to be found", test.path)
				return
			}
			if !test.expectedFound && node != nil {
				t.Errorf("Expected route %s to not be found", test.path)
				return
			}

			if test.expectedFound && len(params) != len(test.expectedParams) {
				t.Errorf("Expected %d parameters, got %d", len(test.expectedParams), len(params))
			}
		})
	}
}

func TestRouterConcurrency(t *testing.T) {
	router := NewRouter()

	// Add some routes
	router.GET("/users/:id", func(c *Context) { c.String(200, "user") })
	router.POST("/users", func(c *Context) { c.String(201, "create") })

	// Test concurrent access
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- true }()

			// Each goroutine performs route lookups
			for j := 0; j < 100; j++ {
				node, params := router.getRoute("GET", "/users/123")
				if node == nil {
					t.Errorf("Goroutine %d: Expected route to be found", id)
					return
				}
				if params["id"] != "123" {
					t.Errorf("Goroutine %d: Expected id = 123, got %s", id, params["id"])
					return
				}
			}
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

func BenchmarkRouterStaticRoute(b *testing.B) {
	router := NewRouter()
	router.GET("/api/v1/users", func(c *Context) { c.String(200, "OK") })

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.getRoute("GET", "/api/v1/users")
	}
}

func BenchmarkRouterParamRoute(b *testing.B) {
	router := NewRouter()
	router.GET("/users/:id/posts/:postId", func(c *Context) { c.String(200, "OK") })

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.getRoute("GET", "/users/123/posts/456")
	}
}

func BenchmarkRouterWildcardRoute(b *testing.B) {
	router := NewRouter()
	router.GET("/files/*filepath", func(c *Context) { c.String(200, "OK") })

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.getRoute("GET", "/files/images/avatars/user123.png")
	}
}

func BenchmarkParsePattern(b *testing.B) {
	pattern := "/api/v1/users/:id/posts/:postId/comments/*commentPath"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parsePattern(pattern)
	}
}

// Integration test with HTTP server
func TestRouterHTTPIntegration(t *testing.T) {
	router := NewRouter()

	router.GET("/hello/:name", func(c *Context) {
		name := c.Param("name")
		c.JSON(200, map[string]string{"message": "Hello " + name})
	})

	// Create a test request
	req := httptest.NewRequest("GET", "/hello/world", nil)
	w := httptest.NewRecorder()

	// Simulate the engine's ServeHTTP behavior
	node, params := router.getRoute(req.Method, req.URL.Path)
	if node == nil {
		t.Fatal("Route should be found")
	}

	// Create context and execute handlers
	c := NewContext(w, req)
	c.params = params
	c.handlers = node.handlers
	c.Next()

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	expected := `{"message":"Hello world"}`
	actual := w.Body.String()
	// Remove trailing newline from JSON encoder
	if actual[len(actual)-1] == '\n' {
		actual = actual[:len(actual)-1]
	}

	if actual != expected {
		t.Errorf("Expected response %s, got %s", expected, actual)
	}
}
