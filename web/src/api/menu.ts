import { request } from '@/utils/http'
import type { MenuNode } from '@/types/api'

export function getMyMenus() {
  return request<MenuNode[]>({ url: '/menus/me', method: 'get' })
}
export function getMenuTree() {
  return request<MenuNode[]>({ url: '/menus', method: 'get' })
}
export function createMenu(body: Record<string, unknown>) {
  return request<MenuNode>({ url: '/menus', method: 'post', data: body })
}
export function updateMenu(id: number, body: Record<string, unknown>) {
  return request<null>({ url: `/menus/${id}`, method: 'put', data: body })
}
export function deleteMenu(id: number) {
  return request<null>({ url: `/menus/${id}`, method: 'delete' })
}
