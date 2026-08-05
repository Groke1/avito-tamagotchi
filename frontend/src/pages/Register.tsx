import { AuthLayout } from '@/app/layouts/AuthLayout'
import { RegisterForm } from '@/features/auth'

export const Register = () => {
  return (
    <AuthLayout>
      <div className="flex flex-col items-center mb-6 text-center max-w-xs">
        <div className="flex items-center gap-1.5 mb-2">
          <span className="text-2xl font-black tracking-tight text-avito-green">Авито</span>
          <span className="bg-avito-blue/15 text-avito-blue px-2.5 py-0.5 rounded-lg text-sm font-semibold">
            Тамагочи
          </span>
        </div>
        <h1 className="text-2xl md:text-3xl font-bold text-on-surface">Создайте питомца</h1>
        <p className="text-sm text-on-surface-variant mt-1">Начните свое приключение в Авито</p>
        <div className="flex items-center gap-2.5 mt-4 w-44">
          <div className="flex-1 h-1.5 bg-surface-high rounded-full overflow-hidden">
            <div className="h-full bg-avito-green w-1/2 rounded-full transition-all duration-300" />
          </div>
          <span className="text-xs font-semibold text-on-surface">Шаг 1 из 2</span>
        </div>
      </div>
      <RegisterForm />
    </AuthLayout>
  )
}
