import { Outlet } from 'react-router-dom'

export const AuthLayout = () => {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center p-4 sm:p-6 md:p-8 bg-surface-bg relative overflow-hidden">
      <div className="hidden sm:block absolute -top-40 -left-20 w-80 h-80 sm:-top-60 sm:-left-32 sm:w-120 sm:h-120 md:-top-80 md:-left-40 md:w-150 md:h-150 bg-avito-green/20 rounded-full blur-2xl pointer-events-none" />
      <div className="hidden md:block absolute left-[10%] -bottom-60 w-100 h-100 lg:left-[20%] lg:-bottom-100 lg:w-200 lg:h-200 bg-avito-red/10 rounded-full blur-2xl pointer-events-none" />
      <div className="hidden sm:block absolute -right-24 bottom-10 w-72 h-72 sm:-right-40 sm:bottom-20 sm:w-100 sm:h-100 md:-right-60 md:bottom-50 md:w-150 md:h-150 bg-avito-blue/20 rounded-full blur-2xl pointer-events-none" />
      <main className="w-full max-w-md sm:max-w-lg z-10 flex flex-col items-center">
        <Outlet />
      </main>
    </div>
  )
}
