# Go-Admin 后台管理系统 — 设计文档

> 日期：2026-06-15
> 状态：设计待评审

## 1. 项目目标与技术选型理由

一个轻量的**权限管理控制台**：超级管理员配置权限，伪超级管理员管理用户，普通用户按角色看到不同界面。核心是把"用户 → 角色 → 菜单/按钮权限"的 RBAC 闭环跑通，**在后台勾选即可配置，无需写 JSON 或改代码**。

**为什么用 Go 而非 Java**：Java 后管（ruoyi/jeecg）生态成熟、可作功能参考，但 Spring 全家桶依赖多、JVM 内存占用大（512MB+）、启动慢。Go 编译为单个二进制、无运行时依赖、常驻内存仅几十 MB、依赖干净（gin/gorm 几个）、部署即扔二进制 + 配置。对"长期常驻、不需要重武器"的后管服务，Go 更省资源、更易维护。

**功能参考但不照搬**：参考 ruoyi-go/ruoyi-web、gin-vue-admin、vben/soybean-admin 的设计**模式**，只做真正需要的功能，砍掉监控/操作日志/定时任务/代码生成/公告/岗位等包袱。

**跨生态选型对比（调研结论）**：
- **后端 Go vs Java vs Node**：Java（JeecgBoot 45.5k★、RuoYi）生态最成熟、模块最多，但 Spring 全家桶 + JVM 重，正是要避开的；Node（NestJS）居中；Go 在"轻量后管"上是最优取舍 → **选 Go**。
- **前端 Vue vs React**：React（Ant Design Pro 38.5k★、react-admin/refine）组件强、AI 语料多，且设计稿本身是 JSX（React 还原更省力）；Vue（vben 31k★、gin-vue-admin、ruoyi）中文后管参考最丰富、权限骨架可大量借鉴 → **选 Vue3**（用户决策）。
- **TS 维持**：作为"代码统一 + 方便 AI 开发/debug"的手段——类型即护栏，能在编译期挡错、给 AI 更多上下文；**禁用 any**。

## 2. 技术栈

### 后端（Go）
| 项 | 选型 | 说明 |
|---|---|---|
| Web 框架 | gin | 主流、轻量 |
| ORM | gorm | 自动迁移建表 |
| 数据库 | **PostgreSQL 16**（本地用 Docker 容器） | |
| 认证 | JWT（无状态） | 登录返回 token，中间件校验 |
| 密码 | bcrypt | 加盐哈希 |
| 配置 | viper + yaml | DB/JWT 等配置外置 |
| 鉴权 | 自研 `Enforcer` 接口 | 查表实现；**不引入 Casbin**，但留可替换接口，未来可换 |
| Redis | 暂不引入 | 验证码/限流/token 黑名单等需要时再加 |

### 前端（Vue3 + TS）
| 项 | 选型 | 说明 |
|---|---|---|
| 框架 | Vue 3 + TypeScript（**禁用 any**） | 每个 API 有请求/响应 type |
| 构建 | Vite | |
| 路由 | Vue Router（动态注册） | 登录后按菜单树建路由 |
| 状态 | Pinia | user / permission / app-theme 三个 store |
| HTTP | axios + 拦截器 | 统一 `{code,msg,data}`；401 自动用 refresh 静默续期，失败才跳登录 |
| UI 实现 | **C 方案**：外壳手写 + Element Plus 重组件 | 见下 |

**前端 C 方案说明**：
- **外壳照搬设计稿**（`/Users/sylar/Downloads/admin`）：侧栏、topbar、卡片式 tab、按钮、徽章、drawer、主题切换（亮/暗、强调色、密度）→ 直接移植已写好的 CSS token，JSX→Vue 模板机械转译，**视觉 100% 还原**。
- **重组件用 Element Plus**：数据表格（排序/分页）、**权限树勾选**、表单校验、日期选择、弹窗 → 用 EP 并套设计稿配色，稳健不踩坑。

## 3. 角色体系

可自定义角色体系（角色为一张表，超管可任意新建角色并勾选权限）。系统预置 3 个角色：

| 角色 | code | 能力 | 约束 |
|---|---|---|---|
| 超级管理员 | `super_admin` | 全部：管菜单、管角色、管用户、管字典 | **写死不可删**，拥有全部权限（鉴权直接放行） |
| 伪超级管理员 | `admin` | 管理普通用户（增删改、分配角色、重置密码） | 不能改菜单、不能配角色权限 |
| 普通用户 | `user` | 按所分配角色看到对应界面 | 仅基础访问 + 个人中心 |

> 超管之外的角色权限完全由超管在后台勾选决定，预置角色只是初始示例。

## 4. 功能模块（一期 8 个）

> 模块取舍依据：对照 RuoYi / JeecgBoot / gin-vue-admin / go-admin / vben / soybean 的全量模块，只保留服务于权限闭环且轻量的功能。

1. **登录认证** — JWT 登录、bcrypt 校验、登出；**防爆破两层**：图形验证码（config 开关，默认开）+ 失败次数锁定（连错 N 次锁账号 M 分钟，超管可手动解锁）
2. **仪表盘** — 登录后落地页，设计稿风格，内容最简（欢迎语 + 统计卡）
3. **用户管理** — 列表/增删改、分配角色、重置密码、启停用（超管 + 伪超管）
4. **角色管理** — 列表/增删改、**勾选菜单+按钮权限**（超管）
5. **菜单/权限管理** — 菜单树增删改，节点含「目录/菜单/按钮」三类型，按钮节点带权限码（超管）
6. **字典管理** — 字典类型 + 字典数据两级，运行时可增删改；init.sql 仅种 2-3 个基础字典示范
7. **在线用户** — 列出当前在线会话（用户/IP/浏览器/登录时间），支持**强制下线（踢人）**
8. **个人中心** — 改本人资料/密码

**头像方案（不上传）**：默认用**首字母头像**——前端组件取昵称/用户名首字，背景色由用户名哈希稳定生成（类 Gmail/Slack），零上传零存储。`sys_user.avatar` 字段保留为可选 URL，留未来上传扩展。

**砍掉的模块**（成熟项目有但本期不要）：部门/数据权限、岗位、参数配置、通知公告、登录日志、定时任务、缓存/数据监控、代码生成、表单构建、Swagger、多租户、工作流、低代码、站内信。

> 注：**i18n（中英）、操作日志审计、服务监控**当初也在此列，后续已补做，见「功能特性」/ README。

**二期留接口**：参数配置等。设计上预留：菜单表节点可扩展、鉴权走可替换 `Enforcer` 接口。（部门 + 数据权限明确不做，见 §10。）

## 5. 数据库设计（PostgreSQL）

```
sys_user            用户
  id, username(uniq), password(bcrypt), nickname, avatar, email, phone,
  status(0停用/1启用), last_login_at,
  login_fail_count, lock_until(防爆破：失败次数 + 锁定到期时间),
  created_at, updated_at, deleted_at

sys_role            角色
  id, name, code(uniq, 如 super_admin), sort, status, remark,
  created_at, updated_at, deleted_at

sys_user_role       用户-角色 关联（多对多）
  user_id, role_id

sys_menu            菜单/权限（树形，自引用）
  id, parent_id(0为根), name, type(M目录/C菜单/F按钮),
  path(路由路径), component(前端组件路径), perm(权限码, 如 user:delete),
  icon, sort, visible(0隐藏/1显示), status,
  created_at, updated_at, deleted_at

sys_role_menu       角色-菜单 关联（多对多）—— RBAC 勾选结果存这里
  role_id, menu_id

sys_dict_type       字典类型
  id, name, type(uniq, 如 sys_status), status, remark, created_at, updated_at

sys_dict_data       字典数据
  id, dict_type(关联 sys_dict_type.type), label, value, css_class/tag_type,
  sort, status, remark, created_at, updated_at

sys_online_session  在线会话（= 在线用户列表，支撑踢人 + refresh 轮换）
  id(=access jti), user_id, username, login_ip, browser, os,
  refresh_id(当前有效 refresh 标识), refresh_expire_at,
  login_at, last_active_at
```
> 说明：`sys_user.avatar` 为可选字段，默认空时前端渲染首字母头像。

**关键关系**：`用户 --(sys_user_role)-- 角色 --(sys_role_menu)-- 菜单/权限`。
登录后：聚合用户所有角色的菜单并集 → 得到该用户的菜单树（建前端路由）+ 权限码集合（控按钮）。超管角色特判：返回全部菜单。

## 6. 权限实现

### 后端鉴权 — 可替换 Enforcer 接口
```go
type Enforcer interface {
    // 该用户（其角色集合）是否拥有某权限码
    Can(userID uint, permCode string) bool
}
// 一期：DBEnforcer 查 user_role -> role_menu -> menu.perm；super_admin 直接 true
// 未来：可写 CasbinEnforcer 实现同接口，中间件零改动
```
gin 中间件：从 JWT 取 userID → 调 `Enforcer.Can(userID, 当前接口要求的permCode)` → 放行或 403。**前端隐藏按钮不等于安全，后端必须独立校验。**

### 认证模型 — 双令牌 + 轮换 + 会话校验（综合业界两派）
整合"无状态 + 刷新令牌"(Auth0/WorkOS) 与"服务端令牌存储"(gin-vue-admin/ruoyi) 两派之长：

- **双令牌**：登录返回 `access`（JWT，短期，默认 30min）+ `refresh`（长期，默认 7 天）。access 携带 `jti`。
- **会话存储**：每次登录在 `sys_online_session` 写一条会话（含 `jti` 与 refresh 标识）。**这张表即"在线用户"列表**。
- **会话校验（即时踢人）**：鉴权中间件校验 access 的 `jti` 会话是否存在（不存在 = 已被踢/已登出 → 401），并刷新 `last_active_at`。**强制下线 = 删除该会话行 → 该 token 下次请求即 401，即时生效**。代价是每请求一次带索引查询（后管量级可忽略）。
- **静默续期 + 轮换**：access 过期时，前端用 refresh 调 `POST /auth/refresh` 换发新 access + **新 refresh 并作废旧 refresh**（轮换，旧 refresh 被重放即视为盗用，拒绝并可告警）。前端 axios 拦截器在 401 时自动续期、失败才跳登录。
- **安全细节**：HS256 + 强密钥（单后端足够）；校验 exp/iss/alg；refresh 建议存 HttpOnly cookie 防 XSS。
- **零 redis**：会话/refresh 一期存 DB，二期多实例时整体迁 redis（TTL 自动过期）。

### 防爆破（两层）
- **图形验证码**：由 `config.captcha.enabled` 控制（默认 true）。开启时 `GET /captcha` 返回图形 + captchaId，登录需带 captchaId + 用户输入校验。
- **失败锁定**：连续输错密码达 `config.login.maxFailCount`（默认 5）→ 锁定该账号 `config.login.lockMinutes`（默认 10）分钟；失败计数与 `lock_until` 存 `sys_user`。登录成功清零；锁定期内拒绝并提示剩余时间；超管可在用户管理手动解锁。
- **登录限流**：IP/账号维度的频率限制中间件，单实例用**内存限流器**（令牌桶/滑动窗口，无需 redis）；超限返回 429。多实例部署时改用 redis 共享计数（二期）。
- **bcrypt 慢哈希**：密码验证本身耗时，拖慢批量试探。

> 三层防爆破（验证码 + 失败锁定 + 内存限流）+ bcrypt，**一期全程零 redis**：失败锁定存 DB、会话存 DB、限流走内存。redis 仅在二期"多实例 + 高并发"时统一升级。

### 前端权限 — 静态公共页 + 动态业务页
仿 ruoyi 的混合路由模型：

- **静态公共页（constantRoutes，前端写死，无需鉴权）**：`登录`、`403`（已登录无权限）、`404`（路由不存在）、`500`（服务异常）。
- **动态业务页（按角色由后端返回）**：登录后 `GET /menus/me` 返回该用户菜单树 → permission store 转成路由并 `addRoute`。这样"超管在后台增删页面"才能生效，无需改前端代码。
- **按钮指令**：`v-permission="'user:delete'"`，从 user store 的权限码集合判断显隐。
- **路由守卫**：无 token → 跳登录（白名单除外）；有 token 无路由 → 先拉菜单；访问无权限路由 → 403；未匹配任何路由 → 404。

## 7. API 清单（一期）

```
认证      GET    /captcha               图形验证码（config 开启时）
          POST   /auth/login            登录 → access + refresh
          POST   /auth/refresh          用 refresh 换新 access + 轮换 refresh
          POST   /auth/logout           登出（删除本会话）
          GET    /auth/me               当前用户信息 + 权限码
          GET    /menus/me              当前用户菜单树（建路由）

在线      GET    /online                 在线会话列表（在线用户）
          DELETE /online/:sessionId     强制下线（踢人）—— 即时生效

用户      GET    /users                 列表（分页/搜索）
          POST   /users                 新增
          PUT    /users/:id             修改
          DELETE /users/:id             删除
          PUT    /users/:id/roles       分配角色
          PUT    /users/:id/password    重置密码
          PUT    /users/:id/status      启停用

角色      GET    /roles                 列表
          POST   /roles                 新增
          PUT    /roles/:id             修改
          DELETE /roles/:id             删除
          PUT    /roles/:id/menus       勾选菜单/按钮权限

菜单      GET    /menus                 菜单树（管理用）
          POST   /menus                 新增
          PUT    /menus/:id             修改
          DELETE /menus/:id             删除

字典      GET    /dict/types            类型列表
          POST   /dict/types            ... 增删改
          GET    /dict/data?type=xxx    某类型下数据
          POST   /dict/data             ... 增删改

个人      PUT    /profile               改本人资料
          PUT    /profile/password      改本人密码
```
统一响应：`{ "code": 0, "msg": "ok", "data": ... }`。

## 8. 目录结构

### 后端
```
go-admin/
  cmd/server/main.go
  configs/config.yaml
  internal/
    handler/     gin handler（路由层）
    service/     业务逻辑
    repository/   gorm 数据访问
    model/       表结构 struct
    middleware/   jwt、enforcer 鉴权、跨域、日志
    auth/        Enforcer 接口 + DBEnforcer 实现
  pkg/           jwt、response、bcrypt 等通用工具
  migrations/    init.sql（种子：3 角色 + 基础菜单 + 示范字典 + 超管账号）
  go.mod
```

### 前端
```
web/
  src/
    api/         按模块的类型化接口（user.ts/role.ts/...）
    components/   外壳组件（Sidebar/Topbar/TabStrip/UserAvatar 首字母头像/...，移植设计稿）
    layout/       整体布局
    views/        8 个业务页（login/dashboard/user/role/menu/dict/online/profile）
    views/error/  静态错误页（403/404/500）
    router/       动态路由 + 守卫
    store/        user / permission / app-theme
    directives/   v-permission
    styles/       设计稿 CSS token
    types/        全局 TS 类型
  vite.config.ts
```

## 9. 开发与运行

- 本地用 Docker 起 PostgreSQL 容器；DB 连接、JWT 密钥、`captcha.enabled` 等走 `config.yaml`。
- 首次启动：gorm AutoMigrate 建表 + 执行 `init.sql` 种子（超管账号、预置角色、基础菜单、示范字典）。
- 前后端分离：后端 `:8080` 提供 API，前端 vite dev server 代理。

## 10. 二期扩展点（设计已预留，本期不实现）

- **Redis（多实例/高并发统一升级）**：把一期用 DB/内存实现的三件事平滑迁过去——会话存储（替 `sys_online_session`）、失败锁定计数（带 TTL 自动解锁）、登录限流（共享计数器）。一期均不依赖。
- **鉴权引擎升级**：实现 `CasbinEnforcer` 替换 DBEnforcer，中间件零改动。

> 注：部门 + 数据权限（data_scope）是 ruoyi 类企业框架的概念，本项目作为轻量控制台不需要，已明确不做（业务真有行级数据归属需求时再按具体场景加）。
