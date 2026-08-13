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
      'size-16 sm:size-24 lg:size-28 border-avito-yellow ring-2 sm:ring-4 ring-amber-300/40 shadow-lg shadow-amber-400/40',
    userName: 'font-extrabold text-xs sm:text-base lg:text-xl max-w-24 sm:max-w-36 lg:max-w-50',
    petName: 'font-semibold text-[10px] sm:text-xs lg:text-sm max-w-24 sm:max-w-36 lg:max-w-50',
    levelBadge: 'bg-amber-50 border-amber-300 text-amber-800 px-2 sm:px-3.5 text-[10px] sm:text-xs',
    pedestal: 'w-22 sm:w-36 lg:w-44 h-24 sm:h-36 lg:h-44 bg-[#FFC800] shadow-md',
    number: 'text-3xl sm:text-5xl lg:text-6xl drop-shadow-sm',
    container: 'z-10',
  },
  2: {
    avatar: 'size-13 sm:size-20 lg:size-24 border-slate-300 shadow-md shadow-slate-400/30',
    userName: 'font-bold text-[11px] sm:text-sm lg:text-lg max-w-20 sm:max-w-32 lg:max-w-42.5',
    petName: 'font-medium text-[10px] sm:text-xs lg:text-sm max-w-20 sm:max-w-32 lg:max-w-42.5',
    levelBadge:
      'bg-slate-100 border-slate-300 text-slate-700 px-1.5 sm:px-3 text-[10px] sm:text-xs',
    pedestal: 'w-18 sm:w-28 lg:w-36 h-18 sm:h-28 lg:h-32 bg-[#C5C8CC] shadow-sm',
    number: 'text-2xl sm:text-4xl lg:text-5xl',
    container: '',
  },
  3: {
    avatar: 'size-13 sm:size-20 lg:size-24 border-amber-700/60 shadow-md shadow-amber-800/30',
    userName: 'font-bold text-[11px] sm:text-sm lg:text-lg max-w-20 sm:max-w-32 lg:max-w-42.5',
    petName: 'font-medium text-[10px] sm:text-xs lg:text-sm max-w-20 sm:max-w-32 lg:max-w-42.5',
    levelBadge:
      'bg-amber-100/70 border-amber-300 text-amber-900 px-1.5 sm:px-3 text-[10px] sm:text-xs',
    pedestal: 'w-18 sm:w-28 lg:w-36 h-14 sm:h-20 lg:h-24 bg-[#C66B27] shadow-sm',
    number: 'text-2xl sm:text-4xl lg:text-5xl',
    container: '',
  },
} as const

export const LeaderboardRankBadge = ({ leader, place, className }: LeaderboardRankBadgeProps) => {
  const userName = leader?.user_name
  const petName = leader?.pet_name
  const level = leader?.level

  const config = PLACE_CONFIG[place]

  return (
    <div
      className={cn('flex flex-col items-center justify-end min-w-0', config.container, className)}
    >
      <div className="relative group">
        <div
          className={cn(
            'rounded-full border-2 sm:border-4 p-0.5 sm:p-1 bg-white flex items-center justify-center overflow-hidden transition-transform duration-300 group-hover:scale-105',
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
      <div className="text-center mt-1.5 sm:mt-3 mb-1.5 sm:mb-3 space-y-0.5 min-w-0 w-full px-0.5">
        <h4 className={cn('text-on-surface truncate mx-auto', config.userName)}>{userName}</h4>
        <p className={cn('text-on-surface-variant truncate mx-auto', config.petName)}>{petName}</p>
        {level !== undefined && (
          <div className="pt-0.5 sm:pt-1">
            <span
              className={cn(
                'inline-block py-0.5 rounded-full border font-bold shadow-xs',
                config.levelBadge,
              )}
            >
              Lvl {level}
            </span>
          </div>
        )}
      </div>
      <div
        className={cn(
          'rounded-t-xl sm:rounded-t-[20px] flex items-center justify-center',
          config.pedestal,
        )}
      >
        <span className={cn('text-white font-extrabold select-none', config.number)}>{place}</span>
      </div>
    </div>
  )
}
