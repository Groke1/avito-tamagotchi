export type RewardStatus = 'active' | 'redeemed' | 'expired'

export interface UserReward {
  reward_id: string
  promo_code: string
  name: string
  description: string
  status: RewardStatus
  expires_at: string | null
  redeemed_at: string | null
}

export interface UserRewardsResponse {
  items: UserReward[]
}

export interface RedeemRewardRequest {
  promo_code: string
}
