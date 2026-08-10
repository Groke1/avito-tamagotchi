import { describe, expect, it } from 'vitest'
import { isApiError, isFetchBaseQueryError } from './guards'

describe('guards', () => {
  describe('isFetchBaseQueryError', () => {
    it('returns true for valid FetchBaseQueryError objects with status property', () => {
      const error = { status: 400, data: 'Bad Request' }
      expect(isFetchBaseQueryError(error)).toBe(true)
    })

    it('returns false for non-object or null errors', () => {
      expect(isFetchBaseQueryError(null)).toBe(false)
      expect(isFetchBaseQueryError(undefined)).toBe(false)
      expect(isFetchBaseQueryError('string error')).toBe(false)
      expect(isFetchBaseQueryError(500)).toBe(false)
    })

    it('returns false for objects without status property', () => {
      expect(isFetchBaseQueryError({ message: 'Error' })).toBe(false)
    })
  })

  describe('isApiError', () => {
    it('returns true for objects containing code and message properties', () => {
      const apiErr = { code: 'PET_FULL', message: 'Питомец сыт' }
      expect(isApiError(apiErr)).toBe(true)
    })

    it('returns true for ApiError with extra properties like retry_after', () => {
      const apiErr = { code: 'COOLDOWN', message: 'Подождите', retry_after: 120 }
      expect(isApiError(apiErr)).toBe(true)
    })

    it('returns false for invalid structures or null', () => {
      expect(isApiError(null)).toBe(false)
      expect(isApiError({ code: 'ERR' })).toBe(false)
      expect(isApiError({ message: 'ERR' })).toBe(false)
      expect(isApiError('error')).toBe(false)
    })
  })
})
