export interface DailyTaskStat {
  id: string
  title: string
  reward_coins: number
  reward_xp: number
  finished_desc: string
  updated_at: string
}

export interface DailyPetStat {
  daily_gained_xp: number
}

export interface DailyRewardStat {
  promo_code: string
  name: string
  description: string
  finished_desc: string
  created_time: string
  status: DailyRewardStatStatus
}

export type DailyRewardStatStatus = 'active' | 'redeemed'

export interface DailyStatResponse {
  user_id: string
  streak: number
  tasks: DailyTaskStat[] | null
  pet: DailyPetStat
  rewards: DailyRewardStat[] | null
}
