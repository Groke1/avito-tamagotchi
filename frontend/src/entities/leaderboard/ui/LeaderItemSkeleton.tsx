import { Skeleton } from '@/shared/ui'

export const LeaderItemSkeleton = () => {
  return (
    <div className="flex items-center justify-between py-3.5 px-4 rounded-2xl border border-transparent">
      <div className="flex items-center gap-4 sm:gap-6">
        <div className="w-8 flex justify-center">
          <Skeleton className="w-5 h-6 rounded-md" />
        </div>
        <Skeleton className="size-12 sm:size-13 rounded-full shrink-0" />
        <div className="space-y-1.5">
          <Skeleton className="h-5 w-28 sm:w-36 rounded-md" />
          <Skeleton className="h-4 w-36 sm:w-44 rounded-md" />
        </div>
      </div>
      <Skeleton className="h-6 w-16 rounded-md shrink-0" />
    </div>
  )
}
