import { setPet, useStrokePetMutation } from '@/entities/pet'
import { isApiError, isFetchBaseQueryError } from '@/shared/lib/guards'
import { useAppDispatch, useAppSelector } from '@/shared/model'
import { toast } from 'sonner'

export const useHandleStroke = () => {
  const [strokePet, { isLoading }] = useStrokePetMutation()
  const dispatch = useAppDispatch()
  const pet = useAppSelector((state) => state.pet.pet)

  const handleStroke = async () => {
    if (pet && pet.happiness >= 100) {
      toast.info('Питомец уже максимально счастлив! (100/100)')
      return
    }

    try {
      const updatedPet = await strokePet().unwrap()
      dispatch(setPet(updatedPet))
      toast.success('Вы погладили питомца! 🖐️ (+5 счастья, +3 XP)')
    } catch (error: unknown) {
      if (isFetchBaseQueryError(error) && isApiError(error.data)) {
        toast.error(error.data.message)
      } else {
        toast.error('Это действие пока недоступно')
      }
    }
  }

  return { isStrokeLoading: isLoading, handleStroke }
}
