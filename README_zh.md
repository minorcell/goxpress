# goxpress

一个快速、直观的 Go Web 框架，灵感来自 Express.js。专为开发者生产力而设计，同时提供出色的性能。

[![Go Report Card](https://goreportcard.com/badge/github.com/minorcell/goxpress)](https://goreportcard.com/report/github.com/minorcell/goxpress)
[![GoDoc](https://godoc.org/github.com/minorcell/goxpress?status.svg)](https://godoc.org/github.com/minorcell/goxpress)
[![Coverage](https://img.shields.io/badge/coverage-88.3%25-brightgreen)](https://github.com/minorcell/goxpress)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## 特性

- 🚀 **类 Express.js API** - 对 Web 开发者友好且直观
- ⚡ **高性能** - 每秒处理超过 100 万请求，路由高效
- 🛡️ **类型安全** - 完整的 Go 类型安全支持，IDE 支持优秀
- 🔧 **中间件支持** - 强大的中间件系统，支持错误处理
- 🗂️ **路由组** - 使用嵌套路由组组织你的 API
- 📦 **零依赖** - 仅基于 Go 标准库构建
- 🧪 **充分测试** - 88.3% 测试覆盖率，全面的性能基准测试

## 快速开始

### 安装

```bash
go mod init your-project
go get github.com/minorcell/goxpress
```

### Hello World

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

## 教程

### 1. 基础 HTTP 服务

#### 简单路由

```go
package main

import "github.com/minorcell/goxpress"

func main() {
    app := goxpress.New()
    
    // 不同的 HTTP 方法
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
        c.JSON(400, map[string]string{"error": "无效的 JSON"})
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

#### 使用参数和查询字符串

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

#### 内置中间件

```go
package main

import "github.com/minorcell/goxpress"

func main() {
    app := goxpress.New()
    
    // 内置中间件
    app.Use(goxpress.Logger())   // 请求日志
    app.Use(goxpress.Recover())  // Panic 恢复
    
    app.GET("/", func(c *goxpress.Context) {
        c.String(200, "Hello with middleware!")
    })
    
    app.Listen(":8080", nil)
}
```

#### 自定义中间件

```go
// 认证中间件
func AuthMiddleware() goxpress.HandlerFunc {
    return func(c *goxpress.Context) {
        token := c.Request.Header.Get("Authorization")
        
        if token == "" {
            c.JSON(401, map[string]string{"error": "缺少 token"})
            c.Abort() // 停止后续处理
            return
        }
        
        // 验证 token（简化版）
        if token != "Bearer valid-token" {
            c.JSON(401, map[string]string{"error": "无效的 token"})
            c.Abort()
            return
        }
        
        // 在上下文中存储用户信息
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
    
    // 全局中间件
    app.Use(goxpress.Logger())
    app.Use(CORSMiddleware())
    
    // 受保护的路由
    protected := app.Route("/api")
    protected.Use(AuthMiddleware()) // 应用到此组中的所有路由
    
    protected.GET("/profile", func(c *goxpress.Context) {
        userID, _ := c.GetString("user_id")
        c.JSON(200, map[string]string{
            "user_id": userID,
            "profile": "用户个人资料数据",
        })
    })
    
    app.Listen(":8080", nil)
}
```

### 3. 上下文和请求处理

#### 请求数据

```go
app.POST("/submit", func(c *goxpress.Context) {
    // JSON 主体解析
    var data struct {
        Name    string `json:"name"`
        Email   string `json:"email"`
        Age     int    `json:"age"`
    }
    
    if err := c.BindJSON(&data); err != nil {
        c.JSON(400, map[string]string{"error": "无效的 JSON 格式"})
        return
    }
    
    // 验证
    if data.Name == "" || data.Email == "" {
        c.JSON(400, map[string]string{"error": "姓名和邮箱是必需的"})
        return
    }
    
    // 路径和查询参数
    category := c.Param("category")
    source := c.Query("source")
    
    // 为其他中间件存储在上下文中
    c.Set("validated_data", data)
    
    c.JSON(200, map[string]interface{}{
        "message":  "数据接收成功",
        "data":     data,
        "category": category,
        "source":   source,
    })
})
```

#### 响应类型

```go
app.GET("/examples", func(c *goxpress.Context) {
    // 字符串响应
    c.String(200, "纯文本响应")
})

app.GET("/json", func(c *goxpress.Context) {
    // JSON 响应
    c.JSON(200, map[string]interface{}{
        "message": "成功",
        "data":    []int{1, 2, 3},
        "meta": map[string]string{
            "version": "1.0",
        },
    })
})

app.GET("/custom", func(c *goxpress.Context) {
    // 自定义标头和状态
    c.Response.Header().Set("X-Custom-Header", "value")
    c.Status(201)
    c.JSON(201, map[string]string{"created": "true"})
})
```

### 4. 路由组和组织

#### 基础路由组

```go
func main() {
    app := goxpress.New()
    
    // API v1 路由
    v1 := app.Route("/api/v1")
    {
        v1.GET("/users", listUsers)
        v1.POST("/users", createUser)
        v1.GET("/users/:id", getUser)
        v1.PUT("/users/:id", updateUser)
        v1.DELETE("/users/:id", deleteUser)
    }
    
    // API v2 路由，不同的实现
    v2 := app.Route("/api/v2")
    {
        v2.GET("/users", listUsersV2)
        v2.POST("/users", createUserV2)
    }
    
    app.Listen(":8080", nil)
}
```

#### 带中间件的嵌套组

```go
func main() {
    app := goxpress.New()
    
    // 全局中间件
    app.Use(goxpress.Logger())
    app.Use(goxpress.Recover())
    
    // 公共 API（无需认证）
    public := app.Route("/api/public")
    public.GET("/health", healthCheck)
    public.POST("/register", registerUser)
    public.POST("/login", loginUser)
    
    // 受保护的 API（需要认证）
    api := app.Route("/api")
    api.Use(AuthMiddleware())
    
    // 用户管理
    users := api.Group("/users")
    users.GET("/", listUsers)
    users.GET("/:id", getUser)
    users.PUT("/:id", updateUser)
    users.DELETE("/:id", deleteUser)
    
    // 仅管理员路由
    admin := api.Group("/admin")
    admin.Use(AdminMiddleware()) // 额外的管理员检查
    admin.GET("/stats", getStats)
    admin.DELETE("/users/:id", adminDeleteUser)
    
    app.Listen(":8080", nil)
}

func AdminMiddleware() goxpress.HandlerFunc {
    return func(c *goxpress.Context) {
        userID, _ := c.GetString("user_id")
        
        // 检查用户是否为管理员（简化版）
        if !isAdmin(userID) {
            c.JSON(403, map[string]string{"error": "需要管理员权限"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

### 5. 错误处理

#### 全局错误处理器

```go
func main() {
    app := goxpress.New()
    
    // 全局错误处理器
    app.UseError(func(err error, c *goxpress.Context) {
        // 记录错误
        fmt.Printf("错误: %v\n", err)
        
        // 返回适当的响应
        c.JSON(500, map[string]string{
            "error":   "内部服务器错误",
            "message": "出了点问题",
        })
    })
    
    app.Use(goxpress.Recover()) // 将 panic 转换为错误
    
    app.GET("/error", func(c *goxpress.Context) {
        // 这将触发错误处理器
        c.Next(fmt.Errorf("出了点问题"))
    })
    
    app.GET("/panic", func(c *goxpress.Context) {
        // 这将被 Recover 中间件捕获
        panic("故意 panic")
    })
    
    app.Listen(":8080", nil)
}
```

### 6. 完整的 REST API 示例

```go
package main

import (
    "fmt"
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
    app.Use(goxpress.Recover())
    app.Use(CORSMiddleware())
    
    // API 路由
    api := app.Route("/api")
    
    // 用户 CRUD
    api.GET("/users", listUsers)
    api.GET("/users/:id", getUser)
    api.POST("/users", createUser)
    api.PUT("/users/:id", updateUser)
    api.DELETE("/users/:id", deleteUser)
    
    app.Listen(":8080", func() {
        fmt.Println("🚀 服务器运行在 http://localhost:8080")
        fmt.Println("📖 试试: curl http://localhost:8080/api/users")
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
        c.JSON(400, map[string]string{"error": "无效的用户 ID"})
        return
    }
    
    for _, user := range users {
        if user.ID == id {
            c.JSON(200, user)
            return
        }
    }
    
    c.JSON(404, map[string]string{"error": "用户未找到"})
}

func createUser(c *goxpress.Context) {
    var newUser User
    if err := c.BindJSON(&newUser); err != nil {
        c.JSON(400, map[string]string{"error": "无效的 JSON"})
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
        c.JSON(400, map[string]string{"error": "无效的用户 ID"})
        return
    }
    
    var updatedUser User
    if err := c.BindJSON(&updatedUser); err != nil {
        c.JSON(400, map[string]string{"error": "无效的 JSON"})
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
    
    c.JSON(404, map[string]string{"error": "用户未找到"})
}

func deleteUser(c *goxpress.Context) {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(400, map[string]string{"error": "无效的用户 ID"})
        return
    }
    
    for i, user := range users {
        if user.ID == id {
            users = append(users[:i], users[i+1:]...)
            c.JSON(200, map[string]string{"message": "用户已删除"})
            return
        }
    }
    
    c.JSON(404, map[string]string{"error": "用户未找到"})
}

func CORSMiddleware() goxpress.HandlerFunc {
    return func(c *goxpress.Context) {
        c.Response.Header().Set("Access-Control-Allow-Origin", "*")
        c.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Response.Header().Set("Access-Control-Allow-Headers", "Content-Type")
        
        if c.Request.Method == "OPTIONS" {
            c.Status(204)
            return
        }
        
        c.Next()
    }
}
```

## 性能

goxpress 在保持开发者生产力的同时提供出色的性能：

### 基准测试

- **简单请求**: ~180万 请求/秒
- **JSON 响应**: ~120万 请求/秒  
- **路径参数**: ~100万 请求/秒
- **路由匹配**: 使用 Radix Tree 算法实现超快匹配

### 真实的性能故事

**90% 的情况下，你的 Web 服务性能不是由你选择的框架决定的。**

大多数应用的真正瓶颈是：

- **数据库查询** - 慢 SQL、缺少索引、N+1 查询
- **外部 API 调用** - 网络延迟、第三方服务限制  
- **业务逻辑** - 复杂计算、低效算法
- **基础设施** - 网络带宽、服务器资源、缓存

**即使切换到快 5 倍的框架，对整体响应时间的影响也微乎其微。**

将优化工作重点放在真正重要的地方：

1. **数据库优化** - 适当的索引、查询优化
2. **缓存策略** - Redis、内存缓存、CDN
3. **API 设计** - 分页、批量操作、高效的数据结构
4. **基础设施** - 负载均衡、适当的资源分配

goxpress 开箱即用就提供出色的性能，所以你可以专注于构建优秀的功能，而不是微优化框架开销。

## API 参考

### 核心类型

```go
type HandlerFunc func(*Context)
type ErrorHandlerFunc func(error, *Context)
```

### Engine 方法

```go
// HTTP 方法
app.GET(pattern string, handlers ...HandlerFunc) *Engine
app.POST(pattern string, handlers ...HandlerFunc) *Engine  
app.PUT(pattern string, handlers ...HandlerFunc) *Engine
app.DELETE(pattern string, handlers ...HandlerFunc) *Engine
app.PATCH(pattern string, handlers ...HandlerFunc) *Engine
app.HEAD(pattern string, handlers ...HandlerFunc) *Engine
app.OPTIONS(pattern string, handlers ...HandlerFunc) *Engine

// 中间件
app.Use(middleware ...HandlerFunc) *Engine
app.UseError(handlers ...ErrorHandlerFunc) *Engine

// 路由组
app.Route(prefix string) *Router

// 服务器
app.Listen(addr string, callback func()) error
app.ListenTLS(addr, certFile, keyFile string, callback func()) error
```

### Context 方法

```go
// 参数和查询
c.Param(key string) string
c.Query(key string) string

// 请求主体
c.BindJSON(obj interface{}) error

// 响应
c.Status(code int)
c.String(code int, format string, values ...interface{}) error
c.JSON(code int, obj interface{}) error

// 流程控制
c.Next(err ...error)
c.Abort()
c.IsAborted() bool

// 数据存储
c.Set(key string, value interface{})
c.Get(key string) (interface{}, bool)
c.GetString(key string) (string, bool)
c.MustGet(key string) interface{}
```

## 高级用法

### 自定义中间件

```go
func TimingMiddleware() goxpress.HandlerFunc {
    return func(c *goxpress.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start)
        c.Response.Header().Set("X-Response-Time", duration.String())
    }
}
```

### 路由模式

```go
// 静态路由
app.GET("/users", handler)

// 参数
app.GET("/users/:id", handler)           // /users/123
app.GET("/users/:id/posts/:postId", handler) // /users/123/posts/456

// 通配符  
app.GET("/files/*filepath", handler)     // /files/css/style.css
```

### 测试

```go
func TestAPI(t *testing.T) {
    app := goxpress.New()
    app.GET("/ping", func(c *goxpress.Context) {
        c.JSON(200, map[string]string{"message": "pong"})
    })
    
    req := httptest.NewRequest("GET", "/ping", nil)
    w := httptest.NewRecorder()
    
    app.ServeHTTP(w, req)
    
    assert.Equal(t, 200, w.Code)
}
```

## 贡献

我们欢迎贡献！请查看我们的[贡献指南](CONTRIBUTING.md)了解详情。

1. Fork 仓库
2. 创建你的特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交你的更改 (`git commit -m 'Add amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开一个 Pull Request

## 许可证

此项目根据 MIT 许可证授权 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 致谢

- 受 [Express.js](https://expressjs.com/) 启发，因其优雅的 API 设计
- 用 ❤️ 为 Go 社区构建
- 特别感谢所有贡献者

---

**使用 goxpress 愉快编码！** 🚀