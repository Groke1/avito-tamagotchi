import { TasksListWidget } from '@/widgets/tasks-list'

export const Tasks = () => {
  return (
    <div className="space-y-6">
      <div className="space-y-1">
        <h1 className="text-3xl font-extrabold text-on-surface">Центр Заданий</h1>
      </div>
      <TasksListWidget />
    </div>
  )
}
