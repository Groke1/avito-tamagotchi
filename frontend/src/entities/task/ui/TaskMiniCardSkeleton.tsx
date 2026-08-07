import { Skeleton } from '@/shared/ui'

export const TaskMiniCardSkeleton = () => {
  return (
    <li className="p-4 rounded-xl border border-outline-variant bg-surface-lowest flex items-center justify-between gap-4">
      <div className="flex items-center gap-3.5 min-w-0 flex-1">
        <Skeleton className="w-10 h-10 rounded-full shrink-0" />
        <div className="min-w-0 flex-1 space-y-1.5">
          <Skeleton className="h-5 w-2/5 rounded-md" />
          <Skeleton className="h-3 w-4/5 rounded-md" />
          <Skeleton className="h-3 rounded-md" />
          <div className="flex items-center gap-2 pt-0.5">
            <Skeleton className="h-4.5 w-16 rounded-full" />
            <Skeleton className="h-4.5 w-16 rounded-full" />
          </div>
        </div>
      </div>
      <div className="shrink-0">
        <Skeleton className="w-24 h-8 rounded-xl" />
      </div>
    </li>
  )
}
