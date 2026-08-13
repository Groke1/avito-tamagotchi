import { getDaysPlural } from '@/shared/lib/utils'
import { Sparkles } from 'lucide-react'
import { StreakDays } from './StreakDays'

export const StreakEntry = () => {
  return (
    <section className="bg-surface-lowest p-6 rounded-card shadow-level-1 flex flex-col gap-4">
      <div className="flex items-center justify-between">
        <h2 className="text-xl sm:text-2xl font-bold text-on-surface flex items-center gap-3">
          <span>Твоя серия:</span>
          <span className="text-avito-green">
            {/* TODO: использовать актуальную ручку */}
            {3} {getDaysPlural(3)}! 🔥
          </span>
        </h2>
      </div>
      <StreakDays streak={3} />
      <div className="flex items-center gap-3 bg-avito-blue-container/10 border border-avito-blue-container/30 text-avito-blue-dark rounded-xl p-4 text-xs sm:text-sm">
        <Sparkles className="size-5.5 shrink-0" />
        <p className="font-semibold">
          Заходи каждый день, чтобы получать бонусные монеты и редкие награды!
        </p>
      </div>
    </section>
  )
}
