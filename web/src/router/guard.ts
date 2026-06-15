import type { Router } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { usePermissionStore } from '@/stores/permission'
import { useAppStore } from '@/stores/app'
import { getMyMenus } from '@/api/menu'
import AppLayout from '@/layout/AppLayout.vue'

export function setupGuard(router: Router): void {
  router.beforeEach(async (to) => {
    const user = useUserStore()
    const perm = usePermissionStore()
    useAppStore().init()

    if (to.meta.public) return true
    if (!user.isLoggedIn) return { path: '/login', query: { redirect: to.fullPath } }

    // 已登录但未拉过菜单 → 拉 me + 菜单，动态挂业务路由
    if (!perm.loaded) {
      try {
        await user.loadMe()
        const menus = await getMyMenus()
        perm.setMenus(menus)
        // 业务路由挂到 Layout 下
        router.addRoute({
          path: '/',
          component: AppLayout,
          redirect: '/dashboard',
          children: [
            ...perm.routes,
            { path: '/profile', name: 'profile', component: () => import('@/views/profile/index.vue'), meta: { title: '个人中心' } },
          ],
        })
        // 按 path 重新进入以命中新挂的路由。不能用 {...to}——首次 to 命中的是
        // catch-all（name: 'not-found'），展开会带上该 name，导致按 name 又跳回 404。
        return { path: to.fullPath, replace: true }
      } catch {
        user.clear()
        return { path: '/login' }
      }
    }
    return true
  })
}
