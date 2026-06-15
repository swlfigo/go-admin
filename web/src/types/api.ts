export interface ApiResult<T> {
  code: number
  msg: string
  data: T
}

export interface Role {
  id: number
  name: string
  code: string
  sort: number
  status: number
  remark: string
}

export interface User {
  id: number
  username: string
  nickname: string
  avatar: string
  email: string
  phone: string
  status: number
  lastLoginAt: string | null
  roles: Role[]
}

export interface MenuNode {
  id: number
  parentId: number
  name: string
  type: string // M | C | F
  path: string
  component: string
  perm: string
  icon: string
  sort: number
  visible?: number
  status?: number
  children?: MenuNode[]
}

export interface LoginResult {
  accessToken: string
  refreshToken: string
}

export interface MeResult {
  user: User
  perms: string[]
}

export interface CaptchaResult {
  enabled: boolean
  captchaId?: string
  img?: string
}

export interface Paged<T> {
  list: T[]
  total: number
}

export interface DictType {
  id: number
  name: string
  type: string
  status: number
  remark: string
}

export interface DictData {
  id: number
  dictType: string
  label: string
  value: string
  tagType: string
  sort: number
  status: number
  remark: string
}

export interface OnlineSession {
  id: string
  userId: number
  username: string
  loginIp: string
  browser: string
  os: string
  loginAt: string
  lastActiveAt: string
}

export interface OperationLog {
  id: number
  username: string
  method: string
  path: string
  ip: string
  status: number
  latencyMs: number
  createdAt: string
}

export interface ServerInfo {
  host: {
    hostname: string
    os: string
    platform: string
    arch: string
  }
  cpu: {
    percent: number
    cores: number
  }
  memory: {
    usedMB: number
    totalMB: number
    percent: number
  }
  disk: {
    usedGB: number
    totalGB: number
    percent: number
  }
  runtime: {
    goVersion: string
    goroutines: number
    numCPU: number
    allocMB: number
  }
  uptimeSeconds: number
}
