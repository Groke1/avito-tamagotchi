import { HeaderDashboard } from '@/widgets/header'
import { PetShowcase, TodayQuestsWidget } from '@/widgets/pet-dashboard'

export const Dashboard = () => {
  return (
    <div className="h-full space-y-6">
      <HeaderDashboard />
      <div className="grid grid-cols-1 lg:grid-cols-10 gap-6 ">
        <div className="lg:col-span-5">
          <PetShowcase />
        </div>
        <div className="lg:col-span-5 space-y-6">
          <TodayQuestsWidget />
        </div>
      </div>
    </div>
  )
}
