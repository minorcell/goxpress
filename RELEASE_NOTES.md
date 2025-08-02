# Release Notes - goxpress v1.1.0

## 🎉 New Features in v1.1.0

This release adds significant enhancements to request/response handling, logging flexibility, and developer experience while maintaining the framework's performance and simplicity.

### 🔧 **Advanced Logger Configuration**

- **Custom Formatters**: Define custom log formats with `Formatter` function
- **Path Filtering**: Skip logging for specific paths with `SkipPaths`
- **Flexible Output**: Direct logs to any `io.Writer`
- **Wildcard Matching**: Use patterns like `/api/*/health` for path filtering

```go
config := goxpress.LoggerConfig{
    SkipPaths: []string{"/health", "/metrics"},
    Output:    logFile,
    Formatter: func(c *goxpress.Context, start time.Time, duration time.Duration) string {
        return fmt.Sprintf("CUSTOM: %s %s\n", c.Request.Method, c.Request.URL.Path)
    },
}
app.Use(goxpress.LoggerWithConfig(config))
```

### 📝 **Enhanced Response Handling**

- **HTML Responses**: `c.HTML(200, "<h1>Hello World</h1>")`
- **HTTP Redirects**: `c.Redirect(302, "https://example.com")`
- **Static File Serving**: `c.File("./public/index.html")`

### 📥 **Improved Request Processing**

- **Form Data Extraction**: `c.PostForm("fieldname")`
- **File Upload Support**: `c.FormFile("avatar")` + `c.SaveUploadedFile()`
- **Enhanced JSON Binding**: Improved error handling in `c.BindJSON()`

### 📚 **Comprehensive Examples**

10+ new examples added with detailed documentation:
- [basic_http] - HTTP methods and parameters
- [context_request] - Form/data handling
- [context_response] - HTML/redirect/file responses
- [custom_middleware] - Custom middleware implementation
- [rest_api] - Complete REST API implementation

### 🧪 **Full Test Coverage**

- HTML response validation
- Redirect functionality testing
- Form data processing tests
- File upload and serving verification
- Custom logger formatter validation

## ✨ Key Features (Carried over from v1.0.0)

## 📦 Installation

```bash
go get github.com/minorcell/goxpress
```

## 🛠️ Built-in Features

### Middleware
- ✅ Request logging (custom formatters)
- ✅ Panic recovery
- ✅ Custom middleware support
- ✅ Error handling chain

### Request/Response
- ✅ JSON binding and responses
- ✅ Query parameter access
- ✅ URL parameter extraction
- ✅ Form data processing
- ✅ File upload handling
- ✅ HTML responses
- ✅ HTTP redirects (301, 302)
- ✅ Static file serving

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
├── 📚 Examples
│   ├── hello_world/     # Basic example
│   ├── basic_http/      # HTTP methods and parameters
│   ├── middleware/      # Built-in middleware usage
│   ├── custom_middleware/ # Custom middleware
│   ├── context_request/ # Request handling
│   ├── context_response/ # Response types
│   ├── route_groups/    # Route organization
│   ├── nested_groups/   # Complex routing
│   ├── error_handling/  # Error management
│   └── rest_api/        # Complete REST API
└── 📄 Project Files
    ├── go.mod           # Module definition
    └── LICENSE          # MIT License
```

## 🔮 Future Roadmap

Future versions may include:

- Built-in CORS middleware  
- Cookie and session support
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

MIT License - See [LICENSE](LICENSE) for details.

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