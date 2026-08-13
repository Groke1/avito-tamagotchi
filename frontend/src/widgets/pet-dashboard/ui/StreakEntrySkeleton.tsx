import { Skeleton } from '@/shared/ui'

export const StreakEntrySkeleton = () => {
  return (
    <section className="bg-surface-lowest p-6 rounded-card shadow-level-1 flex flex-col gap-4">
      <div className="flex items-center justify-between">
        <Skeleton className="h-8 w-56 sm:w-64 rounded-lg" />
      </div>
      <div className="flex flex-wrap w-full lg:grid lg:grid-cols-7 gap-4 justify-between items-center">
        {Array.from({ length: 7 }).map((_, index) => (
          <div key={index} className="flex flex-col gap-1.5 justify-center items-center">
            <Skeleton className="size-8 sm:size-9 rounded-full" />
            <Skeleton className="h-3.5 w-10 rounded-md" />
          </div>
        ))}
      </div>
      <Skeleton className="h-14 w-full rounded-xl" />
    </section>
  )
}
