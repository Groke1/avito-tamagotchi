import { useGetDailyStatQuery } from '@/entities/user'
import { CalendarCheck, Flame, Sparkles } from 'lucide-react'

export const HeaderDailyReport = () => {
  const { data } = useGetDailyStatQuery()

  const streak = data?.streak ?? 0
  const dailyXp = data?.pet.daily_gained_xp ?? 0

  return (
    <header className="flex items-center justify-between pb-6 border-b border-surface-high flex-wrap gap-4">
      <div className="space-y-1">
        <h1 className="text-3xl font-extrabold text-on-surface flex items-center gap-2">
          <span>Ежедневный отчёт</span>
          <CalendarCheck className="size-8 text-avito-blue shrink-0" />
        </h1>
        <p className="text-sm font-semibold text-on-surface-variant">
          Ваша сводная активность и полученные награды за сегодняшний день
        </p>
      </div>
      <div className="bg-surface-low border border-surface-highest px-4.5 py-2 rounded-full shadow-xs">
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2 text-on-surface font-bold text-sm">
            <Flame className="size-5 text-avito-red shrink-0" />
            <span>
              Стрик: {streak} {streak === 1 ? 'день' : streak > 1 && streak < 5 ? 'дня' : 'дней'}
            </span>
          </div>
          <span className="w-px bg-outline-variant h-4" />
          <div className="flex items-center gap-2 text-on-surface-variant font-semibold text-sm">
            <Sparkles className="size-4.5 text-avito-yellow shrink-0" />
            <span>+{dailyXp} XP за день</span>
          </div>
        </div>
      </div>
    </header>
  )
}
