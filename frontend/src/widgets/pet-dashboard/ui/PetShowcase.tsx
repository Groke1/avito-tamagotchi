import { useAppSelector } from '@/shared/model'
import { LazyImage } from '@/shared/ui'
import { PetActionButtons } from './PetActionButtons'
import { PetShowcaseSkeleton } from './PetShowcaseSkeleton'
import { PetVitalityBars } from './PetVitalityBars'

export const PetShowcase = () => {
  const pet = useAppSelector((state) => state.pet.pet)

  if (!pet) return <PetShowcaseSkeleton />

  const { name, level, happiness, satiety, xp, next_level_xp } = pet
  const totalLevelXp = xp + next_level_xp

  return (
    <div className="bg-surface-lowest rounded-card shadow-level-1">
      <div className="flex justify-center bg-surface-bg rounded-tl-card rounded-tr-card p-4 shadow-level-1">
        <div className="size-56 rounded-full flex items-center justify-center shrink-0">
          <LazyImage src="/avito-kot.png" alt="Питомец Авито Тамагочи" className="object-contain" />
        </div>
      </div>
      <div className="p-6">
        <div className="flex flex-col gap-6">
          <div className="space-y-2 pb-4 border-b border-b-surface-high">
            <h2 className="font-extrabold text-5xl text-on-surface truncate">{name}</h2>
            <p className="text-sm text-on-surface-variant">
              Level {level} ({xp}/{totalLevelXp} XP) — До Level {level + 1} осталось {next_level_xp}{' '}
              XP!
            </p>
          </div>
          <PetVitalityBars
            happiness={happiness}
            satiety={satiety}
            totalLevelXp={totalLevelXp}
            xp={xp}
          />
          <PetActionButtons />
        </div>
      </div>
    </div>
  )
}
