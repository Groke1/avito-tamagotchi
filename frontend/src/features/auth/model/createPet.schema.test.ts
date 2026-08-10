import { describe, expect, it } from 'vitest'
import { createPetSchema } from './createPet.schema'

describe('createPetSchema validation', () => {
  it('passes for valid pet name', () => {
    const validData = { name: 'Барсик' }
    const result = createPetSchema.safeParse(validData)

    expect(result.success).toBe(true)
    if (result.success) {
      expect(result.data.name).toBe('Барсик')
    }
  })

  it('fails if name is empty', () => {
    const emptyData = { name: '' }
    const result = createPetSchema.safeParse(emptyData)

    expect(result.success).toBe(false)
  })

  it('fails if name exceeds 25 characters', () => {
    const longName = { name: 'А'.repeat(26) }
    const result = createPetSchema.safeParse(longName)

    expect(result.success).toBe(false)
  })
})
