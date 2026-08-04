import { type PayloadAction, createSlice } from '@reduxjs/toolkit'
import type { AuthResponse, UserState } from './types'

const REFRESH_TOKEN_KEY = 'refresh_token'

const initialState: UserState = {
  user: null,
  accessToken: null,
  isAuthenticated: false,
  isInitialized: false,
}

export const userSlice = createSlice({
  name: 'user',
  initialState,
  reducers: {
    login: (
      state,
      { payload: { user, accessToken, refreshToken } }: PayloadAction<AuthResponse>,
    ) => {
      state.user = user
      state.accessToken = accessToken
      state.isAuthenticated = true
      state.isInitialized = true

      if (refreshToken) {
        localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken)
      }
    },
    logout: (state) => {
      state.user = null
      state.accessToken = null
      state.isAuthenticated = false
      localStorage.removeItem(REFRESH_TOKEN_KEY)
    },
    setAccessToken: (state, action: PayloadAction<string>) => {
      state.accessToken = action.payload
      state.isAuthenticated = true
    },
  },
})

export const { login, logout, setAccessToken } = userSlice.actions

export const getStoredRefreshToken = (): string | null => {
  return localStorage.getItem(REFRESH_TOKEN_KEY)
}
