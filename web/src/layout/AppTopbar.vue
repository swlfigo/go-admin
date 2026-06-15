<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'
import { usePermissionStore } from '@/stores/permission'
import { titleByPath } from '@/i18n/menu'
import type { MenuNode } from '@/types/api'
import UserAvatar from '@/components/UserAvatar.vue'
import Icon from '@/components/Icon.vue'

const emit = defineEmits<{ openSettings: [] }>()

const { t } = useI18n()
const user = useUserStore()
const app = useAppStore()
const perm = usePermissionStore()
const route = useRoute()
const router = useRouter()

const displayName = computed(() => user.user?.nickname || user.user?.username || '')

// 由当前路由推导面包屑：返回原始 path+name，模板里按 path 翻译
interface Crumb {
  parent?: { path: string; name: string }
  current: { path: string; name: string }
}
const crumb = computed<Crumb>(() => {
  const path = route.path
  for (const top of perm.menus) {
    if (top.path === path) return { current: { path: top.path, name: top.name } }
    const child = (top.children || []).find((c: MenuNode) => c.path === path)
    if (child) return { parent: { path: top.path, name: top.name }, current: { path: child.path, name: child.name } }
  }
  const title = (route.meta?.title as string | undefined) || ''
  return { current: { path, name: title } }
})

function toggleTheme() {
  app.setTheme(app.theme === 'light' ? 'dark' : 'light')
}
async function onLogout() {
  await user.logout()
  router.push('/login')
}
function onCommand(c: string) {
  if (c === 'logout') onLogout()
  else if (c === 'profile') router.push('/profile')
}
</script>

<template>
  <header class="topbar">
    <button class="topbar__toggle" :title="t('shell.toggleSidebar')" @click="app.toggleSidebar()">
      <Icon name="panel" :size="19" />
    </button>
    <div class="topdiv" />

    <div class="breadcrumb">
      <span v-if="crumb.parent">{{ titleByPath(crumb.parent.path, crumb.parent.name) }}</span>
      <span v-if="crumb.parent" class="sep"><Icon name="chevron" :size="13" /></span>
      <b>{{ titleByPath(crumb.current.path, crumb.current.name) }}</b>
    </div>

    <div class="topbar__spacer" />

    <div class="topbar__right">
      <button class="iconbtn" :title="t('shell.toggleTheme')" @click="toggleTheme">
        <Icon :name="app.theme === 'dark' ? 'sun' : 'moon'" :size="18" />
      </button>
      <button class="iconbtn" :title="t('shell.themeSettings')" @click="emit('openSettings')">
        <Icon name="palette" :size="18" />
      </button>
      <div class="topdiv" />
      <el-dropdown trigger="click" @command="onCommand">
        <button class="iconbtn" style="width: auto; padding: 0 4px" :title="t('shell.profile')">
          <UserAvatar :name="displayName" :src="user.user?.avatar" :size="30" />
        </button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="profile">{{ t('shell.profile') }}</el-dropdown-item>
            <el-dropdown-item command="logout">{{ t('shell.logout') }}</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </header>
</template>
