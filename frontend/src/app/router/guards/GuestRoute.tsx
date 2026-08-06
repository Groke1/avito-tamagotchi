import { useAppSelector } from '@/app/store/hooks'
import type { FC } from 'react'
import { Navigate, Outlet } from 'react-router-dom'
import { ROUTES_PATHS } from '../paths'

export const GuestRoute: FC = () => {
  const isAuthenticated = useAppSelector((state) => state.user.isAuthenticated)

  if (isAuthenticated) {
    return <Navigate to={ROUTES_PATHS.DASHBOARD} replace />
  }

  return <Outlet />
}
