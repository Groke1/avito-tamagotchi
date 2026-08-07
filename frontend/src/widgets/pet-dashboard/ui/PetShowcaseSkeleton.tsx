import { Skeleton } from '@/shared/ui'

export const PetShowcaseSkeleton = () => (
  <div className="bg-surface-lowest rounded-card shadow-level-1">
    <div className="flex justify-center bg-surface-bg rounded-tl-card rounded-tr-card p-4 shadow-level-1">
      <Skeleton className="size-56 rounded-full" />
    </div>
    <div className="p-6">
      <div className="flex flex-col gap-6">
        <div className="space-y-2 pb-4 border-b border-b-surface-high">
          <Skeleton className="h-10 w-48 rounded-lg" />
          <Skeleton className="h-4 w-72 rounded-md" />
        </div>
        <div className="space-y-4">
          <Skeleton className="h-6 w-full rounded-full" />
          <Skeleton className="h-6 w-full rounded-full" />
          <Skeleton className="h-6 w-full rounded-full" />
        </div>
        <div className="grid grid-cols-2 gap-3 mt-4">
          <Skeleton className="h-9 w-full rounded-md" />
          <Skeleton className="h-9 w-full rounded-md" />
        </div>
      </div>
    </div>
  </div>
)
