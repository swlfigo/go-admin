import { describe, it, expect } from 'vitest'
import { shouldAttemptRefresh } from './http'
import type { RequestConfig } from './http'

/**
 * Tests for shouldAttemptRefresh — the predicate that breaks the self-await deadlock.
 *
 * Root-cause recap:
 *   When /auth/refresh returns 401, the response interceptor fires again.
 *   If it tried to do `await refreshing`, it would deadlock because `refreshing`
 *   is the promise that called /auth/refresh and is waiting for THIS interceptor
 *   to return. The predicate returns false for _skipAuthRefresh requests, causing
 *   the interceptor to take the fast-fail path (call onAuthFail + reject) instead.
 */

describe('shouldAttemptRefresh', () => {
  // ── Happy path: normal API call with 401 ──────────────────────────────────
  it('returns true for a plain 401 with no special config', () => {
    expect(shouldAttemptRefresh(401, {})).toBe(true)
  })

  it('returns true for a 401 when _skipAuthRefresh is explicitly false', () => {
    const cfg: RequestConfig = { _skipAuthRefresh: false }
    expect(shouldAttemptRefresh(401, cfg)).toBe(true)
  })

  it('returns true when config is undefined and status is 401', () => {
    expect(shouldAttemptRefresh(401, undefined)).toBe(true)
  })

  // ── The deadlock-prevention cases ─────────────────────────────────────────
  it('returns false when _skipAuthRefresh is true — prevents self-await deadlock', () => {
    // This is the critical case: /auth/refresh itself returned 401.
    // Without this guard the interceptor would `await refreshing` which is
    // the same promise that depends on this interceptor returning → deadlock.
    const cfg: RequestConfig = { _skipAuthRefresh: true }
    expect(shouldAttemptRefresh(401, cfg)).toBe(false)
  })

  it('returns false for logout() 401 — logout is also flagged to skip refresh', () => {
    const cfg: RequestConfig = { url: '/auth/logout', method: 'post', _skipAuthRefresh: true }
    expect(shouldAttemptRefresh(401, cfg)).toBe(false)
  })

  // ── Non-401 statuses should never trigger refresh ─────────────────────────
  it('returns false for 403 (forbidden, not an auth expiry)', () => {
    expect(shouldAttemptRefresh(403, {})).toBe(false)
  })

  it('returns false for 500', () => {
    expect(shouldAttemptRefresh(500, {})).toBe(false)
  })

  it('returns false for 404', () => {
    expect(shouldAttemptRefresh(404, {})).toBe(false)
  })

  it('returns false when status is undefined (network error before response)', () => {
    expect(shouldAttemptRefresh(undefined, {})).toBe(false)
  })

  // ── Non-401 + _skipAuthRefresh combination ────────────────────────────────
  it('returns false for a non-401 status even when _skipAuthRefresh is false', () => {
    const cfg: RequestConfig = { _skipAuthRefresh: false }
    expect(shouldAttemptRefresh(500, cfg)).toBe(false)
  })
})
