import { request } from '@/utils/http'
import type { Paged, Role } from '@/types/api'

export function listRoles(params: { keyword?: string; page: number; size: number }) {
  return request<Paged<Role>>({ url: '/roles', method: 'get', params })
}
export function getRole(id: number) {
  return request<Role & { menus: { id: number }[] }>({ url: `/roles/${id}`, method: 'get' })
}
export function createRole(body: Record<string, unknown>) {
  return request<Role>({ url: '/roles', method: 'post', data: body })
}
export function updateRole(id: number, body: Record<string, unknown>) {
  return request<null>({ url: `/roles/${id}`, method: 'put', data: body })
}
export function deleteRole(id: number) {
  return request<null>({ url: `/roles/${id}`, method: 'delete' })
}
export function assignRoleMenus(id: number, menuIds: number[]) {
  return request<null>({ url: `/roles/${id}/menus`, method: 'put', data: { menuIds } })
}
