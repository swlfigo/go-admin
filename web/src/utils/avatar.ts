// 取首字母：中文/非拉丁取第一个字；拉丁取首个或两词首字母
export function initials(name: string): string {
  const n = (name || '').trim()
  if (!n) return '?'
  const parts = n.split(/\s+/)
  const isLatin = /^[A-Za-z]/.test(n)
  if (isLatin) {
    if (parts.length >= 2) {
      return (parts[0][0] + parts[1][0]).toUpperCase()
    }
    return parts[0][0].toUpperCase()
  }
  return Array.from(n)[0]
}

// 由用户名哈希出稳定的 HSL 背景色
export function colorOf(name: string): string {
  let hash = 0
  for (let i = 0; i < name.length; i++) {
    hash = (hash << 5) - hash + name.charCodeAt(i)
    hash |= 0
  }
  const hue = Math.abs(hash) % 360
  return `hsl(${hue} 55% 50%)`
}
