import { describe, expect, it } from 'vitest'
import { registerSchema } from './register.schema'

describe('registerSchema validation', () => {
  it('passes for valid registration data', () => {
    const validData = {
      username: 'user_123',
      email: 'user@avito.ru',
      password: 'password123',
      terms: true,
    }

    const result = registerSchema.safeParse(validData)
    expect(result.success).toBe(true)
  })

  it('fails if email is invalid', () => {
    const invalidData = {
      username: 'user_123',
      email: 'not-an-email',
      password: 'password123',
      terms: true,
    }

    const result = registerSchema.safeParse(invalidData)
    expect(result.success).toBe(false)
    if (!result.success) {
      expect(result.error.flatten().fieldErrors.email).toContain('Введите корректный E-mail')
    }
  })

  it('fails if username contains disallowed characters or is too short', () => {
    const invalidChars = {
      username: 'пользователь',
      email: 'user@avito.ru',
      password: 'password123',
      terms: true,
    }

    const result = registerSchema.safeParse(invalidChars)
    expect(result.success).toBe(false)
    if (!result.success) {
      expect(result.error.flatten().fieldErrors.username).toContain(
        'Используйте только латинские буквы, цифры и символ подчеркивания',
      )
    }
  })

  it('fails if password is shorter than 8 characters', () => {
    const shortPassword = {
      username: 'valid_user',
      email: 'user@avito.ru',
      password: '123',
      terms: true,
    }

    const result = registerSchema.safeParse(shortPassword)
    expect(result.success).toBe(false)
    if (!result.success) {
      expect(result.error.flatten().fieldErrors.password).toContain(
        'Пароль должен содержать минимум 8 символов',
      )
    }
  })

  it('fails if terms are not accepted', () => {
    const unacceptedTerms = {
      username: 'valid_user',
      email: 'user@avito.ru',
      password: 'password123',
      terms: false,
    }

    const result = registerSchema.safeParse(unacceptedTerms)
    expect(result.success).toBe(false)
    if (!result.success) {
      expect(result.error.flatten().fieldErrors.terms).toContain(
        'Необходимо согласиться с условиями',
      )
    }
  })
})
