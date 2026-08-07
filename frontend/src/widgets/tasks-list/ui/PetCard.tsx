import { useAppSelector } from '@/shared/model'
import { LazyImage, Progress } from '@/shared/ui'
import { Utensils } from 'lucide-react'

export const PetCard = () => {
  const pet = useAppSelector((state) => state.pet.pet)

  const satietyClamped = Math.min(Math.max(pet?.satiety ?? 0, 0), 100)

  return (
    <div className="bg-surface-lowest shadow-level-1 rounded-card p-4">
      <div className="flex flex-col gap-4 items-center text-center">
        <div className="size-48 rounded-full flex items-center justify-center shrink-0">
          <LazyImage src="/avito-kot.png" alt="Питомец Авито Тамагочи" className="object-contain" />
        </div>
        <div className="space-y-2">
          <h2 className="text-on-surface text-lg font-bold">Ваш питомец ждёт!</h2>
          <p className=" text-on-surface-variant">
            Каждое выполненное задание кормит вашего питомца{' '}
            <span className="text-avito-green-dark">(+15% к Сытости)</span> и продвигает к
            следующему уровню!
          </p>
        </div>
        <div className="w-full bg-surface-highest rounded-xl p-3">
          <div className="space-y-1.5">
            <div className="flex justify-between items-center text-xs font-bold">
              <span className="text-on-surface-variant flex items-center gap-1.5">
                <Utensils className="size-4 text-avito-red" />
                <span>Сытость</span>
              </span>
              <span className="text-on-surface">{satietyClamped}%</span>
            </div>
            <Progress
              value={satietyClamped}
              indicatorClassName="bg-avito-red"
              className="bg-surface-bg"
            />
          </div>
        </div>
      </div>
    </div>
  )
}
