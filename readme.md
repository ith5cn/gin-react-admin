# Gin React Admin

一个用 **Gin + GORM + React + Ant Design** 搭出来的后台管理系统。我做这个项目，不只是为了再写一个 admin 模板，而是想把自己平时做业务后台、权限系统、代码生成、数据维护时积累的经验，沉淀成一套可以学习、可以二开、也可以继续长大的工程。

如果你也在做中后台、低代码、AI 工具、CMS、SaaS 管理端，或者想看一个 Go + React 项目怎么从 0 慢慢长成完整系统，欢迎 Star、Fork、提 Issue。这个仓库会持续迭代，我也会把开发过程中的取舍、坑和重构都尽量留在代码里。

## 项目定位

`gin-react-admin` 是一个偏实战的全栈后台基础工程：

- 后端使用 Go / Gin / GORM，保留清晰的 `api -> service -> model -> router` 分层。
- 前端使用 React / Vite / Ant Design，动态菜单驱动路由。
- 支持 JWT 登录、Redis token 状态、菜单权限、角色权限、字典、配置、附件、日志等后台基础能力。
- 内置首次安装向导，下载代码后可以访问 `/install` 一步步完成环境配置和 SQL 导入。
- 内置代码生成和数据表维护，希望让“后台系统本身”也能帮助开发者更快交付业务。

## 功能特性

- 登录认证：JWT access token / refresh token，Redis 维护 token 状态。
- 权限体系：用户、角色、菜单、按钮权限、动态路由。
- 组织管理：部门、岗位、用户角色绑定。
- 系统配置：配置分组、配置项、批量更新。
- 数据字典：字典类型、字典数据、前端统一渲染。
- 附件管理：附件列表、资源分类、删除确认。
- 日志模块：登录日志、操作日志。
- 代码生成：从数据库表导入字段，生成 Go 后端和 React 前端基础代码。
- 数据表维护：表列表、表结构、碎片清理、优化表、软删除数据回收站。
- 首次安装：`/install` 引导填写 MySQL、Redis、选择 SQL 文件并写入安装锁。

## 技术栈

后端：

- Go
- Gin
- GORM
- MySQL
- Redis
- Zap
- JWT

前端：

- React
- TypeScript
- Vite
- Ant Design
- Zustand
- Axios
- Tailwind CSS
- Lucide React

## 相关版本

这个项目是 Go 版本的后台工程实践。同一套业务思路和后台架构，我也维护了一个 **NestJS 版本**，适合更熟悉 Node.js / TypeScript 技术栈的开发者参考。

- Go 版本：`gin-react-admin`，当前仓库，技术栈为 Gin + GORM + React。
- NestJS 版本：独立仓库，技术栈为 NestJS + TypeScript。仓库地址：`<your-nestjs-repo-url>`

两个版本不会强行保持代码一模一样，而是尽量保持相同的产品思路、权限模型、菜单体系和工程分层。你可以把它们看成同一个后台系统理念在不同技术栈里的两种实现。

## 目录结构

```text
gin-react-admin
├── server                 # Go 后端
│   ├── api                # HTTP 入参、响应处理
│   ├── config             # 环境配置读取
│   ├── database           # 初始化 SQL
│   ├── middleware         # CORS、JWT、Recovery、日志等中间件
│   ├── model              # 数据模型、request、response DTO
│   ├── router             # 路由注册
│   ├── service            # 业务逻辑
│   ├── setup              # MySQL、Redis、Logger 初始化
│   └── utils              # 工具函数
├── web                    # React 前端
│   ├── public
│   └── src
│       ├── api            # 前端接口封装
│       ├── components     # 公共组件
│       ├── pages          # 页面
│       ├── routers        # 静态路由、动态路由转换
│       ├── store          # Zustand 状态
│       └── utils
└── readme.md
```

## 快速开始

### 1. 克隆项目

```shell
git clone <your-repo-url>
cd gin-react-admin
```

### 2. 启动后端

```shell
cd server
go mod tidy
go run .
```

后端默认监听：

```text
http://localhost:8080
```

### 3. 启动前端

```shell
cd web
pnpm install
pnpm run dev
```

前端默认访问：

```text
http://localhost:5173
```

### 4. 首次安装

浏览器访问：

```text
http://localhost:5173/install
```

安装向导会完成：

- 配置 MySQL
- 配置 Redis
- 选择并导入 SQL
- 写入 `server/.env`
- 写入 `server/runtime/install.lock`

初始化 SQL 默认放在：

```text
server/database/ai_system.sql
```

你也可以把自己的完整 SQL 放到：

```text
server/sql/ai_system.sql
```

安装页会自动扫描 `server/sql/*.sql` 和 `server/database/*.sql`。

默认管理员账号以 SQL 文件为准。当前初始化数据通常为：

```text
账号：admin
密码：123456
```

## 环境配置

后端主要读取 `server/.env`：

```env
AI_SYSTEM_MYSQL_HOST=127.0.0.1
AI_SYSTEM_MYSQL_PORT=3306
AI_SYSTEM_MYSQL_USER=root
AI_SYSTEM_MYSQL_PASSWORD=123456
AI_SYSTEM_MYSQL_DB=ai_system
AI_SYSTEM_MYSQL_CONFIG=charset=utf8mb4&parseTime=True&loc=Local

REDIS_MODE=single
REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=
REDIS_DB=0

JWT_SECRET=gin-react-admin-change-me
SERVER_ADDR=:8080
ROUTER_PREFIX=/api/v1
```

项目同时兼容 `/api/v1` 和 `/api` 前缀，方便 Vite 代理和直接访问后端接口。

## 分层约定

这个项目尽量保持后端分层简单清楚：

- `router`：只负责 URL 到 handler 的映射。
- `api`：处理参数绑定、校验、响应格式。
- `service`：写业务逻辑、事务、数据库查询。
- `model`：放数据库模型、请求 DTO、响应 DTO。
- `setup`：初始化数据库、Redis、日志等基础设施。

我比较喜欢这种“够用但不过度设计”的结构：新人能看懂，业务也能继续长。

## 适合谁

- 想学习 Go + React 全栈后台项目的人。
- 想快速启动一个管理端基础工程的人。
- 正在做 CMS、SaaS、AI 工具后台、低代码平台的人。
- 想研究权限、动态菜单、代码生成、安装向导、数据维护这些后台常见能力的人。

## 路线图

- 完善更多业务模块生成模板。
- 增加接口权限校验和操作审计细节。
- 补充 Docker / Docker Compose 一键启动。
- 增加在线文档和截图。
- 优化代码生成后的路由聚合、菜单同步和权限同步。
- 增加更多单元测试和端到端测试。

## 作者的话

我是一个长期折腾后台系统、自动化研发工具和 AI 工作流的开发者。这个项目会记录我对“工程效率”的一些理解：能自动生成的不要重复写，能配置的不要硬编码，能在安装阶段解决的不要丢给使用者猜。

如果你喜欢这种方向，欢迎关注、Star、Fork，也欢迎直接提 Issue 或 PR。  
我希望这个仓库不只是一个代码模板，而是一个能聚集同频开发者的小现场：大家一起把中后台工程做得更顺手、更清楚、更有生命力。

## License

计划以开源协议发布。正式开源前建议补充 `LICENSE` 文件，例如 MIT License。
