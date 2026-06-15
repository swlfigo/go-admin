import { request } from '@/utils/http'
import type { Paged, User } from '@/types/api'

export function listUsers(params: { keyword?: string; page: number; size: number }) {
  return request<Paged<User>>({ url: '/users', method: 'get', params })
}
export function createUser(body: Record<string, unknown>) {
  return request<User>({ url: '/users', method: 'post', data: body })
}
export function updateUser(id: number, body: Record<string, unknown>) {
  return request<null>({ url: `/users/${id}`, method: 'put', data: body })
}
export function deleteUser(id: number) {
  return request<null>({ url: `/users/${id}`, method: 'delete' })
}
export function assignUserRoles(id: number, roleIds: number[]) {
  return request<null>({ url: `/users/${id}/roles`, method: 'put', data: { roleIds } })
}
export function resetUserPassword(id: number, password: string) {
  return request<null>({ url: `/users/${id}/password`, method: 'put', data: { password } })
}
export function unlockUser(id: number) {
  return request<null>({ url: `/users/${id}/unlock`, method: 'put' })
}
