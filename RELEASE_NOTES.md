# Release Notes - goxpress v1.0.0

🎉 **Initial Release**

We're excited to announce the first stable release of goxpress - a fast, intuitive web framework for Go inspired by Express.js!

## 🌟 What is goxpress?

goxpress brings the familiar Express.js development experience to Go while maintaining excellent performance and full type safety. It's designed for developer productivity without sacrificing the performance and reliability that Go is known for.

## ✨ Key Features

### 🚀 **Express.js-like API**
- Familiar routing: `app.GET()`, `app.POST()`, etc.
- Intuitive middleware system
- Chainable method calls
- Route groups and nested routing

### ⚡ **High Performance**
- **1.8M+ requests/sec** for simple routes
- **1.2M+ requests/sec** for JSON responses
- Efficient Radix Tree routing algorithm
- Context object pooling for memory efficiency

### 🛡️ **Type Safe & Go Native**
- Full Go type safety
- Zero external dependencies
- Built on Go standard library
- Excellent IDE support

### 🔧 **Production Ready**
- 90.3% test coverage
- Comprehensive benchmarks
- Panic recovery middleware
- Request logging
- Error handling system

## 📦 Installation

```bash
go get github.com/minorcell/goxpress
```

## 🚀 Quick Start

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

## 🎯 Core Components

### HTTP Routing
- **Static routes**: `/users`
- **Parameters**: `/users/:id`
- **Wildcards**: `/files/*filepath`
- **Route groups**: `/api/v1`

### Middleware System
```go
app.Use(goxpress.Logger())
app.Use(goxpress.Recover())
app.Use(func(c *goxpress.Context) {
    // Custom middleware
    c.Next()
})
```

### Request Handling
```go
app.POST("/users", func(c *goxpress.Context) {
    var user User
    if err := c.BindJSON(&user); err != nil {
        c.JSON(400, map[string]string{"error": "Invalid JSON"})
        return
    }
    c.JSON(201, user)
})
```

### Error Handling
```go
app.UseError(func(err error, c *goxpress.Context) {
    c.JSON(500, map[string]string{"error": err.Error()})
})
```

## 📊 Performance Benchmarks

| Operation | Performance | Memory |
|-----------|-------------|---------|
| Simple requests | 1.8M req/sec | 1129 B/op |
| JSON responses | 1.2M req/sec | 1505 B/op |
| Route parameters | 1.0M req/sec | 1817 B/op |
| Static routing | 7.4M lookups/sec | 120 B/op |
| Parameter extraction | 136M ops/sec | 0 B/op |

## 🔥 Why Choose goxpress?

### The Real Performance Story
90% of web application bottlenecks come from:
- Database queries
- External API calls  
- Business logic complexity
- Infrastructure limitations

**Not the web framework choice.**

goxpress gives you excellent performance out of the box, so you can focus on building great features instead of micro-optimizing framework overhead.

### Developer Experience
- **Familiar API** for Express.js developers
- **Fast development** with intuitive patterns
- **Type safety** catches errors at compile time
- **Zero learning curve** for Go developers

### Production Confidence
- **Comprehensive tests** with 90.3% coverage
- **Performance validated** with extensive benchmarks
- **Memory efficient** with object pooling
- **Error resilient** with panic recovery

## 📚 Documentation

- **English**: [README.md](README.md)
- **中文**: [README_zh.md](README_zh.md)
- **Performance**: [PERFORMANCE.md](PERFORMANCE.md)
- **API Reference**: Available on [pkg.go.dev](https://pkg.go.dev/github.com/minorcell/goxpress)

## 🛠️ Built-in Features

### Middleware
- ✅ Request logging
- ✅ Panic recovery
- ✅ Custom middleware support
- ✅ Error handling chain

### Routing
- ✅ All HTTP methods (GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS)
- ✅ Route parameters and wildcards
- ✅ Route groups and nesting
- ✅ Middleware per route group

### Request/Response
- ✅ JSON binding and responses
- ✅ Query parameter access
- ✅ URL parameter extraction
- ✅ Custom response headers
- ✅ Status code control

## 🎁 What's Included

```
goxpress/
├── 📦 Core Framework
│   ├── goxpress.go      # Main engine
│   ├── router.go        # HTTP routing
│   ├── context.go       # Request context
│   └── middleware.go    # Built-in middleware
├── 📖 Documentation
│   ├── README.md        # English docs
│   ├── README_zh.md     # Chinese docs
│   └── PERFORMANCE.md   # Performance analysis
├── 🧪 Tests & Benchmarks
│   ├── *_test.go        # Unit tests
│   └── benchmark_*.go   # Performance tests
└── 📄 Project Files
    ├── go.mod           # Module definition
    └── LICENSE          # MIT License
```

## 🔮 Future Roadmap

While v1.0.0 covers all essential web framework features, future versions may include:

- Static file serving middleware
- Built-in CORS middleware  
- Cookie and session support
- File upload utilities
- WebSocket support
- Template engine adapters

## 🤝 Contributing

We welcome contributions! goxpress is built with ❤️ for the Go community.

1. Fork the repository
2. Create your feature branch
3. Add tests for new functionality
4. Ensure all tests pass
5. Submit a pull request

## 📄 License

goxpress is released under the MIT License. See [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

- Inspired by [Express.js](https://expressjs.com/) for its elegant API design
- Built with the excellent Go standard library
- Thanks to the Go community for feedback and inspiration

## 🚀 Get Started Today!

```bash
go get github.com/minorcell/goxpress
```

Join thousands of developers building fast, reliable web applications with goxpress!

---

**Happy coding with goxpress!** 🎉

*For questions, issues, or feature requests, please visit our [GitHub repository](https://github.com/minorcell/goxpress).*