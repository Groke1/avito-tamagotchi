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
      className="px-4 py-2.5 rounded-xl font-bold text-xs transition-all flex items-center gap-2 cursor-pointer data-[state=active]:bg-avito-blue data-[state=active]:text-white data-[state=active]:shadow-sm text-on-surface-variant hover:text-on-surface"
    >
      <span>{label}</span>
      {showCount && (
        <span className="px-2 py-0.5 rounded-full text-[10px] font-extrabold bg-surface-high text-on-surface-variant group-data-[state=active]:bg-white/20 group-data-[state=active]:text-white">
          {count}
        </span>
      )}
    </TabsTrigger>
  )
}
