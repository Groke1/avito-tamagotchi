import { logout } from '@/entities/user'
import { type PayloadAction, createSlice } from '@reduxjs/toolkit'
import type { Pet, PetState } from './types'

const initialState: PetState = {
  pet: null,
  hasPet: null,
}

export const petSlice = createSlice({
  name: 'pet',
  initialState,
  reducers: {
    setPet: (state, action: PayloadAction<Pet>) => {
      state.pet = action.payload
      state.hasPet = true
    },
    setHasPet: (state, action: PayloadAction<boolean>) => {
      state.hasPet = action.payload
    },
  },
  extraReducers: (builder) => {
    builder.addCase(logout, (state) => {
      state.pet = null
      state.hasPet = null
    })
  },
})

export const { setPet, setHasPet } = petSlice.actions
