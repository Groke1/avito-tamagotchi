import type { LeaderboardItem } from '@/entities/leaderboard'
import { ROUTES_PATHS } from '@/shared/config'
import { useAppSelector } from '@/shared/model'
import { Button } from '@/shared/ui'
import { Star, Zap } from 'lucide-react'
import type { FC } from 'react'
import { useNavigate } from 'react-router-dom'
import { CurrentUserinfoSkeleton } from './CurrentUserinfoSkeleton'

interface CurrentUserinfoProps {
  user?: LeaderboardItem
  isLoading?: boolean
}

export const CurrentUserinfo: FC<CurrentUserinfoProps> = ({ user, isLoading }) => {
  const navigate = useNavigate()
  const pet = useAppSelector((state) => state.pet.pet)

  if (isLoading || !user) {
    return <CurrentUserinfoSkeleton />
  }

  const rank = user.rank
  const level = user.level
  const xp = pet?.xp ?? 0
  const nextLevelXp = pet?.next_level_xp ?? 120

  return (
    <div className="bg-avito-blue/10 border-2 border-avito-blue/30 rounded-card p-3.5 sm:p-5 flex items-center justify-between flex-wrap sm:flex-nowrap gap-3.5 sm:gap-4 shadow-level-1">
      <div className="flex items-center gap-3 sm:gap-3.5 min-w-0 flex-1">
        <div className="size-11 sm:size-14 rounded-full bg-white border border-avito-blue/30 shadow-xs flex items-center justify-center shrink-0">
          <span className="text-avito-blue-dark font-extrabold text-lg sm:text-2xl tracking-tight">
            #{rank}
          </span>
        </div>
        <div className="min-w-0 flex-1">
          <h3 className="font-extrabold text-base sm:text-xl text-on-surface flex items-center gap-1.5 leading-snug truncate">
            <span>Вы</span>
            <span className="font-bold text-sm sm:text-lg text-on-surface">(Level {level})</span>
          </h3>
          <div className="flex items-center gap-1.5 text-xs sm:text-sm font-bold text-on-surface-variant mt-0.5">
            <div className="size-4 rounded-full bg-avito-blue-dark flex items-center justify-center text-white shrink-0">
              <Star className="size-2.5 fill-white text-white" />
            </div>
            <span>{xp} XP</span>
          </div>
        </div>
      </div>
      <div className="flex flex-col sm:items-end gap-2 w-full sm:w-auto shrink-0 mt-1 sm:mt-0">
        <div className="inline-flex items-center gap-1.5 bg-white/90 backdrop-blur-xs px-2.5 sm:px-3.5 py-1 rounded-full text-[11px] sm:text-xs font-bold text-avito-blue-dark border border-white/80 shadow-2xs max-w-full overflow-hidden text-ellipsis whitespace-nowrap">
          <Zap className="size-3.5 text-amber-500 fill-amber-500 shrink-0" />
          <span className="truncate">До следующего уровня осталось {nextLevelXp} XP!</span>
        </div>
        <Button size="sm" className="w-full sm:w-auto" onClick={() => navigate(ROUTES_PATHS.TASKS)}>
          Получить XP
        </Button>
      </div>
    </div>
  )
}
