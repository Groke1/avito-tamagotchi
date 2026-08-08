import type { RootState } from '@/app/store'
import type { AuthTokens } from '@/entities/user/model/types'
import {
  getStoredRefreshToken,
  logout,
  setAccessToken,
  setStoredRefreshToken,
} from '@/entities/user/model/userSlice'
import {
  type BaseQueryFn,
  type FetchArgs,
  type FetchBaseQueryError,
  createApi,
  fetchBaseQuery,
} from '@reduxjs/toolkit/query/react'

const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

const rawBaseQuery = fetchBaseQuery({
  baseUrl: BASE_URL,
  prepareHeaders: (headers, { getState }) => {
    const token = (getState() as RootState).user.accessToken

    if (token) {
      headers.set('Authorization', `Bearer ${token}`)
    }

    return headers
  },
})

const baseQueryWithReauth: BaseQueryFn<string | FetchArgs, unknown, FetchBaseQueryError> = async (
  args,
  api,
  extraOptions,
) => {
  let result = await rawBaseQuery(args, api, extraOptions)

  if (result.error && result.error.status === 401) {
    const refreshToken = getStoredRefreshToken()

    if (refreshToken) {
      const refreshResult = await rawBaseQuery(
        {
          url: '/auth/refresh',
          method: 'POST',
          body: { refresh_token: refreshToken },
        },
        api,
        extraOptions,
      )

      if (refreshResult.data) {
        const tokens = refreshResult.data as AuthTokens
        api.dispatch(setAccessToken(tokens.access_token))
        if (tokens.refresh_token) {
          setStoredRefreshToken(tokens.refresh_token)
        }

        result = await rawBaseQuery(args, api, extraOptions)
      } else {
        api.dispatch(logout())
      }
    } else {
      api.dispatch(logout())
    }
  }

  return result
}

export const baseApi = createApi({
  reducerPath: 'api',
  baseQuery: baseQueryWithReauth,
  tagTypes: ['User', 'Pet', 'Tasks'],
  endpoints: () => ({}),
})
