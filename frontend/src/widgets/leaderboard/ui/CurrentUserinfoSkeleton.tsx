import { Skeleton } from '@/shared/ui'

export const CurrentUserinfoSkeleton = () => {
  return (
    <div className="bg-avito-blue/10 border-2 border-avito-blue/30 rounded-card p-3.5 sm:p-5 flex items-center justify-between flex-wrap sm:flex-nowrap gap-3.5 sm:gap-4 shadow-xs animate-pulse">
      <div className="flex items-center gap-3 sm:gap-3.5 min-w-0 flex-1">
        <Skeleton className="size-11 sm:size-14 rounded-full shrink-0" />
        <div className="space-y-1.5 flex-1">
          <Skeleton className="h-5 sm:h-6 w-28 sm:w-32 rounded-lg" />
          <Skeleton className="h-4 w-20 sm:w-24 rounded-md" />
        </div>
      </div>
      <div className="flex flex-col sm:items-end gap-2 w-full sm:w-auto shrink-0 mt-1 sm:mt-0">
        <Skeleton className="h-6 w-full sm:w-60 rounded-full" />
        <Skeleton className="h-9 w-full sm:w-32 rounded-xl" />
      </div>
    </div>
  )
}
