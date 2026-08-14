import { clearLatestTrip, useMakeTripMutation } from '@/entities/pet'
import { isApiError, isFetchBaseQueryError } from '@/shared/lib/guards'
import { useAppDispatch, useAppSelector } from '@/shared/model'
import type { PetErrorCode } from '@/shared/model/types'
import { toast } from 'sonner'

export const useHandleJourney = () => {
  const dispatch = useAppDispatch()
  const pet = useAppSelector((state) => state.pet.pet)
  const [makeTrip, { isLoading }] = useMakeTripMutation()

  const handleJourney = async () => {
    if (!pet) {
      toast.error('Сначала создайте питомца')
      return
    }

    try {
      await makeTrip(pet.id).unwrap()
      dispatch(clearLatestTrip())
      toast.success('Питомец отправился в путешествие! 🧭', {
        description: 'Когда путешествие завершится, мы покажем историю и награду.',
      })
    } catch (error: unknown) {
      if (isFetchBaseQueryError(error) && isApiError<PetErrorCode>(error.data)) {
        if ('retry_after' in error.data && error.data.retry_after) {
          toast.info(`${error.data.message}. Попробуйте через ${error.data.retry_after} секунд`)
          return
        }

        toast.error(error.data.message)
        return
      }

      toast.error('Не удалось отправить питомца в путешествие')
    }
  }

  return { handleJourney, isJourneyLoading: isLoading }
}
