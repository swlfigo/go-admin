import { createRouter, createWebHistory, type RouteRecordRaw } from 'vue-router'

export const constantRoutes: RouteRecordRaw[] = [
  { path: '/login', name: 'login', component: () => import('@/views/login/index.vue'), meta: { public: true } },
  { path: '/403', name: '403', component: () => import('@/views/error/403.vue'), meta: { public: true } },
  { path: '/404', name: '404', component: () => import('@/views/error/404.vue'), meta: { public: true } },
  { path: '/500', name: '500', component: () => import('@/views/error/500.vue'), meta: { public: true } },
  // 渲染 404 组件而非 redirect 到公开的 /404——后者会在守卫前被解析，
  // 导致未登录访问 "/" 直接落到公开 404 而非跳登录。非 public，让守卫接管鉴权。
  { path: '/:pathMatch(.*)*', name: 'not-found', component: () => import('@/views/error/404.vue') },
]

const router = createRouter({
  history: createWebHistory(),
  routes: constantRoutes,
})

export default router
