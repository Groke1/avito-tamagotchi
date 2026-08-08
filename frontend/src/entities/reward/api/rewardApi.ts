import { baseApi } from '@/shared/api/baseApi'
import type { RedeemRewardRequest, UserRewardsResponse } from '../model/types'

export const rewardApi = baseApi.injectEndpoints({
  endpoints: (builder) => ({
    getRewards: builder.query<UserRewardsResponse, void>({
      query: () => '/rewards',
      providesTags: ['Rewards'],
    }),
    getActiveRewards: builder.query<UserRewardsResponse, void>({
      query: () => '/rewards/active',
      providesTags: ['Rewards'],
    }),
    redeemReward: builder.mutation<void, RedeemRewardRequest>({
      query: (body) => ({
        url: '/rewards/redeem',
        method: 'POST',
        body,
      }),
      invalidatesTags: ['Rewards', 'User'],
    }),
  }),
})

export const { useGetRewardsQuery, useGetActiveRewardsQuery, useRedeemRewardMutation } = rewardApi
