import { setPet, useFeedPetMutation } from '@/entities/pet'
import { isApiError, isFetchBaseQueryError } from '@/shared/lib/guards'
import { useAppDispatch, useAppSelector } from '@/shared/model'
import { toast } from 'sonner'

export const useHandleFeed = () => {
  const [feedPet, { isLoading }] = useFeedPetMutation()
  const dispatch = useAppDispatch()
  const pet = useAppSelector((state) => state.pet.pet)

  const handleFeed = async () => {
    // TODO: remove the processing logic once the backend fix is in place
    if (pet && pet.satiety >= 100) {
      toast.info('Питомец полностью сыт! (100/100)')
      return
    }

    try {
      const updatedPet = await feedPet().unwrap()
      dispatch(setPet(updatedPet))
      toast.success('Вы покормили питомца! 🍎 (+5 сытости, +2 XP)')
    } catch (error: unknown) {
      if (isFetchBaseQueryError(error) && isApiError(error.data)) {
        if (error.data.code === 'PET_ACTION_UNAVAILABLE') {
          toast.info('Питомец полностью сыт!')
        } else {
          toast.error(error.data.message || 'Это действие пока недоступно')
        }
      } else {
        toast.error('Это действие пока недоступно')
      }
    }
  }

  return { isFeedLoading: isLoading, handleFeed }
}
