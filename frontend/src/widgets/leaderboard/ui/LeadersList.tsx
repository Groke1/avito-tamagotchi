import { LeaderItem, type LeaderboardItem } from '@/entities/leaderboard'
import { Skeleton } from '@/shared/ui'
import type { FC } from 'react'

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
    return (
      <div className="space-y-3 mt-6">
        {Array.from({ length: 5 }).map((_, index) => (
          <div key={index} className="flex items-center justify-between py-3 px-4">
            <div className="flex items-center gap-4">
              <Skeleton className="w-8 h-8 rounded-full" />
              <Skeleton className="size-12 rounded-full" />
              <div className="space-y-2">
                <Skeleton className="h-5 w-32 rounded-md" />
                <Skeleton className="h-4 w-24 rounded-md" />
              </div>
            </div>
            <Skeleton className="h-6 w-16 rounded-md" />
          </div>
        ))}
      </div>
    )
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
            isCurrentUser={!!isCurrent}
          />
        )
      })}
    </ul>
  )
}
