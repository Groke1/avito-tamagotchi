import { useRegisterMutation } from '@/entities/user'
import { zodResolver } from '@hookform/resolvers/zod'
import { type SubmitHandler, useForm, useWatch } from 'react-hook-form'
import { type RegisterFormData, registerSchema } from './register.schema'

export const useRegisterForm = () => {
  const [registerApi, { isLoading }] = useRegisterMutation()

  const {
    register,
    handleSubmit,
    reset,
    control,
    setValue,
    formState: { errors },
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
    mode: 'onTouched',
    defaultValues: {
      username: '',
      email: '',
      password: '',
      terms: false,
    },
  })

  const termsChecked = useWatch({ control, name: 'terms' })

  const onSubmit: SubmitHandler<RegisterFormData> = async (data) => {
    const result = await registerApi(data)

    if (result.data) {
      // TODO: user register slice
      reset()
    }
  }

  return {
    control,
    errors,
    isLoading,
    termsChecked,
    register,
    setValue,
    handleSubmit: handleSubmit(onSubmit),
  }
}
