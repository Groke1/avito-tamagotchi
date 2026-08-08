import { Skeleton } from '@/shared/ui'

export const RewardCardSkeleton = () => {
  return (
    <div className="p-5 rounded-card border border-surface-highest bg-surface-lowest flex flex-col justify-between gap-4 shadow-level-1">
      <div className="flex items-start justify-between gap-2">
        <Skeleton className="w-12 h-12 rounded-full shrink-0" />
        <Skeleton className="h-6 w-20 rounded-full" />
      </div>
      <div className="space-y-1.5">
        <Skeleton className="h-5 w-3/4 rounded-md" />
        <Skeleton className="h-4 w-full rounded-md" />
      </div>
      <div className="p-3 rounded-xl bg-surface-high/60 border border-surface-highest flex items-center justify-between gap-2">
        <div className="space-y-1 flex-1">
          <Skeleton className="h-3 w-14 rounded" />
          <Skeleton className="h-4 w-52 rounded" />
        </div>
        <Skeleton className="h-8 w-8 rounded-lg shrink-0" />
      </div>
      <div className="flex items-center gap-1.5">
        <Skeleton className="h-4 w-28 rounded-md" />
      </div>
      <Skeleton className="w-full h-9 rounded-xl" />
    </div>
  )
}
