<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useAppStore, type TabStyle } from '@/stores/app'
import type { Locale } from '@/i18n'
import Icon from '@/components/Icon.vue'

defineProps<{ open: boolean }>()
const emit = defineEmits<{ close: [] }>()

const { t } = useI18n()
const app = useAppStore()

type Accent = 'go' | 'blue' | 'green' | 'purple' | 'orange' | 'rose' | 'cyan'

const TAB_STYLES: { key: TabStyle; labelKey: string }[] = [
  { key: 'card', labelKey: 'settings.tabCard' },
  { key: 'line', labelKey: 'settings.tabLine' },
  { key: 'plain', labelKey: 'settings.tabPlain' },
]

// 每个强调色对应的 --primary（与 tokens.css 一致），用于色板预览
const ACCENTS: { key: Accent; color: string }[] = [
  { key: 'go', color: '#00ADD8' },
  { key: 'blue', color: 'oklch(0.56 0.16 256)' },
  { key: 'green', color: 'oklch(0.58 0.14 159)' },
  { key: 'purple', color: 'oklch(0.55 0.18 300)' },
  { key: 'orange', color: 'oklch(0.66 0.16 54)' },
  { key: 'rose', color: 'oklch(0.6 0.18 13)' },
  { key: 'cyan', color: 'oklch(0.6 0.12 215)' },
]

const LANGS: { key: Locale; label: string }[] = [
  { key: 'zh-CN', label: '中文' },
  { key: 'en-US', label: 'English' },
]
</script>

<template>
  <template v-if="open">
    <div class="drawer-overlay" @click="emit('close')" />
    <div class="drawer scroll">
      <div class="drawer__head">
        <h3><Icon name="palette" :size="18" style="color: var(--primary)" /> {{ t('settings.title') }}</h3>
        <button class="iconbtn" @click="emit('close')"><Icon name="x" :size="18" /></button>
      </div>

      <div class="drawer__body scroll">
        <!-- 语言 -->
        <div class="drawer__sec">
          <h4>{{ t('settings.language') }}</h4>
          <div class="seg">
            <button
              v-for="l in LANGS"
              :key="l.key"
              :class="{ on: app.locale === l.key }"
              @click="app.setLocale(l.key)"
            >{{ l.label }}</button>
          </div>
        </div>

        <!-- 外观模式 -->
        <div class="drawer__sec">
          <h4>{{ t('settings.appearance') }}</h4>
          <div class="seg">
            <button :class="{ on: app.theme === 'light' }" @click="app.setTheme('light')">
              <Icon name="sun" :size="15" />{{ t('settings.light') }}
            </button>
            <button :class="{ on: app.theme === 'dark' }" @click="app.setTheme('dark')">
              <Icon name="moon" :size="15" />{{ t('settings.dark') }}
            </button>
          </div>
        </div>

        <!-- 强调色 -->
        <div class="drawer__sec">
          <h4>{{ t('settings.accent') }}</h4>
          <div class="swatches">
            <button
              v-for="a in ACCENTS"
              :key="a.key"
              class="swatch"
              :class="{ on: app.accent === a.key }"
              :style="{ background: a.color }"
              :title="t('settings.accents.' + a.key)"
              @click="app.setAccent(a.key)"
            />
          </div>
          <div style="margin-top: 10px; font-size: 12px; color: var(--fg-muted); display: flex; justify-content: space-between">
            <span>{{ t('settings.current') }}</span>
            <b style="color: var(--primary)">{{ t('settings.accents.' + app.accent) }}</b>
          </div>
        </div>

        <!-- 侧边栏 -->
        <div class="drawer__sec">
          <h4>{{ t('settings.sidebar') }}</h4>
          <div class="seg">
            <button :class="{ on: app.sidebarStyle === 'dark' }" @click="app.setSidebarStyle('dark')">{{ t('settings.sidebarDark') }}</button>
            <button :class="{ on: app.sidebarStyle === 'light' }" @click="app.setSidebarStyle('light')">{{ t('settings.sidebarLight') }}</button>
          </div>
          <div style="height: 10px" />
          <div class="optrow">
            <div class="optrow__label">{{ t('settings.collapse') }}</div>
            <button
              class="toggle"
              :class="{ on: app.sidebarCollapsed }"
              role="switch"
              :aria-checked="app.sidebarCollapsed"
              @click="app.toggleSidebar()"
            />
          </div>
        </div>

        <!-- 标签页样式 -->
        <div class="drawer__sec">
          <h4>{{ t('settings.tabs') }}</h4>
          <div class="tabopt-grid">
            <button
              v-for="tab in TAB_STYLES"
              :key="tab.key"
              class="tabopt"
              :class="{ on: app.tabStyle === tab.key }"
              @click="app.setTabStyle(tab.key)"
            >
              <span class="tabopt__preview" :class="'tabopt__preview--' + tab.key">
                <i class="b b1" /><i class="b" /><i class="b" />
              </span>
              <span class="tabopt__name">{{ t(tab.labelKey) }}</span>
            </button>
          </div>
        </div>

        <!-- 界面密度 -->
        <div class="drawer__sec">
          <h4>{{ t('settings.density') }}</h4>
          <div class="seg">
            <button :class="{ on: app.density === 'comfortable' }" @click="app.setDensity('comfortable')">{{ t('settings.comfortable') }}</button>
            <button :class="{ on: app.density === 'compact' }" @click="app.setDensity('compact')">{{ t('settings.compact') }}</button>
          </div>
        </div>
      </div>
    </div>
  </template>
</template>
