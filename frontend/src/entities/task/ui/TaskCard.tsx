import { cn } from '@/shared/lib/utils'
import { Check, CheckCircle2, Coins, Loader2, Sparkles } from 'lucide-react'
import { DefaultTaskIcon, taskIconMap } from '../model/consts'
import type { Task } from '../model/types'

interface TaskCardProps {
  task: Task
  onComplete?: (taskId: string) => void
  isCompleting?: boolean
}

export const TaskCard = ({ task, onComplete, isCompleting = false }: TaskCardProps) => {
  const { id, title, description, reward_coins, reward_xp, status } = task
  const Icon = taskIconMap[title] ?? DefaultTaskIcon
  const isCompleted = status === 'completed'

  return (
    <div
      className={cn(
        'p-6 rounded-2xl border transition-all flex flex-col md:flex-row md:items-center justify-between gap-5 shadow-level-1 hover:shadow-level-2',
        isCompleted
          ? 'bg-avito-green/5 border-avito-green/30'
          : 'bg-surface-lowest border-outline-variant hover:border-avito-blue/30',
      )}
    >
      <div className="flex items-start md:items-center gap-4 flex-1 min-w-0">
        <div
          className={cn(
            'w-12 h-12 rounded-2xl flex items-center justify-center shrink-0 shadow-sm transition-colors',
            isCompleted ? 'bg-avito-green text-white' : 'bg-avito-blue/15 text-avito-blue-dark',
          )}
        >
          {isCompleted ? <Check className="w-6 h-6 stroke-[2.5]" /> : <Icon className="w-6 h-6" />}
        </div>
        <div className="space-y-1.5 flex-1 min-w-0">
          <h4 className="font-extrabold text-base md:text-lg text-on-surface truncate">{title}</h4>
          <p className="text-xs md:text-sm text-on-surface-variant leading-relaxed">
            {description}
          </p>
          <div className="flex items-center gap-2.5 pt-1 flex-wrap">
            <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-extrabold bg-avito-yellow/15 text-yellow-800">
              <Coins className="w-3.5 h-3.5 text-yellow-600" />+{reward_coins} 🪙
            </span>
            <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-extrabold bg-avito-blue/10 text-avito-blue-dark">
              <Sparkles className="w-3.5 h-3.5 text-avito-blue" />+{reward_xp} XP
            </span>
          </div>
        </div>
      </div>
      <div className="shrink-0 flex md:flex-col items-end justify-end">
        {isCompleted ? (
          <div className="inline-flex items-center gap-1.5 px-4 py-2.5 rounded-xl bg-avito-green/15 text-avito-green-dark text-xs font-bold">
            <Check className="w-4 h-4 stroke-[2.5]" />
            <span>Выполнено</span>
          </div>
        ) : (
          <button
            onClick={() => onComplete?.(id)}
            disabled={isCompleting}
            className="w-full md:w-auto px-5 py-2.5 bg-avito-blue hover:bg-avito-blue-dark active:scale-95 text-white font-bold text-xs rounded-xl shadow-level-1 transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2 cursor-pointer"
          >
            {isCompleting ? (
              <Loader2 className="w-4 h-4 animate-spin" />
            ) : (
              <CheckCircle2 className="w-4 h-4" />
            )}
            <span>{isCompleting ? 'Выполнение...' : 'Выполнить'}</span>
          </button>
        )}
      </div>
    </div>
  )
}
