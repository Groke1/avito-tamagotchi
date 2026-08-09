import { useGetPetQuery } from '@/entities/pet'
import { useWebSocket } from '@/shared/lib/websocket/useWebSocket'
import { useAppSelector } from '@/shared/model'
import { Sidebar } from '@/widgets/sidebar'
import { Outlet } from 'react-router-dom'

export const RootLayout = () => {
  const isAuthenticated = useAppSelector((state) => state.user.isAuthenticated)
  useGetPetQuery(undefined, { skip: !isAuthenticated })
  useWebSocket()

  return (
    <div className="flex min-h-screen">
      <div className="shrink-0 basis-64">
        <Sidebar />
      </div>
      <main className="flex-1 overflow-hidden min-w-0 px-8 py-6">
        <Outlet />
      </main>
    </div>
  )
}
