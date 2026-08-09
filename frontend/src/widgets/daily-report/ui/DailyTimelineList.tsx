import type { ActivityItem } from '../model/types'
import { DailyTimelineReward } from './DailyTimelineReward'
import { DailyTimelineTask } from './DailyTimelineTask'

interface DailyTimelineListProps {
  activities: ActivityItem[]
}

export const DailyTimelineList = ({ activities }: DailyTimelineListProps) => {
  return (
    <ul className="space-y-4">
      {activities.map((activity) => {
        if (activity.type === 'task') {
          return <DailyTimelineTask key={activity.id} task={activity} />
        }

        return (
          <DailyTimelineReward
            key={`${activity.promo_code}-${activity.timestamp}`}
            reward={activity}
          />
        )
      })}
    </ul>
  )
}
