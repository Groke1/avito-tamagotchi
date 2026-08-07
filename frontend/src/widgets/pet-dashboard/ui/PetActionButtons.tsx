import { Button } from '@/shared/ui'
import { usePetActionButtons } from '../model/usePetActionButtons'

export const PetActionButtons = () => {
  const {
    handleFeed,
    handleStroke,
    isFeedLoading,
    isFullHappiness,
    isFullSatiety,
    isStrokeLoading,
  } = usePetActionButtons()

  return (
    <div className="grid grid-cols-2 gap-3 mt-4">
      <Button
        variant="outline"
        size="sm"
        isLoading={isFeedLoading}
        disabled={isFeedLoading || isStrokeLoading}
        onClick={handleFeed}
        title={isFullSatiety ? 'Питомец полностью сыт! (100/100)' : ''}
        className="hover:bg-avito-green/10 hover:border-avito-green/60 hover:text-avito-green-dark transition-all duration-200"
      >
        🍎 Покормить (-5 монет, +2 XP)
      </Button>
      <Button
        variant="outline"
        size="sm"
        isLoading={isStrokeLoading}
        disabled={isFeedLoading || isStrokeLoading}
        onClick={handleStroke}
        title={isFullHappiness ? 'Питомец уже максимально счастлив! (100/100)' : ''}
        className="hover:bg-avito-blue/10 hover:border-avito-blue/60 hover:text-avito-blue-dark transition-all duration-200"
      >
        🖐️ Погладить (-7 монет, +3 XP)
      </Button>
    </div>
  )
}
