import { reactive, ref } from 'vue'
import type { Paged } from '@/types/api'

export function usePagedTable<T>(fetcher: (q: { keyword: string; page: number; size: number }) => Promise<Paged<T>>) {
  const rows = ref<T[]>([]) as { value: T[] }
  const total = ref(0)
  const loading = ref(false)
  const query = reactive({ keyword: '', page: 1, size: 10 })

  async function load() {
    loading.value = true
    try {
      const res = await fetcher({ keyword: query.keyword, page: query.page, size: query.size })
      rows.value = res.list
      total.value = res.total
    } finally {
      loading.value = false
    }
  }
  function search() { query.page = 1; return load() }
  function onPage(p: number) { query.page = p; return load() }

  return { rows, total, loading, query, load, search, onPage }
}
