import type { UserReward } from '@/entities/reward'
import { RewardCard } from '@/entities/reward'
import { EmptyState, ErrorState, TabsContent } from '@/shared/ui'
import type { FilterTab } from '../model/types'
import { RewardsListSkeleton } from './RewardsListSkeleton'

interface RewardsContentListProps {
  tabId: FilterTab
  rewards: UserReward[]
  isLoading: boolean
  isError: boolean
  redeemingPromoCode: string | null
  onRedeem: (promoCode: string) => void
  onRetry: () => void
}

export const RewardsContentList = ({
  tabId,
  rewards,
  isLoading,
  isError,
  redeemingPromoCode,
  onRedeem,
  onRetry,
}: RewardsContentListProps) => {
  const getEmptyMessage = (id: FilterTab) => {
    if (id === 'active') return 'У вас пока нет активных промокодов ✨'
    if (id === 'used') return 'Нет использованных или просроченных наград'
    return 'У вас пока нет доступных наград ✨'
  }

  return (
    <TabsContent value={tabId} className="space-y-4">
      {isLoading && <RewardsListSkeleton />}
      {isError && <ErrorState message="Не удалось загрузить список наград" onRetry={onRetry} />}
      {!isLoading && !isError && rewards.length === 0 && (
        <EmptyState message={getEmptyMessage(tabId)} />
      )}
      {!isLoading && !isError && rewards.length > 0 && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {rewards.map((reward) => (
            <RewardCard
              key={reward.reward_id}
              reward={reward}
              onRedeem={onRedeem}
              isRedeeming={redeemingPromoCode === reward.promo_code}
            />
          ))}
        </div>
      )}
    </TabsContent>
  )
}
