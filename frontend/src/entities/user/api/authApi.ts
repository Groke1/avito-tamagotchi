import { baseApi } from '@/shared/api/baseApi'
import type { AuthTokens, LoginDto, RegisterDto, UserResponse } from '../model/types'
import { setUser } from '../model/userSlice'

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
    refreshToken: builder.mutation<AuthTokens, string>({
      query: (token) => ({
        url: '/auth/refresh',
        method: 'POST',
        body: { refresh_token: token },
      }),
    }),
    getProfile: builder.query<UserResponse, void>({
      query: () => ({
        url: '/profile',
        method: 'GET',
      }),
      providesTags: ['User'],
      async onQueryStarted(_, { dispatch, queryFulfilled }) {
        try {
          const { data } = await queryFulfilled
          dispatch(setUser(data))
        } catch {
          //
        }
      },
    }),
  }),
})

export const {
  useLoginMutation,
  useRegisterMutation,
  useRefreshTokenMutation,
  useLazyGetProfileQuery,
} = authApi
