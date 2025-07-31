// Package goxpress provides a fast, intuitive web framework for Go inspired by Express.js.
// It offers Express.js-like API design with excellent performance and full Go type safety.
//
// goxpress is built for developer productivity while maintaining high performance.
// It features a powerful middleware system, efficient routing with Radix Tree algorithm,
// and comprehensive request/response handling capabilities.
//
// Basic Usage:
//
//	app := goxpress.New()
//	app.GET("/", func(c *goxpress.Context) {
//		c.String(200, "Hello, World!")
//	})
//	app.Listen(":8080", nil)
//
// Middleware Support:
//
//	app.Use(goxpress.Logger())
//	app.Use(goxpress.Recover())
//	app.Use(func(c *goxpress.Context) {
//		// Custom middleware logic
//		c.Next()
//	})
//
// Route Groups:
//
//	api := app.Route("/api")
//	api.GET("/users", getUsersHandler)
//	api.POST("/users", createUserHandler)
//
// Error Handling:
//
//	app.UseError(func(err error, c *goxpress.Context) {
//		c.JSON(500, map[string]string{"error": err.Error()})
//	})
package goxpress

import (
	"net/http"
)

// HandlerFunc defines the signature for HTTP request handlers.
// It receives a Context that provides access to request data, response writing,
// and flow control mechanisms.
//
// Example:
//
//	func myHandler(c *goxpress.Context) {
//		id := c.Param("id")
//		c.JSON(200, map[string]string{"user_id": id})
//	}
type HandlerFunc func(*Context)

// ErrorHandlerFunc defines the signature for error handling middleware.
// It receives an error and a Context, allowing custom error processing
// and response generation.
//
// Example:
//
//	func errorHandler(err error, c *goxpress.Context) {
//		c.JSON(500, map[string]string{"error": err.Error()})
//	}
type ErrorHandlerFunc func(error, *Context)

// Engine represents the main goxpress application instance.
// It implements the http.Handler interface and coordinates routing,
// middleware execution, and request processing.
//
// The Engine manages:
//   - HTTP routing system
//   - Global middleware chain
//   - Error handling middleware
//   - HTTP server lifecycle
//
// Create a new Engine instance using New().
type Engine struct {
	router        *Router            // HTTP router for request matching
	middlewares   []HandlerFunc      // Global middleware functions
	errorHandlers []ErrorHandlerFunc // Error handling middleware
}

// New creates and returns a new Engine instance with default configuration.
// The returned Engine is ready to accept route registrations and middleware.
//
// Example:
//
//	app := goxpress.New()
//	app.GET("/", handler)
//	app.Listen(":8080", nil)
func New() *Engine {
	engine := &Engine{
		router:        NewRouter(),
		middlewares:   make([]HandlerFunc, 0),
		errorHandlers: make([]ErrorHandlerFunc, 0),
	}
	return engine
}

// Use registers global middleware functions that will be executed for all requests.
// Middleware are executed in the order they are registered.
// Returns the Engine instance for method chaining.
//
// Middleware functions can:
//   - Perform preprocessing (authentication, logging, etc.)
//   - Modify request/response
//   - Control request flow with c.Next() or c.Abort()
//   - Handle errors by passing them to c.Next(err)
//
// Example:
//
//	app.Use(Logger()).Use(Recover()).Use(func(c *Context) {
//		// Custom middleware logic
//		c.Set("start_time", time.Now())
//		c.Next()
//	})
func (e *Engine) Use(middleware ...HandlerFunc) *Engine {
	e.middlewares = append(e.middlewares, middleware...)
	return e
}

// UseError registers error handling middleware that will be called when
// errors occur during request processing. Error handlers are executed
// in the order they are registered.
// Returns the Engine instance for method chaining.
//
// Error handlers are triggered when:
//   - A handler calls c.Next(err) with a non-nil error
//   - A panic occurs and is recovered by the Recover middleware
//
// Example:
//
//	app.UseError(func(err error, c *Context) {
//		log.Printf("Error: %v", err)
//		c.JSON(500, map[string]string{"error": "Internal Server Error"})
//	})
func (e *Engine) UseError(handler ...ErrorHandlerFunc) *Engine {
	e.errorHandlers = append(e.errorHandlers, handler...)
	return e
}

// GET registers a new route for HTTP GET requests.
// Returns the Engine instance for method chaining.
//
// The pattern supports:
//   - Static paths: "/users"
//   - Parameters: "/users/:id"
//   - Wildcards: "/files/*filepath"
//
// Example:
//
//	app.GET("/users/:id", func(c *Context) {
//		id := c.Param("id")
//		c.JSON(200, map[string]string{"user_id": id})
//	})
func (e *Engine) GET(pattern string, handlers ...HandlerFunc) *Engine {
	e.router.GET(pattern, handlers...)
	return e
}

// POST registers a new route for HTTP POST requests.
// Returns the Engine instance for method chaining.
//
// Example:
//
//	app.POST("/users", func(c *Context) {
//		var user User
//		if err := c.BindJSON(&user); err != nil {
//			c.JSON(400, map[string]string{"error": "Invalid JSON"})
//			return
//		}
//		c.JSON(201, user)
//	})
func (e *Engine) POST(pattern string, handlers ...HandlerFunc) *Engine {
	e.router.POST(pattern, handlers...)
	return e
}

// PUT registers a new route for HTTP PUT requests.
// Returns the Engine instance for method chaining.
//
// Example:
//
//	app.PUT("/users/:id", updateUserHandler)
func (e *Engine) PUT(pattern string, handlers ...HandlerFunc) *Engine {
	e.router.PUT(pattern, handlers...)
	return e
}

// DELETE registers a new route for HTTP DELETE requests.
// Returns the Engine instance for method chaining.
//
// Example:
//
//	app.DELETE("/users/:id", deleteUserHandler)
func (e *Engine) DELETE(pattern string, handlers ...HandlerFunc) *Engine {
	e.router.DELETE(pattern, handlers...)
	return e
}

// PATCH registers a new route for HTTP PATCH requests.
// Returns the Engine instance for method chaining.
//
// Example:
//
//	app.PATCH("/users/:id", patchUserHandler)
func (e *Engine) PATCH(pattern string, handlers ...HandlerFunc) *Engine {
	e.router.PATCH(pattern, handlers...)
	return e
}

// HEAD registers a new route for HTTP HEAD requests.
// Returns the Engine instance for method chaining.
//
// Example:
//
//	app.HEAD("/users/:id", headUserHandler)
func (e *Engine) HEAD(pattern string, handlers ...HandlerFunc) *Engine {
	e.router.HEAD(pattern, handlers...)
	return e
}

// OPTIONS registers a new route for HTTP OPTIONS requests.
// Returns the Engine instance for method chaining.
//
// Example:
//
//	app.OPTIONS("/users", optionsUserHandler)
func (e *Engine) OPTIONS(pattern string, handlers ...HandlerFunc) *Engine {
	e.router.OPTIONS(pattern, handlers...)
	return e
}

// Route creates a new route group with the specified prefix.
// Route groups allow organizing related routes and applying
// group-specific middleware.
//
// Example:
//
//	api := app.Route("/api")
//	api.GET("/users", getUsersHandler)
//	api.POST("/users", createUserHandler)
//
//	v1 := api.Group("/v1")
//	v1.GET("/status", statusHandler)
func (e *Engine) Route(prefix string) *Router {
	return e.router.Group(prefix)
}

// ServeHTTP implements the http.Handler interface, making Engine compatible
// with the standard net/http package. This method handles all incoming HTTP
// requests by:
//
//  1. Creating a Context from the object pool
//  2. Matching the request to a route
//  3. Executing middleware chain and route handlers
//  4. Handling any errors that occur
//  5. Returning the Context to the pool
//
// This method is called automatically by the HTTP server and should not
// be called directly in normal usage.
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Get Context from pool for efficient memory usage
	c := NewContext(w, req)

	// Ensure Context is returned to pool after request processing
	defer func() {
		c.reset()
		contextPool.Put(c)
	}()

	// Find matching route for the request
	node, params := e.router.getRoute(req.Method, req.URL.Path)

	// Set URL parameters if route was found
	if params != nil {
		c.params = params
	}

	// Build handler chain: global middleware + route handlers
	handlers := make([]HandlerFunc, 0)
	handlers = append(handlers, e.middlewares...)

	if node != nil {
		// Route found: add route-specific handlers
		handlers = append(handlers, node.handlers...)
	} else {
		// No route found: add 404 handler
		handlers = append(handlers, func(c *Context) {
			c.Status(http.StatusNotFound)
			c.String(http.StatusNotFound, "404 page not found")
		})
	}

	c.handlers = handlers

	// Execute the handler chain
	c.Next()

	// Process any errors that occurred during request handling
	if c.err != nil && len(e.errorHandlers) > 0 {
		for _, handler := range e.errorHandlers {
			handler(c.err, c)
		}
	}
}

// Listen starts an HTTP server on the specified address.
// The callback function is called after the server starts but before
// it begins accepting connections.
//
// This is a blocking call that will run until the server is stopped
// or encounters an error.
//
// Example:
//
//	app.Listen(":8080", func() {
//		log.Println("Server started on :8080")
//	})
func (e *Engine) Listen(addr string, cb func()) error {
	server := &http.Server{
		Addr:    addr,
		Handler: e,
	}

	if cb != nil {
		cb()
	}

	return server.ListenAndServe()
}

// ListenTLS starts an HTTPS server on the specified address using
// the provided TLS certificate and key files.
// The callback function is called after the server starts but before
// it begins accepting connections.
//
// This is a blocking call that will run until the server is stopped
// or encounters an error.
//
// Example:
//
//	app.ListenTLS(":443", "cert.pem", "key.pem", func() {
//		log.Println("HTTPS Server started on :443")
//	})
func (e *Engine) ListenTLS(addr, certFile, keyFile string, cb func()) error {
	server := &http.Server{
		Addr:    addr,
		Handler: e,
	}

	if cb != nil {
		cb()
	}

	return server.ListenAndServeTLS(certFile, keyFile)
}
