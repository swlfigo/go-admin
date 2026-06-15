import { describe, it, expect } from 'vitest'
import { formatDateTime } from './format'

describe('formatDateTime', () => {
  it('formats an ISO datetime to YYYY-MM-DD HH:mm:ss', () => {
    // 只断言格式（精确值依赖运行机器时区，避免 flaky）
    const out = formatDateTime('2026-06-15T15:48:52.156426+08:00')
    expect(out).toMatch(/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$/)
  })
  it('returns - for empty', () => {
    expect(formatDateTime(null)).toBe('-')
    expect(formatDateTime(undefined)).toBe('-')
    expect(formatDateTime('')).toBe('-')
  })
  it('returns the raw string for unparseable input', () => {
    expect(formatDateTime('not-a-date')).toBe('not-a-date')
  })
})
