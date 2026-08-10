import { TabsTrigger } from '@/shared/ui'
import type { FilterTab } from '../model/types'

interface TasksTabItemProps {
  id: FilterTab
  label: string
  count: number
  showCount: boolean
}

export const TasksTabItem = ({ id, label, count, showCount }: TasksTabItemProps) => {
  return (
    <TabsTrigger
      value={id}
      className="sm:px-4 sm:py-5 px-3.5 py-2.5 rounded-full font-semibold text-xs sm:text-sm text-on-surface bg-surface-lowest outline outline-surface-highest transition-all duration-200 cursor-pointer data-[state=active]:bg-avito-blue-container data-[state=active]:text-[#003A5C] data-[state=active]:outline-transparent whitespace-nowrap"
    >
      <span>{label}</span>
      {showCount && <span>({count})</span>}
    </TabsTrigger>
  )
}
