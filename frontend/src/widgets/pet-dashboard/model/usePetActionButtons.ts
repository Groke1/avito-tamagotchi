import { useHandleFeed, useHandleStroke } from '@/features/pet-actions'
import { useAppSelector } from '@/shared/model'

export const usePetActionButtons = () => {
  const { handleFeed, isFeedLoading } = useHandleFeed()
  const { handleStroke, isStrokeLoading } = useHandleStroke()
  const pet = useAppSelector((state) => state.pet.pet)

  const isFullSatiety = (pet?.satiety ?? 0) >= 100
  const isFullHappiness = (pet?.happiness ?? 0) >= 100

  return {
    handleFeed,
    handleStroke,
    isFullSatiety,
    isFullHappiness,
    isFeedLoading,
    isStrokeLoading,
  }
}
