import { TabsTrigger } from '@/shared/ui'
import type { FilterTab } from '../model/types'

interface RewardsTabItemProps {
  id: FilterTab
  label: string
  count: number
  showCount: boolean
}

export const RewardsTabItem = ({ id, label, count, showCount }: RewardsTabItemProps) => {
  return (
    <TabsTrigger
      value={id}
      className="relative pb-2 pt-1 px-1 font-medium text-base text-on-surface-variant shadow-none! cursor-pointer transition-all hover:text-on-surface  data-[state=active]:text-avito-blue-dark after:absolute after:bottom-0 after:left-0 after:right-0 after:h-[2.5px] after:bg-avito-blue-dark after:opacity-0 data-[state=active]:after:opacity-100 after:transition-opacity whitespace-nowrap"
    >
      <span>
        {label}
        {showCount ? ` (${count})` : ''}
      </span>
    </TabsTrigger>
  )
}
