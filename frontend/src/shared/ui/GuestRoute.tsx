import { ROUTES_PATHS } from '@/app/router/config'
import { useAppSelector } from '@/app/store/hooks'
import type { FC } from 'react'
import { Navigate, Outlet } from 'react-router-dom'

export const GuestRoute: FC = () => {
  const isAuthenticated = useAppSelector((state) => state.user.isAuthenticated)

  if (isAuthenticated) {
    return <Navigate to={ROUTES_PATHS.MAIN} replace />
  }

  return <Outlet />
}
