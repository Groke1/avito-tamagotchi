export interface LeaderboardItem {
  rank: number
  user_name: string
  pet_name: string
  xp: number
  level: number
}

export interface LeaderboardResponse {
  items: LeaderboardItem[]
  current_user: LeaderboardItem
}

export interface LeaderboardParams {
  limit?: number
}
