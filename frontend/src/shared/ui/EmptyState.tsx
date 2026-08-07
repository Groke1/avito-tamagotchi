interface EmptyStateProps {
  message?: string
}

export const EmptyState = ({ message = 'Нет доступных данных ✨' }: EmptyStateProps) => {
  return (
    <div className="p-6 text-center text-xs text-on-surface-variant bg-surface-bg rounded-xl border border-dashed border-outline-variant">
      {message}
    </div>
  )
}
