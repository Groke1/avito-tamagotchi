import { lazy } from 'react'
import { createBrowserRouter } from 'react-router-dom'

const Register = lazy(() => import('@/pages/Register'))
const Login = lazy(() => import('@/pages/Login'))
const CreatePet = lazy(() => import('@/pages/CreatePet'))

export const ROUTES_PATHS = {
  LOGIN: '/login',
  REGISTER: '/register',
  CREATE_PET: '/create-pet',
} as const

export const router = createBrowserRouter([
  { path: ROUTES_PATHS.REGISTER, element: <Register /> },
  { path: ROUTES_PATHS.LOGIN, element: <Login /> },
  { path: ROUTES_PATHS.CREATE_PET, element: <CreatePet /> },
])
