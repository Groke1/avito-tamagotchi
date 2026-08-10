import { cn } from '@/shared/lib/utils'
import { LazyImage } from '@/shared/ui'
import { PawPrint } from 'lucide-react'
import type { LeaderboardItem } from '../model/types'

interface LeaderItemProps {
  leader: LeaderboardItem
  isCurrentUser?: boolean
  className?: string
}

export const LeaderItem = ({ leader, isCurrentUser, className }: LeaderItemProps) => {
  return (
    <li
      className={cn(
        'flex items-center justify-between py-2.5 sm:py-3.5 px-2.5 sm:px-4 rounded-card transition-colors duration-150 gap-2 sm:gap-4',
        isCurrentUser
          ? 'bg-avito-blue/15 border border-avito-blue/40 shadow-level-1'
          : 'hover:bg-surface-low/60 border border-transparent',
        className,
      )}
    >
      <div className="flex items-center gap-2.5 sm:gap-4 min-w-0 flex-1">
        <span className="w-5 sm:w-8 text-center text-base sm:text-2xl font-extrabold text-on-surface shrink-0">
          {leader.rank}
        </span>
        <div className="size-10 sm:size-13 rounded-full border border-surface-container bg-white p-0.5 shadow-level-1 overflow-hidden shrink-0 flex items-center justify-center">
          <LazyImage
            src="/avito-kot.png"
            alt={leader.pet_name || 'Питомец'}
            className="size-full object-contain rounded-full"
          />
        </div>
        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-1.5 min-w-0">
            <h4 className="font-bold text-sm sm:text-lg text-on-surface leading-tight truncate">
              {leader.user_name || 'Игрок'}
            </h4>
            {isCurrentUser && (
              <span className="px-1.5 sm:px-2 py-0.5 rounded-md text-[9px] sm:text-[10px] font-black uppercase tracking-wider bg-avito-blue text-white shadow-level-1 shrink-0">
                ВЫ
              </span>
            )}
          </div>
          <div className="flex items-center gap-1 sm:gap-1.5 text-xs sm:text-sm font-semibold text-on-surface-variant mt-0.5 min-w-0">
            <span className="shrink-0">Lvl {leader.level}</span>
            <span className="shrink-0">•</span>
            <span className="flex items-center gap-1 text-avito-red font-bold min-w-0 truncate">
              <PawPrint className="size-3 sm:size-3.5 fill-avito-red shrink-0" />
              <span className="truncate">{leader.pet_name}</span>
            </span>
          </div>
        </div>
      </div>
      <div className="text-right shrink-0">
        <span className="font-extrabold text-sm sm:text-lg text-avito-blue-dark whitespace-nowrap">
          {leader.xp} XP
        </span>
      </div>
    </li>
  )
}
