import { useHandleFeed, useHandleJourney, useHandleStroke } from '@/features/pet-actions'
import { useAppSelector } from '@/shared/model'

export const usePetActionButtons = () => {
  const { handleFeed, isFeedLoading } = useHandleFeed()
  const { handleStroke, isStrokeLoading } = useHandleStroke()
  const { handleJourney, isJourneyLoading } = useHandleJourney()
  const pet = useAppSelector((state) => state.pet.pet)

  const isFullSatiety = (pet?.satiety ?? 0) >= 100
  const isFullHappiness = (pet?.happiness ?? 0) >= 100

  return {
    handleFeed,
    handleJourney,
    handleStroke,
    isFullSatiety,
    isFullHappiness,
    isFeedLoading,
    isJourneyLoading,
    isStrokeLoading,
  }
}
