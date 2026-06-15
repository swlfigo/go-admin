import { request } from '@/utils/http'
import type { OperationLog, Paged } from '@/types/api'

export function listOperationLogs(params: { keyword?: string; page: number; size: number }) {
  return request<Paged<OperationLog>>({ url: '/logs/operation', method: 'get', params })
}

export function clearOperationLogs() {
  return request<null>({ url: '/logs/operation', method: 'delete' })
}
