import { useAppDispatch } from '@/app/store/hooks'
import { setPet, useFeedPetMutation } from '@/entities/pet'
import { isApiError, isFetchBaseQueryError } from '@/shared/lib/guards'
import { toast } from 'sonner'

export const useHandleFeed = () => {
  const [feedPet, { isLoading }] = useFeedPetMutation()
  const dispatch = useAppDispatch()

  const handleFeed = async () => {
    try {
      const updatedPet = await feedPet().unwrap()
      dispatch(setPet(updatedPet))
      toast.success('Вы покормили питомца! 🍎 (+5 сытости, +2 XP)')
    } catch (error: unknown) {
      if (isFetchBaseQueryError(error) && isApiError(error.data)) {
        toast.error(error.data.message || 'Это действие пока недоступно')
      } else {
        toast.error('Это действие пока недоступно')
      }
    }
  }

  return { isFeedLoading: isLoading, handleFeed }
}
