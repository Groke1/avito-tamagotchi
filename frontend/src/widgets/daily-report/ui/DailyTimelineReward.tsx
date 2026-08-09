import { formatDateStr } from '@/shared/lib/utils'
import { Gift, Tag } from 'lucide-react'
import type { ActivityRewardItem } from '../model/types'

interface DailyTimelineRewardProps {
  reward: ActivityRewardItem
}

export const DailyTimelineReward = ({ reward }: DailyTimelineRewardProps) => {
  const timeStr = formatDateStr(reward.created_time)

  return (
    <div className="bg-surface-lowest border border-surface-highest p-4.5 rounded-xl">
      <div className="flex items-start justify-between gap-4 flex-wrap sm:flex-nowrap">
        <div className="flex items-start gap-3.5 min-w-0">
          <div className="size-10.5 rounded-xl bg-purple-500/15 text-purple-600 border border-purple-500/20 flex items-center justify-center shrink-0 shadow-xs group-hover:scale-105 transition-transform">
            <Gift className="size-5.5" />
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
              <div className="inline-flex items-center gap-2 bg-purple-500/10 text-purple-700 dark:text-purple-300 font-mono text-xs font-bold px-3 py-1.5 rounded-xl border border-purple-500/25 mt-1 shadow-2xs">
                <Tag className="size-3.5 text-purple-500 shrink-0" />
                <span>Промокод: {reward.promo_code}</span>
              </div>
            )}
          </div>
        </div>
        <div className="flex items-center gap-1.5 text-xs font-bold text-on-surface-variant/80 bg-surface-highest/60 px-3 py-1.5 rounded-full border border-surface-highest">
          <span>{timeStr}</span>
        </div>
      </div>
    </div>
  )
}
