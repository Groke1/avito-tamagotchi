import { useAppDispatch } from '@/app/store/hooks'
import { setPet, useCreatePetMutation } from '@/entities/pet'
import { isApiError, isFetchBaseQueryError } from '@/shared/lib/guards'
import type { AuthErrorCode } from '@/shared/model/types'
import { zodResolver } from '@hookform/resolvers/zod'
import { type SubmitHandler, useForm } from 'react-hook-form'
import { toast } from 'sonner'
import { type CreatePetFormData, createPetSchema } from './createPet.schema'

export const useCreatePetForm = () => {
  const [createPet, { isLoading }] = useCreatePetMutation()
  const dispatch = useAppDispatch()

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<CreatePetFormData>({
    resolver: zodResolver(createPetSchema),
    mode: 'onTouched',
    defaultValues: {
      name: '',
    },
  })

  const onSubmit: SubmitHandler<CreatePetFormData> = async (data) => {
    try {
      const pet = await createPet(data).unwrap()
      dispatch(setPet(pet))
      toast.success(`Питомец ${pet.name} успешно создан!`)
    } catch (error: unknown) {
      if (isFetchBaseQueryError(error) && isApiError<AuthErrorCode>(error.data)) {
        toast.error(error.data.message)
      } else {
        toast.error('Произошла ошибка при создании питомца')
      }
    }
  }

  return {
    register,
    errors,
    isLoading,
    handleSubmit: handleSubmit(onSubmit),
  }
}
