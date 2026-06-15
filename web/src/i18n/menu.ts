import i18n from './index'

// 动态菜单名来自后端（中文）。这里按路由 path 映射到 i18n key，
// 命中则翻译显示，未命中则回退后端原名（用户自建的菜单保持原样）。
const PATH_KEY: Record<string, string> = {
  '/dashboard': 'menu.dashboard',
  '/system': 'menu.system',
  '/system/user': 'menu.user',
  '/system/role': 'menu.role',
  '/system/menu': 'menu.menu',
  '/system/dict': 'menu.dict',
  '/monitor/online': 'menu.online',
  '/monitor/log': 'menu.log',
  '/monitor/server': 'menu.server',
  '/profile': 'menu.profile',
}

// 按 path 取标题：命中映射则翻译，否则回退给定的后端原名。
export function titleByPath(path: string, fallback: string): string {
  const key = PATH_KEY[path]
  return key ? i18n.global.t(key) : fallback
}

export function menuTitle(node: { path: string; name: string }): string {
  return titleByPath(node.path, node.name)
}
