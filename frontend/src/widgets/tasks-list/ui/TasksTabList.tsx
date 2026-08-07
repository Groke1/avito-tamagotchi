import { TabsList } from '@/shared/ui'
import { TABS } from '../model/consts'
import type { FilterTab } from '../model/types'
import { TasksTabItem } from './TasksTabItem'

interface TasksTabListProps {
  itemsCount: number
  activeCount: number
  completedCount: number
  showCounts: boolean
}

export const TasksTabList = ({
  itemsCount,
  activeCount,
  completedCount,
  showCounts,
}: TasksTabListProps) => {
  const getTabCount = (tabId: FilterTab) => {
    if (tabId === 'active') return activeCount
    if (tabId === 'completed') return completedCount
    return itemsCount
  }

  return (
    <TabsList className="bg-surface-lowest border border-outline-variant p-1.5 rounded-2xl h-auto gap-1.5 shadow-sm">
      {TABS.map((tab) => (
        <TasksTabItem
          key={tab.id}
          id={tab.id}
          label={tab.label}
          count={getTabCount(tab.id)}
          showCount={showCounts}
        />
      ))}
    </TabsList>
  )
}
