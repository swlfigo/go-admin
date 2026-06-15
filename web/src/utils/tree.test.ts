import { describe, it, expect } from 'vitest'
import { toElTree, collectMenuIds, type ElTreeNode } from './tree'
import type { MenuNode } from '@/types/api'

const menus: MenuNode[] = [
  { id: 1, parentId: 0, name: '系统', type: 'M', path: '', component: '', perm: '', icon: '', sort: 1, children: [
    { id: 2, parentId: 1, name: '用户', type: 'C', path: '', component: '', perm: 'u:list', icon: '', sort: 1, children: [
      { id: 3, parentId: 2, name: '删除', type: 'F', path: '', component: '', perm: 'u:del', icon: '', sort: 1 },
    ] },
  ] },
]

describe('toElTree', () => {
  it('maps id/label/children recursively', () => {
    const tree: ElTreeNode[] = toElTree(menus)
    expect(tree[0].id).toBe(1)
    expect(tree[0].label).toBe('系统')
    expect(tree[0].children?.[0].label).toBe('用户')
    expect(tree[0].children?.[0].children?.[0].label).toBe('删除')
  })
})

describe('collectMenuIds', () => {
  it('returns all ids across all nesting levels', () => {
    const ids = collectMenuIds(menus)
    expect(ids).toContain(1)
    expect(ids).toContain(2)
    expect(ids).toContain(3)
    expect(ids).toHaveLength(3)
  })

  it('returns empty array for empty input', () => {
    expect(collectMenuIds([])).toEqual([])
  })

  it('handles flat (no children) menus', () => {
    const flat: MenuNode[] = [
      { id: 10, parentId: 0, name: 'A', type: 'F', path: '', component: '', perm: '', icon: '', sort: 1 },
      { id: 11, parentId: 0, name: 'B', type: 'F', path: '', component: '', perm: '', icon: '', sort: 2 },
    ]
    expect(collectMenuIds(flat)).toEqual([10, 11])
  })
})
