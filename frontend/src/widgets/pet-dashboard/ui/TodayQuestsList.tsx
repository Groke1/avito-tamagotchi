import { TaskMiniCard } from '@/entities/task'
import { useGetTasksQuery } from '@/entities/task/api/taskApi'
import { EmptyState, ErrorState } from '@/shared/ui'
import { TodayQuestsSkeleton } from './TodayQuestsSkeleton'

export const TodayQuestsList = () => {
  const { data: tasksData, isLoading, isError, refetch } = useGetTasksQuery()

  if (isLoading) return <TodayQuestsSkeleton />
  if (isError) return <ErrorState message="Не удалось загрузить список заданий" onRetry={refetch} />
  if (!tasksData || tasksData.items.length === 0) {
    return <EmptyState message="На сегодня нет доступных заданий ✨" />
  }

  return (
    <ul className="flex flex-col gap-3">
      {tasksData.items.slice(0, 3).map((task) => (
        <TaskMiniCard key={task.id} task={task} />
      ))}
    </ul>
  )
}
