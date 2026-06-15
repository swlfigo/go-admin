<script setup lang="ts">
import { onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { usePagedTable } from '@/composables/usePagedTable'
import { listOperationLogs, clearOperationLogs } from '@/api/log'
import type { OperationLog } from '@/types/api'
import { formatDateTime } from '@/utils/format'

const { t } = useI18n()
const { rows, total, loading, query, load, search, onPage } = usePagedTable<OperationLog>(listOperationLogs)

type TagType = 'success' | 'warning' | 'danger' | 'info'

function methodTag(method: string): TagType {
  switch (method.toUpperCase()) {
    case 'POST': return 'success'
    case 'PUT': return 'warning'
    case 'DELETE': return 'danger'
    default: return 'info'
  }
}

async function clearAll() {
  await ElMessageBox.confirm(t('opLog.confirmClear'), t('common.tip'), { type: 'warning' })
  try {
    await clearOperationLogs()
    ElMessage.success(t('opLog.cleared'))
    load()
  } catch (e) {
    ElMessage.error((e as Error).message)
  }
}

onMounted(load)
</script>

<template>
  <div class="page">
    <div class="page__head">
      <h1>{{ t('opLog.title') }}</h1>
      <div class="page__actions">
        <el-input
          v-model="query.keyword"
          :placeholder="t('opLog.searchPlaceholder')"
          style="width:200px"
          clearable
          @keyup.enter="search"
        />
        <el-button @click="search">{{ t('common.search') }}</el-button>
        <el-button v-permission="'monitor:log:delete'" type="danger" @click="clearAll">{{ t('opLog.clear') }}</el-button>
      </div>
    </div>

    <el-table :data="rows" v-loading="loading" border>
      <el-table-column :label="t('opLog.time')" width="180">
        <template #default="{ row }: { row: OperationLog }">{{ formatDateTime(row.createdAt) }}</template>
      </el-table-column>
      <el-table-column prop="username" :label="t('opLog.username')" width="140" />
      <el-table-column :label="t('opLog.method')" width="100">
        <template #default="{ row }: { row: OperationLog }">
          <el-tag :type="methodTag(row.method)">{{ row.method }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="path" :label="t('opLog.path')" show-overflow-tooltip />
      <el-table-column :label="t('opLog.status')" width="100">
        <template #default="{ row }: { row: OperationLog }">
          <el-tag :type="row.status < 400 ? 'success' : 'danger'">{{ row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column :label="t('opLog.latency')" width="120">
        <template #default="{ row }: { row: OperationLog }">{{ row.latencyMs }} ms</template>
      </el-table-column>
      <el-table-column prop="ip" :label="t('opLog.ip')" width="160" />
    </el-table>

    <el-pagination
      style="margin-top:16px; justify-content:flex-end"
      layout="total, prev, pager, next"
      :total="total" :current-page="query.page" :page-size="query.size"
      @current-change="onPage"
    />
  </div>
</template>

<style scoped>
.page__head{ display:flex; justify-content:space-between; align-items:center; margin-bottom:16px; }
.page__head h1{ font-size:20px; margin:0; }
.page__actions{ display:flex; gap:8px; }
</style>
