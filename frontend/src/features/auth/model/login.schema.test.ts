import { describe, expect, it } from 'vitest'
import { loginSchema } from './login.schema'

describe('loginSchema validation', () => {
  it('passes for valid login data', () => {
    const validData = {
      email: 'alex@avito.ru',
      password: 'supersecretpassword',
    }

    const result = loginSchema.safeParse(validData)
    expect(result.success).toBe(true)
  })

  it('fails for invalid email format', () => {
    const invalidData = {
      email: 'invalid-email',
      password: 'supersecretpassword',
    }

    const result = loginSchema.safeParse(invalidData)
    expect(result.success).toBe(false)
  })

  it('fails for short password (< 8 chars)', () => {
    const invalidData = {
      email: 'alex@avito.ru',
      password: 'short',
    }

    const result = loginSchema.safeParse(invalidData)
    expect(result.success).toBe(false)
  })
})
