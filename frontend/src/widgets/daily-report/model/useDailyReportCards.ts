import type { DailyStatResponse } from '@/entities/user'
import { useMemo } from 'react'

export const useDailyReportCards = (data: DailyStatResponse) => {
  const completedTasksCount = data.tasks?.length ?? 0
  const earnedRewardsCount = useMemo(() => {
    return data.rewards?.filter((r) => r.status === 'active').length ?? 0
  }, [data.rewards])
  const totalCoinsEarned = useMemo(
    () => data.tasks?.reduce((sum, task) => sum + task.reward_coins, 0) ?? 0,
    [data.tasks],
  )
  const totalXp = data.pet?.daily_gained_xp ?? 0

  return { completedTasksCount, earnedRewardsCount, totalCoinsEarned, totalXp }
}
