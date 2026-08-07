import { AlertCircle, RefreshCw } from 'lucide-react'

interface ErrorStateProps {
  title?: string
  message?: string
  onRetry?: () => void
}

export const ErrorState = ({
  title = 'Произошла ошибка',
  message = 'Не удалось загрузить данные',
  onRetry,
}: ErrorStateProps) => {
  return (
    <div className="p-4 rounded-xl border border-red-200 bg-red-50 flex items-center justify-between gap-3 text-red-700">
      <div className="flex items-center gap-2 text-xs font-medium">
        <AlertCircle className="w-4 h-4 shrink-0 text-red-500" />
        <span>{message || title}</span>
      </div>
      {onRetry && (
        <button
          type="button"
          onClick={onRetry}
          className="px-3 py-1.5 bg-white text-red-700 hover:bg-red-100 rounded-lg text-xs font-bold border border-red-200 flex items-center gap-1 transition-colors cursor-pointer shrink-0"
        >
          <RefreshCw className="w-3.5 h-3.5" />
          Повторить
        </button>
      )}
    </div>
  )
}
