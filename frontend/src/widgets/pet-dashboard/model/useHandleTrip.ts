import { useTripPetMutation } from '@/entities/pet'
import { isApiError, isFetchBaseQueryError } from '@/shared/lib/guards'
import { useAppSelector } from '@/shared/model'
import { toast } from 'sonner'

export const useHandleTrip = () => {
  const [tripPet, { isLoading }] = useTripPetMutation()
  const pet = useAppSelector((state) => state.pet.pet)

  const handleTrip = async () => {
    if (!pet) return

    try {
      await tripPet(pet.id).unwrap()
      toast.success('Путешествие создано! ')
    } catch (error: unknown) {
      if (isFetchBaseQueryError(error) && isApiError(error.data)) {
        toast.error(error.data.message)
      } else {
        toast.error('Это действие пока недоступно')
      }
    }
  }

  return { isTripLoading: isLoading, handleTrip }
}
