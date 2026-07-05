---
description: 命名、缩进、import、注释规范（Go + TypeScript）
---

# 代码风格

## Go（server/）

- 缩进：`gofmt` 标准 tab，不手动格式化
- 命名：exported 用 PascalCase，package-private 用 camelCase，常量全大写加下划线
- 错误处理：`if err != nil { return ..., err }` 而非 panic；业务错误用 `errors.Is` 比较，不做字符串匹配
- 日志：用 `zap.Logger`（`loggerInit.Logger.Get()`），禁用 `fmt.Print*` 打日志
- HTTP handler 只做参数绑定和响应，不写业务逻辑；业务逻辑放 service 层
- 注释：exported 函数/类型必须有 godoc；内部函数只在 WHY 不明显时加注释
- 不写 `_ "unused"` 占位注释、不留 TODO 注释进 main 分支

## TypeScript / React（web/）

- 缩进：2 空格（vite 默认）
- 命名：组件 PascalCase，hook `use` 前缀，普通函数/变量 camelCase，类型/接口 PascalCase
- 组件：优先函数组件 + hooks，禁用 class 组件
- import 顺序：外部库 → `@/` 内部路径 → 相对路径，每组空一行
- 类型：优先用 `type`，只在需要 `implements`/`extends` 时用 `interface`
- 不用 `any`，除非对接无类型的第三方库且必须加注释说明原因
- 注释：只在 WHY 非显而易见时写，不写"这段代码做了什么"的注释
