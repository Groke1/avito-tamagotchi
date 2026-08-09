import { formatDateStr } from '@/shared/lib/utils'
import { CheckCircle2, Coins, Sparkles } from 'lucide-react'
import type { ActivityTaskItem } from '../model/types'

interface DailyTimelineTaskProps {
  task: ActivityTaskItem
}

export const DailyTimelineTask = ({ task }: DailyTimelineTaskProps) => {
  const timeStr = formatDateStr(task.updated_at)
  return (
    <div className="bg-surface-lowest border border-surface-highest p-4 rounded-xl space-y-2">
      <div className="flex items-center justify-between flex-wrap gap-2">
        <div className="flex items-center gap-2 font-bold text-on-surface text-base">
          <CheckCircle2 className="size-5 text-avito-green shrink-0" />
          <span>{task.finished_desc || task.title}</span>
        </div>
        <span className="text-xs font-semibold text-on-surface-variant bg-surface-highest px-2.5 py-1 rounded-full">
          {timeStr}
        </span>
      </div>

      <div className="flex items-center gap-3 pt-1">
        {task.reward_coins > 0 && (
          <span className="inline-flex items-center gap-1.5 text-xs font-bold text-on-surface bg-avito-yellow/10 px-2.5 py-1 rounded-md border border-avito-yellow/20">
            <Coins className="size-3.5 text-avito-yellow shrink-0" />+{task.reward_coins} монет
          </span>
        )}
        {task.reward_xp > 0 && (
          <span className="inline-flex items-center gap-1.5 text-xs font-bold text-on-surface bg-avito-green/10 px-2.5 py-1 rounded-md border border-avito-green/20">
            <Sparkles className="size-3.5 text-avito-green shrink-0" />+{task.reward_xp} XP
          </span>
        )}
      </div>
    </div>
  )
}
