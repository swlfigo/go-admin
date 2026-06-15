import { describe, it, expect } from 'vitest'
import { menusToRoutes } from './permission'
import type { MenuNode } from '@/types/api'

const tree: MenuNode[] = [
  { id: 1, parentId: 0, name: '工作台', type: 'C', path: '/dashboard', component: 'dashboard/index', perm: '', icon: 'dash', sort: 1 },
  { id: 2, parentId: 0, name: '系统管理', type: 'M', path: '/system', component: '', perm: '', icon: 'sys', sort: 2, children: [
    { id: 3, parentId: 2, name: '用户管理', type: 'C', path: '/system/user', component: 'system/user/index', perm: 'system:user:list', icon: '', sort: 1 },
  ] },
]

describe('menusToRoutes', () => {
  it('flattens C-type leaves into routes with component path', () => {
    const routes = menusToRoutes(tree)
    const paths = routes.map((r) => r.path)
    expect(paths).toContain('/dashboard')
    expect(paths).toContain('/system/user')
    // 纯目录 M 自身不产出叶子路由
    expect(paths).not.toContain('/system')
  })
  it('attaches component loader meta', () => {
    const routes = menusToRoutes(tree)
    const user = routes.find((r) => r.path === '/system/user')
    expect(user?.meta?.component).toBe('system/user/index')
    expect(user?.meta?.title).toBe('用户管理')
  })

  // Extra tests requested by user
  it('handles empty tree', () => {
    expect(menusToRoutes([])).toEqual([])
  })

  it('flattens a C leaf two levels deep (C under M under M)', () => {
    const deepTree: MenuNode[] = [
      {
        id: 10, parentId: 0, name: '根目录', type: 'M', path: '/root', component: '', perm: '', icon: '', sort: 1,
        children: [
          {
            id: 11, parentId: 10, name: '子目录', type: 'M', path: '/root/sub', component: '', perm: '', icon: '', sort: 1,
            children: [
              { id: 12, parentId: 11, name: '深叶页', type: 'C', path: '/root/sub/leaf', component: 'root/sub/leaf', perm: '', icon: '', sort: 1 },
            ],
          },
        ],
      },
    ]
    const routes = menusToRoutes(deepTree)
    const paths = routes.map((r) => r.path)
    expect(paths).toContain('/root/sub/leaf')
    expect(paths).not.toContain('/root')
    expect(paths).not.toContain('/root/sub')
    const leaf = routes.find((r) => r.path === '/root/sub/leaf')
    expect(leaf?.meta?.title).toBe('深叶页')
    expect(leaf?.meta?.component).toBe('root/sub/leaf')
  })
})
