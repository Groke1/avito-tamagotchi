import { ROUTES_PATHS } from '@/app/router/paths'
import { useAppDispatch } from '@/app/store/hooks'
import { setHasPet, setPet, useLazyGetPetQuery } from '@/entities/pet'
import { login, setAccessToken, useLazyGetProfileQuery, useLoginMutation } from '@/entities/user'
import { isApiError, isFetchBaseQueryError } from '@/shared/lib/guards'
import type { AuthErrorCode } from '@/shared/model/types'
import { zodResolver } from '@hookform/resolvers/zod'
import { type SubmitHandler, useForm } from 'react-hook-form'
import { useNavigate } from 'react-router-dom'
import { toast } from 'sonner'
import { type LoginFormData, loginSchema } from './login.schema'

export const useLoginForm = () => {
  const [loginApi, { isLoading }] = useLoginMutation()
  const [fetchProfile] = useLazyGetProfileQuery()
  const [fetchPet] = useLazyGetPetQuery()
  const dispatch = useAppDispatch()
  const navigate = useNavigate()

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
    try {
      const { access_token, refresh_token } = await loginApi(data).unwrap()
      dispatch(setAccessToken(access_token))

      const user = await fetchProfile().unwrap()

      dispatch(
        login({
          accessToken: access_token,
          refreshToken: refresh_token,
          user,
        }),
      )

      let userHasPet = false
      try {
        const pet = await fetchPet().unwrap()
        dispatch(setPet(pet))
        userHasPet = true
      } catch {
        dispatch(setHasPet(false))
      }

      toast.success('Успешный вход!')
      reset()

      if (userHasPet) {
        navigate(ROUTES_PATHS.MAIN)
      } else {
        navigate(ROUTES_PATHS.CREATE_PET)
      }
    } catch (error: unknown) {
      if (isFetchBaseQueryError(error) && isApiError<AuthErrorCode>(error.data)) {
        toast.error(error.data.message)
      } else {
        toast.error('Произошла ошибка при входе')
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
