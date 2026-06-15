<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
const { t } = useI18n()
const router = useRouter()
</script>
<template>
  <div class="err">
    <div class="err__code">500</div>
    <div class="err__msg">{{ t('error.e500') }}</div>
    <el-button type="primary" @click="router.push('/')">{{ t('error.back') }}</el-button>
  </div>
</template>
<style scoped>
.err{ height:100%; display:flex; flex-direction:column; align-items:center; justify-content:center; gap:12px; }
.err__code{ font-size:64px; font-weight:800; color: var(--primary); }
.err__msg{ color: var(--fg-muted); }
</style>
