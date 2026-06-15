import { request } from '@/utils/http'
import type { DictData, DictType, Paged } from '@/types/api'

export function listDictTypes(params: { keyword?: string; page: number; size: number }) {
  return request<Paged<DictType>>({ url: '/dict/types', method: 'get', params })
}
export function createDictType(body: Record<string, unknown>) {
  return request<DictType>({ url: '/dict/types', method: 'post', data: body })
}
export function updateDictType(id: number, body: Record<string, unknown>) {
  return request<null>({ url: `/dict/types/${id}`, method: 'put', data: body })
}
export function deleteDictType(id: number) {
  return request<null>({ url: `/dict/types/${id}`, method: 'delete' })
}
export function listDictData(type: string) {
  return request<DictData[]>({ url: '/dict/data', method: 'get', params: { type } })
}
export function createDictData(body: Record<string, unknown>) {
  return request<DictData>({ url: '/dict/data', method: 'post', data: body })
}
export function updateDictData(id: number, body: Record<string, unknown>) {
  return request<null>({ url: `/dict/data/${id}`, method: 'put', data: body })
}
export function deleteDictData(id: number) {
  return request<null>({ url: `/dict/data/${id}`, method: 'delete' })
}
