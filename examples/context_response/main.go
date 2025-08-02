package main

// import goxpress
import "github.com/minorcell/goxpress"

// This example demonstrates different types of responses in goxpress
// It shows how to return JSON, plain text, HTML, redirects, and files
func main() {
	// Step 1: Create a new instance of goxpress.Engine
	app := goxpress.New()

	// Step 2: Define routes that return different types of responses

	// JSON response
	app.GET("/api/data", func(c *goxpress.Context) {
		c.JSON(200, map[string]string{"message": "JSON data"})
	})

	// Plain text response
	app.GET("/text", func(c *goxpress.Context) {
		c.String(200, "This is some text")
	})

	// HTML response
	app.GET("/html", func(c *goxpress.Context) {
		c.HTML(200, "<h1>Hello HTML</h1>")
	})

	// Redirect
	app.GET("/redirect", func(c *goxpress.Context) {
		c.Redirect(302, "https://github.com/minorcell/goxpress")
	})

	// File response
	app.GET("/file", func(c *goxpress.Context) {
		c.File("example.txt")
	})

	// Step 3: Start the server
	app.Listen(":8080", func() {
		println("Context response example running at http://localhost:8080")
	})
}