# Go 后端学习指南（以本项目为教材）

> 这份文档以 `server/` 的真实代码为教材，讲清楚：代码怎么流动、为什么这样分层、每个技术点的原理和坑、以及面试常被问到什么。
> 建议对照源码阅读，文中所有路径都相对于 `server/` 目录。

---

## 目录

1. [项目全景与技术栈](#1-项目全景与技术栈)
2. [一次请求的完整生命周期](#2-一次请求的完整生命周期)
3. [分层架构：每一层的职责与纪律](#3-分层架构每一层的职责与纪律)
4. [启动流程 main.go](#4-启动流程-maingo)
5. [配置管理：为什么全走环境变量](#5-配置管理为什么全走环境变量)
6. [Gin 核心概念](#6-gin-核心概念)
7. [中间件：洋葱模型](#7-中间件洋葱模型)
8. [JWT + Redis 双重鉴权（本项目最值得吃透的设计）](#8-jwt--redis-双重鉴权)
9. [密码安全：bcrypt](#9-密码安全bcrypt)
10. [GORM 数据库操作](#10-gorm-数据库操作)
11. [错误处理体系](#11-错误处理体系)
12. [统一响应格式](#12-统一响应格式)
13. [Go 泛型在本项目中的应用](#13-go-泛型在本项目中的应用)
14. [RBAC 权限模型与动态菜单](#14-rbac-权限模型与动态菜单)（含数据权限 dataScope、在线用户）
    - 14.5 [定时任务调度器（robfig/cron）](#145-定时任务调度器robfigcron)
15. [代码生成器的原理](#15-代码生成器的原理)
16. [本项目真实踩过的坑（重构实录）](#16-本项目真实踩过的坑重构实录)
17. [测试](#17-测试)
18. [面试题速查清单](#18-面试题速查清单)

---

## 1. 项目全景与技术栈

这是一个前后端分离的后台管理系统的**后端**：

| 组件 | 选型 | 作用 |
|---|---|---|
| Web 框架 | Gin | 路由、参数绑定、中间件 |
| ORM | GORM | 把 Go 结构体映射到 MySQL 表 |
| 缓存 | Redis (go-redis) | 存 JWT 的 jti，支持 token 主动失效 |
| 认证 | golang-jwt/v5 | 签发/校验 JWT |
| 日志 | zap | 高性能结构化日志 |
| 配置 | godotenv | 从 `.env` 加载环境变量 |
| 密码 | bcrypt | 密码哈希 |

**为什么是这套组合？** 这是 Go Web 开发的"标准班底"——Gin 是使用最广的 Go Web 框架，GORM 是使用最广的 ORM，zap 是 Uber 出品的高性能日志库。学会这套，看懂大多数 Go 业务项目没有障碍。

---

## 2. 一次请求的完整生命周期

以**登录**（`POST /api/v1/base/login`）为例，代码流动路径：

```
浏览器发 JSON { user_name, password }
  │
  ▼
main.go        router.NewRouter().Run(":8080")   ← Gin 监听端口
  │
  ▼
router/router.go        全局中间件依次执行：
  │                     Recovery（兜 panic）→ RequestLogger（记日志）
  │                     → CORS（跨域）→ installGuard（未安装拦截）
  ▼
router/system/base.go   base.POST("/login", systemApi.Login)  ← URL 映射到 handler
  │
  ▼
api/system/user.go      Login(c *gin.Context)：
  │                       1. ShouldBindJSON 把请求体绑定到 LoginRequest 结构体
  │                       2. 调 service 层做业务
  │                       3. 按错误类型组织 HTTP 响应
  ▼
service/system/user.go  Login(userName, password)：
  │                       1. GORM 按 username 查用户
  │                       2. 校验状态和 bcrypt 密码
  │                       3. utils.GenerateToken 签发双 token
  ▼
utils/jwt.go            签 JWT + 把 jti 写入 Redis
  │
  ▼
model/common/response   response.Success(c, tokens)  → 统一 JSON 格式返回
```

再看一个**带鉴权的增删改查**（`PUT /api/v1/system/dept/5`）：

```
请求头 Authorization: Bearer <access_token>
  │
  ▼
middleware/jwt.go     JWTAuth 中间件：解析 token → 查 Redis 确认 jti 有效
  │                   → c.Set("user_id", ...) 把身份放进请求上下文
  ▼
api/system/dept.go    UpdateDept：bindJSON[DeptPayload] 类型化绑定
  ▼
service/system/dept.go  UpdateDept → deptPayloadData 转成更新 map
  ▼
service/system/crud.go  updateWithLevel：维护树形层级 + GORM Updates
  ▼
api/system/common.go    successOrFail：成功返回数据 / 失败按错误类型收敛
```

**新手建议**：拿这两条链路对照源码走一遍，把每个文件真正打开看。理解了"一条请求怎么流"，这个项目就懂了一半。

---

## 3. 分层架构：每一层的职责与纪律

```
server/
├── main.go            程序入口：初始化 + 启动
├── config/            读环境变量 → 配置结构体（不含业务）
├── router/            URL → handler 的映射（不写逻辑）
├── middleware/        横切关注点：鉴权、日志、恢复、跨域
├── api/               HTTP handler：绑参数 + 调 service + 组织响应（薄层）
├── service/           业务逻辑：真正的规则、事务、查询（厚层）
├── model/             数据结构：GORM 模型 + 请求/响应 DTO
├── setup/             基础设施初始化：gorm / redis / logger
├── utils/             无业务依赖的工具：JWT 签发校验
└── database/          初始化 SQL
```

**每层的纪律（本项目 `.claude/rules/backend-api.md` 明文规定）**：

- **handler 禁止写 SQL 和业务判断**。它只做三件事：绑定参数、调用 service、组织响应。看 `api/system/dept.go`，每个函数不超过 10 行。
- **service 禁止直接写 HTTP 响应**。它返回 `(数据, error)`，让 handler 决定怎么回给客户端。
- **model 只描述数据结构**，不写业务方法。

**为什么要分层？**

1. **可测试**：service 不依赖 `*gin.Context`，可以直接用单元测试调用，不需要起 HTTP 服务。
2. **可替换**：明天想加 gRPC 接口，service 层原样复用，只需要新写一层薄的 gRPC handler。
3. **职责单一**：改业务规则只动 service，改接口格式只动 api，互不牵连。

**DTO 与 Model 分离的意义**：`model/system/user.go` 是数据库的样子（有密码哈希、删除时间），`model/system/request/user.go` 是前端能提交的样子（没有 id、没有时间字段）。分开之后，前端多传 `{"password": "xxx", "id": 999}` 这类字段也写不进不该写的地方——这防的是 **Mass Assignment（批量赋值）漏洞**。

> **面试常问**
> - "MVC/三层架构里各层职责是什么？为什么 controller 要薄？"
> - "DTO 和实体类为什么要分开？"（答：安全——防批量赋值；解耦——接口格式和表结构可独立演化）

---

## 4. 启动流程 main.go

```go
loadEnv()                        // 1. 加载 .env（生产环境可无此文件，直接用系统环境变量）
loggerInit.Logger.Initialize()   // 2. 日志最先初始化——之后的报错才有地方输出
defer loggerInit.Logger.Get().Sync()  // 3. 退出前刷掉缓冲区里的日志
initSetup()                      // 4. MySQL、Redis 连接
router.NewRouter().Run(addr)     // 5. 一切就绪才开始接请求
```

**几个值得注意的细节**：

- **初始化有顺序**：日志 → 数据库 → HTTP。如果先启动 HTTP 再连数据库，会有一小段时间接口全 500。
- **`defer Sync()`**：zap 为了性能会缓冲日志，程序退出前必须 flush，否则最后几条日志会丢。`defer` 保证函数返回前一定执行。
- **`Fatal` vs `panic`**：初始化失败用 `Fatal`（打日志后 `os.Exit(1)`），因为基础设施起不来服务就没有存在意义，这叫 **fail fast**。
- **未安装时跳过初始化**：`installService.Installed()` 检查锁文件。没安装时数据库配置还不存在，强行连接必然失败，所以放行安装向导、跳过其余初始化。

> **面试常问**
> - "defer 的执行顺序？"（后进先出，像栈）
> - "defer 常见用途？"（释放资源、解锁、recover、flush 日志）

---

## 5. 配置管理：为什么全走环境变量

看 `config/go_mysql.go` 的模式：

```go
func envOrDefault(key string, defaultValue string) string {
    value := os.Getenv(key)
    if value == "" {
        return defaultValue
    }
    return value
}
```

所有配置都是"环境变量优先，缺省给合理默认值"。**为什么不写在代码或配置文件里？**

1. **一份代码跑所有环境**：开发/测试/生产只是环境变量不同，二进制完全一样（这是 [12-Factor App](https://12factor.net/zh_cn/) 的核心原则之一）。
2. **密钥不进 git**：`JWT_SECRET`、数据库密码写在 `.env` 里，而 `.env` 在 `.gitignore` 中。密钥一旦提交进 git 历史，就应视为泄露。
3. **容器友好**：Docker/K8s 天然通过环境变量注入配置。

**DSN 是什么**：`user:pass@tcp(host:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local` 这串叫数据源名称。`parseTime=True` 很关键——没有它，MySQL 的 datetime 读出来是 `[]byte` 而不是 `time.Time`。

> **面试常问**
> - "配置怎么管理？密钥怎么处理？"（环境变量 + 密钥管理服务；绝不硬编码、绝不进 git）

---

## 6. Gin 核心概念

### 6.1 路由与路由组

```go
// router/router.go
PublicGroup := Router.Group(prefix)      // 公开组：登录、刷新 token
PrivateGroup := Router.Group(prefix)     // 私有组
PrivateGroup.Use(middleware.JWTAuth())   // 私有组统一挂鉴权中间件
```

路由组的价值：**公共前缀 + 公共中间件**只声明一次。新增一个需要登录的接口，挂到 `PrivateGroup` 下就自动有了鉴权，不可能忘。

### 6.2 gin.Context

`c *gin.Context` 是一次请求的"全能上下文"，贯穿整条中间件链：

| 用法 | 例子 |
|---|---|
| 读路径参数 | `c.Param("id")`（对应路由 `/user/:id`） |
| 读查询参数 | `c.Query("tree")`（对应 `?tree=true`） |
| 绑定请求体 | `c.ShouldBindJSON(&payload)` |
| 中间件传值 | `c.Set("user_id", 1)` / `c.Get("user_id")` |
| 写响应 | `c.JSON(200, obj)` |

### 6.3 参数绑定与校验

```go
type LoginRequest struct {
    UserName string `json:"user_name" binding:"required"`
    Password string `json:"password" binding:"required"`
}
```

`binding:"required"` 是声明式校验：字段缺失时 `ShouldBindJSON` 直接返回错误，业务代码里不用写一堆 `if xxx == ""`。本项目在 `api/system/common.go` 里用泛型把绑定收敛成一个函数：

```go
func bindJSON[T any](c *gin.Context) (T, bool) {
    var payload T
    if err := c.ShouldBindJSON(&payload); err != nil {
        response.Fail(c, code.ParamError, err.Error())
        return payload, false
    }
    return payload, true
}
```

每个 handler 从此只需要两行就完成"绑定 + 失败响应"。

**坑**：`ShouldBindJSON` 只能调用一次——请求体是流，读完就没了。

> **面试常问**
> - "Gin 的路由是怎么实现的？"（前缀树 Radix Tree，路径匹配 O(路径长度)）
> - "ShouldBind 和 MustBind(Bind) 的区别？"（Should 返回 error 自己处理；Bind 失败直接写 400 并 Abort）

---

## 7. 中间件：洋葱模型

```
请求 → Recovery → RequestLogger → CORS → installGuard → JWTAuth → handler
响应 ←──────────────────────────────────────────────────────────┘
```

中间件本质是 `func(c *gin.Context)`，调用 `c.Next()` 把控制权交给下一层，`c.Next()` 之后的代码在**响应阶段**执行——像洋葱一样一层层进、一层层出。

看 `middleware/logger.go` 是怎么利用这一点算耗时的：

```go
start := time.Now()
c.Next()                          // ← 这里面把整个 handler 跑完了
latency := time.Since(start)      // ← 回来之后才算得出耗时
```

### Recovery：为什么必须有

Go 里一个 goroutine panic 没人接住，**整个进程会崩**。一次空指针就把服务打挂是不可接受的，所以第一个中间件永远是 Recovery：

```go
defer func() {
    if err := recover(); err != nil {
        // 记堆栈日志 + 返回统一 500
    }
}()
c.Next()
```

**原则**：panic 只用于不可恢复的编程错误；业务问题一律用 `error` 返回值。项目规则明确写着"错误处理用 return err 而非 panic"。

### CORS：浏览器的同源策略

前端 `localhost:5173` 调后端 `localhost:8080` 属于跨域。浏览器发正式请求前会先发 `OPTIONS` **预检请求**，服务端必须应答允许的来源/方法/请求头，浏览器才放行。`middleware/cors.go` 做的就是这个应答。**记住：CORS 是浏览器行为，curl/Postman 不受影响**——所以"Postman 能通、浏览器不通"十有八九是 CORS。

### Abort vs Next

鉴权失败时用 `FailWithAbort`（内部调 `c.Abort()`）——**中断链条**，后续 handler 不再执行。如果只 return 不 Abort，Gin 会继续执行后面的 handler，等于鉴权白做。这是 Gin 新手最容易犯的错误之一。

> **面试常问**
> - "中间件的执行顺序？"（注册顺序进、逆序出）
> - "panic/recover 机制？defer 里 recover 为什么有效？"
> - "CORS 预检请求什么时候触发？"（非简单请求：自定义头、PUT/DELETE、Content-Type: application/json 等）

---

## 8. JWT + Redis 双重鉴权

这是本项目**最值得吃透**的设计，面试出镜率极高。核心文件：`utils/jwt.go`、`middleware/jwt.go`。

### 8.1 JWT 是什么

`Header.Payload.Signature` 三段 Base64。服务端用密钥对前两段做 HMAC-SHA256 签名。**任何人都能解码看到 Payload 内容（所以不能放敏感数据），但没有密钥就无法伪造签名**。

服务端不用存 session，靠签名就能验证身份——这叫**无状态认证**，天然适合多实例横向扩展。

### 8.2 纯 JWT 的致命问题：无法注销

JWT 一旦签发，到期前永远有效。用户点"退出登录"、管理员想踢人下线、token 被盗想作废——纯 JWT 都做不到。

**本项目的解法**：签发时给每个 token 一个随机唯一 ID（`jti`），存入 Redis 并设置和 token 相同的过期时间：

```
签发：Redis SET jwt:access:<jti> = userID， TTL = token 有效期
校验：JWT 签名合法 且 Redis 里 jti 存在 → 才算有效
注销：Redis DEL jwt:access:<jti>          → token 立即失效
```

这等于给无状态的 JWT 加回了一点"状态"，换来了**主动失效能力**。代价是每次请求多一次 Redis 查询（内存级速度，可接受）。

**安全铁律**（`rules/security.md`）：业务鉴权必须走 `ValidateToken`（含 Redis 校验），禁止只用 `ParseToken`（只验签名）——否则撤销机制形同虚设。

### 8.3 为什么要两个 token

| | access token | refresh token |
|---|---|---|
| 用途 | 每个接口都带 | 只用于换新 token |
| 有效期 | 短（默认 2 小时） | 长（默认 24 小时） |
| 暴露频率 | 高 | 极低 |

短命 access token 限制了泄露后的危害窗口；refresh token 只在刷新时用一次，暴露面小。前端拦截器收到 `40101`（access 过期）时自动拿 refresh 换新再重试，用户无感知。

`ValidateToken` 会校验 `token_type`，**refresh token 拿去调业务接口会直接失败**——防止长命 token 被当成通行证。

### 8.4 单端/多端登录

`JWT_LOGIN_MODE=single` 时，Redis 额外记 `jwt:user:<uid>:<type> = 当前jti`。新登录覆盖这个值，旧 token 的 jti 对不上就失效——**新设备登录踢掉旧设备**，一个 key 搞定。

### 8.5 算法安全细节

`ParseToken` 里有一段容易被忽视的代码：

```go
if _, ok := token.Method.(*gojwt.SigningMethodHMAC); !ok {
    return nil, fmt.Errorf("unexpected signing method")
}
```

这是在防**算法混淆攻击**：攻击者把 token 头里的 `alg` 改成 `none`（不签名）或从 RS256 换成 HS256 来绕过校验。服务端必须白名单校验算法，不能信任 token 自己声明的算法。

> **面试常问**
> - "JWT 和 Session 的区别？各自优缺点？"
> - "JWT 如何实现注销/踢人？"（黑名单或本项目的 jti 白名单）
> - "access/refresh 双 token 的意义？刷新流程怎么走？"
> - "JWT 的安全问题？"（alg=none 攻击、密钥弱、Payload 明文、XSS 偷 token）

---

## 9. 密码安全：bcrypt

```go
// 存：bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 验：bcrypt.CompareHashAndPassword(hash, password)
```

**为什么不用 MD5/SHA256？** 它们是为"快"设计的通用哈希——GPU 每秒能算几十亿次，配合彩虹表，弱密码秒破。bcrypt 的特点：

1. **故意慢**（cost 可调，默认 2^10 轮），暴力破解成本高出几个数量级；
2. **自带随机盐**，同一密码两次哈希结果不同，彩虹表失效；
3. 盐就编码在哈希串里，无需单独存储。

**配套设计**：登录失败统一返回"用户名或密码错误"（`ErrLoginFailed` 不区分用户不存在/密码错/账号禁用），防止攻击者**枚举有效用户名**。

> **面试常问**
> - "密码怎么存储？"（bcrypt/argon2，绝不 MD5，绝不明文）
> - "加盐的作用？"（防彩虹表、防相同密码撞哈希）

---

## 10. GORM 数据库操作

### 10.1 模型映射

```go
type AISystemUser struct {
    ID       uint    `json:"id" gorm:"column:id;primaryKey"`
    Username string  `json:"username" gorm:"column:username"`
    Nickname *string `json:"nickname" gorm:"column:nickname"`  // ← 注意是指针
    Roles    []uint  `json:"roles" gorm:"-"`                   // ← gorm:"-" 不映射到表
}
func (AISystemUser) TableName() string { return "ai_system_user" }
```

**为什么可空列要用指针（`*string`）？** 这是 Go + 数据库的经典问题：

- `string` 的零值是 `""`，**无法区分"数据库是 NULL"和"数据库是空字符串"**；
- GORM 把 NULL 扫进非指针字段会直接报错（本项目代码生成器曾踩过这个坑，见第 16 节）；
- 指针为 `nil` 表示 NULL，解引用前必须判空，否则 panic——看 `UserAuthList` 里的 `if user.Nickname != nil`。

同样的思路用在请求 DTO 上：`DeptPayload.Name *string` 为 nil 表示"前端没传这个字段，更新时跳过"，实现**部分更新**语义。

### 10.2 查询：参数化是底线

```go
db.Where("username = ?", userName)           // ✅ 参数化，安全
db.Where("username = '" + userName + "'")    // ❌ 字符串拼接 = SQL 注入
```

`?` 占位符让驱动对值做转义，用户输入永远不会被当成 SQL 执行。**唯一例外**是表名/列名无法参数化——本项目 `service/system/database.go` 对表名先用白名单正则 `^[A-Za-z0-9_]+$` 校验再拼接，`applyFilters` 的列名全部来自代码里写死的映射而非用户输入。

### 10.3 事务

```go
db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Delete(...); err != nil { return err }  // 返回 error → 自动回滚
    if err := tx.Create(...); err != nil { return err }
    return nil                                            // 返回 nil → 自动提交
})
```

什么时候必须用事务？**多条写操作要么全成功要么全失败**时。本项目的例子：给用户重绑角色（先删旧绑定再插新绑定，`BindUserRoles`）、批量保存配置（`BatchUpdateConfig`）。想象删除成功但插入失败——用户就没有任何角色了，这就是"中间状态"。

**坑**：事务回调里必须用 `tx` 而不是外面的 `db`，用错了那条语句就跑在事务外面。

### 10.4 软删除：本项目的特殊做法

GORM 自带软删除（字段类型为 `gorm.DeletedAt` 时 Delete 自动变 UPDATE）。但本项目的表继承自旧系统，用的是普通的 `delete_time` 列，所以**手动实现**：

- 查询：`softDelete(db)` 统一追加 `delete_time IS NULL`；
- 删除：真删用 `db.Delete`（`deleteByID`），软删用 `UPDATE delete_time = NOW()`（`SoftDeleteRecord`）。

**教训**（真实踩坑）：手动软删意味着每条查询都要记得加过滤；而对没有 `delete_time` 列的表加了过滤会直接 SQL 报错。代码生成器就曾在这里翻车，现在按"表里是否真的有这个列"决定。

### 10.5 分页与 N+1

分页三件套：`Count` 总数 → `Offset((page-1)*size)` → `Limit(size)`。

**N+1 问题**：`UserList` 里循环给每个用户查角色——1 次列表查询 + N 次角色查询。数据量小没事，量大时应该改成一次 `WHERE user_id IN (...)` 再在内存分组。源码注释里特意标了这个点，这是面试高频题。

### 10.6 连接池

```go
sqlDB.SetMaxOpenConns(100)   // 最多同时打开的连接
sqlDB.SetMaxIdleConns(10)    // 空闲池里保留的连接
sqlDB.SetConnMaxLifetime(time.Hour)  // 连接最长存活时间
```

为什么要 `ConnMaxLifetime`？MySQL 服务端有 `wait_timeout`，会单方面掐掉久置连接；客户端定期换新连接，避免拿到"半死"连接报错。

> **面试常问**
> - "怎么防 SQL 注入？"（参数化查询；标识符走白名单）
> - "什么是 N+1 查询？怎么解决？"（IN 批查 / JOIN / 预加载）
> - "软删除怎么实现？有什么代价？"（查询都要带过滤、唯一索引要把删除标记算进去）
> - "连接池参数怎么设？"

---

## 11. 错误处理体系

### 11.1 Go 的错误哲学

Go 没有 try/catch，错误是**普通返回值**：

```go
result, err := doSomething()
if err != nil {
    return nil, err   // 处理不了就往上传
}
```

好处是**每个可能失败的点都显式可见**，没有隐藏的跳转路径。

### 11.2 三种错误的分层处理（看 `api/system/common.go` 的 successOrFail）

本项目把错误分成三类，出口行为完全不同：

| 类型 | 例子 | 谁看 | 出口 |
|---|---|---|---|
| 参数错误 | JSON 绑定失败 | 前端开发者 | 40001 + 绑定错误详情 |
| 业务错误 `BizError` | "菜单下存在子菜单，无法删除" | 最终用户 | HTTP 200 + 40002 + 原文案 |
| 内部错误 | SQL 报错、Redis 挂了 | 只有运维/开发 | **记 zap 日志**，返回泛化"系统异常" |

**为什么内部错误不能透传 `err.Error()`？** SQL 报错文本里有表名、列名甚至语句本身——这是给攻击者的**信息泄露**。正确姿势：细节进日志（带 method/path 方便定位），客户端只拿到一句"系统异常"。

### 11.3 sentinel error 与 errors.Is / errors.As

```go
// 定义（service/system/user.go）
var ErrLoginFailed = errors.New("user_name or password is incorrect")

// 判断（api/system/user.go）
if errors.Is(err, systemService.ErrLoginFailed) { ... }
```

**为什么不用 `err.Error() == "xxx"` 比较字符串？** 文案一改代码全崩，而且包装过的错误比不出来。`errors.Is` 沿着包装链（`fmt.Errorf("%w", err)`）逐层比对哨兵值。

`errors.As` 则用于**按类型**匹配并取出错误对象：

```go
var bizErr *systemService.BizError
if errors.As(err, &bizErr) {
    response.Fail(c, code.OperationFailed, bizErr.Error())  // 拿到消息透传
}
```

规则：**固定语义用 sentinel（Is），需要携带数据用自定义类型（As）**。

> **面试常问**
> - "errors.Is / errors.As / errors.Unwrap 的区别？"
> - "error 和 panic 的使用边界？"
> - "怎么设计 API 的错误码体系？"

---

## 12. 统一响应格式

```json
{ "code": 0, "data": {...}, "msg": "操作成功" }
```

所有接口（成功或失败）都长这样（`model/common/response/response.go`）。**业务码和 HTTP 状态码是两套体系**：

- HTTP 状态码给**基础设施**看（网关、监控、浏览器）：401 触发跳登录、500 触发告警；
- 业务码给**前端逻辑**看：`40101` → 拦截器静默刷新 token 重试；`40002` → 弹出 msg。

业务错误（"菜单下有子菜单"）返回 HTTP 200 + code 40002：请求本身处理成功了，只是业务规则说"不行"——这不是服务器错误，不应该污染 5xx 监控指标。

---

## 13. Go 泛型在本项目中的应用

`service/system/crud.go` 是泛型的实战教材：

```go
func createRow[T any](table string, payload map[string]interface{}) (*T, error) {
    ...
    var result T                    // T 在编译期被替换成具体类型
    db.Table(table).Order("id DESC").First(&result)
    return &result, nil
}

// 使用：一行代码获得强类型的 CRUD
func CreatePost(...) (*AISystemPost, error) {
    return createRow[AISystemPost]("ai_system_post", data)
}
```

**没有泛型时的两个坏选择**：每个模型复制一份几乎一样的函数（重复），或用 `interface{}` 传递（丢失类型，到处断言）。泛型让**一份实现服务所有模型，且编译期类型安全**。

另一个例子是 `bindJSON[T any](c)`（api 层）和 `sortTreeChildren[T any]`（对任意树节点类型通用的递归排序）。

> **面试常问**
> - "Go 泛型什么时候该用？"（多类型共享同一算法/容器逻辑；不要为了泛型而泛型）
> - "泛型和 interface{} + 断言的区别？"（编译期检查 vs 运行时检查）

---

## 14. RBAC 权限模型与动态菜单

```
用户 ──< user_role >── 角色 ──< role_menu >── 菜单/按钮
```

权限不直接给用户，而是**给角色，用户挂角色**——新员工入职只需要挂"运营"角色，不用一个个接口去授权。两张中间表（多对多关系）是关系型数据库的标准建模。

**动态菜单**是这个模型的前端呈现：登录后前端调 `/system/user`（`CurrentUserContext`），后端沿着 用户→角色→菜单 查出这个人可见的菜单树 + 按钮权限码，前端据此**动态生成路由和界面**。菜单表里 `type` 字段区分 M（菜单）/B（按钮）——按钮不进路由树，收集成权限码数组（`permissionCodes`），前端用它控制按钮显隐。

**超级管理员**（id=1 或 code=admin）跳过角色过滤直接拿全部菜单，权限码返回 `"*"`。

**接口级校验（真正的防线）**：前端按权限码控制按钮显隐只能"拦君子"——拿到 token 的人完全可以用 curl 直接调接口。所以每条写接口都挂了 `middleware.Perm("system/user/destroy")` 这样的权限中间件（`middleware/permission.go`）：JWT 确认"你是谁"之后，Perm 再查 用户→角色→按钮菜单 确认"你能不能做这件事"，没有权限码返回 HTTP 403 + 业务码 40301。校验失败时**拒绝放行（fail closed）**——权限系统故障时宁可误伤，不能敞开。

> **面试常问**：前端隐藏按钮算不算权限控制？（不算，必须后端拦截；前端只是体验优化）

**树形数据的存储技巧**：表里存 `parent_id`（谁是我爸）+ `level`（祖先路径，如 `"0,1,5"`）。`level` 是空间换时间——查某节点所有后代只需 `level LIKE '0,1,5%'`，不用递归。`normalizeParentAndLevel` 负责维护它。

> **面试常问**
> - "RBAC 是什么？表怎么设计？"
> - "树形结构在数据库里怎么存？"（邻接表 parent_id / 路径枚举 level / 闭包表，各自权衡）

**数据权限（dataScope）——控制"能看到哪些行"**：接口权限回答"你能不能调这个接口"，数据权限回答"调了之后你能看到哪些数据"。两者互补，都在服务端强制。角色表的 `data_scope` 字段取值：

| 值 | 含义 | 过滤方式 |
|---|---|---|
| 1 | 全部数据 | 不加条件（历史数据 NULL 也按此处理，升级不缩小可见范围） |
| 2 | 自定义 | 按 `ai_system_role_dept` 关联表授权的部门过滤 |
| 3 | 本部门 | `dept_id = 我的部门` |
| 4 | 本部门及以下 | 部门树在内存里展开后代，`dept_id IN (...)` |
| 5 | 仅本人 | `id = 我自己` |

实现在 `service/system/datascope.go`：`UserDataScope` 算出操作者的可见范围（多角色取**并集**，最宽松的生效；任何一个角色是"全部"就直接放行），`applyUserDataScope` 把范围拼成 GORM 条件注入查询。目前套在用户列表（`UserList`）上，其他列表要接入时复用这两个函数即可。注意 handler 传的是 **JWT 里的当前用户 ID**，绝不信前端传参——和个人中心同一条纪律。

> **面试常问**："功能权限和数据权限的区别？"（功能权限=接口能不能调，数据权限=行级可见性；RuoYi 等后台框架的 dataScope 就是行级过滤的角色化配置）

**在线用户管理——Redis 会话的又一次复用**：登录/刷新 token 成功后按 access token 的 jti 写一条 `online:token:<jti>` 记录（JSON：用户、IP、浏览器、登录时间），TTL 与 access token 相同——**token 过期记录自动消失，不需要清理任务**。列表用 `SCAN` 遍历（不用 `KEYS`，它会阻塞 Redis）；踢下线=删掉该会话的 access/refresh jti，被踢的 token 下一次请求就过不了 `ValidateToken`。这是"Redis 存 jti"设计的直接红利：状态本来就在服务端，才有"踢人"的抓手。

---

## 14.5 定时任务调度器（robfig/cron）

`service/system/crontab_scheduler.go`，用社区标准库 `robfig/cron/v3`：

```
main.go 启动时 StartCrontabScheduler()
  → 读 ai_tool_crontab 里 status=1 的任务
  → 每条任务按 cron 表达式注册进调度器
任务增删改 → ReloadCrontabScheduler() 全量重载
```

**要点与坑**：

- **表达式解析器**：`cron.SecondOptional` 让 5 段（标准 `分 时 日 月 周`）和 6 段（带秒）都能解析——老数据是 6 段格式。创建/更新时先 `Parse` 校验，非法表达式返回业务错误而不是等运行时爆炸。
- **全量重载而非增量维护**：任务量小（几十个以内），重载简单且不会出"改了没生效"的状态不一致；量大了再优化。
- **闭包捕获循环变量**：`for _, task := range tasks` 里直接把 `task` 塞进闭包，所有任务都会执行最后一个——必须先 `taskCopy := task` 复制（Go 1.22 后 for 循环变量语义已修复，但显式复制意图更清晰）。
- **单个任务不能拖垮进程**：执行函数里 `defer recover()`，HTTP 任务设超时（15s），任何 panic/失败只写执行日志。
- **两种执行方式**：`task_style=2` 是 HTTP 任务（有参数 POST JSON，无参数 GET）；`task_style=1` 是内部任务，`target` 填注册名，从 `crontabTaskRegistry` 查函数执行。新增内部任务用 `RegisterCrontabTask` 注册，内置了 `system/clean-logs`（清理 N 天前的登录/操作日志）。
- **singleton=1 的任务**调度触发一次后自动置停用并移出调度器（手动"执行一次"不算）。

> **面试常问**
> - "定时任务怎么防止单点/重复执行？"（本项目单实例进程内调度即可；多实例部署需要分布式锁（Redis SETNX）或专门的调度中心（xxl-job 等），说清楚量级和取舍）
> - "cron 表达式 6 段和 5 段的区别？"（多出来的第一段是秒；Java Quartz 惯用 6 段，Unix crontab 是 5 段）

---

## 15. 代码生成器的原理

`service/system/codegen.go` + `codegen_templates.go`，流程：

```
1. 读元数据    information_schema.tables/columns 是 MySQL 自带的"描述表结构的表"
2. 存配置     表/字段配置落库，用户在界面上调整（哪些字段进列表/表单/查询、用什么控件）
3. 渲染模板   按配置把字符串模板填充成 Go model/service/api/router + React 页面
4. 写文件     后端写进 server/，前端写进 web/
5. 注册路由   扫描生成目录里的 RegisterXxxRoutes 函数，重写 register.go 汇总表
6. 同步菜单   自动插入菜单和按钮权限记录
```

**要点**：

- 生成的代码**调用稳定契约**（`PageListFiltered`/`CreateRecord` 等）而不是各自复制实现——修一个 bug 所有生成模块同时受益。这些契约函数的签名不能随便改，改了要同步模板和已生成代码。
- 生成的 Go 文件头部有 `// Code generated ... DO NOT EDIT.`——这是 Go 社区惯例，工具和 IDE 认识这个标记；手工改动会在下次生成时被覆盖。
- 新生成的 Go 代码需要**重新编译并重启**才生效（Go 是编译型语言，没有运行时热加载）。
- 模板质量由 `codegen_templates_test.go` 保证：用 `go/parser` 断言生成的代码语法合法。

---

## 16. 本项目真实踩过的坑（重构实录）

这些坑都在 git 历史里有据可查（`0e036bd`、`0e476c9`、`a66269c`、`7fcd713`），每个都是常见的真实教训：

### 坑 1：`map[string]interface{}` 弱类型泛滥
早期 service 全部收 `map[string]interface{}`，配合字符串白名单过滤字段。后果：绑定校验形同虚设、字段名拼错编译器不报错、到处类型断言。**教训：边界处尽早换成强类型结构体，让编译器帮你查错。**

### 坑 2：前后端字段命名不统一（camelCase vs snake_case）
响应是 `parentId`，某些页面提交却是 `parent_id`，后端被迫定义 `ParentID` + `ParentIDSnake` 双字段兜底。更隐蔽的是查询参数：后端读 `query["groupId"]`，前端发 `group_id`——**过滤条件静默失效**（不报错，只是永远查全量）。**教训：命名约定要在项目第一天定死；"静默失效"比报错更可怕。**

### 坑 3：错误详情直接返回给客户端
`err.Error()` 透传把 SQL 报错原文送到前端。**教训：错误要分级——用户看文案、开发者看日志。**

### 坑 4：可空列用非指针类型
生成器给可空的 int 列生成 `int` 字段，扫到 NULL 直接报错。**教训：数据库可空列 ↔ Go 指针类型，一一对应。**

### 坑 5：软删/硬删配置与表结构脱节
配置说"软删除"，但表可能没有 `delete_time` 列——查询直接 SQL 报错；反过来软删表被硬 DELETE，数据一去不返。**教训：任何"配置声明"都要和实际 schema 校验后再生效。**

### 坑 6：无条件生成 `import "time"`
生成的 model 没有日期字段时，未使用的 import 让整个包编译失败（Go 对未使用 import 是**编译错误**，不是警告）。**教训：代码生成器必须对"最小输入"也生成合法代码，用测试守住。**

### 坑 7：一个文件塞下所有模块
848 行的 `modules.go` 装了 8 个模块。**教训：文件按领域拆分，找代码靠目录结构而不是 Ctrl+F。**

---

## 17. 测试

本项目目前唯一的测试是 `service/system/codegen_templates_test.go`，但它示范了几个通用套路：

```go
func TestRenderGoServiceSoftDelete(t *testing.T) {
    soft := renderGoService(fullViewTypeContext(1, 1, true))   // 构造输入
    if !strings.Contains(soft, `SoftDeleteRecord`) {           // 断言输出
        t.Fatalf("软删除表应使用 SoftDeleteRecord:\n%s", soft)
    }
}
```

- **表驱动/子测试**：`t.Run(名字, func)` 一个用例矩阵跑多种场景（软删表/非软删表）；
- **纯函数优先测**：模板渲染不碰数据库，输入结构体、输出字符串，最好测；
- **用真解析器断言**：`go/parser.ParseFile` 验证生成的 Go 代码语法合法，比字符串比对可靠得多；
- 运行：`cd server && go test ./...`，加 `-v` 看每个用例，加 `-run TestXxx` 只跑匹配的。

**下一步建议**（`rules/testing.md` 也这么说）：给 `service/system/user.go`（登录）和 `utils/jwt.go` 补单测，用 `go-sqlmock`/`redismock` 模拟基础设施。这两处是安全核心，最值得测。

> **面试常问**
> - "Go 的表驱动测试是什么？"
> - "依赖数据库的代码怎么做单元测试？"（接口抽象 + mock，或 sqlmock/testcontainers）

---

## 18. 面试题速查清单

**Go 语言**
1. defer 执行顺序与常见用途（含 defer + recover 组合）
2. 指针 vs 值：什么时候用 `*string`（可空语义、避免大结构体拷贝、需要修改原值）
3. `errors.Is` / `errors.As` / `%w` 包装链
4. 泛型 vs `interface{}`：编译期与运行时类型检查
5. panic 会不会打挂整个进程？（会，除非 recover）
6. 未使用的 import/变量是编译错误——Go 的强制整洁

**Web / Gin**
7. 中间件洋葱模型；Abort 和 Next 的区别
8. RESTful 设计：GET/POST/PUT/DELETE 语义
9. CORS 原理与预检请求
10. 参数绑定与声明式校验（binding tag）

**认证 / 安全**
11. JWT 三段结构、签名原理、无状态的意义
12. JWT 注销难题与 jti + Redis 解法（本项目原文）
13. 双 token 机制与自动刷新流程
14. 算法混淆攻击（alg=none）
15. bcrypt vs MD5；加盐的意义
16. SQL 注入与参数化查询；Mass Assignment 与 DTO 白名单
17. 登录错误为什么要模糊化（防用户名枚举）

**数据库 / GORM**
18. 事务 ACID；什么场景必须包事务
19. N+1 查询的识别与解法
20. 软删除的实现方式与代价
21. 连接池三参数的含义
22. NULL 与 Go 零值的区分（指针方案）
23. 树形数据的三种存法（parent_id / 路径 / 闭包表）

**架构**
24. 三层架构各层职责；为什么 handler 要薄
25. RBAC 模型与表设计
26. 统一响应格式；业务码与 HTTP 码的分工
27. 错误分级：用户可见 vs 日志可见
28. 12-Factor：配置与代码分离

---

## 附：推荐阅读顺序

1. `main.go` → `router/router.go` → `router/system/base.go`（看清骨架）
2. `middleware/jwt.go` → `utils/jwt.go`（吃透鉴权，最有含金量）
3. `api/system/user.go` → `service/system/user.go`（一条完整业务链）
4. `api/system/common.go` → `service/system/errors.go`（错误与响应体系）
5. `service/system/dept.go` + `crud.go`（CRUD 套路与泛型）
6. `service/system/codegen.go` + `codegen_templates.go`（进阶：元编程）

每读一个文件，问自己三个问题：**这段代码为什么放在这一层？出错了会走到哪里？我删掉它会坏什么？**——能回答，就是真懂了。
