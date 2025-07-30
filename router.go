// Package relay 是一个类似 Express.js 的 Go Web 框架
// 定位：路由管理模块
// 作用：实现基于 Radix Tree 的高效路由匹配算法，管理路由和路由组
// 使用方法：
//  1. 通过 relay.NewRouter() 创建路由器实例
//  2. 使用 router.GET() 等方法注册路由
//  3. 使用 router.Group() 创建路由组
package relay

import (
	"strings"
	"sync"
)

// 定义一个用于分割路径的strings.Builder池
var builderPool = sync.Pool{
	New: func() interface{} {
		return &strings.Builder{}
	},
}

// Router 管理指定路径前缀的路由和中间件
// 定位：路由管理器
// 作用：
//  1. 管理路由注册和匹配
//  2. 管理路由组和中间件
//  3. 实现 Radix Tree 路由算法
//
// 使用方法：
//  1. 通过 NewRouter() 创建实例
//  2. 使用 Handle() 或 HTTP 方法函数注册路由
//  3. 使用 Group() 创建路由组
//  4. 使用 Use() 注册中间件
type Router struct {
	prefix      string                 // 路由前缀
	middlewares []HandlerFunc          // 路由器级别中间件
	engine      *Engine                // 关联的引擎实例
	subRouters  map[string]*Router     // 子路由器映射
	routes      map[string]*routerTree // 路由树映射，按 HTTP 方法分类
}

// routerTree 表示特定 HTTP 方法的基数树
// 定位：路由树结构
// 作用：为每个 HTTP 方法维护一棵独立的 Radix Tree
// 使用方法：框架内部使用，用于高效路由匹配
type routerTree struct {
	root *routerNode // 树的根节点
}

// routerNode 表示基数树中的节点
// 定位：路由树节点
// 作用：
//  1. 存储路由模式片段
//  2. 维护子节点关系
//  3. 标识是否为通配符节点
//  4. 存储处理函数
//
// 使用方法：框架内部使用，构建和搜索路由树
type routerNode struct {
	pattern  string        // 完整路由模式，例如 /p/:lang/doc
	part     string        // 路由片段，例如 :lang
	children []*routerNode // 子节点列表
	isWild   bool          // 是否为通配符匹配，part 含有 : 或 * 时为 true
	handlers []HandlerFunc // 处理函数列表
}

// NewRouter 创建一个新的 Router 实例
// 定位：Router 结构体的构造函数
// 作用：初始化 Router 实例及其依赖组件
// 使用方法：
//
//	router := NewRouter()
//	通常由 Engine 自动创建，也可以手动创建用于路由组
func NewRouter() *Router {
	return &Router{
		subRouters: make(map[string]*Router),
		routes:     make(map[string]*routerTree),
	}
}

// Use 为路由器注册中间件
// 定位：路由器级别中间件注册方法
// 作用：为当前路由器及其子路由注册中间件
// 使用方法：
//
//	router.Use(loggingMiddleware, authMiddleware)
//	返回 *Router 实例，支持链式调用
func (r *Router) Use(middleware ...HandlerFunc) *Router {
	r.middlewares = append(r.middlewares, middleware...)
	return r
}

// Group 创建具有指定前缀的新路由组
// 定位：路由组创建方法
// 作用：创建具有共同前缀和中间件的新路由组
// 使用方法：
//
//	api := router.Group("/api")
//	api.Use(apiMiddleware)
//	api.GET("/users", getUsers)
//	返回新的 *Router 实例
func (r *Router) Group(prefix string) *Router {
	router := &Router{
		prefix:      r.prefix + prefix,
		middlewares: make([]HandlerFunc, len(r.middlewares)), // 复制父路由器的中间件
		engine:      r.engine,
		subRouters:  make(map[string]*Router),
		routes:      r.routes,
	}

	// 拷贝父路由器的中间件
	copy(router.middlewares, r.middlewares)

	return router
}

// Handle 注册具有指定方法、模式和处理函数的新路由
// 定位：通用路由注册方法
// 作用：为指定 HTTP 方法和路径注册处理函数
// 使用方法：
//
//	router.Handle("GET", "/users/:id", getUser)
func (r *Router) Handle(method, pattern string, handlers ...HandlerFunc) {
	// 合并路由器前缀和路由模式
	fullPattern := r.prefix + pattern
	if r.prefix != "" && pattern == "/" {
		fullPattern = r.prefix
	}

	// 准备处理器链：先添加路由器中间件，再添加路由处理器
	finalHandlers := make([]HandlerFunc, 0)
	finalHandlers = append(finalHandlers, r.middlewares...)
	finalHandlers = append(finalHandlers, handlers...)

	r.addRoute(method, fullPattern, finalHandlers)
}

// HTTP 方法注册函数
// 定位：HTTP 方法注册方法
// 作用：为指定 HTTP 方法和路径注册处理函数
// 使用方法：
//
//	router.GET("/path", handler)
//	返回 *Router 实例，支持链式调用
func (r *Router) GET(pattern string, handlers ...HandlerFunc) *Router {
	r.Handle("GET", pattern, handlers...)
	return r
}

func (r *Router) POST(pattern string, handlers ...HandlerFunc) *Router {
	r.Handle("POST", pattern, handlers...)
	return r
}

func (r *Router) PUT(pattern string, handlers ...HandlerFunc) *Router {
	r.Handle("PUT", pattern, handlers...)
	return r
}

func (r *Router) DELETE(pattern string, handlers ...HandlerFunc) *Router {
	r.Handle("DELETE", pattern, handlers...)
	return r
}

func (r *Router) PATCH(pattern string, handlers ...HandlerFunc) *Router {
	r.Handle("PATCH", pattern, handlers...)
	return r
}

func (r *Router) HEAD(pattern string, handlers ...HandlerFunc) *Router {
	r.Handle("HEAD", pattern, handlers...)
	return r
}

func (r *Router) OPTIONS(pattern string, handlers ...HandlerFunc) *Router {
	r.Handle("OPTIONS", pattern, handlers...)
	return r
}

// parsePattern 将路由模式解析为片段
// 定位：路由模式解析函数
// 作用：将路由模式字符串分割为片段数组
// 使用方法：
//
//	parts := parsePattern("/users/:id")
//	// 返回 ["users", ":id"]
func parsePattern(pattern string) []string {
	// 从池中获取builder
	builder := builderPool.Get().(*strings.Builder)
	defer func() {
		builder.Reset()
		builderPool.Put(builder)
	}()

	// 预估容量以减少重新分配
	parts := make([]string, 0, strings.Count(pattern, "/"))

	start := 0
	for i, char := range pattern {
		if char == '/' {
			if start < i {
				part := pattern[start:i]
				parts = append(parts, part)
				if part[0] == '*' {
					break
				}
			}
			start = i + 1
		}
	}

	// 处理最后一个片段
	if start < len(pattern) {
		part := pattern[start:]
		parts = append(parts, part)
	}

	return parts
}

// addRoute 将新路由添加到基数树
// 定位：路由添加方法
// 作用：将路由模式和处理函数插入到对应的 Radix Tree 中
// 使用方法：由 Handle() 方法内部调用
func (r *Router) addRoute(method, pattern string, handlers []HandlerFunc) {
	// 为方法创建树（如果不存在）
	if r.routes[method] == nil {
		r.routes[method] = &routerTree{root: &routerNode{}}
	}

	parts := parsePattern(pattern)

	// 将模式插入基数树
	r.routes[method].insertRoute(pattern, parts, 0, handlers)
}

// getRoute 查找与给定方法和路径匹配的路由
// 定位：路由查找方法
// 作用：在 Radix Tree 中查找匹配的路由节点
// 使用方法：由 Engine.ServeHTTP() 调用以匹配请求路由
func (r *Router) getRoute(method, path string) (*routerNode, map[string]string) {
	root, ok := r.routes[method]
	if !ok {
		return nil, nil
	}

	searchParts := parsePattern(path)
	params := make(map[string]string)

	node := root.searchRoute(searchParts, 0, params)

	if node != nil {
		return node, params
	}

	return nil, nil
}

// walkMountRoutes 递归遍历路由树，将路由添加到引擎中
// 定位：路由遍历方法
// 作用：递归遍历路由树，将每个节点的路由添加到引擎中
// 使用方法：框架内部使用，用于路由挂载
func (r *Router) walkMountRoutes(node *routerNode, method, mountPrefix string, groupMiddlewares []HandlerFunc, addRoute func(method, pattern string, handlers []HandlerFunc)) {
	// 如果是根节点，递归处理子节点
	if node.pattern == "" {
		for _, child := range node.children {
			r.walkMountRoutes(child, method, mountPrefix, groupMiddlewares, addRoute)
		}
		return
	}

	// 计算完整的路由模式
	pattern := strings.TrimSuffix(mountPrefix, "/") + strings.TrimPrefix(node.pattern, r.prefix)

	// 准备处理器链：先添加路由器中间件，再添加路由处理器
	finalHandlers := make([]HandlerFunc, 0)
	finalHandlers = append(finalHandlers, r.middlewares...)
	finalHandlers = append(finalHandlers, groupMiddlewares...)
	finalHandlers = append(finalHandlers, node.handlers...)

	// 添加路由
	addRoute(method, pattern, finalHandlers)

	// 递归处理子节点
	for _, child := range node.children {
		r.walkMountRoutes(child, method, mountPrefix, groupMiddlewares, addRoute)
	}
}

// insertRoute 将路由模式插入基数树
// 定位：路由插入方法
// 作用：递归地将路由模式插入到 Radix Tree 中
// 使用方法：由 addRoute() 方法内部调用
func (t *routerTree) insertRoute(pattern string, parts []string, height int, handlers []HandlerFunc) {
	// 基本情况：已处理完所有片段
	if len(parts) == height {
		t.root.pattern = pattern
		t.root.handlers = handlers
		return
	}

	part := parts[height]
	child := t.root.matchChild(part)

	if child == nil {
		// 创建新子节点
		child = &routerNode{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		t.root.children = append(t.root.children, child)
	}

	// 递归插入剩余模式
	tree := &routerTree{root: child}
	tree.insertRoute(pattern, parts, height+1, handlers)
}

// searchRoute 在基数树中搜索路由
// 定位：路由搜索方法
// 作用：递归地在 Radix Tree 中搜索匹配的路由节点
// 使用方法：由 getRoute() 方法内部调用
func (t *routerTree) searchRoute(parts []string, height int, params map[string]string) *routerNode {
	// 基本情况：已处理完所有片段或遇到通配符
	if len(parts) == height || strings.HasPrefix(t.root.part, "*") {
		if t.root.pattern == "" {
			return nil
		}
		return t.root
	}

	part := parts[height]

	// 检查所有子节点
	children := t.root.children
	for _, child := range children {
		if child.part == part || child.isWild {
			// 处理参数匹配
			if child.isWild && child.part[0] == ':' {
				params[child.part[1:]] = part
			} else if child.isWild && child.part[0] == '*' {
				// 对于通配符，捕获路径的其余部分
				params[child.part[1:]] = strings.Join(parts[height:], "/")
				return child
			}

			// 递归在子节点中搜索
			tree := &routerTree{root: child}
			result := tree.searchRoute(parts, height+1, params)
			if result != nil {
				return result
			}

			// 如有必要，回溯参数
			if child.isWild && child.part[0] == ':' {
				delete(params, child.part[1:])
			}
		}
	}

	return nil
}

// matchChild 查找与给定片段匹配的子节点
// 定位：子节点匹配方法
// 作用：在当前节点的子节点中查找匹配的节点
// 使用方法：由 insertRoute() 方法内部调用
func (n *routerNode) matchChild(part string) *routerNode {
	for _, child := range n.children {
		// 精确匹配或通配符匹配
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}
