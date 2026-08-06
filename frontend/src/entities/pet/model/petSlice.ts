import { logout } from '@/entities/user'
import { type PayloadAction, createSlice } from '@reduxjs/toolkit'
import type { Pet, PetState } from './types'

const initialState: PetState = {
  pet: null,
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
  },
  extraReducers: (builder) => {
    builder.addCase(logout, (state) => {
      state.pet = null
      state.isInitialized = false
    })
  },
})

export const { setPet } = petSlice.actions
