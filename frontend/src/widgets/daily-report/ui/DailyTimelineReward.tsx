import { formatDateStr } from '@/shared/lib/utils'
import { Gift } from 'lucide-react'
import type { ActivityRewardItem } from '../model/types'

interface DailyTimelineRewardProps {
  reward: ActivityRewardItem
}

export const DailyTimelineReward = ({ reward }: DailyTimelineRewardProps) => {
  const timeStr = formatDateStr(reward.created_time)

  return (
    <div className="bg-surface-lowest border border-surface-highest p-4 rounded-xl space-y-2">
      <div className="flex items-center justify-between flex-wrap gap-2">
        <div className="flex items-center gap-2 font-bold text-on-surface text-base">
          <Gift className="size-5 text-purple-500 shrink-0" />
          <span>{reward.finished_desc || reward.name}</span>
        </div>
        <span className="text-xs font-semibold text-on-surface-variant bg-surface-highest px-2.5 py-1 rounded-full">
          {timeStr}
        </span>
      </div>

      {reward.description && (
        <p className="text-xs font-medium text-on-surface-variant">{reward.description}</p>
      )}

      {reward.promo_code && (
        <div className="inline-block bg-purple-500/10 text-purple-600 font-mono text-xs font-bold px-3 py-1 rounded-md border border-purple-500/20 mt-1">
          Промокод: {reward.promo_code}
        </div>
      )}
    </div>
  )
}
