import type { DailyRewardStat, DailyTaskStat } from '@/entities/user'
import { useDailyActivities } from '../model/useDailyActivities'
import { DailyTimelineList } from './DailyTimelineList'

interface DailyReportTimelineProps {
  tasks: DailyTaskStat[] | null
  rewards: DailyRewardStat[] | null
}

export const DailyReportTimeline = ({ tasks, rewards }: DailyReportTimelineProps) => {
  const activities = useDailyActivities(tasks, rewards)

  return (
    <div className="bg-surface-lowest  rounded-card p-8 space-y-4 shadow-level-1">
      <h2 className="text-2xl font-bold text-on-surface">Детализация активности</h2>
      <DailyTimelineList activities={activities} />
    </div>
  )
}
