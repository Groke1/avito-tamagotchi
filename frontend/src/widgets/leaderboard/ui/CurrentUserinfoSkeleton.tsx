import { Skeleton } from '@/shared/ui'

export const CurrentUserinfoSkeleton = () => {
  return (
    <div className="bg-avito-blue/10 border-2 border-avito-blue/30 rounded-card p-4 sm:p-5 flex items-center justify-between flex-wrap gap-4 shadow-xs animate-pulse">
      <div className="flex items-center gap-3.5">
        <Skeleton className="size-13 sm:size-14 rounded-full" />
        <div className="space-y-2">
          <Skeleton className="h-6 w-32 rounded-lg" />
          <Skeleton className="h-4 w-24 rounded-md" />
        </div>
      </div>
      <div className="flex flex-col items-end gap-2 shrink-0 ml-auto sm:ml-0">
        <Skeleton className="h-6 w-60 rounded-full" />
        <Skeleton className="h-9 w-32 rounded-xl" />
      </div>
    </div>
  )
}
