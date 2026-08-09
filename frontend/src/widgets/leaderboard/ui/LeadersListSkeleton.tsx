import { LeaderItemSkeleton } from '@/entities/leaderboard'

interface LeadersListSkeletonProps {
  count?: number
}

export const LeadersListSkeleton = ({ count = 5 }: LeadersListSkeletonProps) => {
  return (
    <div className="flex flex-col gap-1 mt-4">
      {Array.from({ length: count }).map((_, index) => (
        <LeaderItemSkeleton key={index} />
      ))}
    </div>
  )
}
