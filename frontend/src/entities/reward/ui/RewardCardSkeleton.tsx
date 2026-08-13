import { Skeleton } from '@/shared/ui'

export const RewardCardSkeleton = () => {
  return (
    <div className="p-3.5 sm:p-5 rounded-card border border-surface-highest bg-surface-lowest flex flex-col justify-between gap-3 sm:gap-4 shadow-level-1 min-w-0">
      <div className="flex items-start justify-between gap-2">
        <Skeleton className="w-10 h-10 sm:w-12 sm:h-12 rounded-full shrink-0" />
        <Skeleton className="h-6 w-20 rounded-full" />
      </div>
      <div className="space-y-1.5 min-w-0">
        <Skeleton className="h-5 w-3/4 rounded-md" />
        <Skeleton className="h-4 w-full rounded-md" />
      </div>
      <div className="p-2 sm:p-2.5 rounded-xl bg-surface-high/70 border border-dashed border-surface-highest flex items-center justify-between gap-2 min-w-0">
        <div className="space-y-1 flex-1 min-w-0">
          <Skeleton className="h-2.5 w-14 rounded" />
          <Skeleton className="h-4 w-28 sm:w-36 max-w-full rounded" />
        </div>
        <Skeleton className="size-7 sm:size-8 rounded-lg shrink-0" />
      </div>
      <div className="flex items-center gap-1.5 min-w-0">
        <Skeleton className="h-3.5 w-28 rounded-md" />
      </div>
      <Skeleton className="w-full h-8 sm:h-9 rounded-xl" />
    </div>
  )
}
