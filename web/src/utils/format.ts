// 把后端返回的时间（ISO8601，如 2026-06-15T15:48:52.156426+08:00）格式化为
// 本地可读的 "YYYY-MM-DD HH:mm:ss"。空值返回占位符。
export function formatDateTime(value: string | null | undefined): string {
  if (!value) return '-'
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return String(value)
  const p = (n: number): string => String(n).padStart(2, '0')
  return (
    `${d.getFullYear()}-${p(d.getMonth() + 1)}-${p(d.getDate())} ` +
    `${p(d.getHours())}:${p(d.getMinutes())}:${p(d.getSeconds())}`
  )
}
