import { Sidebar } from '@/widgets/sidebar'
import { Outlet } from 'react-router-dom'

export const RootLayout = () => {
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
