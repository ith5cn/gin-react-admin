# gin-react-admin

Go + React 全栈后台管理系统。后端 Gin/GORM/Redis/JWT，前端 React 18/Vite/Ant Design/Zustand。

## 技术栈

- 后端: Go 1.25, Gin, GORM (MySQL), Redis, JWT (golang-jwt/v5), Zap
- 前端: React 18, Vite, TypeScript, Ant Design 5, Zustand, React Router 6, Tailwind CSS
- 包管理: 后端 `go mod`，前端 `pnpm`

## 常用命令

### 后端（server/）
- 运行: `go run main.go`（需先配置 `.env`）
- 测试: `go test ./...`
- 环境变量: 从 `.env` 加载，参考 `config/go_*.go` 中的 `envOrDefault` 调用

### 前端（web/）
- 安装依赖: `pnpm install`
- 开发运行: `pnpm dev`（代理 `/api` → `http://localhost:8080`）
- 构建: `pnpm build`（先 `tsc -b` 再 `vite build`）
- Lint: `pnpm lint`

## 目录结构

```
gin-react-admin/
├── server/
│   ├── api/system/        # HTTP handler（薄层，只绑参数 + 调 service）
│   ├── service/system/    # 业务逻辑
│   ├── model/system/      # GORM model + request/response DTO
│   ├── router/system/     # 路由注册
│   ├── middleware/        # JWT 鉴权、Recovery、CORS
│   ├── utils/             # JWT 签发/校验、加密工具
│   ├── config/            # 从环境变量读取配置（go_*.go）
│   ├── setup/             # gorm / redis / logger 初始化
│   └── database/          # 初始化 SQL（ai_system.sql）
└── web/
    ├── src/api/           # axios 封装的接口函数
    ├── src/store/         # Zustand 状态（auth.ts, useAppStore.ts）
    ├── src/routers/       # React Router 路由（动态路由 + AuthGuard）
    ├── src/pages/         # 页面组件（system/, login/, install/...）
    └── src/components/    # 通用/业务组件
```

## 后端环境变量（server/.env）

关键变量（无 `.env.example`，需手动创建）：
```
SERVER_ADDR=:8080
ROUTER_PREFIX=/api/v1
DB_AI_SYSTEM_HOST=localhost
DB_AI_SYSTEM_PORT=3306
DB_AI_SYSTEM_USER=root
DB_AI_SYSTEM_PASSWORD=your_password
DB_AI_SYSTEM_DBNAME=ai_system
REDIS_ADDR=localhost:6379
JWT_SECRET=your_secret_here
JWT_ACCESS_EXPIRES_MINUTE=120
JWT_REFRESH_EXPIRES_HOUR=24
JWT_LOGIN_MODE=multi
```

## 项目现状与约束

- **无测试**：server/ 和 web/src/ 均无测试文件；新功能建议补充 service 层单测
- **动态路由**：前端菜单和路由由后端 `/system/user` 接口下发，首次加载由 `Layout` 统一初始化，不要在登录页重复调用 `initUserContext`
- **JWT + Redis 双重校验**：access token 在 Redis 中存 jti；Redis 不可用或 token 被撤销时，所有认证接口会 401，前端拦截器会硬跳回 `/login`

## 规则

以下规则全程生效，已通过 `.claude/rules/` 自动发现加载：
- coding-style.md
- testing.md
- security.md
- git-workflow.md

前端代码（`web/**`）额外受 `rules/frontend.md` 约束，后端 API（`server/**`）额外受 `rules/backend-api.md` 约束，数据库（`server/database/**`、`server/model/**`）额外受 `rules/database.md` 约束。
