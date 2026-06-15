<script setup lang="ts">
import { ref } from 'vue'
import { useAppStore } from '@/stores/app'
import AppSidebar from './AppSidebar.vue'
import AppTopbar from './AppTopbar.vue'
import AppTabs from './AppTabs.vue'
import SettingsDrawer from './SettingsDrawer.vue'

const app = useAppStore()
const settingsOpen = ref(false)
</script>

<template>
  <div class="app" :data-collapsed="app.sidebarCollapsed">
    <AppSidebar />
    <div class="main">
      <AppTopbar @open-settings="settingsOpen = true" />
      <AppTabs />
      <main class="content scroll">
        <router-view />
      </main>
    </div>
    <SettingsDrawer :open="settingsOpen" @close="settingsOpen = false" />
  </div>
</template>
