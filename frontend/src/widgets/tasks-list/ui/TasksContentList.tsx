import type { Task } from '@/entities/task'
import { TaskCard } from '@/entities/task'
import { EmptyState, ErrorState, TabsContent } from '@/shared/ui'
import type { FilterTab } from '../model/types'
import { TasksListSkeleton } from './TasksListSkeleton'

interface TasksContentListProps {
  tabId: FilterTab
  tasks: Task[]
  isLoading: boolean
  isError: boolean
  completingTaskId: string | null
  onComplete: (taskId: string) => void
  onRetry: () => void
}

export const TasksContentList = ({
  tabId,
  tasks,
  isLoading,
  isError,
  completingTaskId,
  onComplete,
  onRetry,
}: TasksContentListProps) => {
  const getEmptyMessage = (id: FilterTab) => {
    if (id === 'completed') return 'У вас пока нет выполненных заданий 🎯'
    if (id === 'active') return 'Все задания выполнены! Отличная работа 🎉'
    return 'На сегодня нет доступных заданий ✨'
  }

  return (
    <TabsContent value={tabId} className="mt-6 space-y-4">
      {isLoading && <TasksListSkeleton />}
      {isError && <ErrorState message="Не удалось загрузить список заданий" onRetry={onRetry} />}
      {!isLoading && !isError && tasks.length === 0 && (
        <EmptyState message={getEmptyMessage(tabId)} />
      )}
      {!isLoading && !isError && tasks.length > 0 && (
        <div className="flex flex-col gap-4">
          {tasks.map((task) => (
            <TaskCard
              key={task.id}
              task={task}
              onComplete={onComplete}
              isCompleting={completingTaskId === task.id}
            />
          ))}
        </div>
      )}
    </TabsContent>
  )
}
