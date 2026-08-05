import { ROUTES_PATHS } from '@/app/router/config'
import { useAppSelector } from '@/app/store/hooks'
import type { FC } from 'react'
import { Navigate, Outlet, useLocation } from 'react-router-dom'

export const ProtectedRoute: FC = () => {
  const isAuthenticated = useAppSelector((state) => state.user.isAuthenticated)
  const location = useLocation()

  if (!isAuthenticated) {
    return <Navigate to={ROUTES_PATHS.LOGIN} state={{ from: location }} replace />
  }

  return <Outlet />
}
