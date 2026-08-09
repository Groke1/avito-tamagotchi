import type { DailyRewardStat, DailyTaskStat } from '@/entities/user'
import { useMemo } from 'react'
import type { ActivityItem, ActivityRewardItem, ActivityTaskItem } from './types'

export const useDailyActivities = (
  tasks: DailyTaskStat[] | null,
  rewards: DailyRewardStat[] | null,
): ActivityItem[] => {
  return useMemo<ActivityItem[]>(() => {
    const taskItems: ActivityTaskItem[] = (tasks ?? []).map((task) => ({
      ...task,
      type: 'task',
      timestamp: new Date(task.updated_at).getTime(),
    }))

    const rewardItems: ActivityRewardItem[] = (rewards ?? []).map((reward) => ({
      ...reward,
      type: 'reward',
      timestamp: new Date(reward.created_time).getTime(),
    }))

    return [...taskItems, ...rewardItems].sort((a, b) => b.timestamp - a.timestamp)
  }, [tasks, rewards])
}
