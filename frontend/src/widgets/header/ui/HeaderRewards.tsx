import { Coins, Gift } from 'lucide-react'
import { useHeaderRewards } from '../model/useHeaderRewards'

export const HeaderRewards = () => {
  const { coins, activeCount, totalCount } = useHeaderRewards()

  return (
    <header className="flex items-center justify-between pb-4 sm:pb-6 border-b border-surface-high flex-wrap gap-3 sm:gap-4">
      <h1 className="text-2xl sm:text-3xl font-extrabold text-on-surface">Мои Награды 🎁</h1>
      <div className="bg-surface-low border border-surface-highest px-3.5 sm:px-4.5 py-1.5 sm:py-2 rounded-2xl sm:rounded-full">
        <div className="flex items-center gap-3 sm:gap-4 flex-wrap">
          <div className="flex items-center gap-1.5 sm:gap-2 text-on-surface font-bold text-xs sm:text-sm">
            <Coins className="size-4 sm:size-5 text-avito-yellow shrink-0" />
            <span>{coins} монет</span>
          </div>
          <span className="hidden sm:block w-px bg-outline-variant h-4" />
          <div className="flex items-center gap-1.5 sm:gap-2 text-on-surface-variant font-semibold text-xs sm:text-sm">
            <Gift className="size-4 sm:size-4.5 text-avito-blue shrink-0" />
            <span>
              Активных промокодов: {activeCount} из {totalCount}
            </span>
          </div>
        </div>
      </div>
    </header>
  )
}
