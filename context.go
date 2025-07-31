// Package goxpress provides a fast, intuitive web framework for Go inspired by Express.js.
// This file contains the Context implementation that wraps HTTP request/response
// and provides convenient methods for handling web requests.
package goxpress

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// contextPool is a sync.Pool for Context objects to reduce GC pressure
// and improve performance by reusing Context instances.
var contextPool = sync.Pool{
	New: func() interface{} {
		return &Context{
			params: make(map[string]string),
			store:  make(map[string]interface{}),
			index:  -1,
		}
	},
}

// Context represents the context of the current HTTP request.
// It wraps the http.Request and http.ResponseWriter and provides
// convenient methods for handling request data, generating responses,
// and managing request flow through middleware chains.
//
// Context provides:
//   - URL parameter extraction
//   - Query parameter access
//   - JSON request/response handling
//   - Middleware flow control
//   - Request-scoped data storage
//   - Error handling
//
// Context instances are pooled for efficient memory usage and should
// not be stored beyond the scope of a single request.
type Context struct {
	// Embedded standard context for cancellation and deadlines
	context.Context

	// HTTP request and response
	Request  *http.Request       // Original HTTP request
	Response http.ResponseWriter // HTTP response writer

	// URL parameters extracted from route patterns
	params map[string]string

	// Middleware chain management
	handlers []HandlerFunc // Chain of handlers to execute
	index    int           // Current position in handler chain

	// Request flow control
	aborted bool // Whether request processing should be aborted

	// Response state tracking
	statusCodeWritten bool // Whether response status has been written

	// Error handling
	err error // Error that occurred during request processing

	// Request-scoped data storage
	store map[string]interface{} // Key-value store for request data
}

// NewContext creates a new Context instance from the pool and initializes it
// with the given HTTP request and response writer.
//
// This function is used internally by the framework and typically should not
// be called directly by application code.
func NewContext(w http.ResponseWriter, req *http.Request) *Context {
	// Get Context from pool or create new one
	c := contextPool.Get().(*Context)

	// Initialize request-related fields
	c.Context = req.Context()
	c.Request = req
	c.Response = w

	// Reset state fields
	c.index = -1
	c.aborted = false
	c.statusCodeWritten = false
	c.err = nil

	return c
}

// reset clears the Context state and prepares it for return to the pool.
// This method is called internally to clean up Context instances before
// they are returned to the pool for reuse.
func (c *Context) reset() {
	// Clear maps
	for k := range c.params {
		delete(c.params, k)
	}

	for k := range c.store {
		delete(c.store, k)
	}

	// Reset other fields
	c.Context = nil
	c.Request = nil
	c.Response = nil
	c.handlers = nil
	c.index = -1
	c.aborted = false
	c.statusCodeWritten = false
	c.err = nil
}

// Param returns the value of the URL parameter with the given name.
// URL parameters are extracted from route patterns like "/users/:id".
//
// Example:
//
//	// Route: "/users/:id"
//	// Request: "/users/123"
//	id := c.Param("id") // Returns "123"
func (c *Context) Param(key string) string {
	return c.params[key]
}

// Query returns the value of the URL query parameter with the given name.
// Returns an empty string if the parameter doesn't exist.
//
// Example:
//
//	// Request: "/search?q=golang&page=1"
//	query := c.Query("q")    // Returns "golang"
//	page := c.Query("page")  // Returns "1"
//	empty := c.Query("foo")  // Returns ""
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// BindJSON parses the request body as JSON and stores the result
// in the value pointed to by obj. The request body is consumed
// during this operation.
//
// Example:
//
//	var user struct {
//		Name  string `json:"name"`
//		Email string `json:"email"`
//	}
//	if err := c.BindJSON(&user); err != nil {
//		c.JSON(400, map[string]string{"error": "Invalid JSON"})
//		return
//	}
func (c *Context) BindJSON(obj interface{}) error {
	return json.NewDecoder(c.Request.Body).Decode(obj)
}

// Status sets the HTTP status code for the response.
// This method can only be called once per request; subsequent calls are ignored.
//
// Example:
//
//	c.Status(201) // Set status to 201 Created
func (c *Context) Status(code int) {
	if !c.statusCodeWritten {
		c.Response.WriteHeader(code)
		c.statusCodeWritten = true
	}
}

// JSON serializes the given object to JSON and writes it to the response
// with the specified status code. It automatically sets the Content-Type
// header to "application/json".
//
// Example:
//
//	c.JSON(200, map[string]interface{}{
//		"message": "success",
//		"data":    user,
//	})
func (c *Context) JSON(code int, obj interface{}) error {
	if !c.statusCodeWritten {
		c.Response.Header().Set("Content-Type", "application/json")
		c.Response.WriteHeader(code)
		c.statusCodeWritten = true
	}
	return json.NewEncoder(c.Response).Encode(obj)
}

// String writes a formatted string to the response with the specified status code.
// It automatically sets the Content-Type header to "text/plain; charset=utf-8".
//
// Example:
//
//	c.String(200, "Hello %s", name)
//	c.String(404, "Page not found")
func (c *Context) String(code int, format string, values ...interface{}) error {
	if !c.statusCodeWritten {
		c.Response.Header().Set("Content-Type", "text/plain; charset=utf-8")
		c.Response.WriteHeader(code)
		c.statusCodeWritten = true
	}
	_, err := c.Response.Write([]byte(fmt.Sprintf(format, values...)))
	return err
}

// Next executes the next handler in the middleware chain.
// If an error is provided, it will be stored in the context
// for later processing by error handlers.
//
// This method should be called by middleware to continue processing
// the request. If not called, the request processing stops.
//
// Example:
//
//	func MyMiddleware(c *Context) {
//		// Pre-processing
//		log.Println("Before")
//
//		c.Next() // Execute subsequent handlers
//
//		// Post-processing
//		log.Println("After")
//	}
func (c *Context) Next(err ...error) {
	// Store error if provided
	if len(err) > 0 && err[0] != nil {
		c.err = err[0]
	}

	// Advance to next handler and execute it
	c.index++
	for c.index < len(c.handlers) {
		// Stop if request has been aborted
		if c.aborted {
			return
		}

		handler := c.handlers[c.index]
		handler(c)
		c.index++
	}
}

// Abort prevents any pending handlers from being called.
// Note that this will not stop the current handler from executing.
// For instance, if you have an authorization middleware that validates
// each request, you might want to call Abort if the authorization fails.
//
// Example:
//
//	func AuthMiddleware(c *Context) {
//		if !isAuthorized(c) {
//			c.JSON(401, map[string]string{"error": "Unauthorized"})
//			c.Abort()
//			return
//		}
//		c.Next()
//	}
func (c *Context) Abort() {
	c.aborted = true
}

// IsAborted returns true if the current context was aborted.
// This can be used to check if request processing should continue.
//
// Example:
//
//	if c.IsAborted() {
//		return // Stop processing
//	}
func (c *Context) IsAborted() bool {
	return c.aborted
}

// Set stores a key-value pair in the context's data store.
// This data is available throughout the request lifecycle and
// can be accessed by subsequent middleware and handlers.
//
// Example:
//
//	c.Set("user_id", "123")
//	c.Set("start_time", time.Now())
func (c *Context) Set(key string, value interface{}) {
	c.store[key] = value
}

// Get retrieves a value from the context's data store.
// Returns the value and a boolean indicating whether the key exists.
//
// Example:
//
//	if user, exists := c.Get("user"); exists {
//		// Use user data
//		fmt.Println("User:", user)
//	}
func (c *Context) Get(key string) (interface{}, bool) {
	value, exists := c.store[key]
	return value, exists
}

// MustGet retrieves a value from the context's data store.
// It panics if the key doesn't exist. Use this only when you're
// certain the key exists.
//
// Example:
//
//	user := c.MustGet("user").(User)
func (c *Context) MustGet(key string) interface{} {
	if value, exists := c.store[key]; exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

// GetString retrieves a string value from the context's data store.
// Returns the string value and a boolean indicating success.
// Returns false if the key doesn't exist or the value is not a string.
//
// Example:
//
//	if name, ok := c.GetString("user_name"); ok {
//		fmt.Println("User name:", name)
//	}
func (c *Context) GetString(key string) (string, bool) {
	if val, ok := c.Get(key); ok && val != nil {
		if str, ok := val.(string); ok {
			return str, true
		}
	}
	return "", false
}
