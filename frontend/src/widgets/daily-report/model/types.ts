import type { DailyRewardStat, DailyTaskStat } from '@/entities/user'

export type ActivityTaskItem = DailyTaskStat & { type: 'task'; timestamp: number }
export type ActivityRewardItem = DailyRewardStat & { type: 'reward'; timestamp: number }
export type ActivityItem = ActivityTaskItem | ActivityRewardItem
