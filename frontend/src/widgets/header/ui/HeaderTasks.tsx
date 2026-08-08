import { Calendar, CheckCircle2, Coins } from 'lucide-react'
import { useHeaderTasks } from '../model/useHeaderTasks'

export const HeaderTasks = () => {
  const { coins, completedCount, totalCount, formattedDate } = useHeaderTasks()

  return (
    <header className="flex items-center justify-between pb-6 border-b border-surface-high flex-wrap gap-4">
      <div className="space-y-1">
        <h1 className="text-3xl font-extrabold text-on-surface">Центр Заданий 🎯</h1>
        {formattedDate && (
          <p className="text-sm font-semibold text-on-surface-variant flex items-center gap-1.5">
            <Calendar className="size-3.5 text-avito-blue shrink-0" />
            <span>Задания на {formattedDate}</span>
          </p>
        )}
      </div>
      <div className="bg-surface-low border border-surface-highest px-4.5 py-2 rounded-full">
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2 text-on-surface font-bold text-sm">
            <Coins className="size-5 text-avito-yellow shrink-0" />
            <span>{coins} монет</span>
          </div>
          <span className="w-px bg-outline-variant h-4" />
          <div className="flex items-center gap-2 text-on-surface-variant font-semibold text-sm">
            <CheckCircle2 className="size-4.5 text-avito-green shrink-0" />
            <span>
              Выполнено сегодня: {completedCount} из {totalCount}
            </span>
          </div>
        </div>
      </div>
    </header>
  )
}
