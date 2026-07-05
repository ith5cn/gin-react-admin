---
description: 前端组件、状态管理、路由约定（React + Zustand + React Router 6）
paths: web/**
---

# 前端规范

## 状态管理（Zustand）

- 全局状态放 `src/store/`：`auth.ts`（认证/菜单/权限）、`useAppStore.ts`（布局/主题）
- 不在组件内直接调用 `useAuthStore.getState().xxx()`；在 event handler 中可用，但 React 渲染期间只用 hook selector
- `initUserContext` 只由 `Layout` 组件调用，登录页不得重复调用，避免竞态和跳转闪烁

## 路由

- 动态路由由 `/system/user` 接口下发，经 `normalizeBackendMenuTree` → `transformMenuToRoutes` 生成
- `Layout` 组件是认证门卫（token 为空跳 `/login`，`isUserInitialized=false` 时显示 Spin）
- 新增静态布局内路由：在 `src/routers/staticRoutes.tsx` 加，设 `isLayout: true`
- 新增公开路由（不需登录）：在 `src/routers/publicRouters.tsx` 加

## 组件约定

- 页面组件放 `src/pages/<模块>/index.tsx`，组件目录名对应后端 `component` 字段路径
- 通用组件放 `src/components/`，业务组件放对应模块的 `components/` 子目录
- 禁止在页面组件里写内联样式；优先 Tailwind CSS 类名，其次 Ant Design token

## API 调用

- 所有接口封装在 `src/api/` 目录下，组件不直接用 axios
- 错误处理：业务错误由 `src/utils/request.ts` 的拦截器统一弹 `message.error`，组件层 catch 后直接 `return`，不重复弹窗
- token 刷新由拦截器自动处理（40101 → refresh → 重试），组件无需感知
