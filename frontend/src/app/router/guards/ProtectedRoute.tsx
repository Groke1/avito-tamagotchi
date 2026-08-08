import { ROUTES_PATHS } from '@/shared/config'
import { useAppSelector } from '@/shared/model'
import { PageLoader } from '@/shared/ui'
import type { FC } from 'react'
import { Navigate, Outlet, useLocation } from 'react-router-dom'

export const ProtectedRoute: FC = () => {
  const isAuthenticated = useAppSelector((state) => state.user.isAuthenticated)
  const { pet, isInitialized } = useAppSelector((state) => state.pet)
  const location = useLocation()

  if (!isAuthenticated) {
    return <Navigate to={ROUTES_PATHS.LOGIN} state={{ from: location }} replace />
  }

  if (!isInitialized) {
    return <PageLoader />
  }

  const hasPet = Boolean(pet)

  if (location.pathname === ROUTES_PATHS.CREATE_PET && hasPet) {
    return <Navigate to={ROUTES_PATHS.DASHBOARD} replace />
  }

  if (location.pathname !== ROUTES_PATHS.CREATE_PET && !hasPet) {
    return <Navigate to={ROUTES_PATHS.CREATE_PET} replace />
  }

  return <Outlet />
}
