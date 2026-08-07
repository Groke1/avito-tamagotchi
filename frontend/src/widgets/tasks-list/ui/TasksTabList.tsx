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
    <TabsList className="gap-2">
      {TABS.map(({ id, label }) => (
        <TasksTabItem
          key={id}
          id={id}
          label={label}
          count={getTabCount(id)}
          showCount={showCounts}
        />
      ))}
    </TabsList>
  )
}
