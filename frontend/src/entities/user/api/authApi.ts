import { baseApi } from '@/shared/api/baseApi'
import type { AuthTokens, LoginDto, RegisterDto, UserResponse } from '../model/types'

export const authApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    login: builder.mutation<AuthTokens, LoginDto>({
      query: (credentials) => ({
        url: '/auth/login',
        method: 'POST',
        body: credentials,
      }),
    }),
    register: builder.mutation<AuthTokens, RegisterDto>({
      query: (userData) => ({
        url: '/auth/register',
        method: 'POST',
        body: userData,
      }),
    }),
    refreshToken: builder.mutation<AuthTokens, { refresh_token: string }>({
      query: (refreshToken) => ({
        url: '/auth/refresh',
        method: 'POST',
        body: { refresh_token: refreshToken },
      }),
    }),
    getProfile: builder.query<UserResponse, void>({
      query: () => ({
        url: '/profile',
        method: 'GET',
      }),
      providesTags: ['User'],
    }),
  }),
})

export const { useLoginMutation, useRegisterMutation, useRefreshTokenMutation } = authApi
