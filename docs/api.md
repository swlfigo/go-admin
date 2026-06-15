# Go Admin · API 文档

后端基址:`http://localhost:8080/api`(前端通过 Vite 代理 `/api` 访问)。

## 通用约定

**统一响应体**

```json
{ "code": 0, "msg": "ok", "data": {} }
```

- 成功:`code = 0`,业务数据在 `data`。
- 失败:HTTP 状态码 ≠ 200,响应 `code` 复用 HTTP 状态码,`msg` 为错误信息。
- **客户端应按 HTTP 状态码判断成败**(不要只看 `code`)。

**认证**:除登录/验证码/续期外,所有接口需带

```
Authorization: Bearer <accessToken>
```

**常见状态码**

| 状态 | 含义 |
|---|---|
| 200 | 成功 |
| 400 | 参数错误 / 验证码错误 / 业务校验失败 |
| 401 | 未登录 / token 失效 / 被强制下线 |
| 403 | 已登录但无该接口权限 |
| 429 | 登录请求过于频繁(限流) |

**分页**:列表接口接受查询参数 `keyword`、`page`(默认 1)、`size`(默认 10),返回 `{ "list": [...], "total": n }`。

---

## 认证

### `GET /captcha` — 图形验证码
返回 `{ "enabled": true, "captchaId": "...", "img": "data:image/png;base64,..." }`(`captcha.enabled=false` 时仅返回 `{ "enabled": false }`)。

### `POST /auth/login` — 登录
请求:`{ "username", "password", "captchaId", "captchaCode" }`
响应:`{ "accessToken", "refreshToken" }`
失败:验证码错误 / 用户名密码错误 / 账号锁定 / 账号停用。

### `POST /auth/refresh` — 刷新令牌
请求:`{ "refreshToken" }` → 响应:`{ "accessToken", "refreshToken" }`(旧 refresh 轮换作废)。

### `POST /auth/logout` — 登出 🔒
删除当前会话(该 token 立即失效)。

### `GET /auth/me` — 当前用户 🔒
响应:`{ "user": { ...,"roles":[...] }, "perms": ["system:user:list", ...] }`。

### `GET /menus/me` — 当前用户菜单树 🔒
响应:`MenuNode[]`(已过滤按钮类型,用于前端动态建路由)。

---

## 个人中心 🔒（仅需登录，无需特定权限）

| 方法 | 路径 | 请求体 | 说明 |
|---|---|---|---|
| PUT | `/profile` | `{ nickname, email, phone }` | 改本人资料 |
| PUT | `/profile/password` | `{ oldPassword, newPassword }` | 改本人密码(校验旧密码) |

---

## 用户管理 🔒

| 方法 | 路径 | 权限码 | 请求体 / 说明 |
|---|---|---|---|
| GET | `/users` | `system:user:list` | 查询 `keyword/page/size` → `{list,total}` |
| POST | `/users` | `system:user:create` | `{ username, password, nickname, email, phone, roleIds }` |
| PUT | `/users/:id` | `system:user:update` | `{ nickname, email, phone, status }` |
| DELETE | `/users/:id` | `system:user:delete` | 软删除 |
| PUT | `/users/:id/roles` | `system:user:assignrole` | `{ roleIds: number[] }` |
| PUT | `/users/:id/password` | `system:user:resetpwd` | `{ password }` 重置密码 |
| PUT | `/users/:id/unlock` | `system:user:update` | 解锁(清失败计数) |

## 角色管理 🔒

| 方法 | 路径 | 权限码 | 请求体 / 说明 |
|---|---|---|---|
| GET | `/roles` | `system:role:list` | `{list,total}` |
| GET | `/roles/:id` | `system:role:list` | 角色详情(含已分配菜单) |
| POST | `/roles` | `system:role:create` | `{ name, code, sort, status, remark }` |
| PUT | `/roles/:id` | `system:role:update` | 同上(code 不可改) |
| DELETE | `/roles/:id` | `system:role:delete` | 内置 `super_admin` 不可删 |
| PUT | `/roles/:id/menus` | `system:role:assignmenu` | `{ menuIds: number[] }` 全量替换 |

## 菜单管理 🔒

| 方法 | 路径 | 权限码 | 请求体 / 说明 |
|---|---|---|---|
| GET | `/menus` | `system:menu:list` | 完整菜单树(含按钮) |
| POST | `/menus` | `system:menu:create` | `{ parentId, name, type(M/C/F), path, component, perm, icon, sort, visible, status }` |
| PUT | `/menus/:id` | `system:menu:update` | 同上 |
| DELETE | `/menus/:id` | `system:menu:delete` | 有子节点不可删 |

> `type`:`M` 目录 · `C` 菜单(对应一个页面/路由)· `F` 按钮(带权限码,不进路由)。

## 字典管理 🔒

| 方法 | 路径 | 权限码 | 请求体 / 说明 |
|---|---|---|---|
| GET | `/dict/types` | `system:dict:list` | `{list,total}` |
| POST | `/dict/types` | `system:dict:create` | `{ name, type, status, remark }` |
| PUT | `/dict/types/:id` | `system:dict:update` | |
| DELETE | `/dict/types/:id` | `system:dict:delete` | |
| GET | `/dict/data?type=xxx` | `system:dict:list` | 某类型下的字典项 `DictData[]` |
| POST | `/dict/data` | `system:dict:create` | `{ dictType, label, value, tagType, sort, status, remark }` |
| PUT | `/dict/data/:id` | `system:dict:update` | |
| DELETE | `/dict/data/:id` | `system:dict:delete` | |

## 在线用户 🔒

| 方法 | 路径 | 权限码 | 说明 |
|---|---|---|---|
| GET | `/online` | `monitor:online:list` | 当前在线会话 `OnlineSession[]` |
| DELETE | `/online/:id` | `monitor:online:kick` | 强制下线(`:id` 为会话 id,即 access 的 jti) |

---

## 数据结构

```ts
User    { id, username, nickname, avatar, email, phone, status, lastLoginAt, roles[] }
Role    { id, name, code, sort, status, remark }
MenuNode{ id, parentId, name, type, path, component, perm, icon, sort, visible, status, children[] }
DictType{ id, name, type, status, remark }
DictData{ id, dictType, label, value, tagType, sort, status, remark }
OnlineSession { id, userId, username, loginIp, browser, os, loginAt, lastActiveAt }
```

> 🔒 = 需要 `Authorization: Bearer <accessToken>`;带权限码的接口还需对应角色权限,否则返回 403。超级管理员(`super_admin`)绕过所有权限校验。
