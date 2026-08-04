import { useLoginMutation } from '@/entities/user'
import { zodResolver } from '@hookform/resolvers/zod'
import { type SubmitHandler, useForm } from 'react-hook-form'
import { type LoginFormData, loginSchema } from './login.schema'

export const useLoginForm = () => {
  const [loginApi, { isLoading }] = useLoginMutation()

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
    mode: 'onTouched',
  })

  const onSubmit: SubmitHandler<LoginFormData> = async (data) => {
    const result = await loginApi(data)

    if (result.data) {
      // TODO: user login slice
      // dispatch(login(...))
      reset()
    }
  }

  return {
    register,
    errors,
    isLoading,
    handleSubmit: handleSubmit(onSubmit),
  }
}
