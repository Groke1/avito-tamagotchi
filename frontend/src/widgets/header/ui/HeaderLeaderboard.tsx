import { useGetLeaderboardQuery } from '@/entities/leaderboard'
import { Award, Trophy, Zap } from 'lucide-react'

export const HeaderLeaderboard = () => {
  const { data } = useGetLeaderboardQuery({ limit: 20 })
  const currentUser = data?.current_user

  return (
    <header className="flex items-center justify-between pb-6 border-b border-surface-high flex-wrap gap-4">
      <div className="space-y-1">
        <h1 className="text-3xl font-extrabold text-on-surface flex items-center gap-2">
          <span>Таблица Лидеров</span>
          <Trophy className="size-8 text-avito-yellow shrink-0" />
        </h1>
        <p className="text-sm font-semibold text-on-surface-variant">
          Соревнуйтесь с другими владельцами питомцев Авито и занимайте высшие места!
        </p>
      </div>
      {currentUser && (
        <div className="bg-surface-low border border-surface-highest px-4.5 py-2 rounded-full shadow-xs">
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-2 text-on-surface font-bold text-sm">
              <Award className="size-5 text-avito-blue shrink-0" />
              <span>Ваше место: #{currentUser.rank}</span>
            </div>
            <span className="w-px bg-outline-variant h-4" />
            <div className="flex items-center gap-2 text-on-surface-variant font-semibold text-sm">
              <Zap className="size-4.5 text-avito-green shrink-0" />
              <span>Уровень: {currentUser.level}</span>
            </div>
          </div>
        </div>
      )}
    </header>
  )
}
