import type { TasksResponse } from '@/entities/task'
import { useState } from 'react'
import type { FilterTab } from './types'

export const useTasksListTabs = (tasksData?: TasksResponse) => {
  const [activeTab, setActiveTab] = useState<FilterTab>('all')

  const items = tasksData?.items ?? []
  const activeCount = items.filter((t) => t.status === 'active').length
  const completedCount = items.filter((t) => t.status === 'completed').length

  const getFilteredTasks = (tab: FilterTab) => {
    if (tab === 'active') return items.filter((t) => t.status === 'active')
    if (tab === 'completed') return items.filter((t) => t.status === 'completed')

    return items
  }

  const handleSetTab = (tab: unknown) => {
    setActiveTab(tab as FilterTab)
  }

  return { items, activeTab, activeCount, completedCount, getFilteredTasks, handleSetTab }
}
