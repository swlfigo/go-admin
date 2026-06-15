import { request } from '@/utils/http'
import type { ServerInfo } from '@/types/api'

export function getServerInfo() {
  return request<ServerInfo>({ url: '/monitor/server', method: 'get' })
}
