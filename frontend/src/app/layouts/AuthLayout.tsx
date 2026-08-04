import type { FC, PropsWithChildren } from 'react'

export const AuthLayout: FC<PropsWithChildren> = ({ children }) => {
  return (
    <div className="min-h-screen flex flex-col">
      <header className="flex items-center justify-between px-6 py-4">
        <h2>
          Авито <span>Тамагочи</span>
        </h2>
        <div>
          <div>Шаг 1 из 2</div>
          <div className="relative p-1 rounded-full bg-surface-high">
            <span className="absolute inset-0 w-1/2 bg-avito-green rounded-full"></span>
          </div>
        </div>
      </header>
      <div className="flex justify-center items-center ">{children}</div>
    </div>
  )
}
