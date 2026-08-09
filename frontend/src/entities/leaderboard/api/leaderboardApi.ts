import { baseApi } from '@/shared/api/baseApi'
import type { LeaderboardParams, LeaderboardResponse } from '../model/types'

const LEADERBOARD_URL =
  import.meta.env.VITE_API_LEADERBOARD_URL || 'http://localhost:8082/api/v1/leaderboard'

export const leaderboardApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getLeaderboard: builder.query<LeaderboardResponse, LeaderboardParams | void>({
      query: (params) => ({
        url: LEADERBOARD_URL,
        params: params?.limit ? { limit: params.limit } : undefined,
      }),
      providesTags: ['Leaderboard'],
    }),
  }),
})

export const { useGetLeaderboardQuery } = leaderboardApi
