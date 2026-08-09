import { LeaderItem, type LeaderboardItem } from '@/entities/leaderboard'
import type { FC } from 'react'
import { LeadersListSkeleton } from './LeadersListSkeleton'

export interface LeadersListProps {
  leaders?: LeaderboardItem[]
  currentUserRank?: number
  currentUserName?: string
  isLoading?: boolean
}

export const LeadersList: FC<LeadersListProps> = ({
  leaders,
  currentUserRank,
  currentUserName,
  isLoading,
}) => {
  if (isLoading) {
    return <LeadersListSkeleton count={10} />
  }

  return (
    <ul className="flex flex-col gap-1 mt-4">
      {leaders?.map((leader) => {
        const isCurrent =
          currentUserRank === leader.rank ||
          (currentUserName &&
            (leader.user_name || '').toLowerCase() === currentUserName.toLowerCase())

        return (
          <LeaderItem
            key={`${leader.rank}-${leader.user_name}`}
            leader={leader}
            isCurrentUser={Boolean(isCurrent)}
          />
        )
      })}
    </ul>
  )
}
