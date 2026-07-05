---
description: 测试约定（当前项目无测试，补测时遵循本规范）
---

# 测试规范

## 现状

server/ 和 web/ 均无测试文件。优先补充 service 层单测，风险最高的路径是：
- `service/system/user.go` — 登录、权限、CurrentUserContext
- `utils/jwt.go` — token 签发、校验、Redis jti 验证

## Go 测试（server/）

- 文件命名：`xxx_test.go`，与被测文件同包（优先 `package xxx`，需测内部逻辑时用 `package xxx_test`）
- 运行：`go test ./...`（在 server/ 目录执行）
- Service 层测试：mock Redis 和 DB（用 `go-redis/redismock`、`go-sqlmock`），不依赖真实基础设施
- 不允许测试中硬编码真实密码、token 或数据库连接串

## TypeScript 测试（web/）

- 目前无测试框架配置；若补充，推荐 vitest（与 vite 生态一致）
- 组件测试放 `src/components/__tests__/`，工具函数测试放 `src/utils/__tests__/`
- store 测试（zustand）独立测 state 变更，不 mock 路由
