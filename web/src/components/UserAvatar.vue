<script setup lang="ts">
import { computed } from 'vue'
import { initials, colorOf } from '@/utils/avatar'

const props = defineProps<{ name: string; src?: string; size?: number }>()
const px = computed(() => `${props.size ?? 32}px`)
const bg = computed(() => colorOf(props.name))
const text = computed(() => initials(props.name))
</script>

<template>
  <img v-if="src" :src="src" :style="{ width: px, height: px, borderRadius: '50%' }" alt="" />
  <span
    v-else
    :style="{
      width: px, height: px, background: bg, color: '#fff',
      display: 'inline-flex', alignItems: 'center', justifyContent: 'center',
      borderRadius: '50%', fontWeight: 600, fontSize: 'calc(' + px + ' * 0.42)',
    }"
  >{{ text }}</span>
</template>
