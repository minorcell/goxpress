# goxpress

一个快速、直观的 Go Web 框架，灵感来自 Express.js。专为开发者生产力而设计，同时提供出色的性能。

[![Go Report Card](https://goreportcard.com/badge/github.com/minorcell/goxpress)](https://goreportcard.com/report/github.com/minorcell/goxpress)
[![GoDoc](https://godoc.org/github.com/minorcell/goxpress?status.svg)](https://godoc.org/github.com/minorcell/goxpress)
[![Coverage](https://img.shields.io/badge/coverage-90.3%25-brightgreen)](https://github.com/minorcell/goxpress)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

[English](README.md) | 中文

## 前言

事实上，Go 生态中已有许多优秀且成熟的 Web 框架，如 Gin、Fiber、Echo 等，它们拥有广泛的社区支持和丰富的生态资源。那么，为什么我们还需要另一个框架 —— goxpress？这是否只是又一次“造轮子”？

我认为这个问题的答案可以从以下几个核心设计理念中找到：

- **坚持 Go 的“少即是多”哲学，提供极简 API**
  goxpress 提供零配置、快速上手的开发体验。只需三行代码，即可启动一个基本的 Web 服务。简洁的 API 不仅意味着更低的学习成本，也保留了高度的灵活性，让开发者能够更自由地掌控业务逻辑，而不是被过度封装所束缚。

- **借鉴 Express 的 API 风格，降低学习曲线**
  选择与 Express 接近的编程风格有两个主要考虑：一是对前端开发者（尤其是熟悉 Node.js 的开发者）更友好，几乎可以“无痛”迁移到 Go；二是 Express 本身是一款久经验证的框架，其设计理念值得借鉴。站在巨人的肩膀上，可以走得更稳、更远。

- **发挥 Go 的语言优势**
  goxpress 并非简单的“Go 版 Express”，而是结合了 Go 的语言特性加以优化。例如原生的高并发能力、强大的类型系统、稳定且功能强大的标准库支持，这些都是 Node.js 生态所难以比拟的基础能力。

## 特性

- 🚀 **Express.js 风格的 API** - 熟悉的语法，上手即用
- ⚡ **高性能** - 100 万+ QPS，路由匹配极速
- 🛡️ **类型安全** - 完整的 Go 类型支持，IDE 友好
- 🔧 **中间件支持** - 强大的中间件生态，错误处理完善
- 🗂️ **路由分组** - 优雅的 API 组织方式
- 📦 **零依赖** - 仅基于 Go 标准库构建
- 🧪 **测试完善** - 90.3% 的测试覆盖率

## 快速开始

### 安装

说实话，安装过程再简单不过了：

```bash
go mod init your-project
go get github.com/minorcell/goxpress
```

### Hello World

三行代码，一个完整的 Web 服务就起来了：

```go
package main

import "github.com/minorcell/goxpress"

func main() {
    app := goxpress.New()

    app.GET("/", func(c *goxpress.Context) {
        c.String(200, "Hello, World!")
    })

    app.Listen(":8080", func() {
        println("服务器运行在 http://localhost:8080")
    })
}
```

### 链式调用

如果你喜欢链式的写法（很多人都喜欢），那也没问题：

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

## 性能

### 基准测试

说到性能，我们还是很有信心的。在我们的基准测试中：

- **吞吐量**: 超过 100 万 QPS（在一个 8 核 CPU 的机器上）
- **内存使用**: 每个请求约 2.5KB 内存分配
- **延迟**: P99 延迟低于 1ms

### 但框架的性能真的重要吗？

这是个有趣的问题。说实话，对于大多数应用来说，框架的性能可能不是瓶颈。数据库查询、网络 I/O、业务逻辑的复杂度往往比框架本身的开销大得多。

不过，高性能的框架确实有几个好处：

- **更低的资源消耗** - 意味着更低的云服务器成本
- **更好的用户体验** - 响应时间更快
- **更高的并发处理能力** - 可以支持更多的用户

所以，性能虽然不是唯一考虑因素，但有总比没有好，对吧？

## API 参考

### 核心类型

```go
type HandlerFunc func(*Context)
type ErrorHandlerFunc func(*Context, error)
```

### Engine 方法

- `New() *Engine` - 创建新的引擎实例
- `GET(path, handler)` - 注册 GET 路由
- `POST(path, handler)` - 注册 POST 路由
- `PUT(path, handler)` - 注册 PUT 路由
- `DELETE(path, handler)` - 注册 DELETE 路由
- `PATCH(path, handler)` - 注册 PATCH 路由
- `HEAD(path, handler)` - 注册 HEAD 路由
- `OPTIONS(path, handler)` - 注册 OPTIONS 路由
- `Use(handlers...)` - 注册全局中间件
- `Route(prefix)` - 创建路由组
- `Listen(addr, callback)` - 启动服务器
- `SetErrorHandler(handler)` - 设置全局错误处理器

### Context 方法

- `Param(key) string` - 获取路径参数
- `Query(key) string` - 获取查询参数
- `PostForm(key) string` - 获取表单数据
- `BindJSON(obj)` - 绑定 JSON 数据
- `JSON(code, obj)` - 返回 JSON 响应
- `String(code, text)` - 返回文本响应
- `HTML(code, html)` - 返回 HTML 响应
- `Status(code)` - 设置状态码
- `Redirect(code, url)` - 重定向
- `Set(key, value)` - 在上下文中存储数据
- `GetString(key)` - 从上下文中获取字符串
- `Next()` - 调用下一个中间件
- `Abort()` - 中止请求处理
- `Error(err)` - 触发错误处理
- `StatusCode() int` - 获取当前状态码

## 贡献

欢迎大家贡献代码！如果你有任何想法或者发现了 bug，随时提 issue 或者 pull request。

在提交代码之前，请确保：

- 代码通过了所有测试
- 新功能有对应的测试用例
- 遵循 Go 的代码规范

## 许可证

MIT License. 详见 [LICENSE](LICENSE) 文件。

## 致谢

感谢所有为这个项目做出贡献的开发者，以及 Express.js 团队的优秀设计理念。
