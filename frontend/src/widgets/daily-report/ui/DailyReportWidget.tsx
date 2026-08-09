import { useGetDailyStatQuery } from '@/entities/user'
import { ErrorState } from '@/shared/ui'
import { DailyReportCards } from './DailyReportCards'
import { DailyReportSkeleton } from './DailyReportSkeleton'
import { DailyReportTimeline } from './DailyReportTimeline'

export const DailyReportWidget = () => {
  const { data, isLoading, isError, refetch } = useGetDailyStatQuery()

  if (isLoading) return <DailyReportSkeleton />

  if (isError || !data) {
    return (
      <ErrorState
        title="Не удалось загрузить отчёт"
        message="Произошла ошибка при получении дневной статистики. Попробуйте обновить страницу."
        onRetry={refetch}
      />
    )
  }

  return (
    <div className="space-y-6">
      <DailyReportCards data={data} />
      <DailyReportTimeline rewards={data.rewards} tasks={data.tasks} />
    </div>
  )
}
