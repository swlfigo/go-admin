import axios, {
  type AxiosInstance,
  type AxiosRequestConfig,
  type AxiosResponse,
  type InternalAxiosRequestConfig,
} from 'axios'
import type { ApiResult } from '@/types/api'

// Extended config type that allows skipping the auth-refresh interceptor.
// Used by auth endpoints (refreshToken, logout) to prevent the 401 self-await deadlock.
export interface RequestConfig extends AxiosRequestConfig {
  _skipAuthRefresh?: boolean
}

// 这两个回调由 user store 在初始化时注入，避免循环依赖。
let tokenGetter: () => string | null = () => null
let refreshHandler: () => Promise<string | null> = async () => null
let onAuthFail: () => void = () => {}

export function configureHttp(opts: {
  getToken: () => string | null
  refresh: () => Promise<string | null>
  onFail: () => void
}): void {
  tokenGetter = opts.getToken
  refreshHandler = opts.refresh
  onAuthFail = opts.onFail
}

/**
 * Pure predicate: should the interceptor attempt a token refresh for this response?
 *
 * Returns false (skip refresh) when:
 *  - status is not 401, OR
 *  - the request itself was an auth endpoint flagged with _skipAuthRefresh.
 *
 * This breaks the self-await deadlock: when the /auth/refresh call returns 401
 * (expired refresh token), _skipAuthRefresh is true, so we never await the
 * in-flight `refreshing` promise that is waiting on this very interceptor.
 */
export function shouldAttemptRefresh(
  status: number | undefined,
  config: RequestConfig | undefined,
): boolean {
  if (status !== 401) return false
  if (config?._skipAuthRefresh) return false
  return true
}

const instance: AxiosInstance = axios.create({ baseURL: '/api', timeout: 15000 })

instance.interceptors.request.use((config: InternalAxiosRequestConfig) => {
  const token = tokenGetter()
  if (token) config.headers.Authorization = `Bearer ${token}`
  return config
})

let refreshing: Promise<string | null> | null = null

instance.interceptors.response.use(
  (resp: AxiosResponse) => {
    const body = resp.data as ApiResult<unknown>
    // Return as AxiosResponse to satisfy interceptor type; request<T> casts at call site
    return body.data as AxiosResponse
  },
  async (error: unknown) => {
    const axiosError = error as {
      response?: { status?: number; data?: { msg?: string } }
      config?: RequestConfig & { _retried?: boolean }
      message?: string
    }
    const status = axiosError.response?.status
    const original = axiosError.config

    // If this is a flagged auth endpoint (refreshToken / logout) returning 401,
    // do NOT attempt refresh — that would deadlock by awaiting the promise whose
    // resolution depends on this interceptor returning. Instead, fail fast.
    if (!shouldAttemptRefresh(status, original)) {
      if (status === 401 && original?._skipAuthRefresh) {
        // Refresh token itself is expired/invalid → log out immediately.
        onAuthFail()
      }
      const msg = axiosError.response?.data?.msg ?? axiosError.message ?? '请求失败'
      return Promise.reject(new Error(msg))
    }

    if (original && !original._retried) {
      original._retried = true
      // 合并并发刷新
      if (!refreshing) refreshing = refreshHandler().finally(() => { refreshing = null })
      const newToken = await refreshing
      if (newToken) {
        original.headers = original.headers ?? {}
        ;(original.headers as Record<string, string>).Authorization = `Bearer ${newToken}`
        return instance(original)
      }
      onAuthFail()
    }
    const msg = axiosError.response?.data?.msg ?? axiosError.message ?? '请求失败'
    return Promise.reject(new Error(msg))
  },
)

// 泛型包装：调用处拿到的就是 data 的类型
export function request<T>(config: RequestConfig): Promise<T> {
  return instance(config) as unknown as Promise<T>
}
