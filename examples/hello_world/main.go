package main

import "github.com/minorcell/goxpress"

func main() {
	app := goxpress.New()

	app.GET("/", func(c *goxpress.Context) {
		c.String(200, "Hello, World!")
	})

	app.Listen(":8080", func() {
		println("Server running at http://localhost:8080")
	})
}
