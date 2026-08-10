import { logout } from '@/entities/user'
import { describe, expect, it } from 'vitest'
import { petSlice, setPet } from './petSlice'
import type { Pet, PetState } from './types'

describe('petSlice', () => {
  const initialState: PetState = {
    pet: null,
    isInitialized: false,
  }

  const mockPet: Pet = {
    id: 1,
    name: 'Кот',
    level: 3,
    xp: 120,
    next_level_xp: 200,
    satiety: 80,
    happiness: 90,
  }

  it('should handle setPet', () => {
    const state = petSlice.reducer(initialState, setPet(mockPet))

    expect(state.pet).toEqual(mockPet)
    expect(state.isInitialized).toBe(true)
  })

  it('should reset pet state on logout', () => {
    const activeState: PetState = {
      pet: mockPet,
      isInitialized: true,
    }

    const state = petSlice.reducer(activeState, logout())

    expect(state.pet).toBeNull()
    expect(state.isInitialized).toBe(false)
  })
})
