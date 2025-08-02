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
- `StatusCode() int` - Get current status code

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
