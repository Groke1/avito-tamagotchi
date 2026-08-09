import { baseApi } from '@/shared/api/baseApi'
import type { DailyStatResponse } from '../model/dailyStatTypes'

export const dailyStatApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getDailyStat: builder.query<DailyStatResponse, void>({
      query: () => '/daily-stat',
      providesTags: ['User', 'Pet', 'Tasks', 'Rewards'],
    }),
  }),
})

export const { useGetDailyStatQuery, useLazyGetDailyStatQuery } = dailyStatApi
