import { logout } from '@/entities/user'
import { describe, expect, it } from 'vitest'
import { clearLatestTrip, petSlice, setLatestTrip, setPet } from './petSlice'
import type { Pet, PetState, TripResult } from './types'

describe('petSlice', () => {
  const initialState: PetState = {
    pet: null,
    latestTrip: null,
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

  const mockTrip: TripResult = {
    story: 'Я сходил в путешествие и вернулся с находкой.',
    coins: 35,
    xp: 60,
    reward: {
      id: 'reward-1',
      promo_code: 'TRIP-001',
      name: 'Подарок от путешествия',
      description: 'Скидка на заказ',
      status: 'active',
      expires_at: '2026-08-31T00:00:00Z',
      earned_reason: 'trip',
      redeemed_at: null,
    },
  }

  it('should handle setPet', () => {
    const state = petSlice.reducer(initialState, setPet(mockPet))

    expect(state.pet).toEqual(mockPet)
    expect(state.isInitialized).toBe(true)
  })

  it('should handle setLatestTrip', () => {
    const state = petSlice.reducer(initialState, setLatestTrip(mockTrip))

    expect(state.latestTrip).toEqual(mockTrip)
  })

  it('should clear latest trip', () => {
    const activeState: PetState = {
      pet: mockPet,
      latestTrip: mockTrip,
      isInitialized: true,
    }

    const state = petSlice.reducer(activeState, clearLatestTrip())

    expect(state.latestTrip).toBeNull()
  })

  it('should reset pet state on logout', () => {
    const activeState: PetState = {
      pet: mockPet,
      latestTrip: mockTrip,
      isInitialized: true,
    }

    const state = petSlice.reducer(activeState, logout())

    expect(state.pet).toBeNull()
    expect(state.latestTrip).toBeNull()
    expect(state.isInitialized).toBe(false)
  })
})
