import { Skeleton } from '@/shared/ui'

export const DailyReportSkeleton = () => {
  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {Array.from({ length: 4 }).map((_, index) => (
          <div
            key={index}
            className="bg-surface-lowest border border-surface-highest rounded-card p-5 flex items-center gap-4 shadow-level-1"
          >
            <Skeleton className="size-14 rounded-full shrink-0" />
            <div className="space-y-2 flex-1">
              <Skeleton className="h-3 w-20" />
              <Skeleton className="h-5 w-28" />
            </div>
          </div>
        ))}
      </div>
      <div className="bg-surface-lowest border border-surface-highest rounded-card p-8 space-y-6 shadow-level-1">
        <Skeleton className="h-7 w-56" />
        <div className="space-y-4">
          {Array.from({ length: 3 }).map((_, index) => (
            <div
              key={index}
              className="bg-surface-lowest border border-surface-highest p-4 rounded-xl space-y-3"
            >
              <div className="flex items-center justify-between">
                <Skeleton className="h-5 w-64" />
                <Skeleton className="h-4 w-16 rounded-full" />
              </div>
              <Skeleton className="h-4 w-3/4" />
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}
