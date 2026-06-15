import { request } from '@/utils/http'
import type { OnlineSession } from '@/types/api'

export function listOnline() {
  return request<OnlineSession[]>({ url: '/online', method: 'get' })
}
export function kickOnline(id: string) {
  return request<null>({ url: `/online/${id}`, method: 'delete' })
}
