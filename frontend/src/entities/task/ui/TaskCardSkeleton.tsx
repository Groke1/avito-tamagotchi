import { Skeleton } from '@/shared/ui'

export const TaskCardSkeleton = () => {
  return (
    <div className="p-6 rounded-2xl border border-outline-variant bg-surface-lowest flex flex-col md:flex-row md:items-center justify-between gap-5 shadow-level-1">
      <div className="flex items-start md:items-center gap-4 flex-1 min-w-0">
        <Skeleton className="w-12 h-12 rounded-2xl shrink-0" />
        <div className="space-y-2 flex-1 min-w-0">
          <Skeleton className="h-5 w-2/5 rounded-md" />
          <Skeleton className="h-4 w-4/5 rounded-md" />
          <div className="flex items-center gap-2.5 pt-1">
            <Skeleton className="h-6 w-20 rounded-full" />
            <Skeleton className="h-6 w-20 rounded-full" />
          </div>
        </div>
      </div>
      <div className="shrink-0">
        <Skeleton className="w-full md:w-32 h-10 rounded-xl" />
      </div>
    </div>
  )
}
