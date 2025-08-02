# goxpress

A fast, intuitive Go web framework inspired by Express.js. Built for developer productivity with excellent performance.

[![Go Report Card](https://goreportcard.com/badge/github.com/minorcell/goxpress)](https://goreportcard.com/report/github.com/minorcell/goxpress)
[![GoDoc](https://godoc.org/github.com/minorcell/goxpress?status.svg)](https://godoc.org/github.com/minorcell/goxpress)
[![Coverage](https://img.shields.io/badge/coverage-90.3%25-brightgreen)](https://github.com/minorcell/goxpress)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

English | [中文](README_zh.md)

## Why goxpress?

In fact, the Go ecosystem already has many excellent and mature web frameworks such as Gin, Fiber, Echo, and others, all of which enjoy broad community support and rich ecosystem resources. So why do we still need another framework — goxpress? Is it just another case of “reinventing the wheel”?

I believe the answer to this question lies in the following core design principles:

- ** Adhering to Go’s “less is more” philosophy by providing a minimal API **
  goxpress offers a zero-configuration, quick-start development experience. You can launch a basic web service with just three lines of code. A simple API not only lowers the learning curve but also retains high flexibility, allowing developers to freely control business logic instead of being constrained by over-encapsulation.

- ** Borrowing Express’s API style to reduce the learning curve **
  Choosing a programming style similar to Express is based on two main considerations: first, it is more friendly to front-end developers, especially those familiar with Node.js, enabling an almost painless transition to Go; second, Express itself is a battle-tested framework whose design philosophy is worth learning from. Standing on the shoulders of giants lets us move more steadily and further.

- ** Leveraging Go’s language strengths **
  goxpress is not simply a “Go version of Express,” but rather an optimization that takes advantage of Go’s language features. For example, native high concurrency capabilities, a strong type system, and a stable, powerful standard library — all foundational strengths that the Node.js ecosystem struggles to match.

## Features

- 🚀 **Express.js-style API** - Familiar syntax, ready to use
- ⚡ **High Performance** - 1M+ QPS with lightning-fast routing
- 🛡️ **Type Safe** - Full Go type support, IDE friendly
- 🔧 **Middleware Support** - Powerful middleware ecosystem with comprehensive error handling
- 🗂️ **Route Groups** - Elegant API organization
- 📦 **Zero Dependencies** - Built only on Go standard library
- 🧪 **Well Tested** - 90.3% test coverage

## Quick Start

### Installation

Honestly, the installation process couldn't be simpler:

```bash
go mod init your-project
go get github.com/minorcell/goxpress
```

### Hello World

Three lines of code, and you have a complete web service up and running:

```go
package main

import "github.com/minorcell/goxpress"

func main() {
    app := goxpress.New()

    app.GET("/", func(c *goxpress.Context) {
        c.String(200, "Hello, World!")
    })

    app.Listen(":8080", func() {
        println("Server running on http://localhost:8080")
    })
}
```

### Method Chaining

If you prefer the chaining style (many people do), no problem:

```go
package main

import "github.com/minorcell/goxpress"

func main() {
    goxpress.New().
        Use(goxpress.Logger(), goxpress.Recover()).
        GET("/", func(c *goxpress.Context) {
            c.String(200, "Hello World")
        }).
        Listen(":8080", func() {
            println("Server started on port 8080")
        })
}
```

## Tutorial

### 1. Basic HTTP Server

#### Simple Routes

For routing, we basically copied Express's style completely, so if you've used Express, there should be zero learning curve here:

```go
package main

import "github.com/minorcell/goxpress"

func main() {
    app := goxpress.New()

    // Various HTTP methods, use them however you want
    app.GET("/users", getUsers)
    app.POST("/users", createUser)
    app.PUT("/users/:id", updateUser)
    app.DELETE("/users/:id", deleteUser)

    app.Listen(":8080", nil)
}

func getUsers(c *goxpress.Context) {
    c.JSON(200, map[string]interface{}{
        "users": []string{"Alice", "Bob", "Charlie"},
    })
}

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

func updateUser(c *goxpress.Context) {
    id := c.Param("id")
    c.JSON(200, map[string]string{
        "message": "User " + id + " updated",
    })
}

func deleteUser(c *goxpress.Context) {
    id := c.Param("id")
    c.JSON(200, map[string]string{
        "message": "User " + id + " deleted",
    })
}
```

#### Parameters and Query Strings

Getting path parameters and query parameters is also intuitive, pretty much self-explanatory:

```go
app.GET("/users/:id", func(c *goxpress.Context) {
    // Path parameter
    userID := c.Param("id")

    // Query parameters
    page := c.Query("page")
    limit := c.Query("limit")

    c.JSON(200, map[string]string{
        "user_id": userID,
        "page":    page,
        "limit":   limit,
    })
})

// GET /users/123?page=1&limit=10
// Returns: {"user_id": "123", "page": "1", "limit": "10"}
```

### 2. Middleware

For middleware, we provide some commonly used built-in ones, and of course you can easily write your own.

#### Built-in Middleware

```go
package main

import "github.com/minorcell/goxpress"

func main() {
    app := goxpress.New()

    // Built-in middleware, ready to use
    app.Use(goxpress.Logger())   // Request logging
    app.Use(goxpress.Recover())  // Panic recovery, prevents crashes

    app.GET("/", func(c *goxpress.Context) {
        c.String(200, "Hello with middleware!")
    })

    app.Listen(":8080", nil)
}
```

#### Custom Middleware

Writing custom middleware is actually pretty simple - just a function that returns a `HandlerFunc`:

```go
// Authentication middleware
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

// CORS middleware
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

func main() {
    app := goxpress.New()

    // Global middleware, applies to all routes
    app.Use(goxpress.Logger())
    app.Use(CORSMiddleware())

    // Protected route group
    protected := app.Route("/api")
    protected.Use(AuthMiddleware()) // Only applies to routes in this group

    protected.GET("/profile", func(c *goxpress.Context) {
        userID, _ := c.GetString("user_id")
        c.JSON(200, map[string]string{
            "user_id": userID,
            "profile": "Here's the user profile data",
        })
    })

    app.Listen(":8080", nil)
}
```

### 3. Context and Request Handling

Context is the core of goxpress - all request handling revolves around it.

#### Request Data

Getting various request data is convenient:

```go
app.POST("/users", func(c *goxpress.Context) {
    // JSON data binding
    var user struct {
        Name  string `json:"name"`
        Email string `json:"email"`
        Age   int    `json:"age"`
    }

    if err := c.BindJSON(&user); err != nil {
        c.JSON(400, map[string]string{"error": "JSON format has issues"})
        return
    }

    // Form data
    name := c.PostForm("name")
    email := c.PostForm("email")

    // File upload
    file, err := c.FormFile("avatar")
    if err == nil {
        // Handle file upload...
        c.SaveUploadedFile(file, "./uploads/"+file.Filename)
    }

    c.JSON(200, map[string]interface{}{
        "message": "Data received successfully",
        "user":    user,
    })
})
```

#### Response Types

Support various response formats - return whatever you want:

```go
app.GET("/api/data", func(c *goxpress.Context) {
    // JSON response
    c.JSON(200, map[string]string{"message": "JSON data"})
})

app.GET("/text", func(c *goxpress.Context) {
    // Plain text response
    c.String(200, "This is some text")
})

app.GET("/html", func(c *goxpress.Context) {
    // HTML response
    c.HTML(200, "<h1>Hello HTML</h1>")
})

app.GET("/redirect", func(c *goxpress.Context) {
    // Redirect
    c.Redirect(302, "https://github.com/minorcell/goxpress")
})
```

### 4. Route Groups and Organization

When projects get bigger, route organization becomes important. We provide very flexible route grouping functionality.

#### Basic Route Groups

```go
func main() {
    app := goxpress.New()

    // API v1 group
    v1 := app.Route("/api/v1")
    v1.GET("/users", func(c *goxpress.Context) {
        c.JSON(200, map[string]string{"version": "v1", "users": "User list"})
    })
    v1.GET("/posts", func(c *goxpress.Context) {
        c.JSON(200, map[string]string{"version": "v1", "posts": "Post list"})
    })

    // API v2 group
    v2 := app.Route("/api/v2")
    v2.GET("/users", func(c *goxpress.Context) {
        c.JSON(200, map[string]string{"version": "v2", "users": "User list (new version)"})
    })

    app.Listen(":8080", nil)
}
```

#### Nested Groups with Middleware

```go
func main() {
    app := goxpress.New()

    // Global middleware
    app.Use(goxpress.Logger())

    // API group with its own middleware
    api := app.Route("/api")
    api.Use(CORSMiddleware())

    // Public APIs
    public := api.Route("/public")
    public.GET("/health", func(c *goxpress.Context) {
        c.JSON(200, map[string]string{"status": "OK"})
    })

    // Protected APIs requiring authentication
    protected := api.Route("/protected")
    protected.Use(AdminMiddleware())
    protected.GET("/admin", func(c *goxpress.Context) {
        c.JSON(200, map[string]string{"message": "Admin-only endpoint"})
    })
    protected.DELETE("/users/:id", func(c *goxpress.Context) {
        id := c.Param("id")
        c.JSON(200, map[string]string{"message": "User " + id + " deleted"})
    })

    app.Listen(":8080", nil)
}

func AdminMiddleware() goxpress.HandlerFunc {
    return func(c *goxpress.Context) {
        // Here we should check if user is admin
        // For demo purposes, we simplify
        role := c.Request.Header.Get("User-Role")
        if role != "admin" {
            c.JSON(403, map[string]string{"error": "Admin privileges required"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### 5. Error Handling

For error handling, we provide a global error handler so you can handle various errors uniformly.

#### Global Error Handler

```go
func main() {
    app := goxpress.New()

    // Set global error handler
    app.SetErrorHandler(func(c *goxpress.Context, err error) {
        // Log the error
        println("Error occurred:", err.Error())

        // Return different responses based on error type
        if err.Error() == "unauthorized" {
            c.JSON(401, map[string]string{"error": "Unauthorized"})
        } else {
            c.JSON(500, map[string]string{"error": "Internal server error"})
        }
    })

    app.GET("/error", func(c *goxpress.Context) {
        // Trigger an error
        panic("This is a test error")
    })

    app.GET("/auth-error", func(c *goxpress.Context) {
        // Return an authentication error
        c.Error(fmt.Errorf("unauthorized"))
    })

    app.Listen(":8080", nil)
}
```

### 6. Complete REST API Example

Here's a complete example showing how a standard REST API should be written:

```go
package main

import (
    "strconv"
    "github.com/minorcell/goxpress"
)

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

var users = []User{
    {ID: 1, Name: "Alice", Email: "alice@example.com"},
    {ID: 2, Name: "Bob", Email: "bob@example.com"},
}
var nextID = 3

func main() {
    app := goxpress.New()

    // Middleware
    app.Use(goxpress.Logger())
    app.Use(CORSMiddleware())

    // API routes
    api := app.Route("/api")
    api.GET("/users", listUsers)           // Get user list
    api.GET("/users/:id", getUser)         // Get single user
    api.POST("/users", createUser)         // Create user
    api.PUT("/users/:id", updateUser)      // Update user
    api.DELETE("/users/:id", deleteUser)   // Delete user

    app.Listen(":8080", func() {
        println("API server running on http://localhost:8080")
    })
}

func listUsers(c *goxpress.Context) {
    c.JSON(200, map[string]interface{}{"users": users})
}

func getUser(c *goxpress.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(400, map[string]string{"error": "Wrong ID format"})
        return
    }

    for _, user := range users {
        if user.ID == id {
            c.JSON(200, user)
            return
        }
    }

    c.JSON(404, map[string]string{"error": "User not found"})
}

func createUser(c *goxpress.Context) {
    var newUser User
    if err := c.BindJSON(&newUser); err != nil {
        c.JSON(400, map[string]string{"error": "Request data format error"})
        return
    }

    newUser.ID = nextID
    nextID++
    users = append(users, newUser)

    c.JSON(201, newUser)
}

func updateUser(c *goxpress.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(400, map[string]string{"error": "Wrong ID format"})
        return
    }

    var updatedUser User
    if err := c.BindJSON(&updatedUser); err != nil {
        c.JSON(400, map[string]string{"error": "Request data format error"})
        return
    }

    for i, user := range users {
        if user.ID == id {
            updatedUser.ID = id
            users[i] = updatedUser
            c.JSON(200, updatedUser)
            return
        }
    }

    c.JSON(404, map[string]string{"error": "User not found"})
}

func deleteUser(c *goxpress.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(400, map[string]string{"error": "Wrong ID format"})
        return
    }

    for i, user := range users {
        if user.ID == id {
            users = append(users[:i], users[i+1:]...)
            c.JSON(200, map[string]string{"message": "User deleted successfully"})
            return
        }
    }

    c.JSON(404, map[string]string{"error": "User not found"})
}

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
```

## Performance

### Benchmarks

When it comes to performance, we're pretty confident. In our benchmarks:

- **Throughput**: Over 1M QPS (on an 8-core CPU machine)
- **Memory usage**: About 2.5KB memory allocation per request
- **Latency**: P99 latency under 1ms

### But does framework performance really matter?

This is an interesting question. Honestly, for most applications, framework performance might not be the bottleneck. Database queries, network I/O, and business logic complexity often have much more overhead than the framework itself.

However, high-performance frameworks do have several benefits:

- **Lower resource consumption** - means lower cloud server costs
- **Better user experience** - faster response times
- **Higher concurrency handling** - can support more users

So while performance isn't the only consideration, it's better to have it than not, right?

## API Reference

### Core Types

```go
type HandlerFunc func(*Context)
type ErrorHandlerFunc func(*Context, error)
```

### Engine Methods

- `New() *Engine` - Create new engine instance
- `GET(path, handler)` - Register GET route
- `POST(path, handler)` - Register POST route
- `PUT(path, handler)` - Register PUT route
- `DELETE(path, handler)` - Register DELETE route
- `PATCH(path, handler)` - Register PATCH route
- `HEAD(path, handler)` - Register HEAD route
- `OPTIONS(path, handler)` - Register OPTIONS route
- `Use(handlers...)` - Register global middleware
- `Route(prefix)` - Create route group
- `Listen(addr, callback)` - Start server
- `SetErrorHandler(handler)` - Set global error handler

### Context Methods

- `Param(key) string` - Get path parameter
- `Query(key) string` - Get query parameter
- `PostForm(key) string` - Get form data
- `BindJSON(obj)` - Bind JSON data
- `JSON(code, obj)` - Return JSON response
- `String(code, text)` - Return text response
- `HTML(code, html)` - Return HTML response
- `Status(code)` - Set status code
- `Redirect(code, url)` - Redirect
- `Set(key, value)` - Store data in context
- `GetString(key)` - Get string from context
- `Next()` - Call next middleware
- `Abort()` - Abort request processing
- `Error(err)` - Trigger error handling

## Advanced Usage

### Custom Middleware

Write a timing middleware to monitor request processing time:

```go
func TimingMiddleware() goxpress.HandlerFunc {
    return func(c *goxpress.Context) {
        start := time.Now()

        c.Next()

        duration := time.Since(start)
        println("Request processing time:", duration.String())
    }
}
```

### Route Patterns

Support various route patterns:

- `/users/:id` - Single parameter
- `/files/*filepath` - Wildcard matching
- `/api/v:version/users` - Custom parameter names

### Testing

Testing your API is also simple:

```go
func TestAPI(t *testing.T) {
    app := goxpress.New()
    app.GET("/test", func(c *goxpress.Context) {
        c.JSON(200, map[string]string{"message": "test"})
    })

    // Here you can use any HTTP testing tool
    // like the httptest package
}
```

## Contributing

Everyone's welcome to contribute code! If you have any ideas or find bugs, feel free to create issues or pull requests.

Before submitting code, please ensure:

- Code passes all tests
- New features have corresponding test cases
- Follow Go coding standards

## License

MIT License. See [LICENSE](LICENSE) file for details.

## Acknowledgments

Thanks to all developers who contributed to this project, and the Express.js team for their excellent design philosophy.
