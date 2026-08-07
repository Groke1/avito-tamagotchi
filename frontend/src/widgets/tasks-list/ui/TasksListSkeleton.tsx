import { TaskCardSkeleton } from '@/entities/task'

export const TasksListSkeleton = () => {
  return (
    <div className="flex flex-col gap-4">
      <TaskCardSkeleton />
      <TaskCardSkeleton />
      <TaskCardSkeleton />
    </div>
  )
}
