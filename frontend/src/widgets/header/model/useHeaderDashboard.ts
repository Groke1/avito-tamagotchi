import { useGetProfileQuery } from '@/entities/user'
import { useAppSelector } from '@/shared/model'
import { useHandleLogout } from './useHandleLogout'

export const useHeaderDashboard = () => {
  useGetProfileQuery()
  const user = useAppSelector((state) => state.user.user)
  const { handleLogout, isLogoutLoading } = useHandleLogout()

  return {
    user,
    coins: user?.coins ?? 0,
    username: user?.username || 'Пользователь',
    handleLogout,
    isLogoutLoading,
  }
}
