export const TaskMiniCardSkeleton = () => {
  return (
    <li className="p-4 rounded-xl border border-outline-variant bg-surface-lowest flex items-center justify-between gap-4 animate-pulse">
      <div className="flex items-center gap-3.5 flex-1">
        <div className="w-10 h-10 rounded-full bg-surface-high shrink-0" />
        <div className="flex-1 space-y-2">
          <div className="h-4 bg-surface-high rounded w-3/5" />
          <div className="h-3 bg-surface-high rounded w-4/5" />
          <div className="flex gap-2 pt-1">
            <div className="h-4 w-16 bg-surface-high rounded-full" />
            <div className="h-4 w-16 bg-surface-high rounded-full" />
          </div>
        </div>
      </div>
      <div className="w-20 h-8 bg-surface-high rounded-xl shrink-0" />
    </li>
  )
}
