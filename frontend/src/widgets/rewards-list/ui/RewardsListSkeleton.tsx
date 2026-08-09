import { RewardCardSkeleton } from '@/entities/reward'

export const RewardsListSkeleton = () => {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <RewardCardSkeleton />
      <RewardCardSkeleton />
      <RewardCardSkeleton />
    </div>
  )
}
