## goxpress Design 文档

goxpress 是一个用 Go 开发的 Web 框架，目标是成为“Go 版 Express”，既保留 Express 的开发体验，又充分发挥 Go 的高性能和可预测性。

---

### 1. 核心定位与设计目标

* **Go 化 Express**：外观、用法与 Express 保持高度一致，前端或 Node.js 开发者切换时几乎零学习成本。
* **大道至简 (Simplicity from Go)**：接口最小化、无魔法、拥抱标准库 `net/http`，所有依赖与上下文显式传递。
* **灵活组合 (Flexibility from Express)**：所有功能通过可插拔中间件实现，支持参数化路由、路由分组、本地中间件。
* **高性能**：底层使用 Radix Tree 路由查找，最小化每次请求的内存分配。
* **易于测试**：对 `http.Handler` 的兼容性使得可使用 `httptest` 轻松编写单元测试与集成测试。
* **快速上手**：仅需几行代码即可启动 "Hello, World" 服务。
* **生产力**：丰富的 Context 方法减少模板代码，让开发者专注业务。

---

### 2. 核心哲学

1. **最小化接口 & 无魔法**

   * 所有核心组件（Engine、Router、Context）仅提供正交且可预测的功能。
   * 避免全局变量与隐式链条，全靠显式注册与调用。
2. **拥抱 `net/http`**

   * `Engine` 实现 `http.Handler`，可以与标准库或第三方中间件无缝集成。
   * 用户也可将标准 `http.Handler` 包装为 goxpress 中间件。
3. **中间件驱动**

   * 框架灵魂在于中间件（Middleware）。所有功能（日志、认证、CORS、压缩）均通过插件形式提供。
4. **路由自由**

   * 参数化路由、通配符、路由分组、子路由模块化。
5. **无绑定**

   * 不强制使用模板引擎、ORM 或特定项目结构，自由组合。

---

### 3. 整体架构 (High-Level)

```
Request → Server (Engine) → Router → [Middleware Chain] → Handler → Response
```

1. **Engine (引擎/应用)**

   * `app := goxpress.New()`，实现 `http.Handler`。
   * 管理全局中间件、路由实例。
   * 提供：`Use()`, `UseError()`, `GET/POST/...`, `Route()`, `Mount()`, `Listen()`。
2. **Router (路由器)**

   * 每个 HTTP 方法对应一棵 Radix Tree，`O(k)` 路径查找。
   * 支持静态段、动态段（`:id`）与通配符（`*path`）。
   * 提供路由分组与挂载（`Mount`）。
3. **Middleware (中间件)**

   * `type HandlerFunc func(*Context)`；中间件和最终处理器共用此签名。
   * 支持链式调用：中间件通过 `c.Next(err?)` 或 `c.Abort()` 控制流程。
4. **Context (上下文)**

   * 封装 `http.ResponseWriter`、`*http.Request`、标准 `context.Context`。
   * 附带 URL 参数、存储区（`Set/Get`）、中间件链管理。
   * 提供 `Param()`, `Query()`, `BindJSON()`, `JSON()`, `String()`, `Status()` 等便捷方法。

---

### 4. Express 式 API 设计

| Express                      | goxpress                                    |
| ---------------------------- | ---------------------------------------- |
| `const app = express()`      | `app := goxpress.New()`                     |
| `app.use(mw)`                | `app.Use(mw...)`                         |
| `app.get(path, h...)`        | `app.GET(path, h...)`                    |
| `app.post(path, h...)`       | `app.POST(path, h...)`                   |
| `app.listen(3000, cb)`       | `app.Listen(":3000", cb?)`               |
| `app.route("/x").get()`      | `app.Route("/x").GET(...).POST(...)`     |
| `const r = express.Router()` | `r := goxpress.NewRouter()`                 |
| `app.use('/api', r)`         | `app.Mount("/api", r)`                   |
| `next(err)`                  | `c.Next(err)` + `app.UseError(errMw...)` |

* **链式调用**：所有方法返回 `*Engine` 或 `*Router`，可连续写。

---

### 5. 内置中间件与插件

* **BodyParser**：`app.Use(goxpress.BodyParser())`，支持 JSON 与 URL‑encoded。
* **CookieParser**：`app.Use(goxpress.CookieParser())`。
* **CORS**：`app.Use(goxpress.CORS())`。
* **Static**：`app.Static("/assets", "./public")`（同 Express `express.static`）。
* **Logger & Recover**：默认请求日志与 Panic 捕获。
* **Error Handler**：`app.UseError(handler)`，捕获 `c.Next(err)` 或 Panic，统一返回。

---

### 6. 路由策略与优先级

1. **多树 vs 单树**：按方法分多棵 Radix Tree。
2. **静态 > 动态 > 通配**：匹配优先级，首匹配即走。
3. **嵌套路由**：`Mount(prefix, Router)` 会自动拼路径并继承中间件。
4. **参数冲突处理**：在插入与查找时区分静态段与动态段节点。
5. **通配符与正则扩展**：节点设计预留扩展点，支持 `*filepath` 及未来正则。

---

### 7. 错误流程与控制流

* **错误中间件**：签名 `func(error, *Context)`，在 `c.Next(err)` 或 Panic 后触发。
* **Abort / IsAborted**：`c.Abort()` 立即停止后续常规中间件。
* **Recover**：框架层面自动捕获 Panic 并转换为错误中间件调用。
* **超时/取消**：与标准 `context.Context` 对接，支持请求超时与手动取消。

---

### 8. 测试与生产准备

* **单元测试**：使用 `httptest.NewRecorder()` + `app.ServeHTTP()` 断言状态码和响应体。
* **Graceful Shutdown**：对外暴露底层 `http.Server` 或提供 `app.Shutdown()`，优雅停服。
* **TLS/HTTP2 支持**：`app.ListenTLS()` 或用户传入自定义 `http.Server`。
* **Metrics**：可选 Prometheus 插件，自动记录请求时长、状态分布等。

---

### 9. 开发者体验示例

```go
app := goxpress.New().
    Use(goxpress.Logger()).
    Use(goxpress.Recover()).
    Static("/public", "./public").

api := goxpress.NewRouter().
    Use(auth).
    GET("/users/:id", getUser)

app.Mount("/api", api).
    GET("/", func(c *goxpress.Context) {
        c.String(200, "Hello, World!")
    }).
    Listen(":8080", nil)
```