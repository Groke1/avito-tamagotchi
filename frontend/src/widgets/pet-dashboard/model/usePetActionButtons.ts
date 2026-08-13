import { useHandleFeed, useHandleStroke } from '@/features/pet-actions'
import { useAppSelector } from '@/shared/model'
import { useHandleTrip } from './useHandleTrip'

export const usePetActionButtons = () => {
  const { handleFeed, isFeedLoading } = useHandleFeed()
  const { handleStroke, isStrokeLoading } = useHandleStroke()
  const { handleTrip, isTripLoading } = useHandleTrip()

  const pet = useAppSelector((state) => state.pet.pet)

  const isFullSatiety = (pet?.satiety ?? 0) >= 100
  const isFullHappiness = (pet?.happiness ?? 0) >= 100

  return {
    handleFeed,
    handleStroke,
    handleTrip,
    isFullSatiety,
    isFullHappiness,
    isFeedLoading,
    isStrokeLoading,
    isTripLoading,
  }
}
