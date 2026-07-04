---
description: 后端 API 分层规范、错误处理、中间件约定（Gin + GORM）
paths: server/**
---

# 后端 API 规范

## 分层约定

```
handler (api/)   → 参数绑定 + 调 service + 组织响应，≤30 行
service/         → 业务逻辑，可依赖 GORM model 和 utils
model/           → GORM 结构体（system/）+ DTO（request/ response/）
```

- handler 禁止写 SQL 或业务判断；service 禁止直接写 HTTP 响应
- 新接口必须在 `router/system/base.go` 注册，公开接口挂 `PublicGroup`，需鉴权的挂 `PrivateGroup`

## 错误处理

- service 层用具名 sentinel error（`var ErrXxx = errors.New("...")`），handler 用 `errors.Is` 区分
- HTTP 响应统一用 `response.Success` / `response.Fail` / `response.FailWithHTTP`
- 登录相关失败统一返回泛化错误，不区分"用户名不存在"和"密码错误"

## JWT 鉴权

- 私有路由通过 `middleware.JWTAuth()` 注册，该中间件调用 `utils.ValidateToken`（含 Redis jti 校验）
- 中间件将 `user_id`（uint）和 `username`（string）写入 Gin context，handler 通过 `c.Get("user_id")` 读取
- token 过期返回 `code.AccessTokenExpired`（40101），前端拦截器会自动刷新后重试

## 路由前缀

- 默认 `ROUTER_PREFIX=/api/v1`，所有路由带此前缀
- 前端 vite proxy：`/api` → `http://localhost:8080`，并将 `/api` 重写为 `/api/v1`
- 实际访问路径：`/api/v1/base/login`，前端调用 `/api/base/login`

## 多数据库

- 目前只有 `ai_system` 库，通过 `gormInit.Gorm.Get("ai_system")` 获取连接
- 新增业务库：在 `config/go_mysql.go` 中添加配置，在 `setup/gorm/` 中注册

## 接口权限

- 需要权限控制的接口在路由注册时挂 `middleware.Perm("<权限码>")`，权限码规范为 `system/<模块>/<动作>`（如 `system/user/create`），必须与菜单表 type='B' 记录的 code 一致
- 仅登录即可的接口（当前用户上下文、下拉数据、上传、个人中心）不挂 Perm
- 新增权限码时同步：菜单表按钮记录（`ai_system.sql` 种子 + 已有库数据）、前端页面 auth 数组

## 数据权限（dataScope）

- 行级可见性过滤在 service 层做：`UserDataScope(operatorID)` 算范围，`applyUserDataScope` 注入 GORM 查询（见 `service/system/datascope.go`）
- operatorID 必须来自 JWT context（handler 用 `currentUserID(c)`），禁止接收前端传的用户 id
- 角色 `data_scope` 为 NULL 时按 1（全部）处理，避免升级缩小老角色可见范围；多角色取并集
- 新列表接入数据权限时复用上述两个函数，不要各自手写部门过滤

## 定时任务

- 任务增删改后必须调 `ReloadCrontabScheduler()` 同步调度器（service 层已封装，别绕过 service 直接写表）
- cron 表达式创建/更新时用 `validateCrontabRule` 先校验；支持 5 段和 6 段（带秒）格式
- 新增系统内部任务用 `RegisterCrontabTask(name, fn)` 注册，target 填注册名；任务函数自己兜底错误，不允许 panic 外逸
