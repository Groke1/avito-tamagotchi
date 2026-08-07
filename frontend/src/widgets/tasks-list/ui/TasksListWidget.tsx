import { TaskCard, TaskCardSkeleton } from '@/entities/task'
import { useCompleteTaskMutation, useGetTasksQuery } from '@/entities/task/api/taskApi'
import { cn } from '@/shared/lib/utils'
import { EmptyState, ErrorState } from '@/shared/ui'
import { useState } from 'react'
import { toast } from 'sonner'

type FilterTab = 'all' | 'active' | 'completed'

const TABS: { id: FilterTab; label: string }[] = [
  { id: 'all', label: 'Все задания' },
  { id: 'active', label: 'Активные' },
  { id: 'completed', label: 'Выполненные' },
]

export const TasksListWidget = () => {
  const [activeTab, setActiveTab] = useState<FilterTab>('all')
  const [completingTaskId, setCompletingTaskId] = useState<string | null>(null)

  const { data: tasksData, isLoading, isError, refetch } = useGetTasksQuery()
  const [completeTask] = useCompleteTaskMutation()

  const handleComplete = async (taskId: string) => {
    try {
      setCompletingTaskId(taskId)
      await completeTask(taskId).unwrap()
      toast.success('Задание успешно выполнено! Награда получена 🪙✨')
    } catch (err: unknown) {
      const errorObj = err as { status?: number }
      if (errorObj?.status === 409) {
        toast.error('Награда за это задание уже получена!')
      } else {
        toast.error('Не удалось выполнить задание. Попробуйте позже.')
      }
    } finally {
      setCompletingTaskId(null)
    }
  }

  const items = tasksData?.items ?? []
  const filteredTasks = items.filter((task) => {
    if (activeTab === 'active') return task.status === 'active'
    if (activeTab === 'completed') return task.status === 'completed'

    return true
  })

  const activeCount = items.filter((t) => t.status === 'active').length
  const completedCount = items.filter((t) => t.status === 'completed').length

  return (
    <section className="space-y-6">
      <div className="flex items-center gap-2 p-1.5 bg-surface-lowest border border-outline-variant rounded-2xl w-fit shadow-sm">
        {TABS.map((tab) => {
          const count =
            tab.id === 'all' ? items.length : tab.id === 'active' ? activeCount : completedCount

          return (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              className={cn(
                'px-4 py-2 rounded-xl font-bold text-xs transition-all flex items-center gap-2 cursor-pointer',
                activeTab === tab.id
                  ? 'bg-avito-blue text-white shadow-sm'
                  : 'text-on-surface-variant hover:text-on-surface hover:bg-surface-bg',
              )}
            >
              <span>{tab.label}</span>
              {!isLoading && !isError && (
                <span
                  className={cn(
                    'px-2 py-0.5 rounded-full text-[10px] font-extrabold',
                    activeTab === tab.id
                      ? 'bg-white/20 text-white'
                      : 'bg-surface-high text-on-surface-variant',
                  )}
                >
                  {count}
                </span>
              )}
            </button>
          )
        })}
      </div>
      {isLoading && (
        <div className="flex flex-col gap-4">
          <TaskCardSkeleton />
          <TaskCardSkeleton />
          <TaskCardSkeleton />
        </div>
      )}
      {isError && <ErrorState message="Не удалось загрузить список заданий" onRetry={refetch} />}
      {!isLoading && !isError && filteredTasks.length === 0 && (
        <EmptyState
          message={
            activeTab === 'completed'
              ? 'У вас пока нет выполненных заданий 🎯'
              : activeTab === 'active'
                ? 'Все задания выполнены! Отличная работа 🎉'
                : 'На сегодня нет доступных заданий ✨'
          }
        />
      )}
      {!isLoading && !isError && filteredTasks.length > 0 && (
        <div className="flex flex-col gap-4">
          {filteredTasks.map((task) => (
            <TaskCard
              key={task.id}
              task={task}
              onComplete={handleComplete}
              isCompleting={completingTaskId === task.id}
            />
          ))}
        </div>
      )}
    </section>
  )
}
