package main

// import goxpress
import "github.com/minorcell/goxpress"

// This example demonstrates basic HTTP server functionality with goxpress
// It shows how to create routes for different HTTP methods and handle requests
func main() {
	// Step 1: Create a new instance of goxpress.Engine
	app := goxpress.New()

	// Step 2: Define routes for different HTTP methods
	// Various HTTP methods, use them however you want
	app.GET("/users", getUsers)
	app.POST("/users", createUser)
	app.PUT("/users/:id", updateUser)
	app.DELETE("/users/:id", deleteUser)

	// Step 3: Start the server
	app.Listen(":8080", func() {
		println("Basic HTTP server running at http://localhost:8080")
	})
}

// Handler for GET /users - returns a list of users
func getUsers(c *goxpress.Context) {
	c.JSON(200, map[string]interface{}{
		"users": []string{"Alice", "Bob", "Charlie"},
	})
}

// Handler for POST /users - creates a new user from JSON data
func createUser(c *goxpress.Context) {
	var user struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.BindJSON(&user); err != nil {
		c.JSON(400, map[string]string{"error": "Oops, JSON format is wrong"})
		return
	}

	c.JSON(201, map[string]interface{}{
		"message": "User created successfully",
		"user":    user,
	})
}

// Handler for PUT /users/:id - updates a user by ID
func updateUser(c *goxpress.Context) {
	id := c.Param("id")
	c.JSON(200, map[string]string{
		"message": "User " + id + " updated",
	})
}

// Handler for DELETE /users/:id - deletes a user by ID
func deleteUser(c *goxpress.Context) {
	id := c.Param("id")
	c.JSON(200, map[string]string{
		"message": "User " + id + " deleted",
	})
}