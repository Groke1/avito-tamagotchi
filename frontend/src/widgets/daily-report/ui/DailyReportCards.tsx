import type { DailyStatResponse } from '@/entities/user'
import { getDaysPlural } from '@/shared/lib/utils'
import { CheckSquare, Flame, Gift, Sparkles } from 'lucide-react'
import { useDailyReportCards } from '../model/useDailyReportCards'

interface DailyReportCardsProps {
  data: DailyStatResponse
}

export const DailyReportCards = ({ data }: DailyReportCardsProps) => {
  const { completedTasksCount, earnedRewardsCount, totalCoinsEarned, totalXp } =
    useDailyReportCards(data)

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <div className="bg-surface-lowest rounded-card p-5 flex items-center gap-4 shadow-level-1">
        <div className="size-14 flex items-center justify-center shrink-0 bg-avito-yellow/20 text-avito-yellow rounded-full">
          <Sparkles className="size-6" />
        </div>
        <div>
          <p className="text-xs font-semibold text-on-surface-variant uppercase tracking-wider">
            Заработано
          </p>
          <p className="text-2xl font-extrabold text-on-surface mt-0.5">+{totalXp} XP</p>
        </div>
      </div>
      <div className="bg-surface-lowest rounded-card p-5 flex items-center gap-4 shadow-level-1">
        <div className="size-14 flex items-center justify-center shrink-0 bg-avito-blue/20 text-avito-blue rounded-full">
          <CheckSquare className="size-6" />
        </div>
        <div>
          <p className="text-xs font-semibold text-on-surface-variant uppercase tracking-wider">
            Выполнено заданий
          </p>
          <p className="text-2xl font-extrabold text-on-surface mt-0.5">
            {completedTasksCount}{' '}
            {totalCoinsEarned > 0 && (
              <span className="text-xs font-bold text-avito-yellow">+{totalCoinsEarned} монет</span>
            )}
          </p>
        </div>
      </div>
      <div className="bg-surface-lowest rounded-card p-5 flex items-center gap-4 shadow-level-1">
        <div className="size-14 flex items-center justify-center shrink-0 bg-purple-500/20 text-purple-500 rounded-full">
          <Gift className="size-6" />
        </div>
        <div>
          <p className="text-xs font-semibold text-on-surface-variant uppercase tracking-wider">
            Получено наград
          </p>
          <p className="text-2xl font-extrabold text-on-surface mt-0.5">{earnedRewardsCount}</p>
        </div>
      </div>
      <div className="bg-surface-lowest rounded-card p-5 flex items-center gap-4 shadow-level-1">
        <div className="size-14 flex items-center justify-center shrink-0 bg-avito-red/20 text-avito-red rounded-full">
          <Flame className="size-6" />
        </div>
        <div>
          <p className="text-xs font-semibold text-on-surface-variant uppercase tracking-wider">
            Серия дней
          </p>
          <p className="text-2xl font-extrabold text-on-surface mt-0.5">
            {data.streak} {getDaysPlural(data.streak)}
          </p>
        </div>
      </div>
    </div>
  )
}
