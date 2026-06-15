import { defineStore } from 'pinia'
import i18n, { type Locale, LOCALE_KEY } from '@/i18n'

type Theme = 'light' | 'dark'
type Accent = 'go' | 'blue' | 'green' | 'purple' | 'orange' | 'rose' | 'cyan'
type Density = 'comfortable' | 'compact'
type SidebarStyle = 'dark' | 'light'
export type TabStyle = 'card' | 'line' | 'plain'

function apply(theme: Theme, accent: Accent, density: Density, sidebar: SidebarStyle, tab: TabStyle) {
  const el = document.documentElement
  el.dataset.theme = theme
  el.dataset.accent = accent
  el.dataset.density = density
  el.dataset.sidebar = sidebar
  el.dataset.tabstyle = tab
}

export const useAppStore = defineStore('app', {
  state: () => ({
    theme: (localStorage.getItem('ga_theme') as Theme) || 'light',
    accent: (localStorage.getItem('ga_accent') as Accent) || 'go',
    density: (localStorage.getItem('ga_density') as Density) || 'comfortable',
    sidebarStyle: (localStorage.getItem('ga_sidebar') as SidebarStyle) || 'dark',
    tabStyle: (localStorage.getItem('ga_tabstyle') as TabStyle) || 'line',
    locale: (localStorage.getItem(LOCALE_KEY) as Locale) || 'zh-CN',
    sidebarCollapsed: false,
  }),
  actions: {
    init() { apply(this.theme, this.accent, this.density, this.sidebarStyle, this.tabStyle) },
    setTheme(t: Theme) { this.theme = t; localStorage.setItem('ga_theme', t); this.init() },
    setAccent(a: Accent) { this.accent = a; localStorage.setItem('ga_accent', a); this.init() },
    setDensity(d: Density) { this.density = d; localStorage.setItem('ga_density', d); this.init() },
    setSidebarStyle(s: SidebarStyle) { this.sidebarStyle = s; localStorage.setItem('ga_sidebar', s); this.init() },
    setTabStyle(t: TabStyle) { this.tabStyle = t; localStorage.setItem('ga_tabstyle', t); this.init() },
    setLocale(l: Locale) { this.locale = l; localStorage.setItem(LOCALE_KEY, l); i18n.global.locale.value = l },
    toggleSidebar() { this.sidebarCollapsed = !this.sidebarCollapsed },
  },
})
