import { HeaderTasks } from '@/widgets/header'
import { TasksListWidget } from '@/widgets/tasks-list'

export const Tasks = () => {
  return (
    <div className="space-y-8">
      <HeaderTasks />
      <TasksListWidget />
    </div>
  )
}
