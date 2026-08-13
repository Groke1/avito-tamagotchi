import { cn } from '@/shared/lib/utils'
import { Check, Gift, Lock } from 'lucide-react'
import type { FC } from 'react'
import { useStreakDays } from '../model/useStreakDays'

interface StreakDaysProps {
  streak: number
}

export const StreakDays: FC<StreakDaysProps> = ({ streak }) => {
  const { daysCycle } = useStreakDays(streak)

  return (
    <div className="flex flex-wrap w-full lg:grid lg:grid-cols-7 gap-4">
      {daysCycle.map(({ label, dayStep, isDone, isNextGoal, isFuture, isBonus, isToday }) => {
        return (
          <div
            key={dayStep}
            className="flex flex-col gap-1.5 justify-center items-center text-center"
          >
            <div
              className={cn(
                'flex items-center justify-center rounded-full size-8 sm:size-9 border-2 transition-all',
                isDone && 'bg-avito-green border-avito-green text-avito-green-dark',
                isNextGoal &&
                  'bg-transparent border-avito-green-dark text-avito-green-dark font-bold',
                isFuture && !isBonus && 'bg-surface-high border-transparent text-on-surface/30',
                isBonus && !isDone && 'bg-avito-red/10 border-transparent text-avito-red',
                !isDone && !isNextGoal && !isFuture && 'bg-surface-low border-surface-container',
              )}
            >
              {isDone && <Check className="size-4 shrink-0 stroke-3 text-avito-green-dark" />}
              {isNextGoal && !isBonus && dayStep}
              {isFuture && !isBonus && (
                <Lock className="size-4 shrink-0 stroke-[2.5] text-on-surface/30" />
              )}
              {isBonus && !isDone && (
                <Gift className="size-4.5 shrink-0 stroke-[2.5] text-avito-red/60" />
              )}
            </div>
            <span
              className={cn(
                'text-xs font-semibold',
                isToday ? 'text-on-surface font-bold' : 'text-on-surface-variant/70',
              )}
            >
              {label}
            </span>
          </div>
        )
      })}
    </div>
  )
}
