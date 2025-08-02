package main

// import goxpress
import "github.com/minorcell/goxpress"

// This example demonstrates how to use built-in middleware in goxpress
// Built-in middleware provides common functionality like logging and panic recovery
func main() {
	// Step 1: Create a new instance of goxpress.Engine
	app := goxpress.New()

	// Step 2: Register built-in middleware
	// Logger middleware logs HTTP requests
	// Recover middleware recovers from panics to prevent server crashes
	app.Use(goxpress.Logger())  // Request logging
	// or you can use logger whit config
	// app.Use(goxpress.LoggerWithConfig(goxpress.LoggerConfig{
	// 	Output: ios.Stdout,
	//  Formatter: goxpress.DefaultLogFormatter,
	//  SkipPaths: []string{"/health"}
	// }))
	app.Use(goxpress.Recover()) // Panic recovery

	// Step 3: Define routes
	app.GET("/", func(c *goxpress.Context) {
		c.String(200, "Hello with middleware!")
	})

	// Step 4: Start the server
	app.Listen(":8080", func() {
		println("Middleware example running at http://localhost:8080")
	})
}