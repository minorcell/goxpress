// Package goxpress provides a fast, intuitive web framework for Go inspired by Express.js.
// This file contains the HTTP routing system implementation using Radix Tree algorithm
// for efficient route matching and parameter extraction.
package goxpress

import (
	"strings"
	"sync"
)

// builderPool is a sync.Pool for strings.Builder to reduce memory allocations
// during path parsing operations.
var builderPool = sync.Pool{
	New: func() interface{} {
		return &strings.Builder{}
	},
}

// Router represents the HTTP router that manages route registration and matching.
// It uses a Radix Tree data structure for efficient route lookup and supports:
//   - Static routes: "/users"
//   - Parameter routes: "/users/:id"
//   - Wildcard routes: "/files/*filepath"
//   - Route groups with shared prefixes and middleware
//
// The Router is safe for concurrent read access after route registration is complete.
type Router struct {
	prefix      string                 // Route group prefix
	middlewares []HandlerFunc          // Group-specific middleware
	engine      *Engine                // Reference to parent engine
	subRouters  map[string]*Router     // Nested route groups
	routes      map[string]*routerTree // HTTP method -> route tree mapping
}

// routerTree implements a Radix Tree for efficient route matching.
// Each HTTP method has its own tree to avoid conflicts between
// different HTTP verbs on the same path.
type routerTree struct {
	root *routerNode // Root node of the tree
}

// routerNode represents a single node in the Radix Tree.
// Each node can represent part of a URL path and may contain
// handlers if it represents a complete route.
type routerNode struct {
	pattern  string        // Complete route pattern (e.g., "/users/:id")
	part     string        // Path segment for this node (e.g., ":id")
	children []*routerNode // Child nodes
	isWild   bool          // True if this node represents a parameter or wildcard
	handlers []HandlerFunc // Route handlers (only set for terminal nodes)
}

// NewRouter creates and returns a new Router instance.
// The router is initialized with empty route trees for all HTTP methods.
//
// Example:
//
//	router := NewRouter()
//	router.GET("/users", getUsersHandler)
func NewRouter() *Router {
	return &Router{
		subRouters: make(map[string]*Router),
		routes:     make(map[string]*routerTree),
	}
}

// Use registers middleware functions for this router group.
// Middleware registered on a router will only apply to routes
// defined on that router and its sub-groups.
// Returns the Router instance for method chaining.
//
// Example:
//
//	api := app.Route("/api")
//	api.Use(AuthMiddleware()).Use(LoggingMiddleware())
func (r *Router) Use(middleware ...HandlerFunc) *Router {
	r.middlewares = append(r.middlewares, middleware...)
	return r
}

// Group creates a new sub-router with the given prefix.
// The sub-router inherits middleware from its parent and can
// define additional middleware that only applies to its routes.
//
// Example:
//
//	api := app.Route("/api")
//	v1 := api.Group("/v1")  // Routes will have "/api/v1" prefix
//	v1.GET("/users", handler)  // Handles "/api/v1/users"
func (r *Router) Group(prefix string) *Router {
	router := &Router{
		prefix:      r.prefix + prefix,
		middlewares: make([]HandlerFunc, len(r.middlewares)), // Copy parent middleware
		engine:      r.engine,
		subRouters:  make(map[string]*Router),
		routes:      r.routes, // Share route trees with parent
	}

	// Copy parent middleware to new router
	copy(router.middlewares, r.middlewares)

	r.subRouters[prefix] = router
	return router
}

// Handle registers a new route with the specified HTTP method and pattern.
// This is the core route registration method used by all HTTP method helpers.
//
// The method combines the router's prefix with the pattern and prepares
// the final handler chain including group middleware.
func (r *Router) Handle(method, pattern string, handlers ...HandlerFunc) {
	// Combine router prefix with route pattern
	fullPattern := r.prefix + pattern
	if r.prefix != "" && pattern == "/" {
		fullPattern = r.prefix
	}

	// Build final handler chain: group middleware + route handlers
	finalHandlers := make([]HandlerFunc, 0)
	finalHandlers = append(finalHandlers, r.middlewares...)
	finalHandlers = append(finalHandlers, handlers...)

	// Register the route
	r.addRoute(method, fullPattern, finalHandlers)
}

// GET registers a new route for HTTP GET requests.
// Returns the Router instance for method chaining.
//
// Example:
//
//	router.GET("/users/:id", getUserHandler)
func (r *Router) GET(pattern string, handlers ...HandlerFunc) *Router {
	r.Handle("GET", pattern, handlers...)
	return r
}

// POST registers a new route for HTTP POST requests.
// Returns the Router instance for method chaining.
func (r *Router) POST(pattern string, handlers ...HandlerFunc) *Router {
	r.Handle("POST", pattern, handlers...)
	return r
}

// PUT registers a new route for HTTP PUT requests.
// Returns the Router instance for method chaining.
func (r *Router) PUT(pattern string, handlers ...HandlerFunc) *Router {
	r.Handle("PUT", pattern, handlers...)
	return r
}

// DELETE registers a new route for HTTP DELETE requests.
// Returns the Router instance for method chaining.
func (r *Router) DELETE(pattern string, handlers ...HandlerFunc) *Router {
	r.Handle("DELETE", pattern, handlers...)
	return r
}

// PATCH registers a new route for HTTP PATCH requests.
// Returns the Router instance for method chaining.
func (r *Router) PATCH(pattern string, handlers ...HandlerFunc) *Router {
	r.Handle("PATCH", pattern, handlers...)
	return r
}

// HEAD registers a new route for HTTP HEAD requests.
// Returns the Router instance for method chaining.
func (r *Router) HEAD(pattern string, handlers ...HandlerFunc) *Router {
	r.Handle("HEAD", pattern, handlers...)
	return r
}

// OPTIONS registers a new route for HTTP OPTIONS requests.
// Returns the Router instance for method chaining.
func (r *Router) OPTIONS(pattern string, handlers ...HandlerFunc) *Router {
	r.Handle("OPTIONS", pattern, handlers...)
	return r
}

// parsePattern splits a URL pattern into path segments, removing empty segments.
// It uses a pool of strings.Builder for efficient string operations.
//
// Examples:
//
//	"/users/:id" -> ["users", ":id"]
//	"/api/v1/users" -> ["api", "v1", "users"]
//	"/files/*filepath" -> ["files", "*filepath"]
func parsePattern(pattern string) []string {
	// Get builder from pool for efficient string operations
	builder := builderPool.Get().(*strings.Builder)
	defer func() {
		builder.Reset()
		builderPool.Put(builder)
	}()

	// Pre-allocate slice with estimated capacity
	parts := make([]string, 0, strings.Count(pattern, "/"))

	// Split by '/' and filter out empty parts
	segments := strings.Split(pattern, "/")
	for _, segment := range segments {
		if segment != "" {
			parts = append(parts, segment)
		}
	}

	return parts
}

// addRoute adds a new route to the appropriate route tree.
// It creates the tree for the HTTP method if it doesn't exist,
// then inserts the route pattern into the Radix Tree.
func (r *Router) addRoute(method, pattern string, handlers []HandlerFunc) {
	// Create route tree for method if it doesn't exist
	if r.routes[method] == nil {
		r.routes[method] = &routerTree{root: &routerNode{}}
	}

	parts := parsePattern(pattern)

	// Insert pattern into the Radix Tree
	r.routes[method].insertRoute(pattern, parts, 0, handlers)
}

// getRoute finds a matching route for the given HTTP method and path.
// Returns the matching node and extracted URL parameters, or nil if no match.
//
// The method performs efficient tree traversal to find the best match,
// extracting parameters along the way.
func (r *Router) getRoute(method, path string) (*routerNode, map[string]string) {
	root, ok := r.routes[method]
	if !ok {
		return nil, nil
	}

	searchParts := parsePattern(path)
	params := make(map[string]string)

	node := root.searchRoute(searchParts, 0, params)

	return node, params
}

// walkMountRoutes recursively walks through route tree nodes to mount routes
// from sub-routers. This is used internally for route group management.
func (r *Router) walkMountRoutes(node *routerNode, method, mountPrefix string, groupMiddlewares []HandlerFunc, addRoute func(method, pattern string, handlers []HandlerFunc)) {
	// If this is a root node, recursively process children
	if node.pattern == "" {
		for _, child := range node.children {
			r.walkMountRoutes(child, method, mountPrefix, groupMiddlewares, addRoute)
		}
		return
	}

	// Calculate the full route pattern
	pattern := strings.TrimSuffix(mountPrefix, "/") + strings.TrimPrefix(node.pattern, r.prefix)

	// Combine group middleware with route handlers
	finalHandlers := make([]HandlerFunc, 0)
	finalHandlers = append(finalHandlers, groupMiddlewares...)
	finalHandlers = append(finalHandlers, node.handlers...)

	// Add the route
	addRoute(method, pattern, finalHandlers)
}

// insertRoute recursively inserts a route pattern into the Radix Tree.
// It builds the tree structure by creating nodes for each path segment
// and handles parameter and wildcard matching.
func (t *routerTree) insertRoute(pattern string, parts []string, height int, handlers []HandlerFunc) {
	// Base case: all segments processed
	if len(parts) == height {
		t.root.pattern = pattern
		t.root.handlers = handlers
		return
	}

	part := parts[height]
	child := t.root.matchChild(part)

	if child == nil {
		// Create new child node
		child = &routerNode{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		t.root.children = append(t.root.children, child)
	}

	// Recursively insert remaining parts
	childTree := &routerTree{root: child}
	childTree.insertRoute(pattern, parts, height+1, handlers)
}

// searchRoute performs recursive search through the Radix Tree to find
// a matching route. It extracts URL parameters during traversal.
func (t *routerTree) searchRoute(parts []string, height int, params map[string]string) *routerNode {
	// Base case: all parts processed or wildcard encountered
	if len(parts) == height || strings.HasPrefix(t.root.part, "*") {
		if t.root.pattern == "" {
			return nil
		}
		return t.root
	}

	part := parts[height]
	// Check all children for matches
	for _, child := range t.root.children {
		if child.part == part || child.isWild {
			// Handle parameter matching
			if child.isWild && child.part[0] == ':' {
				params[child.part[1:]] = part
			} else if child.isWild && child.part[0] == '*' {
				// For wildcard, capture the rest of the path
				params[child.part[1:]] = strings.Join(parts[height:], "/")
				return child
			}

			// Recursively search in child node
			childTree := &routerTree{root: child}
			result := childTree.searchRoute(parts, height+1, params)
			if result != nil {
				return result
			}

			// Backtrack parameters if necessary
			if child.isWild && child.part[0] == ':' {
				delete(params, child.part[1:])
			}
		}
	}

	return nil
}

// matchChild finds a direct child node that matches the given part.
// Returns nil if no exact match is found.
func (n *routerNode) matchChild(part string) *routerNode {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}
