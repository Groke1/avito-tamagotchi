import type { UserRewardsResponse } from '@/entities/reward'
import { useState } from 'react'
import type { FilterTab } from './types'

export const useRewardsListTabs = (rewardsData?: UserRewardsResponse) => {
  const [activeTab, setActiveTab] = useState<FilterTab>('all')

  const items = rewardsData?.items ?? []
  const activeCount = items.filter((r) => r.status === 'active').length
  const usedCount = items.filter((r) => r.status === 'redeemed').length

  const getFilteredRewards = (tab: FilterTab) => {
    if (tab === 'active') return items.filter((r) => r.status === 'active')
    if (tab === 'used') return items.filter((r) => r.status === 'redeemed')
    return items
  }

  const handleSetTab = (tab: unknown) => {
    setActiveTab(tab as FilterTab)
  }

  return { items, activeTab, activeCount, usedCount, getFilteredRewards, handleSetTab }
}
