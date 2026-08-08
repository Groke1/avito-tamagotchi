import { TaskMiniCardSkeleton } from '@/entities/task'

export const TodayQuestsSkeleton = () => {
  return (
    <ul className="flex flex-col gap-3">
      <TaskMiniCardSkeleton />
      <TaskMiniCardSkeleton />
      <TaskMiniCardSkeleton />
    </ul>
  )
}
