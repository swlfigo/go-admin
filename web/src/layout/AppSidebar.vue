<script setup lang="ts">
import { computed, reactive, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { usePermissionStore } from '@/stores/permission'
import { menuTitle } from '@/i18n/menu'
import type { MenuNode } from '@/types/api'
import Icon, { type IconName } from '@/components/Icon.vue'

const { t } = useI18n()

const perm = usePermissionStore()
const route = useRoute()
const router = useRouter()

// 仅渲染可见节点（排除按钮 type F）
const tree = computed(() => perm.menus.filter((m) => m.type !== 'F'))

// name / path → 图标 的小型查表，找不到时回退到通用图标
function iconFor(node: MenuNode): IconName {
  const byName: Record<string, IconName> = {
    工作台: 'grid',
    系统管理: 'settings',
    用户管理: 'users',
    角色管理: 'shield',
    菜单管理: 'list',
    字典管理: 'book',
    在线用户: 'activity',
  }
  if (byName[node.name]) return byName[node.name]
  const p = node.path || ''
  if (p === '/dashboard') return 'grid'
  if (p === '/system') return 'settings'
  if (p.includes('/user')) return 'users'
  if (p.includes('/role')) return 'shield'
  if (p.includes('/menu')) return 'list'
  if (p.includes('/dict')) return 'book'
  if (p.includes('/monitor/log')) return 'book'
  if (p.includes('/monitor/server')) return 'activity'
  if (p.includes('/monitor')) return 'activity'
  return 'grid'
}

function childrenOf(node: MenuNode): MenuNode[] {
  return (node.children || []).filter((c) => c.type !== 'F')
}

function isActive(path: string): boolean {
  return !!path && route.path === path
}

function groupActive(node: MenuNode): boolean {
  return childrenOf(node).some((c) => isActive(c.path))
}

// 展开状态：默认把含当前路由的分组展开
const open = reactive<Record<number, boolean>>({})
function syncOpen() {
  for (const node of tree.value) {
    if (node.type === 'M' && groupActive(node) && open[node.id] === undefined) {
      open[node.id] = true
    }
  }
}
watch(() => [perm.menus, route.path], syncOpen, { immediate: true })

function toggleGroup(id: number) {
  open[id] = !open[id]
}

function go(path: string) {
  if (path) router.push(path)
}
</script>

<template>
  <aside class="sidebar">
    <div class="sidebar__logo">
      <img class="sidebar__mark" src="/logo.svg" alt="Go Admin" width="32" height="32" />
      <div class="sidebar__name">Go Admin<small>{{ t('shell.brandTagline') }}</small></div>
    </div>

    <nav class="sidebar__nav scroll">
      <div class="navsection">
        <div class="navsection__title">{{ t('shell.nav') }}</div>

        <template v-for="node in tree" :key="node.id">
          <!-- M 型目录：可展开分组 -->
          <div
            v-if="node.type === 'M'"
            class="navgroup"
            :class="{ 'navgroup--open': open[node.id], 'navgroup--active': groupActive(node) }"
          >
            <button class="navgroup__header" :title="menuTitle(node)" @click="toggleGroup(node.id)">
              <span class="navgroup__icon"><Icon :name="iconFor(node)" :size="20" /></span>
              <span class="navgroup__label">{{ menuTitle(node) }}</span>
              <span class="navgroup__chevron"><Icon name="chevron" :size="15" /></span>
            </button>

            <div class="navgroup__items">
              <div>
                <div class="navgroup__inner">
                  <button
                    v-for="child in childrenOf(node)"
                    :key="child.id"
                    class="navitem"
                    :class="{ 'navitem--active': isActive(child.path) }"
                    @click="go(child.path)"
                  >
                    <span class="navitem__label">{{ menuTitle(child) }}</span>
                  </button>
                </div>
              </div>
            </div>

            <!-- 收起态悬浮菜单 -->
            <div class="navgroup__flyout">
              <div class="navgroup__flyout-title">{{ menuTitle(node) }}</div>
              <button
                v-for="child in childrenOf(node)"
                :key="child.id"
                class="navitem"
                :class="{ 'navitem--active': isActive(child.path) }"
                @click="go(child.path)"
              >
                <span class="navitem__label">{{ menuTitle(child) }}</span>
              </button>
            </div>
          </div>

          <!-- 顶层 C 型菜单：直接导航 -->
          <button
            v-else
            class="navitem navitem--top"
            :class="{ 'navitem--active': isActive(node.path) }"
            :title="menuTitle(node)"
            @click="go(node.path)"
          >
            <span class="navitem__icon"><Icon :name="iconFor(node)" :size="20" /></span>
            <span class="navitem__label">{{ menuTitle(node) }}</span>
          </button>
        </template>
      </div>
    </nav>
  </aside>
</template>
