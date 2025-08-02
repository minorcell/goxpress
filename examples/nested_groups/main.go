package main

// import goxpress
import "github.com/minorcell/goxpress"

// This example demonstrates how to use nested route groups with middleware in goxpress
// Nested groups allow for more complex route organization and selective middleware application
func main() {
	// Step 1: Create a new instance of goxpress.Engine
	app := goxpress.New()

	// Step 2: Register global middleware
	// This middleware applies to all routes
	app.Use(goxpress.Logger())

	// Step 3: Create API group with its own middleware
	api := app.Route("/api")
	api.Use(CORSMiddleware())

	// Step 4: Create public APIs subgroup
	public := api.Group("/public")
	public.GET("/health", func(c *goxpress.Context) {
		c.JSON(200, map[string]string{"status": "OK"})
	})

	// Step 5: Create protected APIs subgroup with authentication middleware
	protected := api.Group("/protected")
	protected.Use(AdminMiddleware())
	protected.GET("/admin", func(c *goxpress.Context) {
		c.JSON(200, map[string]string{"message": "Admin-only endpoint"})
	})
	protected.DELETE("/users/:id", func(c *goxpress.Context) {
		id := c.Param("id")
		c.JSON(200, map[string]string{"message": "User " + id + " deleted"})
	})

	// Step 6: Start the server
	app.Listen(":8080", func() {
		println("Nested groups example running at http://localhost:8080")
	})
}

// CORSMiddleware adds Cross-Origin Resource Sharing headers to responses
func CORSMiddleware() goxpress.HandlerFunc {
	return func(c *goxpress.Context) {
		c.Response.Header().Set("Access-Control-Allow-Origin", "*")
		c.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Response.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.Status(204)
			return
		}

		c.Next()
	}
}

// AdminMiddleware checks if the user has admin privileges
func AdminMiddleware() goxpress.HandlerFunc {
	return func(c *goxpress.Context) {
		// Here we should check if user is admin
		// For demo purposes, we simplify
		role := c.Request.Header.Get("User-Role")
		if role != "admin" {
			c.JSON(403, map[string]string{"error": "Admin privileges required"})
			c.Abort()
			return
		}
		c.Next()
	}
}