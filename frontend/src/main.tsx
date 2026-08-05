import { Toaster } from '@/shared/ui/sonner'
import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { RouterProvider } from 'react-router-dom'
import { StoreProvider } from './app/providers'
import { router } from './app/router/config'
import './index.css'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <StoreProvider>
      <Toaster />
      <RouterProvider router={router} />
    </StoreProvider>
  </StrictMode>,
)
