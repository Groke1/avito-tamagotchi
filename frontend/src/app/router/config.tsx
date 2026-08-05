import { CreatePet, Login, Register } from '@/pages'
import { PageLoader } from '@/shared/ui'
import { GuestRoute } from '@/shared/ui/GuestRoute'
import { ProtectedRoute } from '@/shared/ui/ProtectedRoute'
import { createBrowserRouter } from 'react-router-dom'
import { rootLoader } from '../loaders/rootLoader'
import { ROUTES_PATHS } from './paths'

export const router = createBrowserRouter([
  {
    loader: rootLoader,
    hydrateFallbackElement: <PageLoader />,
    children: [
      {
        element: <GuestRoute />,
        children: [
          { path: ROUTES_PATHS.REGISTER, element: <Register /> },
          { path: ROUTES_PATHS.LOGIN, element: <Login /> },
        ],
      },
      {
        element: <ProtectedRoute />,
        children: [
          { path: ROUTES_PATHS.CREATE_PET, element: <CreatePet /> },
          { path: ROUTES_PATHS.MAIN, element: <div>Main</div> },
        ],
      },
    ],
  },
])
