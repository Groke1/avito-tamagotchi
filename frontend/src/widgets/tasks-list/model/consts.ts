import type { FilterTab } from './types'

export const TABS: { id: FilterTab; label: string }[] = [
  { id: 'all', label: 'Все задания' },
  { id: 'active', label: 'Активные' },
  { id: 'completed', label: 'Выполненные' },
]
