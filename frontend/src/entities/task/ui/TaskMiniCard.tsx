import { Check, Coins, Loader2, Sparkles } from 'lucide-react'
import { DefaultTaskIcon, taskIconMap } from '../model/consts'
import type { Task } from '../model/types'

interface TaskMiniCardProps {
  task: Task
  onComplete?: (taskId: string) => void
  isCompleting?: boolean
}

export const TaskMiniCard = ({ task, onComplete, isCompleting = false }: TaskMiniCardProps) => {
  const { id, title, description, reward_coins, reward_xp, status } = task
  const Icon = taskIconMap[title] ?? DefaultTaskIcon
  const isCompleted = status === 'completed'

  return (
    <li
      className={`p-4 rounded-xl border transition-all flex items-center justify-between gap-4 ${
        isCompleted
          ? 'bg-avito-green/5 border-avito-green/30'
          : 'bg-surface-lowest border-outline-variant hover:border-avito-blue/40 shadow-sm'
      }`}
    >
      <div className="flex items-center gap-3.5 min-w-0 flex-1">
        <div
          className={`w-10 h-10 rounded-full flex items-center justify-center shrink-0 shadow-sm transition-colors ${
            isCompleted ? 'bg-avito-green text-white' : 'bg-avito-blue/15 text-avito-blue-dark'
          }`}
        >
          {isCompleted ? <Check className="w-5 h-5 stroke-[2.5]" /> : <Icon className="w-5 h-5" />}
        </div>

        <div className="min-w-0 flex-1 space-y-1">
          <div className="flex items-center gap-2 flex-wrap">
            <h5 className="font-bold text-sm text-on-surface truncate">{title}</h5>
          </div>

          <p className="text-xs text-on-surface-variant line-clamp-1">{description}</p>

          <div className="flex items-center gap-2 pt-0.5">
            <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-bold bg-avito-yellow/15 text-yellow-700">
              <Coins className="w-3 h-3 text-yellow-600" />+{reward_coins} 🪙
            </span>
            <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-bold bg-avito-blue/10 text-avito-blue-dark">
              <Sparkles className="w-3 h-3 text-avito-blue" />+{reward_xp} XP
            </span>
          </div>
        </div>
      </div>
      <div className="shrink-0">
        {isCompleted ? (
          <div className="flex items-center gap-1 px-3 py-1.5 rounded-lg bg-avito-green/15 text-avito-green-dark text-xs font-bold">
            <Check className="w-4 h-4" />
            Выполнено
          </div>
        ) : (
          <button
            onClick={() => onComplete?.(id)}
            disabled={isCompleting}
            className="px-3.5 py-2 bg-avito-blue hover:bg-avito-blue-dark active:scale-95 text-white font-bold text-xs rounded-xl shadow-level-1 transition-all disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none flex items-center gap-1.5 cursor-pointer"
          >
            {isCompleting && <Loader2 className="w-3.5 h-3.5 animate-spin" />}
            {isCompleting ? 'Загрузка...' : 'Выполнить'}
          </button>
        )}
      </div>
    </li>
  )
}
