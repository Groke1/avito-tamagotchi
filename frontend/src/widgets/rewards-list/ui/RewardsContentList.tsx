import type { UserReward } from '@/entities/reward'
import { RewardCard } from '@/entities/reward'
import { ErrorState, TabsContent } from '@/shared/ui'
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
  // const getEmptyMessage = (id: FilterTab) => {
  //   if (id === 'active') return 'У вас пока нет активных промокодов ✨'
  //   if (id === 'used') return 'Нет использованных или просроченных наград'
  //   return 'У вас пока нет полученных наград ✨'
  // }

  const actualRewards: UserReward[] = rewards.length
    ? rewards
    : [
        {
          reward_id: '3fa85f64-5717-4562-b3fc-2c963f66afa6',
          promo_code: 'DELIVERYDISC-X8K2M7Q9WP',
          name: 'Скидка 10% на доставку',
          description: 'Скидка 10% на одну покупку с Авито Доставкой',
          status: 'active',
          expires_at: '2026-09-06T19:07:12Z',
          redeemed_at: null,
        },
      ]

  return (
    <TabsContent value={tabId} className="space-y-4">
      {isLoading && <RewardsListSkeleton />}
      {isError && <ErrorState message="Не удалось загрузить список наград" onRetry={onRetry} />}
      {/* {!isLoading && !isError && rewards.length === 0 && (
        <EmptyState message={getEmptyMessage(tabId)} />
      )} */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {actualRewards.map((reward) => (
          <RewardCard
            key={reward.reward_id}
            reward={reward}
            onRedeem={onRedeem}
            isRedeeming={redeemingPromoCode === reward.promo_code}
          />
        ))}
      </div>
    </TabsContent>
  )
}
