package main

// import goxpress
import "github.com/minorcell/goxpress"

func main() {
	// Step 1: Create a new instance of goxpress.Engine
	app := goxpress.New()

	// Step 2: Define routes, goxexpress supports all HTTP methods
	app.GET("/", func(c *goxpress.Context) {
		c.String(200, "Hello, World!")
	})

	// Step 3: Start the server, you can specify the running port and the fallback
	app.Listen(":8080", func() {
		println("Server running at http://localhost:8080")
	})
}
