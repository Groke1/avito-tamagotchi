import { TaskMiniCard, type TasksResponse } from '@/entities/task'
// import { useGetTasksQuery } from '@/entities/task/api/taskApi'
import { ROUTES_PATHS } from '@/shared/config'
// import { EmptyState, ErrorState } from '@/shared/ui'
import { NavLink } from 'react-router-dom'

// import { TodayQuestsSkeleton } from './TodayQuestsSkeleton'

const tasksData: TasksResponse = {
  date: 'asd',
  items: [
    {
      id: '3fa85f64-5717-4562-b3fc-2c963f66afa6',
      title: 'Лояльный продавец',
      description: 'Открой три разных объявления в выбранной категории',
      reward_coins: 20,
      reward_xp: 50,
      status: 'active',
      completed_at: null,
    },
  ],
}
export const TodayQuestsWidget = () => {
  // const { data: tasksData, isLoading, isError, refetch } = useGetTasksQuery()

  return (
    <section className="bg-surface-lowest p-6 rounded-card shadow-level-1 flex flex-col gap-4">
      <div className="flex items-center justify-between gap-4">
        <h4 className="text-on-surface font-bold text-2xl">Задания на сегодня</h4>
        <NavLink
          to={ROUTES_PATHS.TASKS}
          className="font-bold text-xs text-avito-blue-dark hover:underline cursor-pointer"
        >
          Все задания
        </NavLink>
      </div>
      {/* {isLoading && <TodayQuestsSkeleton />}
      {isError && <ErrorState message="Не удалось загрузить список заданий" onRetry={refetch} />}
      {!isLoading && !isError && tasksData?.items.length === 0 && (
        <EmptyState message="На сегодня нет доступных заданий ✨" />
      )} */}
      {
        <ul className="flex flex-col gap-3">
          {tasksData.items.slice(0, 3).map((task) => (
            <TaskMiniCard key={task.id} task={task} />
          ))}
        </ul>
      }
    </section>
  )
}
