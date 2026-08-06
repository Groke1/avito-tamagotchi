import { Progress } from '@/shared/ui'
import { Smile, Sparkles, Utensils } from 'lucide-react'
import type { FC } from 'react'
import { usePetVitalityBars } from '../model/usePetVitalityBars'

interface PetVitalityBarsProps {
  satiety: number
  happiness: number
  totalLevelXp: number
  xp?: number
}

export const PetVitalityBars: FC<PetVitalityBarsProps> = ({
  satiety,
  happiness,
  totalLevelXp,
  xp = 0,
}) => {
  const { satietyClamped, happinessClamped, xpPercentage } = usePetVitalityBars({
    satiety,
    happiness,
    xp,
    totalLevelXp,
  })

  return (
    <div className="flex flex-col gap-4">
      <div className="space-y-1.5">
        <div className="flex justify-between items-center text-xs font-bold">
          <span className="text-on-surface-variant flex items-center gap-1.5">
            <Utensils className="size-4 text-avito-green" />
            <span>Сытость</span>
          </span>
          <span className="text-on-surface">{satietyClamped}%</span>
        </div>
        <Progress value={satietyClamped} indicatorClassName="bg-avito-green" />
      </div>
      <div className="space-y-1.5">
        <div className="flex justify-between items-center text-xs font-bold">
          <span className="text-on-surface-variant flex items-center gap-1.5">
            <Smile className="size-4 text-avito-blue" />
            <span>Счастье</span>
          </span>
          <span className="text-on-surface">{happinessClamped}%</span>
        </div>
        <Progress value={happinessClamped} indicatorClassName="bg-avito-blue" />
      </div>
      <div className="space-y-1.5">
        <div className="flex justify-between items-center text-xs font-bold">
          <span className="text-on-surface-variant flex items-center gap-1.5">
            <Sparkles className="size-4 text-avito-yellow" />
            <span>Опыт уровня</span>
          </span>
          <span className="text-on-surface">{xpPercentage}%</span>
        </div>
        <Progress value={xpPercentage} indicatorClassName="bg-avito-yellow" />
      </div>
    </div>
  )
}
