import { useTripPetMutation } from '@/entities/pet'
import { isApiError, isFetchBaseQueryError } from '@/shared/lib/guards'
import { useAppSelector } from '@/shared/model'
import { toast } from 'sonner'

export const useHandleTrip = () => {
  const [tripPet, { isLoading }] = useTripPetMutation()
  const pet = useAppSelector((state) => state.pet.pet)
  const user = useAppSelector((state) => state.user.user)

  const handleTrip = async () => {
    if (!pet) {
      toast.error('Сначала создайте питомца')
      return
    }

    if (user && user.coins < 100) {
      toast.error('Недостаточно монет для путешествия (нужно 100 монет)')
      return
    }

    try {
      await tripPet(pet.id).unwrap()

      toast.success('Питомец отправился в путешествие! 🧭', {
        description: 'Когда путешествие завершится, мы покажем историю и награду.',
      })
    } catch (error: unknown) {
      if (isFetchBaseQueryError(error) && isApiError(error.data)) {
        if ('retry_after' in error.data && error.data.retry_after) {
          toast.info(`${error.data.message}. Попробуйте через ${error.data.retry_after} секунд`)
          return
        }

        if (error.data.code === 'INSUFFICIENT_COINS') {
          toast.error(error.data.message)
          return
        }

        toast.info('Питомец уже в путешествии, попробуйте позже.')
      } else {
        toast.error('Это действие пока недоступно')
      }
    }
  }

  return { isTripLoading: isLoading, handleTrip }
}
