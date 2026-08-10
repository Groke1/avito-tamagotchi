import { cn } from '@/shared/lib/utils'
import { Button } from '@/shared/ui'
import { Check, Coins, Sparkles } from 'lucide-react'
import type { Task, TaskType } from '../model/types'

interface TaskCardProps {
  task: Task
  onComplete?: (taskId: string) => void
  isCompleting?: boolean
}

const typeMap: Record<TaskType, { badgeClass: string; emoji: string }> = {
  Отзывы: {
    badgeClass: 'bg-amber-100 text-amber-800',
    emoji: '⭐',
  },
  Поиск: {
    badgeClass: 'bg-purple-100 text-purple-800',
    emoji: '🔍',
  },
  Сообщения: {
    badgeClass: 'bg-blue-100 text-avito-blue-dark',
    emoji: '💬',
  },
  Категории: {
    badgeClass: 'bg-emerald-100 text-emerald-800',
    emoji: '🏷️',
  },
  Покупки: {
    badgeClass: 'bg-rose-100 text-rose-800',
    emoji: '🛍️',
  },
  Избранное: {
    badgeClass: 'bg-pink-100 text-pink-800',
    emoji: '❤️',
  },
}

export const TaskCard = ({ task, onComplete, isCompleting = false }: TaskCardProps) => {
  const { id, title, description, reward_coins, reward_xp, status, task_type } = task
  const isCompleted = status === 'completed'

  const category = typeMap[task_type] ?? {
    badgeClass: 'bg-blue-100 text-avito-blue-dark',
    emoji: '✨',
  }

  return (
    <div
      className={cn(
        'p-3.5 sm:p-4 rounded-card border bg-surface-lowest transition-all flex flex-col justify-between gap-3 shadow-level-1 relative min-w-0',
        isCompleted ? 'border-avito-green/40 bg-avito-green/5' : 'border-surface-highest',
      )}
    >
      <div className="flex items-center justify-between gap-2">
        <span
          className={cn(
            'px-2.5 py-0.5 rounded-full text-xs font-bold',
            isCompleted ? 'bg-surface-high text-on-surface-variant' : category.badgeClass,
          )}
        >
          {task_type}
        </span>
        {isCompleted && (
          <div className="inline-flex items-center p-1 rounded-full bg-avito-green/20 text-surface-lowest">
            <Check className="size-4 stroke-[2.5]" />
          </div>
        )}
      </div>

      <div className="space-y-1 min-w-0">
        <h4 className="font-extrabold text-sm sm:text-base text-on-surface leading-snug wrap-break-word">
          {category.emoji} {title}
        </h4>
        <p className="text-xs text-on-surface-variant leading-relaxed line-clamp-2 wrap-break-word">
          {description}
        </p>
      </div>

      <div className="flex items-center gap-2 pt-0.5 flex-wrap">
        <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-avito-yellow/15 text-yellow-700">
          <Coins className="w-3 h-3 text-yellow-600 shrink-0" />+{reward_coins} 🪙
        </span>
        <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-avito-blue/10 text-avito-blue-dark">
          <Sparkles className="w-3 h-3 text-avito-blue shrink-0" />+{reward_xp} XP
        </span>
      </div>

      {isCompleted ? (
        <div className="w-full py-2 px-3 rounded-xl bg-avito-green/15 text-avito-green-dark font-bold text-xs flex items-center justify-center gap-1.5">
          <span>Выполнено</span>
        </div>
      ) : (
        <Button
          size="sm"
          onClick={() => onComplete?.(id)}
          disabled={isCompleting}
          isLoading={isCompleting}
          className="w-full"
        >
          Выполнить
        </Button>
      )}
    </div>
  )
}
