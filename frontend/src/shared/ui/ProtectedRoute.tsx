import { ROUTES_PATHS } from '@/app/router/paths'
import { useAppSelector } from '@/app/store/hooks'
import type { FC } from 'react'
import { Navigate, Outlet, useLocation } from 'react-router-dom'

export const ProtectedRoute: FC = () => {
  const isAuthenticated = useAppSelector((state) => state.user.isAuthenticated)
  const hasPet = useAppSelector((state) => state.pet.hasPet)
  const location = useLocation()

  if (!isAuthenticated) {
    return <Navigate to={ROUTES_PATHS.LOGIN} state={{ from: location }} replace />
  }

  if (location.pathname === ROUTES_PATHS.CREATE_PET && hasPet) {
    return <Navigate to={ROUTES_PATHS.MAIN} replace />
  }

  if (location.pathname !== ROUTES_PATHS.CREATE_PET && !hasPet) {
    return <Navigate to={ROUTES_PATHS.CREATE_PET} replace />
  }

  return <Outlet />
}
