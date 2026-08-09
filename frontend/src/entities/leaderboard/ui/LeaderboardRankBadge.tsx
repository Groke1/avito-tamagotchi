import { cn } from '@/shared/lib/utils'
import { LazyImage } from '@/shared/ui'
import type { LeaderboardItem } from '../model/types'

interface LeaderboardRankBadgeProps {
  leader?: LeaderboardItem
  place: 1 | 2 | 3
  className?: string
}

const PLACE_CONFIG = {
  1: {
    avatar:
      'size-24 sm:size-28 border-avito-yellow ring-4 ring-amber-300/40 shadow-lg shadow-amber-400/40',
    userName: 'font-extrabold text-lg sm:text-xl max-w-37.5 sm:max-w-50',
    petName: 'font-semibold text-xs sm:text-sm max-w-37.5 sm:max-w-50',
    levelBadge: 'bg-amber-50 border-amber-300 text-amber-800 px-3.5',
    pedestal: 'w-32 sm:w-44 h-36 sm:h-44 bg-[#FFC800] shadow-md',
    number: 'text-6xl drop-shadow-sm',
    container: 'z-10',
  },
  2: {
    avatar: 'size-20 sm:size-24 border-slate-300 shadow-md shadow-slate-400/30',
    userName: 'font-bold text-base sm:text-lg max-w-32.5 sm:max-w-42.5',
    petName: 'font-medium text-xs sm:text-sm max-w-32.5 sm:max-w-42.5',
    levelBadge: 'bg-slate-100 border-slate-300 text-slate-700 px-3',
    pedestal: 'w-28 sm:w-36 h-28 sm:h-32 bg-[#C5C8CC] shadow-sm',
    number: 'text-5xl',
    container: '',
  },
  3: {
    avatar: 'size-20 sm:size-24 border-amber-700/60 shadow-md shadow-amber-800/30',
    userName: 'font-bold text-base sm:text-lg max-w-32.5 sm:max-w-42.5',
    petName: 'font-medium text-xs sm:text-sm max-w-32.5 sm:max-w-42.5',
    levelBadge: 'bg-amber-100/70 border-amber-300 text-amber-900 px-3',
    pedestal: 'w-28 sm:w-36 h-20 sm:h-24 bg-[#C66B27] shadow-sm',
    number: 'text-5xl',
    container: '',
  },
} as const

export const LeaderboardRankBadge = ({ leader, place, className }: LeaderboardRankBadgeProps) => {
  const userName = leader?.user_name
  const petName = leader?.pet_name
  const level = leader?.level

  const config = PLACE_CONFIG[place]

  return (
    <div className={cn('flex flex-col items-center justify-end', config.container, className)}>
      <div className="relative group">
        <div
          className={cn(
            'rounded-full border-4 p-1 bg-white flex items-center justify-center overflow-hidden transition-transform duration-300 group-hover:scale-105',
            config.avatar,
          )}
        >
          <LazyImage
            src="/avito-kot.png"
            alt={petName}
            className="size-full object-contain rounded-full"
          />
        </div>
      </div>
      <div className="text-center mt-3 mb-3 space-y-0.5">
        <h4 className={cn('text-on-surface truncate', config.userName)}>{userName}</h4>
        <p className={cn('text-on-surface-variant truncate', config.petName)}>{petName}</p>
        {level !== undefined && (
          <div className="pt-1">
            <span
              className={cn(
                'inline-block py-0.5 rounded-full border font-bold text-xs shadow-xs',
                config.levelBadge,
              )}
            >
              Lvl {level}
            </span>
          </div>
        )}
      </div>
      <div className={cn('rounded-t-[20px] flex items-center justify-center', config.pedestal)}>
        <span className={cn('text-white font-extrabold select-none', config.number)}>{place}</span>
      </div>
    </div>
  )
}
