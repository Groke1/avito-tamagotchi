import { cn } from '@/shared/lib/utils'
import { Button } from '@/shared/ui'
import { Check, Coins, Sparkles } from 'lucide-react'
import type { Task } from '../model/types'

interface TaskCardProps {
  task: Task
  onComplete?: (taskId: string) => void
  isCompleting?: boolean
}

const categoryMap: Record<string, { label: string; badgeClass: string; emoji: string }> = {
  'Первая продажа месяца': {
    label: 'Продажи',
    badgeClass: 'bg-emerald-100 text-emerald-800',
    emoji: '🏷️',
  },
  'Лояльный продавец': {
    label: 'Ежедневное',
    badgeClass: 'bg-blue-100 text-avito-blue-dark',
    emoji: '❤️',
  },
  'Быстрый ответ покупателю': {
    label: 'Доставка',
    badgeClass: 'bg-pink-100 text-pink-800',
    emoji: '📦',
  },
}

export const TaskCard = ({ task, onComplete, isCompleting = false }: TaskCardProps) => {
  const { id, title, description, reward_coins, reward_xp, status } = task
  const isCompleted = status === 'completed'
  const category = categoryMap[title] ?? {
    label: 'Задание',
    badgeClass: 'bg-blue-100 text-avito-blue-dark',
    emoji: '✨',
  }

  return (
    <div
      className={cn(
        'p-5 rounded-card border bg-surface-lowest transition-all flex flex-col justify-between gap-4 shadow-level-1 relative',
        isCompleted ? 'border-avito-green/40 bg-avito-green/5' : 'border-surface-highest ',
      )}
    >
      <div className="flex items-center justify-between gap-2">
        <span
          className={cn(
            'px-3 py-1 rounded-full text-xs font-bold',
            isCompleted ? 'bg-surface-high text-on-surface-variant' : category.badgeClass,
          )}
        >
          {category.label}
        </span>
        {isCompleted && (
          <div className="inline-flex items-center p-1.5 rounded-full bg-avito-green/20 text-surface-lowest">
            <Check className="size-5 stroke-[2.5]" />
          </div>
        )}
      </div>
      <div className="space-y-1.5">
        <h4 className="font-extrabold text-base text-on-surface leading-tight">
          {category.emoji} {title}
        </h4>
        <p className="text-xs text-on-surface-variant leading-relaxed line-clamp-2">
          {description}
        </p>
      </div>
      <div className="flex items-center gap-2 pt-0.5">
        <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-avito-yellow/15 text-yellow-700">
          <Coins className="w-3 h-3 text-yellow-600" />+{reward_coins} 🪙
        </span>
        <span className="inline-flex items-center gap-1 px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-avito-blue/10 text-avito-blue-dark">
          <Sparkles className="w-3 h-3 text-avito-blue" />+{reward_xp} XP
        </span>
      </div>
      {isCompleted ? (
        <div className="w-full py-2.5 px-4 rounded-xl bg-avito-green/15 text-avito-green-dark font-bold text-xs flex items-center justify-center gap-1.5">
          <span>Выполнено</span>
        </div>
      ) : (
        <Button
          size={'sm'}
          onClick={() => onComplete?.(id)}
          disabled={isCompleting}
          isLoading={isCompleting}
        >
          Выполнить
        </Button>
      )}
    </div>
  )
}
