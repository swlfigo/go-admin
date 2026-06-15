import { defineStore } from 'pinia'

export interface TabItem {
  path: string
  title: string
  pinned: boolean
}

const STORAGE_KEY = 'ga_tabs'
const DASHBOARD: TabItem = { path: '/dashboard', title: '工作台', pinned: true }

function loadTabs(): TabItem[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (raw) {
      const parsed = JSON.parse(raw) as TabItem[]
      if (Array.isArray(parsed) && parsed.length) {
        // 工作台恒置首且固定
        const rest = parsed.filter((t) => t.path !== '/dashboard')
        return [{ ...DASHBOARD }, ...rest]
      }
    }
  } catch {
    /* ignore corrupt storage */
  }
  return [{ ...DASHBOARD }]
}

export const useTabsStore = defineStore('tabs', {
  state: () => ({
    tabs: loadTabs(),
  }),
  actions: {
    persist() {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(this.tabs))
    },
    // 固定的排前面（保持各自相对顺序）
    reorder() {
      this.tabs = [...this.tabs.filter((t) => t.pinned), ...this.tabs.filter((t) => !t.pinned)]
    },
    addTab(t: TabItem) {
      if (this.tabs.some((x) => x.path === t.path)) return
      this.tabs.push(t)
      this.persist()
    },
    // 移除一个 tab（固定的不移除），返回移除后应导航到的路径。
    removeTab(path: string): string {
      const idx = this.tabs.findIndex((x) => x.path === path)
      if (idx === -1 || this.tabs[idx].pinned) return path
      this.tabs.splice(idx, 1)
      this.persist()
      return this.tabs[Math.min(idx, this.tabs.length - 1)]?.path ?? '/dashboard'
    },
    togglePin(path: string) {
      const t = this.tabs.find((x) => x.path === path)
      if (!t || t.path === '/dashboard') return // 工作台恒固定
      t.pinned = !t.pinned
      this.reorder()
      this.persist()
    },
    // 以下批量关闭都保留固定的 tab，返回导航目标
    closeOthers(path: string): string {
      this.tabs = this.tabs.filter((t) => t.pinned || t.path === path)
      this.persist()
      return this.tabs.some((t) => t.path === path) ? path : '/dashboard'
    },
    closeLeft(path: string): string {
      const idx = this.tabs.findIndex((t) => t.path === path)
      if (idx < 0) return path
      this.tabs = this.tabs.filter((t, i) => t.pinned || i >= idx)
      this.persist()
      return path
    },
    closeRight(path: string): string {
      const idx = this.tabs.findIndex((t) => t.path === path)
      if (idx < 0) return path
      this.tabs = this.tabs.filter((t, i) => t.pinned || i <= idx)
      this.persist()
      return path
    },
    closeAll(): string {
      this.tabs = this.tabs.filter((t) => t.pinned)
      if (!this.tabs.length) this.tabs = [{ ...DASHBOARD }]
      this.persist()
      return this.tabs[this.tabs.length - 1]?.path ?? '/dashboard'
    },
    reset() {
      this.tabs = [{ ...DASHBOARD }]
      this.persist()
    },
  },
})
