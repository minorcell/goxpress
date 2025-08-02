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

## 教程

### 1. 基础 HTTP 服务

#### 简单路由

路由这块儿，我们基本上完全照搬了 Express 的风格，所以如果你用过 Express，这里应该毫无学习成本：

```go
package main

import "github.com/minorcell/goxpress"

func main() {
    app := goxpress.New()

    // 各种 HTTP 方法，想怎么用就怎么用
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
        c.JSON(400, map[string]string{"error": "JSON 格式不对哦"})
        return
    }

    c.JSON(201, map[string]interface{}{
        "message": "用户创建成功",
        "user":    user,
    })
}

func updateUser(c *goxpress.Context) {
    id := c.Param("id")
    c.JSON(200, map[string]string{
        "message": "用户 " + id + " 已更新",
    })
}

func deleteUser(c *goxpress.Context) {
    id := c.Param("id")
    c.JSON(200, map[string]string{
        "message": "用户 " + id + " 已删除",
    })
}
```

#### 参数和查询字符串

获取路径参数和查询参数也很直观，基本上一看就懂：

```go
app.GET("/users/:id", func(c *goxpress.Context) {
    // 路径参数
    userID := c.Param("id")

    // 查询参数
    page := c.Query("page")
    limit := c.Query("limit")

    c.JSON(200, map[string]string{
        "user_id": userID,
        "page":    page,
        "limit":   limit,
    })
})

// GET /users/123?page=1&limit=10
// 返回: {"user_id": "123", "page": "1", "limit": "10"}
```

### 2. 中间件

中间件这块儿，我们提供了一些常用的内置中间件，当然你也可以很容易地写自己的。

#### 内置中间件

```go
package main

import "github.com/minorcell/goxpress"

func main() {
    app := goxpress.New()

    // 内置中间件，开箱即用
    app.Use(goxpress.Logger())   // 请求日志
    app.Use(goxpress.Recover())  // Panic 恢复，避免程序崩溃

    app.GET("/", func(c *goxpress.Context) {
        c.String(200, "Hello with middleware!")
    })

    app.Listen(":8080", nil)
}
```

#### 自定义中间件

写个自定义中间件其实挺简单的，就是一个返回 `HandlerFunc` 的函数：

```go
// 认证中间件
func AuthMiddleware() goxpress.HandlerFunc {
    return func(c *goxpress.Context) {
        token := c.Request.Header.Get("Authorization")

        if token == "" {
            c.JSON(401, map[string]string{"error": "哎呀，忘记带 token 了"})
            c.Abort() // 停止后续处理
            return
        }

        // 验证 token（这里简化了，实际项目中你可能需要 JWT 或其他方式）
        if token != "Bearer valid-token" {
            c.JSON(401, map[string]string{"error": "token 不对哦"})
            c.Abort()
            return
        }

        // 在上下文中存储用户信息，后面的处理器可以用到
        c.Set("user_id", "12345")
        c.Next() // 继续到下一个中间件/处理器
    }
}

// CORS 中间件
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

    // 全局中间件，对所有路由生效
    app.Use(goxpress.Logger())
    app.Use(CORSMiddleware())

    // 受保护的路由组
    protected := app.Route("/api")
    protected.Use(AuthMiddleware()) // 只对这个组里的路由生效

    protected.GET("/profile", func(c *goxpress.Context) {
        userID, _ := c.GetString("user_id")
        c.JSON(200, map[string]string{
            "user_id": userID,
            "profile": "这里是用户个人资料数据",
        })
    })

    app.Listen(":8080", nil)
}
```

### 3. 上下文和请求处理

Context 是 goxpress 的核心，所有的请求处理都围绕它展开。

#### 请求数据

获取各种请求数据都很方便：

```go
app.POST("/users", func(c *goxpress.Context) {
    // JSON 数据绑定
    var user struct {
        Name  string `json:"name"`
        Email string `json:"email"`
        Age   int    `json:"age"`
    }

    if err := c.BindJSON(&user); err != nil {
        c.JSON(400, map[string]string{"error": "JSON 格式有问题"})
        return
    }

    // 表单数据
    name := c.PostForm("name")
    email := c.PostForm("email")

    // 文件上传
    file, err := c.FormFile("avatar")
    if err == nil {
        // 处理文件上传...
        c.SaveUploadedFile(file, "./uploads/"+file.Filename)
    }

    c.JSON(200, map[string]interface{}{
        "message": "数据接收成功",
        "user":    user,
    })
})
```

#### 响应类型

支持各种响应格式，想返回什么就返回什么：

```go
app.GET("/api/data", func(c *goxpress.Context) {
    // JSON 响应
    c.JSON(200, map[string]string{"message": "JSON 数据"})
})

app.GET("/text", func(c *goxpress.Context) {
    // 纯文本响应
    c.String(200, "这是一段文本")
})

app.GET("/html", func(c *goxpress.Context) {
    // HTML 响应
    c.HTML(200, "<h1>Hello HTML</h1>")
})

app.GET("/redirect", func(c *goxpress.Context) {
    // 重定向
    c.Redirect(302, "https://github.com/minorcell/goxpress")
})
```

### 4. 路由组和组织

当项目变大的时候，路由组织就变得很重要了。我们提供了很灵活的路由分组功能。

#### 基础路由组

```go
func main() {
    app := goxpress.New()

    // API v1 组
    v1 := app.Route("/api/v1")
    v1.GET("/users", func(c *goxpress.Context) {
        c.JSON(200, map[string]string{"version": "v1", "users": "用户列表"})
    })
    v1.GET("/posts", func(c *goxpress.Context) {
        c.JSON(200, map[string]string{"version": "v1", "posts": "文章列表"})
    })

    // API v2 组
    v2 := app.Route("/api/v2")
    v2.GET("/users", func(c *goxpress.Context) {
        c.JSON(200, map[string]string{"version": "v2", "users": "用户列表（新版）"})
    })

    app.Listen(":8080", nil)
}
```

#### 带中间件的嵌套组

```go
func main() {
    app := goxpress.New()

    // 全局中间件
    app.Use(goxpress.Logger())

    // API 组，有自己的中间件
    api := app.Route("/api")
    api.Use(CORSMiddleware())

    // 公开的 API
    public := api.Route("/public")
    public.GET("/health", func(c *goxpress.Context) {
        c.JSON(200, map[string]string{"status": "OK"})
    })

    // 需要认证的 API
    protected := api.Route("/protected")
    protected.Use(AdminMiddleware())
    protected.GET("/admin", func(c *goxpress.Context) {
        c.JSON(200, map[string]string{"message": "管理员专用接口"})
    })
    protected.DELETE("/users/:id", func(c *goxpress.Context) {
        id := c.Param("id")
        c.JSON(200, map[string]string{"message": "用户 " + id + " 已删除"})
    })

    app.Listen(":8080", nil)
}

func AdminMiddleware() goxpress.HandlerFunc {
    return func(c *goxpress.Context) {
        // 这里应该检查用户是否是管理员
        // 为了演示，我们简化处理
        role := c.Request.Header.Get("User-Role")
        if role != "admin" {
            c.JSON(403, map[string]string{"error": "需要管理员权限"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### 5. 错误处理

错误处理这块儿，我们提供了全局错误处理器，让你可以统一处理各种错误。

#### 全局错误处理器

```go
func main() {
    app := goxpress.New()

    // 设置全局错误处理器
    app.SetErrorHandler(func(c *goxpress.Context, err error) {
        // 记录错误日志
        println("发生错误:", err.Error())

        // 根据错误类型返回不同的响应
        if err.Error() == "unauthorized" {
            c.JSON(401, map[string]string{"error": "未授权"})
        } else {
            c.JSON(500, map[string]string{"error": "服务器内部错误"})
        }
    })

    app.GET("/error", func(c *goxpress.Context) {
        // 触发一个错误
        panic("这是一个测试错误")
    })

    app.GET("/auth-error", func(c *goxpress.Context) {
        // 返回一个认证错误
        c.Error(fmt.Errorf("unauthorized"))
    })

    app.Listen(":8080", nil)
}
```

### 6. 完整的 REST API 示例

来个完整的例子，展示一个标准的 REST API 应该怎么写：

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

    // 中间件
    app.Use(goxpress.Logger())
    app.Use(CORSMiddleware())

    // API 路由
    api := app.Route("/api")
    api.GET("/users", listUsers)           // 获取用户列表
    api.GET("/users/:id", getUser)         // 获取单个用户
    api.POST("/users", createUser)         // 创建用户
    api.PUT("/users/:id", updateUser)      // 更新用户
    api.DELETE("/users/:id", deleteUser)   // 删除用户

    app.Listen(":8080", func() {
        println("API 服务器运行在 http://localhost:8080")
    })
}

func listUsers(c *goxpress.Context) {
    c.JSON(200, map[string]interface{}{
        "users": users,
        "total": len(users),
    })
}

func getUser(c *goxpress.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(400, map[string]string{"error": "ID 格式不对"})
        return
    }

    for _, user := range users {
        if user.ID == id {
            c.JSON(200, user)
            return
        }
    }

    c.JSON(404, map[string]string{"error": "用户不存在"})
}

func createUser(c *goxpress.Context) {
    var newUser User
    if err := c.BindJSON(&newUser); err != nil {
        c.JSON(400, map[string]string{"error": "请求数据格式错误"})
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
        c.JSON(400, map[string]string{"error": "ID 格式不对"})
        return
    }

    var updatedUser User
    if err := c.BindJSON(&updatedUser); err != nil {
        c.JSON(400, map[string]string{"error": "请求数据格式错误"})
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

    c.JSON(404, map[string]string{"error": "用户不存在"})
}

func deleteUser(c *goxpress.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(400, map[string]string{"error": "ID 格式不对"})
        return
    }

    for i, user := range users {
        if user.ID == id {
            users = append(users[:i], users[i+1:]...)
            c.JSON(200, map[string]string{"message": "用户删除成功"})
            return
        }
    }

    c.JSON(404, map[string]string{"error": "用户不存在"})
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

## 高级用法

### 自定义中间件

写个计时中间件来监控请求处理时间：

```go
func TimingMiddleware() goxpress.HandlerFunc {
    return func(c *goxpress.Context) {
        start := time.Now()

        c.Next()

        duration := time.Since(start)
        println("请求处理时间:", duration.String())
    }
}
```

### 路由模式

支持各种路由模式：

- `/users/:id` - 单个参数
- `/files/*filepath` - 通配符匹配
- `/api/v:version/users` - 自定义参数名

### 测试

测试你的 API 也很简单：

```go
func TestAPI(t *testing.T) {
    app := goxpress.New()
    app.GET("/test", func(c *goxpress.Context) {
        c.JSON(200, map[string]string{"message": "test"})
    })

    // 这里你可以用任何 HTTP 测试工具
    // 比如 httptest 包
}
```

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
