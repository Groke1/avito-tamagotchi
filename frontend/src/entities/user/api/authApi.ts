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
    logout: builder.mutation<void, string | null>({
      query: (token) => ({
        url: '/auth/logout',
        method: 'POST',
        body: { refresh_token: token },
      }),
    }),
    refreshToken: builder.mutation<AuthTokens, string>({
      query: (token) => ({
        url: '/auth/refresh',
        method: 'POST',
        body: { refresh_token: token },
      }),
    }),
    sendUserAction: builder.mutation<void, void>({
      query: () => ({
        url: '/action',
        method: 'POST',
      }),
    }),
    getProfile: builder.query<UserResponse, void>({
      query: () => ({
        url: '/profile',
        method: 'GET',
      }),
      providesTags: ['User'],
      async onQueryStarted(_, { dispatch, queryFulfilled }) {
        const { data } = await queryFulfilled
        dispatch(setUser(data))
      },
    }),
  }),
})

export const {
  useLoginMutation,
  useRegisterMutation,
  useLogoutMutation,
  useRefreshTokenMutation,
  useSendUserActionMutation,
  useGetProfileQuery,
  useLazyGetProfileQuery,
} = authApi
