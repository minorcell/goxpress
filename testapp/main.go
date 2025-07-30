package main

import (
	"fmt"

	"relay"
)

func main() {
	app := relay.New().
		Use(relay.Logger()).
		Use(relay.Recover())

	app.GET("/", func(c *relay.Context) {
		c.String(200, "Hello, Relay!")
	})

	app.GET("/users/:id", func(c *relay.Context) {
		id := c.Param("id")
		c.JSON(200, map[string]interface{}{
			"id":   id,
			"name": "User " + id,
		})
	})

	api := app.Route("/api").
		Use(func(c *relay.Context) {
			c.Set("version", "v1")
			c.Next()
		})

	api.GET("/status", func(c *relay.Context) {
		version, _ := c.GetString("version")
		c.JSON(200, map[string]interface{}{
			"status":  "ok",
			"version": version,
		})
	})

	app.UseError(func(err error, c *relay.Context) {
		c.JSON(500, map[string]interface{}{
			"error": err.Error(),
		})
	})

	fmt.Println("Server starting on :8080")
	err := app.Listen(":8080", nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
