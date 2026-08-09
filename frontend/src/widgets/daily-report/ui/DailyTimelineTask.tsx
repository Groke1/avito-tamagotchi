import { formatDateStr } from '@/shared/lib/utils'
import { CheckCircle2, Coins, Sparkles } from 'lucide-react'
import type { ActivityTaskItem } from '../model/types'

interface DailyTimelineTaskProps {
  task: ActivityTaskItem
}

export const DailyTimelineTask = ({ task }: DailyTimelineTaskProps) => {
  const timeStr = formatDateStr(task.updated_at)

  return (
    <div className="bg-surface-lowest border border-surface-highest  p-4.5 rounded-xl">
      <div className="flex items-start justify-between gap-4 flex-wrap sm:flex-nowrap">
        <div className="flex items-start gap-3.5 min-w-0">
          <div className="size-10.5 rounded-xl bg-emerald-500/15 text-emerald-600 border border-emerald-500/20 flex items-center justify-center shrink-0 shadow-xs group-hover:scale-105 transition-transform">
            <CheckCircle2 className="size-5.5" />
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

        <div className="flex items-center gap-1.5 text-xs font-bold text-on-surface-variant/80 bg-surface-highest/60 px-3 py-1.5 rounded-full border border-surface-highest shrink-0 self-start">
          <span>{timeStr}</span>
        </div>
      </div>
    </div>
  )
}
