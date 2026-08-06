import { useAppDispatch } from '@/app/store/hooks'
import { setPet, useStrokePetMutation } from '@/entities/pet'
import { isApiError, isFetchBaseQueryError } from '@/shared/lib/guards'
import { toast } from 'sonner'

export const useHandleStroke = () => {
  const [strokePet, { isLoading }] = useStrokePetMutation()
  const dispatch = useAppDispatch()

  const handleStroke = async () => {
    try {
      const updatedPet = await strokePet().unwrap()
      dispatch(setPet(updatedPet))
      toast.success('Вы погладили питомца! 🖐️ (+5 счастья, +3 XP)')
    } catch (error: unknown) {
      if (isFetchBaseQueryError(error) && isApiError(error.data)) {
        toast.error(error.data.message || 'Это действие пока недоступно')
      } else {
        toast.error('Это действие пока недоступно')
      }
    }
  }

  return { isStrokeLoading: isLoading, handleStroke }
}
