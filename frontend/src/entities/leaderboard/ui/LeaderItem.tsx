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
        'flex items-center justify-between py-3.5 px-4 rounded-2xl transition-colors duration-150',
        isCurrentUser
          ? 'bg-avito-blue/15 border border-avito-blue/40 shadow-level-1'
          : 'hover:bg-surface-low/60 border border-transparent',
        className,
      )}
    >
      <div className="flex items-center gap-4 sm:gap-6">
        <span className="w-8 text-center text-xl sm:text-2xl font-extrabold text-on-surface">
          {leader.rank}
        </span>
        <div className="size-12 sm:size-13 rounded-full border border-surface-container bg-white p-0.5 shadow-level-1 overflow-hidden shrink-0 flex items-center justify-center">
          <LazyImage
            src="/avito-kot.png"
            alt={leader.pet_name || 'Питомец'}
            className="size-full object-contain rounded-full"
          />
        </div>
        <div>
          <div className="flex items-center gap-2">
            <h4 className="font-bold text-base sm:text-lg text-on-surface leading-tight">
              {leader.user_name || 'Игрок'}
            </h4>
            {isCurrentUser && (
              <span className="px-2 py-0.5 rounded-md text-[10px] font-black uppercase tracking-wider bg-avito-blue text-white shadow-level-1">
                ВЫ
              </span>
            )}
          </div>
          <div className="flex items-center gap-1.5 text-xs sm:text-sm font-semibold text-on-surface-variant mt-0.5">
            <span>Lvl {leader.level}</span>
            <span>•</span>
            <span className="flex items-center gap-1 text-avito-red font-bold">
              <PawPrint className="size-3.5 fill-avito-red shrink-0" />
              <span>{leader.pet_name}</span>
            </span>
          </div>
        </div>
      </div>
      <div className="text-right">
        <span className="font-extrabold text-base sm:text-lg text-avito-blue-dark">
          {leader.xp} XP
        </span>
      </div>
    </li>
  )
}
