<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { getServerInfo } from '@/api/monitor'
import type { ServerInfo } from '@/types/api'

const { t } = useI18n()

const info = ref<ServerInfo | null>(null)
const loading = ref(false)
let timer: ReturnType<typeof setInterval> | null = null

async function load(showError = true) {
  loading.value = true
  try {
    info.value = await getServerInfo()
  } catch (e) {
    if (showError) ElMessage.error((e as Error).message || t('monitor.loadFailed'))
  } finally {
    loading.value = false
  }
}

function progressType(percent: number): 'success' | 'warning' | 'exception' {
  if (percent >= 90) return 'exception'
  if (percent >= 70) return 'warning'
  return 'success'
}

function round(n: number): number {
  return Math.round(n * 10) / 10
}

const uptimeText = computed<string>(() => {
  const total = info.value?.uptimeSeconds ?? 0
  const d = Math.floor(total / 86400)
  const h = Math.floor((total % 86400) / 3600)
  const m = Math.floor((total % 3600) / 60)
  return `${d}${t('monitor.unitDay')} ${h}${t('monitor.unitHour')} ${m}${t('monitor.unitMinute')}`
})

onMounted(() => {
  load()
  timer = setInterval(() => load(false), 5000)
})
onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<template>
  <div class="page">
    <div class="page__head">
      <h1>{{ t('monitor.title') }}</h1>
      <el-button :loading="loading" @click="load()">{{ t('common.refresh') }}</el-button>
    </div>

    <div v-if="info" class="grid">
      <!-- 主机 -->
      <div class="card">
        <div class="card__head"><h3>{{ t('monitor.host') }}</h3></div>
        <div class="card__body">
          <div class="row"><span class="row__k">{{ t('monitor.hostname') }}</span><span class="row__v">{{ info.host.hostname }}</span></div>
          <div class="row"><span class="row__k">{{ t('monitor.os') }}</span><span class="row__v">{{ info.host.os }}</span></div>
          <div class="row"><span class="row__k">{{ t('monitor.platform') }}</span><span class="row__v">{{ info.host.platform }}</span></div>
          <div class="row"><span class="row__k">{{ t('monitor.arch') }}</span><span class="row__v">{{ info.host.arch }}</span></div>
        </div>
      </div>

      <!-- CPU -->
      <div class="card">
        <div class="card__head"><h3>{{ t('monitor.cpu') }}</h3></div>
        <div class="card__body">
          <el-progress :percentage="round(info.cpu.percent)" :status="progressType(info.cpu.percent)" :stroke-width="14" />
          <div class="row"><span class="row__k">{{ t('monitor.cores') }}</span><span class="row__v">{{ info.cpu.cores }}</span></div>
        </div>
      </div>

      <!-- 内存 -->
      <div class="card">
        <div class="card__head"><h3>{{ t('monitor.memory') }}</h3></div>
        <div class="card__body">
          <el-progress :percentage="round(info.memory.percent)" :status="progressType(info.memory.percent)" :stroke-width="14" />
          <div class="row"><span class="row__k">{{ t('monitor.used') }}</span><span class="row__v">{{ round(info.memory.usedMB) }} / {{ round(info.memory.totalMB) }} MB</span></div>
        </div>
      </div>

      <!-- 磁盘 -->
      <div class="card">
        <div class="card__head"><h3>{{ t('monitor.disk') }}</h3></div>
        <div class="card__body">
          <el-progress :percentage="round(info.disk.percent)" :status="progressType(info.disk.percent)" :stroke-width="14" />
          <div class="row"><span class="row__k">{{ t('monitor.used') }}</span><span class="row__v">{{ round(info.disk.usedGB) }} / {{ round(info.disk.totalGB) }} GB</span></div>
        </div>
      </div>

      <!-- Go 运行时 -->
      <div class="card">
        <div class="card__head"><h3>{{ t('monitor.goRuntime') }}</h3></div>
        <div class="card__body">
          <div class="row"><span class="row__k">{{ t('monitor.goVersion') }}</span><span class="row__v">{{ info.runtime.goVersion }}</span></div>
          <div class="row"><span class="row__k">{{ t('monitor.goroutines') }}</span><span class="row__v">{{ info.runtime.goroutines }}</span></div>
          <div class="row"><span class="row__k">{{ t('monitor.numCPU') }}</span><span class="row__v">{{ info.runtime.numCPU }}</span></div>
          <div class="row"><span class="row__k">{{ t('monitor.allocMB') }}</span><span class="row__v">{{ round(info.runtime.allocMB) }} MB</span></div>
        </div>
      </div>

      <!-- 运行时长 -->
      <div class="card">
        <div class="card__head"><h3>{{ t('monitor.uptime') }}</h3></div>
        <div class="card__body">
          <div class="uptime">{{ uptimeText }}</div>
        </div>
      </div>
    </div>

    <el-empty v-else-if="!loading" :description="t('monitor.loadFailed')" />
  </div>
</template>

<style scoped>
.page__head{ display:flex; justify-content:space-between; align-items:center; margin-bottom:16px; }
.page__head h1{ font-size:20px; margin:0; }
.grid{ display:grid; grid-template-columns:repeat(3,1fr); gap:16px; }
@media (max-width:1100px){ .grid{ grid-template-columns:repeat(2,1fr); } }
@media (max-width:680px){ .grid{ grid-template-columns:1fr; } }
.row{ display:flex; justify-content:space-between; align-items:center; gap:12px; padding:6px 0; font-size:13px; }
.row + .row{ border-top:1px solid var(--border); }
.row__k{ color:var(--fg-muted); }
.row__v{ font-weight:600; font-family:var(--font-mono); text-align:right; word-break:break-all; }
.card__body .el-progress{ margin-bottom:8px; }
.uptime{ font-size:22px; font-weight:800; font-family:var(--font-mono); letter-spacing:-.01em; }
</style>
