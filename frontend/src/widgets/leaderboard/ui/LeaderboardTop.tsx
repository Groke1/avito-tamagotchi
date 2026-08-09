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
    <div className="bg-surface-lowest rounded-card p-6 sm:p-10 shadow-level-1 overflow-hidden">
      <h2 className="text-2xl sm:text-3xl font-extrabold text-on-surface text-center mb-8 sm:mb-10 tracking-tight">
        Зал славы
      </h2>
      <div className="flex justify-center items-end gap-3 sm:gap-10 pb-2">
        <LeaderboardRankBadge place={2} leader={secondPlaceItem} />
        <LeaderboardRankBadge place={1} leader={firstPlaceItem} />
        <LeaderboardRankBadge place={3} leader={thirdPlaceItem} />
      </div>
    </div>
  )
}
