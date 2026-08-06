import { ROUTES_PATHS } from '@/app/router/paths'
import { useAppDispatch } from '@/app/store/hooks'
import { setHasPet } from '@/entities/pet'
import { login, setAccessToken, useLazyGetProfileQuery, useRegisterMutation } from '@/entities/user'
import { isApiError, isFetchBaseQueryError } from '@/shared/lib/guards'
import type { AuthErrorCode } from '@/shared/model/types'
import { zodResolver } from '@hookform/resolvers/zod'
import { type SubmitHandler, useForm, useWatch } from 'react-hook-form'
import { useNavigate } from 'react-router-dom'
import { toast } from 'sonner'
import { type RegisterFormData, registerSchema } from './register.schema'

export const useRegisterForm = () => {
  const [registerApi, { isLoading }] = useRegisterMutation()
  const [fetchProfile] = useLazyGetProfileQuery()
  const dispatch = useAppDispatch()
  const navigate = useNavigate()

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

  const onSubmit: SubmitHandler<RegisterFormData> = async ({ email, username, password }) => {
    try {
      const { access_token, refresh_token } = await registerApi({
        email,
        username,
        password,
      }).unwrap()

      dispatch(setAccessToken(access_token))

      const user = await fetchProfile().unwrap()

      dispatch(
        login({
          accessToken: access_token,
          refreshToken: refresh_token,
          user,
        }),
      )
      dispatch(setHasPet(false))

      toast.success('Регистрация прошла успешно!')
      reset(
        { username: '', email: '', password: '', terms: false },
        { keepErrors: false, keepTouched: false },
      )

      navigate(ROUTES_PATHS.CREATE_PET)
    } catch (error: unknown) {
      if (isFetchBaseQueryError(error) && isApiError<AuthErrorCode>(error.data)) {
        toast.error(error.data.message)
      } else {
        toast.error('Произошла ошибка при регистрации')
      }
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
