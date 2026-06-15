import type { App, DirectiveBinding } from 'vue'
import { useUserStore } from '@/stores/user'

// 用法：v-permission="'system:user:delete'"  无权限则移除元素
export function setupPermissionDirective(app: App): void {
  app.directive('permission', {
    mounted(el: HTMLElement, binding: DirectiveBinding<string>) {
      const user = useUserStore()
      const code = binding.value
      // 超管或拥有该权限码放行；否则移除
      const isSuper = user.user?.roles?.some((r) => r.code === 'super_admin')
      if (!isSuper && code && !user.hasPerm(code)) {
        el.parentNode?.removeChild(el)
      }
    },
  })
}
