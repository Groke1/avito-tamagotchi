import { HeaderDashboard } from '@/widgets/header'
import {
  PetShowcase,
  StreakEntry,
  TodayQuestsWidget,
  TripResultWidget,
} from '@/widgets/pet-dashboard'
import { PetHungry } from '@/widgets/pet-hungry'

export const Dashboard = () => {
  return (
    <div className="h-full space-y-6">
      <HeaderDashboard />
      <PetHungry />
      <div className="grid grid-cols-1 lg:grid-cols-10 gap-6 ">
        <div className="lg:col-span-5 space-y-6">
          <PetShowcase />
          <TripResultWidget />
        </div>
        <div className="lg:col-span-5 space-y-6">
          <StreakEntry />
          <TodayQuestsWidget />
        </div>
      </div>
    </div>
  )
}
