import { defineStore } from 'pinia'
import type { RouteRecordRaw } from 'vue-router'
import type { MenuNode } from '@/types/api'

export interface FlatRoute {
  path: string
  name: string
  meta: { title: string; component: string; icon: string }
}

// 把菜单树拍平成"可渲染页面"的叶子路由（仅 type C，带组件路径）。
export function menusToRoutes(tree: MenuNode[]): FlatRoute[] {
  const out: FlatRoute[] = []
  const walk = (nodes: MenuNode[]): void => {
    for (const n of nodes) {
      if (n.type === 'C' && n.component) {
        out.push({
          path: n.path,
          name: n.path,
          meta: { title: n.name, component: n.component, icon: n.icon },
        })
      }
      if (n.children?.length) walk(n.children)
    }
  }
  walk(tree)
  return out
}

// 动态 import 所有页面组件（Plan 4 会补齐 views/**）。
const modules = import.meta.glob('@/views/**/*.vue')

export const usePermissionStore = defineStore('permission', {
  state: () => ({
    menus: [] as MenuNode[],
    routes: [] as RouteRecordRaw[],
    loaded: false,
  }),
  actions: {
    setMenus(tree: MenuNode[]) {
      this.menus = tree
      this.routes = menusToRoutes(tree).map((r) => {
        const loader = modules[`/src/views/${r.meta.component}.vue`]
        return {
          path: r.path,
          name: r.name,
          component: loader ?? (() => import('@/views/error/404.vue')),
          meta: r.meta,
        } as RouteRecordRaw
      })
      this.loaded = true
    },
    reset() {
      this.menus = []
      this.routes = []
      this.loaded = false
    },
  },
})
