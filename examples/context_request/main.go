package main

import (
	"fmt"
	"github.com/minorcell/goxpress"
)

// This example demonstrates how to handle request data in goxpress
// It shows how to bind JSON data, access form data, and handle file uploads
func main() {
	// Step 1: Create a new instance of goxpress.Engine
	app := goxpress.New()

	// Step 2: Define routes for handling different types of request data
	app.POST("/users", func(c *goxpress.Context) {
		// 1. JSON data binding with enhanced error handling
		var user struct {
			Name  string `json:"name"`
			Email string `json:"email"`
			Age   int    `json:"age"`
		}

		if err := c.BindJSON(&user); err != nil {
			c.JSON(400, map[string]string{
				"error":   "Invalid JSON format",
				"details": err.Error(),
			})
			return
		}

		// 2. Form data extraction
		name := c.PostForm("name")
		email := c.PostForm("email")

		// 3. File upload handling with more detailed information
		file, err := c.FormFile("avatar")
		var fileInfo map[string]interface{}
		if err == nil {
			fileInfo = map[string]interface{}{
				"filename":  file.Filename,
				"size":      file.Size,
				"header":    file.Header,
			}
			
			// In a real application, you would save the file:
			// if err := c.SaveUploadedFile(file, "./uploads/"+file.Filename); err != nil {
			// 	c.JSON(500, map[string]string{"error": "Failed to save file"})
			// 	return
			// }
		} else {
			fileInfo = nil
		}

		// 4. Send comprehensive response
		c.JSON(200, map[string]interface{}{
			"message": "Data received successfully",
			"user": map[string]interface{}{
				"name":  user.Name,
				"email": user.Email,
				"age":   user.Age,
			},
			"form": map[string]string{
				"name":  name,
				"email": email,
			},
			"file": fileInfo,
		})
	})

	// Step 3: Start the server
	app.Listen(":8080", func() {
		fmt.Println("Context request example running at http://localhost:8080")
	})
}