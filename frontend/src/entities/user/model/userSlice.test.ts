import { describe, expect, it, beforeEach } from 'vitest'
import { userSlice, login, logout, setUser, setAccessToken } from './userSlice'
import type { UserState, AuthResponse } from './types'

describe('userSlice', () => {
  let initialState: UserState

  beforeEach(() => {
    localStorage.clear()
    initialState = {
      user: null,
      accessToken: null,
      isAuthenticated: false,
      isInitialized: false,
    }
  })

  it('should return initial state by default', () => {
    const state = userSlice.reducer(undefined, { type: 'unknown' })
    expect(state).toEqual(initialState)
  })

  it('should handle login', () => {
    const authPayload: AuthResponse = {
      user: { user_id: 'u-123', email: 'test@avito.ru', username: 'Test User', coins: 100 },
      accessToken: 'token-abc',
      refreshToken: 'refresh-xyz',
    }

    const state = userSlice.reducer(initialState, login(authPayload))

    expect(state.user).toEqual(authPayload.user)
    expect(state.accessToken).toBe('token-abc')
    expect(state.isAuthenticated).toBe(true)
    expect(state.isInitialized).toBe(true)
    expect(localStorage.getItem('refresh_token')).toBe('refresh-xyz')
  })

  it('should handle logout', () => {
    const loggedInState: UserState = {
      user: { user_id: 'u-123', email: 'test@avito.ru', username: 'Test User', coins: 100 },
      accessToken: 'token-abc',
      isAuthenticated: true,
      isInitialized: true,
    }
    localStorage.setItem('refresh_token', 'refresh-xyz')

    const state = userSlice.reducer(loggedInState, logout())

    expect(state.user).toBeNull()
    expect(state.accessToken).toBeNull()
    expect(state.isAuthenticated).toBe(false)
    expect(state.isInitialized).toBe(true)
    expect(localStorage.getItem('refresh_token')).toBeNull()
  })

  it('should handle setUser and setAccessToken', () => {
    const user = { user_id: 'u-1', email: 'a@b.com', username: 'Alex', coins: 50 }
    let state = userSlice.reducer(initialState, setUser(user))
    expect(state.user).toEqual(user)
    expect(state.isAuthenticated).toBe(true)

    state = userSlice.reducer(state, setAccessToken('new-access-token'))
    expect(state.accessToken).toBe('new-access-token')
  })
})
