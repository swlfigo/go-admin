import { request } from '@/utils/http'
import type { RequestConfig } from '@/utils/http'
import type { CaptchaResult, LoginResult, MeResult } from '@/types/api'

export function getCaptcha() {
  return request<CaptchaResult>({ url: '/captcha', method: 'get' })
}

export function login(body: {
  username: string
  password: string
  captchaId?: string
  captchaCode?: string
}) {
  return request<LoginResult>({ url: '/auth/login', method: 'post', data: body })
}

export function refreshToken(token: string) {
  const config: RequestConfig = {
    url: '/auth/refresh',
    method: 'post',
    data: { refreshToken: token },
    _skipAuthRefresh: true,
  }
  return request<LoginResult>(config)
}

export function logout() {
  const config: RequestConfig = {
    url: '/auth/logout',
    method: 'post',
    _skipAuthRefresh: true,
  }
  return request<null>(config)
}

export function getMe() {
  return request<MeResult>({ url: '/auth/me', method: 'get' })
}
