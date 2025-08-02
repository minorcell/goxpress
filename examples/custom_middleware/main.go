package main

// import goxpress
import "github.com/minorcell/goxpress"

// This example demonstrates how to create custom middleware in goxpress
// Custom middleware allows you to add your own functionality to the request processing pipeline
func main() {
	// Step 1: Create a new instance of goxpress.Engine
	app := goxpress.New()

	// Step 2: Register global middleware
	// Logger middleware logs all requests
	// CORS middleware adds Cross-Origin Resource Sharing headers
	app.Use(goxpress.Logger())
	app.Use(CORSMiddleware())

	// Step 3: Create a route group with additional middleware
	// Protected routes require authentication
	protected := app.Route("/api")
	protected.Use(AuthMiddleware()) // Only applies to routes in this group

	// Step 4: Define routes
	protected.GET("/profile", func(c *goxpress.Context) {
		userID, _ := c.GetString("user_id")
		c.JSON(200, map[string]string{
			"user_id": userID,
			"profile": "Here's the user profile data",
		})
	})

	// Step 5: Start the server
	app.Listen(":8080", func() {
		println("Custom middleware example running at http://localhost:8080")
		println("Try: curl -H \"Authorization: Bearer valid-token\" http://localhost:8080/api/profile")
	})
}

// AuthMiddleware is a custom authentication middleware
// It checks for a valid Authorization header
func AuthMiddleware() goxpress.HandlerFunc {
	return func(c *goxpress.Context) {
		token := c.Request.Header.Get("Authorization")

		if token == "" {
			c.JSON(401, map[string]string{"error": "Oops, forgot to bring the token"})
			c.Abort() // Stop further processing
			return
		}

		// Validate token (simplified here, in real projects you might need JWT or other methods)
		if token != "Bearer valid-token" {
			c.JSON(401, map[string]string{"error": "Wrong token"})
			c.Abort()
			return
		}

		// Store user info in context for later handlers to use
		c.Set("user_id", "12345")
		c.Next() // Continue to next middleware/handler
	}
}

// CORSMiddleware is a custom CORS (Cross-Origin Resource Sharing) middleware
// It adds the necessary headers to allow cross-origin requests
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