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
      className="px-4 py-5 rounded-full font-semibold text-md text-on-surface bg-surface-lowest outline  outline-surface-highest shadow-level-1 transition-all duration-200 cursor-pointer data-[state=active]:bg-avito-blue-container data-[state=active]:text-[#003A5C] data-[state=active]:outline-transparent"
    >
      <span>{label}</span>
      {showCount && <span>({count})</span>}
    </TabsTrigger>
  )
}
