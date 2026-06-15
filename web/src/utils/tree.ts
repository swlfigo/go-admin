import type { MenuNode } from '@/types/api'

export interface ElTreeNode {
  id: number
  label: string
  children?: ElTreeNode[]
}

// 菜单树 → Element Plus el-tree 数据（id/label/children）
export function toElTree(menus: MenuNode[]): ElTreeNode[] {
  return menus.map((m) => ({
    id: m.id,
    label: m.name,
    children: m.children?.length ? toElTree(m.children) : undefined,
  }))
}

// 把扁平菜单（含 children）的全部 id 收集出来（用于"全选"或回显基准）
export function collectMenuIds(menus: MenuNode[]): number[] {
  const ids: number[] = []
  const walk = (ns: MenuNode[]) => {
    for (const n of ns) {
      ids.push(n.id)
      if (n.children?.length) walk(n.children)
    }
  }
  walk(menus)
  return ids
}
