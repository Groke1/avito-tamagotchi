import { describe, expect, it } from 'vitest'
import { clearLatestTrip, petSlice, setPet } from './petSlice'
import type { Pet, PetState, TripResult } from './types'

describe('petSlice', () => {
  const initialState: PetState = {
    pet: null,
    latestTrip: null,
    dismissedStory: null,
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
    story: 'Кот исследовал город',
    coins: 50,
    xp: 40,
    reward: null,
  }

  it('should handle setPet', () => {
    const state = petSlice.reducer(initialState, setPet(mockPet))

    expect(state.pet).toEqual(mockPet)
    expect(state.isInitialized).toBe(true)
  })

  it('should handle clearLatestTrip and save dismissedStory', () => {
    const stateWithTrip: PetState = {
      ...initialState,
      latestTrip: mockTrip,
    }

    const state = petSlice.reducer(stateWithTrip, clearLatestTrip())

    expect(state.latestTrip).toBeNull()
    expect(state.dismissedStory).toBe(mockTrip.story)
  })
})
