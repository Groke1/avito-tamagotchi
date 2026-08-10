import { type LeaderboardItem, LeaderboardRankBadge } from '@/entities/leaderboard'
import { type FC, useMemo } from 'react'

interface LeaderboardTopProps {
  leaders?: LeaderboardItem[]
}

export const LeaderboardTop: FC<LeaderboardTopProps> = ({ leaders }) => {
  const firstPlaceItem = useMemo(() => leaders?.find((l) => l.rank === 1), [leaders])
  const secondPlaceItem = useMemo(() => leaders?.find((l) => l.rank === 2), [leaders])
  const thirdPlaceItem = useMemo(() => leaders?.find((l) => l.rank === 3), [leaders])

  return (
    <div className="bg-surface-lowest rounded-card p-4 sm:p-8 lg:p-10 shadow-level-1 overflow-hidden">
      <h2 className="text-xl sm:text-2xl lg:text-3xl font-extrabold text-on-surface text-center mb-5 sm:mb-8 lg:mb-10 tracking-tight">
        Зал славы
      </h2>
      <div className="flex justify-center items-end gap-1.5 sm:gap-6 lg:gap-10 pb-2">
        <LeaderboardRankBadge place={2} leader={secondPlaceItem} />
        <LeaderboardRankBadge place={1} leader={firstPlaceItem} />
        <LeaderboardRankBadge place={3} leader={thirdPlaceItem} />
      </div>
    </div>
  )
}
