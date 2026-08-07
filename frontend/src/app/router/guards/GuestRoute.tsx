import { ROUTES_PATHS } from '@/shared/config'
import { useAppSelector } from '@/shared/model'
import type { FC } from 'react'
import { Navigate, Outlet } from 'react-router-dom'

export const GuestRoute: FC = () => {
  const isAuthenticated = useAppSelector((state) => state.user.isAuthenticated)

  if (isAuthenticated) {
    return <Navigate to={ROUTES_PATHS.DASHBOARD} replace />
  }

  return <Outlet />
}
