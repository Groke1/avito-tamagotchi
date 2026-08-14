import { logout } from '@/entities/user'
import { type PayloadAction, createSlice } from '@reduxjs/toolkit'
import { petApi } from '../api/petApi'
import type { Pet, PetState, TripResult } from './types'

const initialState: PetState = {
  pet: null,
  latestTrip: null,
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
    setLatestTrip: (state, action: PayloadAction<TripResult | null>) => {
      state.latestTrip = action.payload
    },
    clearLatestTrip: (state) => {
      state.latestTrip = null
    },
  },
  extraReducers: (builder) => {
    builder.addCase(logout, (state) => {
      state.pet = null
      state.latestTrip = null
      state.isInitialized = false
    })
    builder.addMatcher(petApi.endpoints.getPet.matchFulfilled, (state, action) => {
      state.pet = action.payload
      state.isInitialized = true
    })
    builder.addMatcher(petApi.endpoints.getLastTrip.matchFulfilled, (state, action) => {
      state.latestTrip = action.payload
    })
  },
})

export const { setPet, setLatestTrip, clearLatestTrip } = petSlice.actions
