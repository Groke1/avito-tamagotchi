import { HeaderLeaderboard } from '@/widgets/header'
import { LeaderboardWidget } from '@/widgets/leaderboard'

export const Leaderboard = () => {
  return (
    <div className="space-y-8">
      <HeaderLeaderboard />
      <LeaderboardWidget />
    </div>
  )
}
