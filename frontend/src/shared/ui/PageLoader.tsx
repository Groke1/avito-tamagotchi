import type { FC } from 'react'

interface PageLoaderProps {
  label?: string
}

export const PageLoader: FC<PageLoaderProps> = ({ label = 'Загрузка...' }) => {
  return (
    <div className="min-h-screen flex items-center justify-center bg-surface-bg">
      <div className="flex flex-col items-center gap-3">
        <div className="w-10 h-10 border-4 border-avito-green border-t-transparent rounded-full animate-spin" />
        <span className="text-sm font-semibold text-on-surface-variant">{label}</span>
      </div>
    </div>
  )
}
