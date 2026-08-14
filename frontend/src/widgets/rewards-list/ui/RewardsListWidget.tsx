import { useGetRewardsQuery } from '@/entities/reward'
import { Tabs } from '@/shared/ui'
import { TABS } from '../model/consts'
import { useRedeemReward } from '../model/useRedeemReward'
import { useRewardsListTabs } from '../model/useRewardsListTabs'
import { RewardsContentList } from './RewardsContentList'
import { RewardsTabList } from './RewardsTabList'

export const RewardsListWidget = () => {
  const { data: rewardsData, isLoading, isError, refetch } = useGetRewardsQuery(undefined, {
    refetchOnMountOrArgChange: true,
  })

  const { redeemingPromoCode, handleRedeem } = useRedeemReward()
  const { items, activeCount, usedCount, activeTab, getFilteredRewards, handleSetTab } =
    useRewardsListTabs(rewardsData)

  return (
    <section>
      <Tabs value={activeTab} onValueChange={handleSetTab} className="w-full">
        <RewardsTabList
          itemsCount={items.length}
          activeCount={activeCount}
          usedCount={usedCount}
          showCounts={!isLoading && !isError}
        />
        <div className="mt-6">
          {TABS.map(({ id }) => (
            <RewardsContentList
              key={id}
              tabId={id}
              rewards={getFilteredRewards(id)}
              isLoading={isLoading}
              isError={isError}
              redeemingPromoCode={redeemingPromoCode}
              onRedeem={handleRedeem}
              onRetry={refetch}
            />
          ))}
        </div>
      </Tabs>
    </section>
  )
}
