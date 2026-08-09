import { DailyReportWidget } from '@/widgets/daily-report'
import { HeaderDailyReport } from '@/widgets/header'

export const DailyReport = () => {
  return (
    <div className="space-y-8">
      <HeaderDailyReport />
      <DailyReportWidget />
    </div>
  )
}
