import { useGetTasksQuery } from '@/entities/task/api/taskApi'
import { useGetProfileQuery } from '@/entities/user'
import { formatTaskDate } from '@/shared/lib/utils'
import { useAppSelector } from '@/shared/model'

export const useHeaderTasks = () => {
  useGetProfileQuery()
  const user = useAppSelector((state) => state.user.user)
  const { data: tasksData } = useGetTasksQuery()

  const totalCount = tasksData?.items.length ?? 0
  const completedCount = tasksData?.items.filter((t) => t.status === 'completed').length ?? 0
  const formattedDate = formatTaskDate({ dateStr: tasksData?.date })

  return {
    coins: user?.coins ?? 0,
    totalCount,
    completedCount,
    formattedDate,
  }
}
