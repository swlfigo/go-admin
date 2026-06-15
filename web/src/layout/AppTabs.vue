<script setup lang="ts">
import { onMounted, onUnmounted, reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { usePermissionStore } from '@/stores/permission'
import { useTabsStore } from '@/stores/tabs'
import { titleByPath } from '@/i18n/menu'
import type { MenuNode } from '@/types/api'
import type { TabItem } from '@/stores/tabs'
import Icon from '@/components/Icon.vue'

const { t } = useI18n()

const route = useRoute()
const router = useRouter()
const perm = usePermissionStore()
const tabs = useTabsStore()

// 不为这些路由建立 tab。
const SKIP_PATHS = new Set(['/login', '/403', '/404', '/500'])

function titleFromMenus(path: string): string | undefined {
  for (const top of perm.menus) {
    if (top.path === path) return top.name
    const child = (top.children ?? []).find((c: MenuNode) => c.path === path)
    if (child) return child.name
  }
  return undefined
}

function resolveTitle(): string | undefined {
  const metaTitle = route.meta?.title
  if (typeof metaTitle === 'string' && metaTitle) return metaTitle
  return titleFromMenus(route.path)
}

watch(
  () => route.path,
  (path) => {
    if (SKIP_PATHS.has(path)) return
    const title = resolveTitle()
    if (!title) return
    tabs.addTab({ path, title, pinned: path === '/dashboard' })
  },
  { immediate: true },
)

function onTabClick(tab: TabItem) {
  if (tab.path !== route.path) router.push(tab.path)
}

// 关闭/批量关闭后：若当前激活页签已被关掉，导航到目标路径
function navigateIfClosed(next: string) {
  if (!tabs.tabs.some((t) => t.path === route.path)) router.push(next)
}

function doClose(tab: TabItem) {
  navigateIfClosed(tabs.removeTab(tab.path))
}

// ---- 右键上下文菜单 ----
const menu = reactive({ open: false, x: 0, y: 0, tab: null as TabItem | null })

function openMenu(e: MouseEvent, tab: TabItem) {
  menu.open = true
  menu.x = e.clientX
  menu.y = e.clientY
  menu.tab = tab
}
function closeMenu() {
  menu.open = false
  menu.tab = null
}
function act(fn: (path: string) => string) {
  if (!menu.tab) return
  navigateIfClosed(fn(menu.tab.path))
  closeMenu()
}
function actPin() {
  if (menu.tab) tabs.togglePin(menu.tab.path)
  closeMenu()
}

onMounted(() => window.addEventListener('click', closeMenu))
onUnmounted(() => window.removeEventListener('click', closeMenu))
</script>

<template>
  <nav class="tabstrip scroll">
    <button
      v-for="tab in tabs.tabs"
      :key="tab.path"
      class="tab"
      :class="{ 'tab--active': route.path === tab.path, 'tab--pinned': tab.pinned }"
      type="button"
      @click="onTabClick(tab)"
      @contextmenu.prevent="openMenu($event, tab)"
    >
      <Icon v-if="tab.pinned" name="pin" :size="12" class="tab__pin" />
      <span class="tab__label">{{ titleByPath(tab.path, tab.title) }}</span>
      <span v-if="!tab.pinned" class="tab__close" :title="t('shell.closeTab')" @click.stop="doClose(tab)">
        <Icon name="x" :size="13" />
      </span>
    </button>
  </nav>

  <!-- 右键菜单 -->
  <div
    v-if="menu.open && menu.tab"
    class="ctxmenu"
    :style="{ left: menu.x + 'px', top: menu.y + 'px' }"
    @click.stop
  >
    <button class="ctxmenu__item" :disabled="menu.tab.path === '/dashboard'" @click="actPin">
      <Icon name="pin" :size="14" />{{ menu.tab.pinned ? t('shell.unpinTab') : t('shell.pinTab') }}
    </button>
    <button class="ctxmenu__item ctxmenu__item--danger" :disabled="menu.tab.pinned" @click="act(tabs.removeTab)">
      <Icon name="x" :size="14" />{{ t('shell.closeTab') }}
    </button>
    <div class="ctxmenu__sep" />
    <button class="ctxmenu__item" @click="act(tabs.closeLeft)">{{ t('shell.closeLeft') }}</button>
    <button class="ctxmenu__item" @click="act(tabs.closeRight)">{{ t('shell.closeRight') }}</button>
    <button class="ctxmenu__item" @click="act(tabs.closeOthers)">{{ t('shell.closeOthers') }}</button>
    <button class="ctxmenu__item" @click="act(() => tabs.closeAll())">{{ t('shell.closeAll') }}</button>
  </div>
</template>
