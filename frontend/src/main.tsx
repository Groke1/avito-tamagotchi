import { AuthInitializer, StoreProvider } from './app/providers'
import { router } from './app/router/config'
import { Toaster } from '@/shared/ui/sonner'
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { RouterProvider } from 'react-router-dom'
import './index.css'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <StoreProvider>
      <AuthInitializer>
        <Toaster />
        <RouterProvider router={router} />
      </AuthInitializer>
    </StoreProvider>
  </StrictMode>,
)

