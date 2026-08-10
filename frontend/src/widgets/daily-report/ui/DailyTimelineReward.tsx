import { formatDateStr } from '@/shared/lib/utils'
import { Gift, Tag } from 'lucide-react'
import type { ActivityRewardItem } from '../model/types'

interface DailyTimelineRewardProps {
  reward: ActivityRewardItem
}

export const DailyTimelineReward = ({ reward }: DailyTimelineRewardProps) => {
  const timeStr = formatDateStr(reward.created_time)

  return (
    <div className="bg-purple-500/10 border border-purple-500/20 p-4.5 rounded-card">
      <div className="flex items-start justify-between gap-4 flex-wrap sm:flex-nowrap">
        <div className="flex items-start gap-3.5 min-w-0">
          <div className="size-10 rounded-full bg-purple-500 text-white flex items-center justify-center shrink-0 shadow-xs">
            <Gift className="size-5" />
          </div>
          <div className="space-y-1.5 min-w-0">
            <h3 className="font-extrabold text-on-surface text-base leading-snug">
              {reward.finished_desc || reward.name}
            </h3>
            {reward.description && (
              <p className="text-xs font-semibold text-on-surface-variant leading-relaxed">
                {reward.description}
              </p>
            )}
            {reward.promo_code && (
              <div className="inline-flex items-center gap-2 bg-purple-500/15 text-purple-800/40 text-xs font-bold px-3 py-1.5 rounded-xl border border-purple-500/20 mt-1 shadow-2xs">
                <Tag className="size-3.5 text-purple-600 shrink-0" />
                <span>Промокод: {reward.promo_code}</span>
              </div>
            )}
          </div>
        </div>
        <div className="flex items-center gap-1.5 text-xs font-bold text-purple-800 dark:text-purple-300 bg-purple-500/15 px-3 py-1.5 rounded-full border border-purple-500/20 shrink-0 self-start">
          <span>{timeStr}</span>
        </div>
      </div>
    </div>
  )
}
