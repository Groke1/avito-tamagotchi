interface EmptyStateProps {
  message?: string
}

export const EmptyState = ({ message = 'Нет доступных данных ✨' }: EmptyStateProps) => {
  return (
    <div className="p-6 text-center  text-on-surface-variant bg-surface-bg rounded-xl border-2 border-dashed border-surface-highest">
      {message}
    </div>
  )
}
