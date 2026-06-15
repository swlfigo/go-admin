import { describe, it, expect } from 'vitest'
import { initials, colorOf } from './avatar'

describe('avatar', () => {
  it('takes first CJK char', () => {
    expect(initials('超级管理员')).toBe('超')
  })
  it('takes up to two initials for latin words', () => {
    expect(initials('john doe')).toBe('JD')
    expect(initials('alice')).toBe('A')
  })
  it('falls back to ? for empty', () => {
    expect(initials('')).toBe('?')
  })
  it('color is stable per name and is an hsl string', () => {
    const c1 = colorOf('admin')
    const c2 = colorOf('admin')
    expect(c1).toBe(c2)
    expect(c1).toMatch(/^hsl\(/)
  })
  it('initials of a name with leading/trailing spaces returns first letter uppercased', () => {
    expect(initials('  bob  ')).toBe('B')
  })
  it('colorOf returns different hues for two clearly different names', () => {
    const c1 = colorOf('admin')
    const c2 = colorOf('张三')
    expect(c1).not.toBe(c2)
  })
})
