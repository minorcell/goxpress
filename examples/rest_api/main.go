package main

// import goxpress
import (
	"strconv"
	"github.com/minorcell/goxpress"
)

// User represents a user entity in our system
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// In-memory storage for users
var users = []User{
	{ID: 1, Name: "Alice", Email: "alice@example.com"},
	{ID: 2, Name: "Bob", Email: "bob@example.com"},
}
var nextID = 3

// This example demonstrates a complete REST API implementation with goxpress
// It shows how to build a standard CRUD API for a resource (users in this case)
func main() {
	// Step 1: Create a new instance of goxpress.Engine
	app := goxpress.New()

	// Step 2: Register middleware
	// Logger middleware for request logging
	// CORS middleware for Cross-Origin Resource Sharing
	app.Use(goxpress.Logger())
	app.Use(CORSMiddleware())

	// Step 3: Create API routes
	api := app.Route("/api")
	api.GET("/users", listUsers)           // Get user list
	api.GET("/users/:id", getUser)         // Get single user
	api.POST("/users", createUser)         // Create user
	api.PUT("/users/:id", updateUser)      // Update user
	api.DELETE("/users/:id", deleteUser)   // Delete user

	// Step 4: Start the server
	app.Listen(":8080", func() {
		println("REST API example running at http://localhost:8080")
		println("API endpoints:")
		println("  GET    /api/users      - List all users")
		println("  GET    /api/users/:id  - Get a user by ID")
		println("  POST   /api/users      - Create a new user")
		println("  PUT    /api/users/:id  - Update a user by ID")
		println("  DELETE /api/users/:id  - Delete a user by ID")
	})
}

// listUsers returns a list of all users
func listUsers(c *goxpress.Context) {
	c.JSON(200, map[string]interface{}{"users": users})
}

// getUser returns a single user by ID
func getUser(c *goxpress.Context) {
	// Parse the user ID from the URL parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, map[string]string{"error": "Wrong ID format"})
		return
	}

	// Find the user with the specified ID
	for _, user := range users {
		if user.ID == id {
			c.JSON(200, user)
			return
		}
	}

	// User not found
	c.JSON(404, map[string]string{"error": "User not found"})
}

// createUser creates a new user from the request body
func createUser(c *goxpress.Context) {
	// Parse the user data from the request body
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(400, map[string]string{"error": "Request data format error"})
		return
	}

	// Assign an ID and save the user
	newUser.ID = nextID
	nextID++
	users = append(users, newUser)

	// Return the created user
	c.JSON(201, newUser)
}

// updateUser updates an existing user by ID
func updateUser(c *goxpress.Context) {
	// Parse the user ID from the URL parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, map[string]string{"error": "Wrong ID format"})
		return
	}

	// Parse the updated user data from the request body
	var updatedUser User
	if err := c.BindJSON(&updatedUser); err != nil {
		c.JSON(400, map[string]string{"error": "Request data format error"})
		return
	}

	// Find and update the user with the specified ID
	for i, user := range users {
		if user.ID == id {
			updatedUser.ID = id
			users[i] = updatedUser
			c.JSON(200, updatedUser)
			return
		}
	}

	// User not found
	c.JSON(404, map[string]string{"error": "User not found"})
}

// deleteUser removes a user by ID
func deleteUser(c *goxpress.Context) {
	// Parse the user ID from the URL parameter
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, map[string]string{"error": "Wrong ID format"})
		return
	}

	// Find and remove the user with the specified ID
	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			c.JSON(200, map[string]string{"message": "User deleted successfully"})
			return
		}
	}

	// User not found
	c.JSON(404, map[string]string{"error": "User not found"})
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