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
        'p-4 rounded-xl border transition-all flex items-center justify-between gap-4',
        isCompleted
          ? 'bg-avito-green/5 border-avito-green/40'
          : 'bg-surface-lowest border-outline-variant',
      )}
    >
      <div className="flex items-center gap-3.5 min-w-0 flex-1">
        <div
          className={cn(
            'w-10 h-10 rounded-full flex items-center justify-center shrink-0 transition-colors shadow-sm',
            isCompleted ? 'bg-avito-green text-white' : 'bg-avito-blue/15 text-avito-blue-dark',
          )}
        >
          <Icon className="w-5 h-5" />
        </div>
        <div className="min-w-0 flex-1 space-y-1">
          <div className="flex items-center gap-2 flex-wrap">
            <h5 className="font-bold text-sm text-on-surface truncate">{title}</h5>
          </div>
          <p className="text-xs text-on-surface-variant line-clamp-2">{description}</p>
          <div className="flex items-center gap-2 pt-0.5">
            <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-avito-yellow/15 text-yellow-700">
              <Coins className="w-3 h-3 text-yellow-600" />+{reward_coins} 🪙
            </span>
            <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-avito-blue/10 text-avito-blue-dark">
              <Sparkles className="w-3 h-3 text-avito-blue" />+{reward_xp} XP
            </span>
          </div>
        </div>
      </div>
      <div className="shrink-0">
        {isCompleted ? (
          <div className="w-8 h-8 rounded-full bg-avito-green text-white shadow-sm flex items-center justify-center">
            <Check className="w-4 h-4 stroke-[2.5]" />
          </div>
        ) : (
          <div className="w-8 h-8 rounded-full border-2 border-outline-variant text-outline-variant flex items-center justify-center">
            <Hourglass className="w-3.5 h-3.5" />
          </div>
        )}
      </div>
    </li>
  )
}
