import { formatDateStr } from '@/shared/lib/utils'
import { Check, Coins, Sparkles } from 'lucide-react'
import type { ActivityTaskItem } from '../model/types'

interface DailyTimelineTaskProps {
  task: ActivityTaskItem
}

export const DailyTimelineTask = ({ task }: DailyTimelineTaskProps) => {
  const timeStr = formatDateStr(task.updated_at)

  return (
    <div className="bg-emerald-500/10 border border-emerald-500/20 p-4.5 rounded-card">
      <div className="flex items-start justify-between gap-4 flex-wrap sm:flex-nowrap">
        <div className="flex items-start gap-3.5 min-w-0 flex-wrap sm:flex-nowrap">
          <div className="size-10 rounded-full bg-emerald-500 text-white flex items-center justify-center shrink-0 shadow-xs">
            <Check className="size-5 stroke-[2.5]" />
          </div>
          <div className="space-y-1.5 min-w-0">
            <h3 className="font-extrabold text-on-surface text-base leading-snug">
              {task.finished_desc || task.title}
            </h3>
            {task.finished_desc && task.title && task.finished_desc !== task.title && (
              <p className="text-xs font-semibold text-on-surface-variant truncate">{task.title}</p>
            )}

            <div className="flex items-center gap-2.5 pt-1 flex-wrap">
              {task.reward_coins > 0 && (
                <div className="flex items-center gap-2 pt-0.5">
                  <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-avito-yellow/15 text-yellow-700">
                    <Coins className="w-3 h-3 text-yellow-600" />+{task.reward_coins} 🪙
                  </span>
                </div>
              )}
              {task.reward_xp > 0 && (
                <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-avito-blue/10 text-avito-blue-dark">
                  <Sparkles className="w-3 h-3 text-avito-blue" />+{task.reward_xp} XP
                </span>
              )}
            </div>
          </div>
        </div>
        <div className="flex items-center gap-1.5 text-xs font-bold text-emerald-800 dark:text-emerald-300 bg-emerald-500/15 px-3 py-1.5 rounded-full border border-emerald-500/20 shrink-0 self-start">
          <span>{timeStr}</span>
        </div>
      </div>
    </div>
  )
}
