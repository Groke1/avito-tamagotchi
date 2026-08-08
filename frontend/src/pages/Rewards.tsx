import { HeaderRewards } from '@/widgets/header'
import { RewardsListWidget } from '@/widgets/rewards-list'

export const Rewards = () => {
  return (
    <div className="space-y-8">
      <HeaderRewards />
      <RewardsListWidget />
    </div>
  )
}
