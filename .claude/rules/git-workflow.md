---
description: 分支、commit、PR 规范
---

# Git 工作流

## 现状

当前 git 历史使用中文 commit（"前端代码"、"1"），尚未启用 Conventional Commits。

## Commit 规范

采用 Conventional Commits 格式（推荐逐步落地）：

```
<type>(<scope>): <subject>

[optional body]
```

- type: `feat` / `fix` / `refactor` / `chore` / `docs` / `test`
- scope: `auth` / `user` / `menu` / `role` / `frontend` / `backend` 等
- subject: 动词开头，不加句号，中英文均可
- 例：`fix(auth): 移除登录页重复调用 initUserContext 导致的跳转闪烁`

## 分支命名

- `feat/<短描述>` — 新功能
- `fix/<短描述>` — bug 修复
- `chore/<短描述>` — 构建/依赖/脚手架变更

## 提交原则

- 一个 commit 只做一件事
- 不提交 `.env`、`node_modules/`、`dist/`、`vendor/`
- 不跳过 pre-commit hook（`--no-verify`）
- 涉及数据库 schema 变更时，同步更新 `server/database/ai_system.sql`
