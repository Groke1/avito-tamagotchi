import { clearLatestTrip, useLazyGetLastTripQuery } from '@/entities/pet'
import { useAppDispatch, useAppSelector } from '@/shared/model'
import { Button } from '@/shared/ui'
import { Coins, Gift, Map, Sparkles, X } from 'lucide-react'
import { useEffect } from 'react'

export const TripResultWidget = () => {
  const dispatch = useAppDispatch()
  const pet = useAppSelector((state) => state.pet.pet)
  const latestTrip = useAppSelector((state) => state.pet.latestTrip)
  const [fetchLastTrip] = useLazyGetLastTripQuery()

  useEffect(() => {
    if (!pet) {
      return
    }

    void fetchLastTrip()
  }, [fetchLastTrip, pet])

  if (!latestTrip) {
    return null
  }

  const reward = latestTrip.reward

  return (
    <section className="bg-surface-lowest rounded-card shadow-level-1 p-6 space-y-5">
      <div className="flex items-start justify-between gap-4">
        <div className="space-y-1">
          <div className="inline-flex items-center gap-2 text-avito-blue-dark text-sm font-semibold">
            <Map className="size-4" />
            Последнее путешествие
          </div>
          <h3 className="text-xl font-extrabold text-on-surface">Питомец вернулся домой</h3>
        </div>
        <Button
          type="button"
          variant="ghost"
          size="icon"
          onClick={() => dispatch(clearLatestTrip())}
          className="shrink-0"
          aria-label="Скрыть результат путешествия"
        >
          <X className="size-4" />
        </Button>
      </div>

      <div className="flex flex-wrap items-center gap-3 text-sm">
        <div className="inline-flex items-center gap-2 rounded-xl bg-avito-green/10 px-3 py-2 font-semibold text-avito-green-dark">
          <Coins className="size-4" />
          +{latestTrip.coins} монет
        </div>
        <div className="inline-flex items-center gap-2 rounded-xl bg-avito-blue/10 px-3 py-2 font-semibold text-avito-blue-dark">
          <Sparkles className="size-4" />
          +{latestTrip.xp} XP
        </div>
        {reward && (
          <div className="inline-flex items-center gap-2 rounded-xl bg-surface-high px-3 py-2 font-semibold text-on-surface">
            <Gift className="size-4" />
            {reward.name}
          </div>
        )}
      </div>

      <div className="space-y-2">
        <p className="text-sm font-semibold text-on-surface">История</p>
        <p className="text-sm leading-7 text-on-surface-variant">{latestTrip.story}</p>
      </div>

      {reward ? (
        <div className="space-y-2 border-t border-surface-high pt-4">
          <p className="text-sm font-semibold text-on-surface">Награда из путешествия</p>
          <div className="space-y-1">
            <p className="text-sm font-bold text-on-surface">{reward.name}</p>
            <p className="text-sm text-on-surface-variant">{reward.description}</p>
          </div>
          <div className="text-xs font-mono font-semibold text-on-surface-variant">
            Промокод: {reward.promo_code}
          </div>
        </div>
      ) : null}
    </section>
  )
}
