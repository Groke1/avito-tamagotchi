import { CreatePet, Dashboard, Login, Register, Rewards, Tasks } from '@/pages'
import { ROUTES_PATHS } from '@/shared/config'
import { PageLoader } from '@/shared/ui'
import { createBrowserRouter } from 'react-router-dom'
import { AuthLayout, RootLayout } from '../layouts'
import { GuestRoute, ProtectedRoute } from './guards'
import { rootLoader } from './loaders'

export const router = createBrowserRouter([
  {
    loader: rootLoader,
    hydrateFallbackElement: <PageLoader />,
    children: [
      {
        element: <GuestRoute />,
        children: [
          {
            element: <AuthLayout />,
            children: [
              { path: ROUTES_PATHS.REGISTER, element: <Register /> },
              { path: ROUTES_PATHS.LOGIN, element: <Login /> },
            ],
          },
        ],
      },
      {
        element: <ProtectedRoute />,
        children: [
          {
            element: <AuthLayout />,
            children: [{ path: ROUTES_PATHS.CREATE_PET, element: <CreatePet /> }],
          },
          {
            element: <RootLayout />,
            children: [
              { path: ROUTES_PATHS.DASHBOARD, element: <Dashboard /> },
              { path: ROUTES_PATHS.TASKS, element: <Tasks /> },
              { path: ROUTES_PATHS.REWARDS, element: <Rewards /> },
              { path: ROUTES_PATHS.SHOP, element: <div>Shop</div> },
              { path: ROUTES_PATHS.LEADERBOARD, element: <div>Leaderboard</div> },
              { path: ROUTES_PATHS.DAILY_REPORT, element: <div>Daily report</div> },
            ],
          },
        ],
      },
    ],
  },
])
