<script setup lang="ts">
import { computed } from 'vue'

// Inline stroke line-icons, 24x24 viewBox, currentColor stroke.
// Ported from the design mockup (admin/icons.jsx).
export type IconName =
  | 'grid'
  | 'settings'
  | 'users'
  | 'shield'
  | 'list'
  | 'book'
  | 'activity'
  | 'sun'
  | 'moon'
  | 'bell'
  | 'palette'
  | 'panel'
  | 'chevron'
  | 'search'
  | 'plus'
  | 'x'
  | 'logout'
  | 'user'
  | 'refresh'
  | 'pin'

const ICON_PATHS: Record<IconName, string> = {
  // grid / dashboard
  grid: '<rect x="3" y="3" width="7" height="9" rx="1.5"/><rect x="14" y="3" width="7" height="5" rx="1.5"/><rect x="14" y="12" width="7" height="9" rx="1.5"/><rect x="3" y="16" width="7" height="5" rx="1.5"/>',
  // settings / sliders
  settings: '<path d="M4 6h11"/><path d="M19 6h1"/><circle cx="17" cy="6" r="2"/><path d="M4 12h2"/><path d="M10 12h10"/><circle cx="8" cy="12" r="2"/><path d="M4 18h11"/><path d="M19 18h1"/><circle cx="17" cy="18" r="2"/>',
  users: '<circle cx="9" cy="8" r="3.4"/><path d="M3.5 19.5a5.5 5.5 0 0 1 11 0"/><path d="M16 5.2a3.4 3.4 0 0 1 0 6.6"/><path d="M18 19.5a5.5 5.5 0 0 0-3-4.9"/>',
  shield: '<path d="M12 3 5 6v5c0 4.4 3 7.6 7 9 4-1.4 7-4.6 7-9V6l-7-3z"/><path d="m9.5 12 1.8 1.8 3.5-3.6"/>',
  list: '<path d="M8 6h12"/><path d="M8 12h12"/><path d="M8 18h12"/><circle cx="4" cy="6" r="1"/><circle cx="4" cy="12" r="1"/><circle cx="4" cy="18" r="1"/>',
  book: '<path d="M5 4.5A1.5 1.5 0 0 1 6.5 3H19v15H6.5A1.5 1.5 0 0 0 5 19.5V4.5z"/><path d="M5 19.5A1.5 1.5 0 0 0 6.5 21H19"/><path d="M9 7.5h6"/><path d="M9 11h6"/>',
  activity: '<path d="M3 12h4l3 7 4-14 3 7h4"/>',
  // ui
  sun: '<circle cx="12" cy="12" r="4.2"/><path d="M12 2v2.5M12 19.5V22M2 12h2.5M19.5 12H22M4.6 4.6 6.4 6.4M17.6 17.6l1.8 1.8M19.4 4.6 17.6 6.4M6.4 17.6 4.6 19.4"/>',
  moon: '<path d="M20 14.5A8.2 8.2 0 0 1 9.5 4 8.3 8.3 0 1 0 20 14.5z"/>',
  bell: '<path d="M18 8.5a6 6 0 1 0-12 0c0 6-2.5 7.5-2.5 7.5h17S18 14.5 18 8.5z"/><path d="M10.3 20a2 2 0 0 0 3.4 0"/>',
  palette: '<path d="M12 3a9 9 0 1 0 0 18c1.2 0 2-.9 2-2 0-.6-.2-1-.6-1.4-.3-.4-.5-.8-.5-1.3 0-1 .8-1.8 1.8-1.8H17a4 4 0 0 0 4-4c0-4.4-4-7.7-9-7.7z"/><circle cx="7.5" cy="11" r="1.1"/><circle cx="11" cy="7.2" r="1.1"/><circle cx="15.5" cy="8" r="1.1"/>',
  panel: '<rect x="3" y="4" width="18" height="16" rx="2"/><path d="M9 4v16"/>',
  chevron: '<path d="m9 6 6 6-6 6"/>',
  search: '<circle cx="11" cy="11" r="7"/><path d="m20 20-3.2-3.2"/>',
  plus: '<path d="M12 5v14M5 12h14"/>',
  x: '<path d="M6 6l12 12M18 6 6 18"/>',
  logout: '<path d="M9 4.5H6a1.5 1.5 0 0 0-1.5 1.5v12A1.5 1.5 0 0 0 6 19.5h3"/><path d="M15.5 12H21"/><path d="m18 8.5 3.5 3.5L18 15.5"/>',
  user: '<circle cx="12" cy="8" r="4"/><path d="M4.5 20a7.5 7.5 0 0 1 15 0"/>',
  refresh: '<path d="M21 12a9 9 0 1 1-3-6.7"/><path d="M21 4v4h-4"/>',
  pin: '<path d="M9 3.5h6l-1 5 2.5 2.5v1.5H7.5V11L10 8.5 9 3.5z"/><path d="M12 12.5v8"/>',
}

const props = withDefaults(
  defineProps<{ name: IconName; size?: number; stroke?: number }>(),
  { size: 18, stroke: 1.75 },
)

const inner = computed(() => ICON_PATHS[props.name])
</script>

<template>
  <svg
    class="icon"
    :width="size"
    :height="size"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    :stroke-width="stroke"
    stroke-linecap="round"
    stroke-linejoin="round"
    aria-hidden="true"
    v-html="inner"
  />
</template>
