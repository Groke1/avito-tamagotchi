import { Header } from '@/widgets/header'
import {
  DailyProgressWidget,
  LiveActivityBanner,
  PetShowcase,
  TodayQuestsWidget,
} from '@/widgets/pet-dashboard'

export const Dashboard = () => {
  return (
    <div className="h-full space-y-6">
      <Header />
      <LiveActivityBanner />
      <div className="grid grid-cols-1 lg:grid-cols-12 gap-6 items-center">
        <div className="lg:col-span-7">
          <PetShowcase />
        </div>
        <div className="lg:col-span-5 space-y-6">
          <DailyProgressWidget />
          <TodayQuestsWidget />
        </div>
      </div>
    </div>
  )
}
