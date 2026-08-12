import { useGetPetQuery } from '@/entities/pet'
import { useWebSocket } from '@/shared/lib/websocket/useWebSocket'
import { useAppSelector } from '@/shared/model'
import { MobileNav, Sidebar } from '@/widgets/sidebar'
import { Outlet } from 'react-router-dom'

export const RootLayout = () => {
  const isAuthenticated = useAppSelector((state) => state.user.isAuthenticated)
  useGetPetQuery(undefined, { skip: !isAuthenticated })
  useWebSocket()

  return (
    <div className="flex flex-col lg:flex-row min-h-screen bg-surface-bg">
      <MobileNav />
      <div className="hidden lg:block shrink-0 basis-64 h-screen sticky top-0">
        <Sidebar />
      </div>
      <main className="flex-1 min-w-0 px-4 py-5 pb-24 lg:px-8 lg:py-6 lg:pb-8">
        <Outlet />
      </main>
    </div>
  )
}
