package main

// import goxpress
import "github.com/minorcell/goxpress"

// This example demonstrates how to use route groups in goxpress
// Route groups help organize routes and apply middleware to specific groups
func main() {
	// Step 1: Create a new instance of goxpress.Engine
	app := goxpress.New()

	// Step 2: Create API v1 group
	v1 := app.Route("/api/v1")
	v1.GET("/users", func(c *goxpress.Context) {
		c.JSON(200, map[string]string{"version": "v1", "users": "User list"})
	})
	v1.GET("/posts", func(c *goxpress.Context) {
		c.JSON(200, map[string]string{"version": "v1", "posts": "Post list"})
	})

	// Step 3: Create API v2 group
	v2 := app.Route("/api/v2")
	v2.GET("/users", func(c *goxpress.Context) {
		c.JSON(200, map[string]string{"version": "v2", "users": "User list (new version)"})
	})

	// Step 4: Start the server
	app.Listen(":8080", func() {
		println("Route groups example running at http://localhost:8080")
	})
}