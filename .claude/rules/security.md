---
description: 安全禁止事项、密钥处理、敏感文件规范
---

# 安全规范

## 密钥与凭据

- 禁止在代码中硬编码 `JWT_SECRET`、数据库密码、Redis 密码
- 所有密钥通过环境变量传入，生产环境使用系统环境变量而非 `.env` 文件
- `.env` 文件必须在 `.gitignore` 中；不得将 `.env` 提交到 git（当前项目已有 `.env` 但无 `.env.example`，需补充）

## JWT / 认证

- access token 在 Redis 中存 jti；禁止绕过 `ValidateToken`（含 Redis 校验）直接使用 `ParseToken`（只校验签名）做业务鉴权
- refresh token 只能用于 `/base/token/refresh` 接口，禁止用于普通业务接口
- 登录失败统一返回泛化错误，不暴露"用户名不存在"和"密码错误"的区别（当前已正确实现）

## 输入校验

- 所有接口入参通过 `ShouldBindJSON` + struct tag 校验，handler 层禁止直接读 raw body
- SQL 通过 GORM 参数化查询，禁止字符串拼接 SQL
- 前端表单的 XSS：渲染用户输入时不得使用 `dangerouslySetInnerHTML`

## 敏感文件

确保 `.gitignore` 包含：
```
server/.env
*.env
*.pem
*.key
```
