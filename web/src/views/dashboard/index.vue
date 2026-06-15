<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/stores/user'
import { listUsers } from '@/api/user'
import { listRoles } from '@/api/role'
import { getMenuTree } from '@/api/menu'
import { listOnline } from '@/api/online'
import type { MenuNode } from '@/types/api'
import type { IconName } from '@/components/Icon.vue'
import Icon from '@/components/Icon.vue'

const { t } = useI18n()
const user = useUserStore()

interface StatCard {
  key: string
  labelKey: string
  icon: IconName
  tint: number // hue for the icon tile tint
  value: string // '…' loading, '—' error, or a number string
}

const cards = ref<StatCard[]>([
  { key: 'users', labelKey: 'dashboard.statUsers', icon: 'users', tint: 255, value: '…' },
  { key: 'roles', labelKey: 'dashboard.statRoles', icon: 'shield', tint: 150, value: '…' },
  { key: 'menus', labelKey: 'dashboard.statMenus', icon: 'list', tint: 30, value: '…' },
  { key: 'online', labelKey: 'dashboard.statOnline', icon: 'activity', tint: 320, value: '…' },
])

function setValue(key: string, value: string) {
  const card = cards.value.find((c) => c.key === key)
  if (card) card.value = value
}

function icoStyle(tint: number): Record<string, string> {
  return {
    background: `color-mix(in oklab, oklch(0.62 0.18 ${tint}) 14%, transparent)`,
    color: `oklch(0.55 0.18 ${tint})`,
  }
}

function countNodes(nodes: MenuNode[]): number {
  let total = 0
  for (const n of nodes) {
    total += 1
    if (n.children?.length) total += countNodes(n.children)
  }
  return total
}

async function loadUsers() {
  try {
    const res = await listUsers({ page: 1, size: 1 })
    setValue('users', String(res.total))
  } catch {
    setValue('users', '—')
  }
}
async function loadRoles() {
  try {
    const res = await listRoles({ page: 1, size: 1 })
    setValue('roles', String(res.total))
  } catch {
    setValue('roles', '—')
  }
}
async function loadMenus() {
  try {
    const tree = await getMenuTree()
    setValue('menus', String(countNodes(tree)))
  } catch {
    setValue('menus', '—')
  }
}
async function loadOnline() {
  try {
    const list = await listOnline()
    setValue('online', String(list.length))
  } catch {
    setValue('online', '—')
  }
}

onMounted(() => {
  void loadUsers()
  void loadRoles()
  void loadMenus()
  void loadOnline()
})
</script>

<template>
  <div class="page">
    <div class="page__head">
      <div>
        <h1 class="page__title">{{ t('dashboard.title') }}</h1>
        <p class="page__sub">
          {{ t('dashboard.welcome', { name: user.user?.nickname || user.user?.username || '' }) }}
        </p>
      </div>
    </div>

    <div class="stats">
      <div v-for="c in cards" :key="c.key" class="statcard">
        <div class="statcard__top">
          <div class="statcard__ico" :style="icoStyle(c.tint)">
            <Icon :name="c.icon" :size="20" />
          </div>
        </div>
        <div class="statcard__label">{{ t(c.labelKey) }}</div>
        <div class="statcard__value">{{ c.value }}</div>
      </div>
    </div>
  </div>
</template>
