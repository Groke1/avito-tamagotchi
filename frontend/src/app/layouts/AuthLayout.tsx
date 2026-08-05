import type { FC, PropsWithChildren } from 'react'

export const AuthLayout: FC<PropsWithChildren> = ({ children }) => {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center p-4 md:p-8 bg-surface-bg relative overflow-hidden">
      <div className="hidden xl:block absolute -top-80 -left-40 w-150 h-150 bg-avito-green/20 rounded-full blur-2xl pointer-events-none" />
      <div className="hidden xl:block absolute left-[20%] -bottom-100 w-200 h-200 bg-avito-red/10 rounded-full blur-2xl pointer-events-none" />
      <div className="hidden xl:block absolute -right-60 bottom-50 w-150 h-150 bg-avito-blue/20 rounded-full blur-2xl pointer-events-none" />
      <main className="w-full max-w-lg z-10 flex flex-col items-center">{children}</main>
    </div>
  )
}
