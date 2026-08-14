import { Button } from '@/shared/ui'
import { usePetActionButtons } from '../model/usePetActionButtons'

export const PetActionButtons = () => {
  const {
    handleFeed,
    handleJourney,
    handleStroke,
    isFeedLoading,
    isFullHappiness,
    isFullSatiety,
    isJourneyLoading,
    isStrokeLoading,
  } = usePetActionButtons()

  return (
    <div className="grid grid-cols-1 xl:grid-cols-2 gap-3 mt-4">
      <Button
        variant="outline"
        size="sm"
        isLoading={isFeedLoading}
        disabled={isFeedLoading || isStrokeLoading}
        onClick={handleFeed}
        title={isFullSatiety ? 'Питомец полностью сыт! (100/100)' : ''}
        className="w-full hover:bg-avito-green/10 hover:border-avito-green/60 hover:text-avito-green-dark transition-all duration-200"
      >
        🍎 Покормить (-5 монет)
      </Button>
      <Button
        variant="outline"
        size="sm"
        isLoading={isStrokeLoading}
        disabled={isFeedLoading || isStrokeLoading || isJourneyLoading}
        onClick={handleStroke}
        title={isFullHappiness ? 'Питомец уже максимально счастлив! (100/100)' : ''}
        className="w-full hover:bg-avito-blue/10 hover:border-avito-blue/60 hover:text-avito-blue-dark transition-all duration-200"
      >
        🖐️ Погладить
      </Button>
      <Button
        size="sm"
        isLoading={isJourneyLoading}
        disabled={isFeedLoading || isStrokeLoading || isJourneyLoading}
        onClick={handleJourney}
        className="col-span-1 xl:col-span-2"
      >
        🧭 Отправить в путешествие
      </Button>
    </div>
  )
}
