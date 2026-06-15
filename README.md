<div align="center">

<img src="web/public/logo.svg" width="84" height="84" alt="Go Admin" />

# Go Admin

**基于 Go + Vue 的全栈后台管理框架 · 轻量 · 现代 · 开箱即用**

`Go · Vue · 管理系统`

[![CI](https://github.com/swlfigo/go-admin/actions/workflows/ci.yml/badge.svg)](https://github.com/swlfigo/go-admin/actions/workflows/ci.yml)
[![Go](https://img.shields.io/badge/Go-1.26-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![Vue](https://img.shields.io/badge/Vue-3-42B883?logo=vuedotjs&logoColor=white)](https://vuejs.org)
[![TypeScript](https://img.shields.io/badge/TypeScript-strict-3178C6?logo=typescript&logoColor=white)](https://www.typescriptlang.org)
[![Element Plus](https://img.shields.io/badge/Element%20Plus-UI-409EFF)](https://element-plus.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white)](https://www.postgresql.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

**简体中文** · [English](README.en.md)

<br/>

<img src="snapshot/dashboard.png" alt="Go Admin 预览" width="880" />

</div>

---

## 简介

Go Admin 是一个**权限管理控制台**:超级管理员配置权限、伪超级管理员管理用户、普通用户按角色看到不同界面。核心是把 **用户 → 角色 → 菜单/按钮权限** 的 RBAC 闭环跑通——**在后台勾选即可配置,无需写 JSON 或改代码**。

选用 Go 而非 Java:编译为单个二进制、无运行时依赖、常驻内存仅几十 MB、依赖干净,适合"长期常驻、轻量"的后台服务。

## ✨ 功能特性

**🔐 认证与安全**
- 双令牌(access + refresh)+ **refresh 轮换** + 服务端会话校验
- 四层防爆破:图形验证码 · 失败锁定 · 登录限流 · bcrypt 慢哈希
- 一期零 Redis(会话/锁定存 DB,限流走内存),鉴权走可替换 `Enforcer` 接口(预留 Casbin 升级位)

**👥 RBAC 权限**
- 用户管理:增删改 · 分配角色 · 重置密码 · 解锁 · 启停用
- 角色管理:**菜单/按钮权限树勾选**(超管写死不可删)
- 菜单管理:目录/菜单/按钮三级树,管理动态路由 + 权限码
- 字典管理:类型 + 数据两级,运行时可配
- **操作日志审计**:中间件自动记录所有增删改(用户/方法/路径/状态/耗时/IP)+ 日志查询页
- **服务监控**:CPU / 内存 / 磁盘 / Go 运行时(协程数/版本/运行时长)实时看板
- 接口级授权中间件:前端隐藏按钮 + 后端独立校验

**📊 体验**
- 动态菜单(后端按角色下发 → 前端动态注册路由)+ `v-permission` 按钮级控制
- **在线用户** + 即时踢人 · 个人中心(改资料/改密码)
- 标签页:打开/关闭/**固定(pin)** + 右键菜单 + 持久化
- 主题:明暗模式 · 7 种强调色(默认 Go 青)· 界面密度 · 侧栏/标签页风格,实时切换 + 持久化
- **多语言**:中文 / English 实时切换(含 Element Plus 组件、动态菜单名),设置中可选
- 首字母头像(零上传,哈希配色)

## 🎨 品牌

Go Admin 的品牌色取自两种技术的官方颜色,天然协调:

| | 颜色 | 用途 |
|---|---|---|
| **Go 青** | `#00ADD8` | 系统主色 · 链接 · 高亮 |
| **Vue 绿** | `#42B883` | 成功状态 · 数据趋势 |
| **品牌渐变** | `#00ADD8 → #42B883` | Logo · 侧栏标识(固定,不随强调色变) |

Logo 为 **Gopher V**——Go 官方吉祥物地鼠的几何提炼:大眼睛最强识别,两颗门牙即 Vue 的双 V 标志,16px favicon 仍认得出是脸。

## 🏗 技术栈

| 层 | 选型 |
|---|---|
| **后端** | Go 1.26 · gin · gorm · PostgreSQL 16 · golang-jwt · bcrypt · viper · base64Captcha |
| **前端** | Vue 3 + TypeScript(strict,禁用 any)· Vite · Pinia · Vue Router · Element Plus · vitest |
| **架构** | 分层(handler / service / repository / model)· 设计 token 主题系统 · 动态路由 |

## 📁 目录结构

```
go-admin/
├── server/              Go 后端
│   ├── cmd/server/      程序入口
│   ├── internal/        handler / service / repository / model / middleware / auth
│   ├── pkg/             jwt / password / captcha / response 通用工具
│   └── configs/         config.yaml
├── web/                 Vue 前端
│   └── src/             api / stores / layout / views / components / router
└── docs/                设计文档
```

## 🚀 快速开始

**前置依赖**:Go 1.26+ · Node 18+ · PostgreSQL 16(本地用 Docker 起即可,或任意可连的 PostgreSQL)

**1. 启动 PostgreSQL**

```bash
docker run -d --name go-admin-pg \
  -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=go_admin \
  -p 5432:5432 postgres:16
```

**2. 后端**(`:8080`)

```bash
cd server
make run          # 自动建表 + 种子数据(预置超管/角色/菜单/字典)
```

**3. 前端**(`:5173`)

```bash
cd web
npm install
npm run dev
```

浏览器打开 **http://localhost:5173**,默认账号 **`admin` / `admin123`**。

## 🐳 Docker 一键部署

本机只要有 Docker,不用装 Go/Node:

```bash
docker compose up -d --build
```

访问 **http://localhost:8088**,默认账号 `admin` / `admin123`。

- 三个服务:`db`(PostgreSQL)+ `server`(Go 后端,首次自动建表 + 种子)+ `web`(Caddy 托管 SPA 并反代 `/api`)
- **镜像全部本地构建**,只拉取公共基础镜像(postgres / golang / node / caddy)——**无需 DockerHub 账号**
- 生产前改 `docker-compose.yml` 里的 `GOADMIN_JWT_SECRET`(及可选 `GOADMIN_ADMIN_PASSWORD`)
- 停止:`docker compose down`(数据保留在 `pgdata` 卷)

> 配置支持**环境变量覆盖**(`GOADMIN_<段>_<键>`,如 `GOADMIN_DATABASE_HOST`、`GOADMIN_JWT_SECRET`),无需改 `config.yaml`,便于注入密钥。

## ⚙️ 配置说明

全部配置集中在 `server/configs/config.yaml`:

| 配置块 | 说明 |
|---|---|
| `database` | PostgreSQL 连接(host/port/user/password/name/sslmode) |
| `jwt` | `secret`(⚠️ 生产务必换强随机串)· access/refresh 有效期 |
| `admin` | **初始超管账号/密码**(仅首次建库写入,之后改密不被覆盖) |
| `captcha` | 图形验证码开关 |
| `login` | 失败锁定阈值 · 锁定时长 · 登录限流 |

> **安全提示**:上线前务必更换 `jwt.secret` 与 `admin.password`。

## 🧪 测试

```bash
cd server && make test     # 后端:集成测试(串行,需 PostgreSQL)
cd web && npm run test      # 前端:vitest 单元测试
```

## 📄 许可证

[MIT](LICENSE)

---

<div align="center">
<sub>Go <code>#00ADD8</code> × Vue <code>#42B883</code> · Manrope · Built with Go Admin</sub>
</div>
