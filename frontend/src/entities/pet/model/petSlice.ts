import { logout } from '@/entities/user'
import { type PayloadAction, createSlice } from '@reduxjs/toolkit'
import { petApi } from '../api/petApi'
import type { Pet, PetState } from './types'

const initialState: PetState = {
  pet: null,
  latestTrip: null,
  dismissedStory: null,
  isInitialized: false,
}

export const petSlice = createSlice({
  name: 'pet',
  initialState,
  reducers: {
    setPet: (state, action: PayloadAction<Pet | null>) => {
      state.pet = action.payload
      state.isInitialized = true
    },
    clearLatestTrip: (state) => {
      if (state.latestTrip?.story) {
        state.dismissedStory = state.latestTrip.story
      }
      state.latestTrip = null
    },
  },
  extraReducers: (builder) => {
    builder.addCase(logout, (state) => {
      state.pet = null
      state.latestTrip = null
      state.dismissedStory = null
      state.isInitialized = false
    })
    builder.addMatcher(petApi.endpoints.getPet.matchFulfilled, (state, action) => {
      state.pet = action.payload
      state.isInitialized = true
    })
    builder.addMatcher(petApi.endpoints.getPetTripLast.matchFulfilled, (state, action) => {
      if (action.payload?.story && action.payload.story === state.dismissedStory) {
        state.latestTrip = null
      } else {
        state.latestTrip = action.payload
      }
    })
  },
})

export const { setPet, clearLatestTrip } = petSlice.actions
