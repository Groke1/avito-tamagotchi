import { Skeleton } from '@/shared/ui'

export const TaskCardSkeleton = () => {
  return (
    <div className="p-4 rounded-card border border-surface-highest bg-surface-lowest flex flex-col justify-between gap-4 shadow-level-1">
      <div className="flex items-center justify-between gap-2">
        <Skeleton className="h-5 w-20 rounded-full" />
      </div>
      <div className="space-y-1">
        <Skeleton className="h-5 w-3/6 rounded-md" />
        <Skeleton className="h-3.5 w-5/6 rounded-md" />
      </div>
      <div className="flex items-center gap-2 pt-0.5">
        <Skeleton className="h-5 w-16 rounded-full" />
        <Skeleton className="h-5 w-16 rounded-full" />
      </div>
      <Skeleton className="w-full h-8 rounded-xl" />
    </div>
  )
}
