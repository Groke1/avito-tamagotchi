import { describe, expect, it } from 'vitest'
import { formatTaskDate, getDaysPlural } from './utils'

describe('utils', () => {
  describe('getDaysPlural', () => {
    it('returns correct Russian form for 1 day', () => {
      expect(getDaysPlural(1)).toBe('день')
      expect(getDaysPlural(21)).toBe('день')
      expect(getDaysPlural(101)).toBe('день')
    })

    it('returns correct Russian form for 2-4 days', () => {
      expect(getDaysPlural(2)).toBe('дня')
      expect(getDaysPlural(3)).toBe('дня')
      expect(getDaysPlural(4)).toBe('дня')
      expect(getDaysPlural(22)).toBe('дня')
    })

    it('returns correct Russian form for 5+ days and exceptions 11-14', () => {
      expect(getDaysPlural(0)).toBe('дней')
      expect(getDaysPlural(5)).toBe('дней')
      expect(getDaysPlural(11)).toBe('дней')
      expect(getDaysPlural(12)).toBe('дней')
      expect(getDaysPlural(14)).toBe('дней')
      expect(getDaysPlural(20)).toBe('дней')
    })
  })

  describe('formatTaskDate', () => {
    it('returns empty string for null or undefined', () => {
      expect(formatTaskDate({ dateStr: null })).toBe('')
      expect(formatTaskDate({ dateStr: undefined })).toBe('')
    })

    it('formats ISO date string into Russian date', () => {
      const formatted = formatTaskDate({ dateStr: '2026-08-10T12:00:00Z' })
      expect(formatted).toContain('2026')
      expect(formatted).toContain('августа')
    })

    it('supports object argument with hasMinutes option', () => {
      const formatted = formatTaskDate({ dateStr: '2026-08-10T12:30:00Z', hasMinutes: true })
      expect(formatted).toContain('2026')
      expect(formatted).toContain('августа')
      expect(formatted).toMatch(/30/)
    })
  })
})
