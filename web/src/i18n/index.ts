import { createI18n } from 'vue-i18n'
import zhCN from './locales/zh-CN'
import enUS from './locales/en-US'

export type Locale = 'zh-CN' | 'en-US'

export const LOCALE_KEY = 'ga_locale'

function initialLocale(): Locale {
  const saved = localStorage.getItem(LOCALE_KEY)
  return saved === 'en-US' ? 'en-US' : 'zh-CN'
}

const i18n = createI18n({
  legacy: false, // Composition API 模式
  locale: initialLocale(),
  fallbackLocale: 'zh-CN',
  messages: { 'zh-CN': zhCN, 'en-US': enUS },
})

export default i18n
