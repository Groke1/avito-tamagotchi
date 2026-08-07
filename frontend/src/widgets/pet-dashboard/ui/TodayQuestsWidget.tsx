import { TaskMiniCard, type TasksResponse } from '@/entities/task'
// import { useGetTasksQuery } from '@/entities/task/api/taskApi'
import { useMemo } from 'react'

export const TodayQuestsWidget = () => {
  // const { data: tasks } = useGetTasksQuery()

  const tasks: TasksResponse | undefined = useMemo(() => {
    return {
      date: '2026-08-07',
      items: [
        {
          id: 'e9b11c75-392c-473d-815a-52ef389d3110',
          title: 'Первая продажа месяца',
          description:
            'Успешно завершите сделку по любому объявлению с подключенной Авито Доставкой.',
          reward_coins: 500,
          reward_xp: 1000,
          status: 'active',
          completed_at: null,
        },
        {
          id: 'cd51eda4-c593-4037-a0c1-c5bf3e7da86f',
          title: 'Лояльный продавец',
          description: 'Получите новый отзыв с оценкой 5 звезд от верифицированного покупателя.',
          reward_coins: 300,
          reward_xp: 500,
          status: 'active',
          completed_at: null,
        },
      ],
    }
  }, [])

  const activeTasks = useMemo(() => {
    return tasks?.items.filter(({ status }) => status === 'active')
  }, [tasks])

  return (
    <div className="bg-surface-lowest p-6 rounded-card shadow-level-1">
      <div className="flex flex-col gap-4">
        <div className="flex items-center justify-between gap-4">
          <h4 className="text-on-surface font-bold text-2xl">Задания на сегодня</h4>
          <div className="font-bold text-sm text-avito-blue-dark">Все задания</div>
        </div>
        <ul className="flex flex-col gap-2">
          {activeTasks?.map((task) => (
            <TaskMiniCard key={task.id} {...task} />
          ))}
        </ul>
      </div>
    </div>
  )
}
