import { useAppDispatch } from '@/app/store/hooks'
import { getStoredRefreshToken, logout, useLogoutMutation } from '@/entities/user'
import { baseApi } from '@/shared/api/baseApi'
import { toast } from 'sonner'

export const useHandleLogout = () => {
  const [logoutApi, { isLoading: isLogoutLoading }] = useLogoutMutation()
  const dispatch = useAppDispatch()

  const handleLogout = async () => {
    try {
      const token = getStoredRefreshToken()

      if (token) {
        await logoutApi(token).unwrap()
      }
    } finally {
      dispatch(logout())
      dispatch(baseApi.util.resetApiState())
      toast.success('Вы успешно вышли из системы')
    }
  }

  return { handleLogout, isLogoutLoading }
}
