import { cn } from '@/shared/lib/utils'
import { Check, Coins, Hourglass, Sparkles } from 'lucide-react'
import { DefaultTaskIcon, taskIconMap } from '../model/consts'
import type { Task } from '../model/types'

interface TaskMiniCardProps {
  task: Task
}

export const TaskMiniCard = ({ task }: TaskMiniCardProps) => {
  const { title, description, reward_coins, reward_xp, status } = task
  const Icon = taskIconMap[title] ?? DefaultTaskIcon
  const isCompleted = status === 'completed'

  return (
    <li
      className={cn(
        'p-3.5 sm:p-4 rounded-xl border transition-all flex items-center justify-between gap-3 sm:gap-4 min-w-0',
        isCompleted
          ? 'bg-avito-green/5 border-avito-green/40'
          : 'bg-surface-lowest border-surface-highest',
      )}
    >
      <div className="flex items-center gap-3 sm:gap-3.5 min-w-0 flex-1">
        <div
          className={cn(
            'size-9 sm:size-10 rounded-full flex items-center justify-center shrink-0 transition-colors shadow-level-1',
            isCompleted
              ? 'bg-avito-green/20 text-surface-lowest'
              : 'bg-avito-blue/15 text-avito-blue-dark',
          )}
        >
          <Icon className="size-4.5 sm:size-5" />
        </div>
        <div className="min-w-0 flex-1 space-y-1">
          <div className="flex items-center gap-2 flex-wrap">
            <h5 className="font-bold text-xs sm:text-sm text-on-surface truncate">{title}</h5>
          </div>
          <p className="text-xs text-on-surface-variant line-clamp-2 wrap-break-word">{description}</p>
          <div className="flex items-center gap-1.5 sm:gap-2 pt-0.5 flex-wrap">
            <span className="inline-flex items-center gap-1 px-2 sm:px-2.5 py-0.5 rounded-full text-[10px] sm:text-[11px] font-bold bg-avito-yellow/15 text-yellow-700">
              <Coins className="size-3 text-yellow-600 shrink-0" />+{reward_coins} 🪙
            </span>
            <span className="inline-flex items-center gap-1 px-2 sm:px-2.5 py-0.5 rounded-full text-[10px] sm:text-[11px] font-bold bg-avito-blue/10 text-avito-blue-dark">
              <Sparkles className="size-3 text-avito-blue shrink-0" />+{reward_xp} XP
            </span>
          </div>
        </div>
      </div>
      <div className="shrink-0">
        {isCompleted ? (
          <div className="size-7 sm:size-8 rounded-full bg-avito-green/20 text-surface-lowest shadow-level-1 flex items-center justify-center">
            <Check className="size-3.5 sm:size-4 stroke-[2.5]" />
          </div>
        ) : (
          <div className="size-7 sm:size-8 rounded-full border-2 border-outline-variant text-outline-variant flex items-center justify-center">
            <Hourglass className="size-3 sm:size-3.5" />
          </div>
        )}
      </div>
    </li>
  )
}
