import type { DailyRewardStat, DailyTaskStat } from '@/entities/user'

interface DailyReportTimelineProps {
  tasks: DailyTaskStat[] | null
  rewards: DailyRewardStat[] | null
}

export const DailyReportTimeline = ({ tasks, rewards }: DailyReportTimelineProps) => {
  return (
    <div className="bg-surface-low border border-surface-high rounded-2xl p-6 space-y-6 shadow-xs"></div>
  )
}
