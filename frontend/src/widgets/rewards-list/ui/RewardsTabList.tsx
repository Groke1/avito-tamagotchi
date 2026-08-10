import { TabsList } from '@/shared/ui'
import { TABS } from '../model/consts'
import type { FilterTab } from '../model/types'
import { RewardsTabItem } from './RewardsTabItem'

interface RewardsTabListProps {
  itemsCount: number
  activeCount: number
  usedCount: number
  showCounts: boolean
}

export const RewardsTabList = ({
  itemsCount,
  activeCount,
  usedCount,
  showCounts,
}: RewardsTabListProps) => {
  const getTabCount = (tabId: FilterTab) => {
    if (tabId === 'active') return activeCount
    if (tabId === 'used') return usedCount

    return itemsCount
  }

  return (
    <div className="w-full overflow-x-auto scrollbar-none pb-1">
      <TabsList className="gap-2 sm:gap-4 flex-nowrap sm:flex-wrap min-w-max">
        {TABS.map(({ id, label }) => (
          <RewardsTabItem
            key={id}
            id={id}
            label={label}
            count={getTabCount(id)}
            showCount={showCounts}
          />
        ))}
      </TabsList>
    </div>
  )
}
