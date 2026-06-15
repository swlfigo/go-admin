import { request } from '@/utils/http'

export function updateProfile(body: { nickname: string; email: string; phone: string }) {
  return request<null>({ url: '/profile', method: 'put', data: body })
}
export function changePassword(body: { oldPassword: string; newPassword: string }) {
  return request<null>({ url: '/profile/password', method: 'put', data: body })
}
