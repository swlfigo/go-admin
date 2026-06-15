<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as onlineApi from '@/api/online'
import type { OnlineSession } from '@/types/api'
import { formatDateTime } from '@/utils/format'

const { t } = useI18n()

const rows = ref<OnlineSession[]>([])
const loading = ref(false)

async function load() {
  loading.value = true
  try { rows.value = await onlineApi.listOnline() } finally { loading.value = false }
}
async function kick(row: OnlineSession) {
  await ElMessageBox.confirm(t('online.confirmKick', { name: row.username }), t('common.tip'), { type: 'warning' })
  await onlineApi.kickOnline(row.id); ElMessage.success(t('online.kicked')); load()
}

onMounted(load)
</script>

<template>
  <div class="page">
    <div class="page__head"><h1>{{ t('online.title') }}</h1><el-button @click="load">{{ t('common.refresh') }}</el-button></div>
    <el-table :data="rows" v-loading="loading" border>
      <el-table-column prop="username" :label="t('online.user')" />
      <el-table-column prop="loginIp" :label="t('online.ip')" />
      <el-table-column prop="browser" :label="t('online.browser')" />
      <el-table-column prop="os" :label="t('online.os')" />
      <el-table-column :label="t('online.loginTime')" width="180">
        <template #default="{ row }: { row: OnlineSession }">{{ formatDateTime(row.loginAt) }}</template>
      </el-table-column>
      <el-table-column :label="t('common.operations')" width="120">
        <template #default="{ row }">
          <el-button v-permission="'monitor:online:kick'" link type="danger" @click="kick(row)">{{ t('online.kick') }}</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<style scoped>
.page__head{ display:flex; justify-content:space-between; align-items:center; margin-bottom:16px; }
.page__head h1{ font-size:20px; margin:0; }
</style>
