import { defineStore } from 'pinia'
import type { User } from '@/types/api'
import * as authApi from '@/api/auth'

const ACCESS_KEY = 'ga_access'
const REFRESH_KEY = 'ga_refresh'

export const useUserStore = defineStore('user', {
  state: () => ({
    accessToken: localStorage.getItem(ACCESS_KEY),
    refreshToken: localStorage.getItem(REFRESH_KEY),
    user: null as User | null,
    perms: [] as string[],
  }),
  getters: {
    isLoggedIn: (s) => !!s.accessToken,
    hasPerm: (s) => (code: string) => s.perms.includes(code),
  },
  actions: {
    setTokens(access: string, refresh: string) {
      this.accessToken = access
      this.refreshToken = refresh
      localStorage.setItem(ACCESS_KEY, access)
      localStorage.setItem(REFRESH_KEY, refresh)
    },
    async login(payload: { username: string; password: string; captchaId?: string; captchaCode?: string }) {
      const res = await authApi.login(payload)
      this.setTokens(res.accessToken, res.refreshToken)
    },
    async loadMe() {
      const me = await authApi.getMe()
      this.user = me.user
      this.perms = me.perms
    },
    async doRefresh(): Promise<string | null> {
      if (!this.refreshToken) return null
      try {
        const res = await authApi.refreshToken(this.refreshToken)
        this.setTokens(res.accessToken, res.refreshToken)
        return res.accessToken
      } catch {
        return null
      }
    },
    clear() {
      this.accessToken = null
      this.refreshToken = null
      this.user = null
      this.perms = []
      localStorage.removeItem(ACCESS_KEY)
      localStorage.removeItem(REFRESH_KEY)
    },
    async logout() {
      try { await authApi.logout() } catch { /* ignore */ }
      this.clear()
    },
  },
})
