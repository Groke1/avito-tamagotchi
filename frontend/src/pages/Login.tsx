import { AuthLayout } from '@/app/layouts/AuthLayout'
import { LoginForm } from '@/features/auth'

const Login = () => {
  return (
    <AuthLayout>
      <div className="flex flex-col items-center mb-6 text-center max-w-sm">
        <div className="flex items-center gap-1.5 mb-2">
          <span className="text-2xl font-black tracking-tight text-avito-green">Авито</span>
          <span className="bg-avito-blue/15 text-avito-blue px-2.5 py-0.5 rounded-lg text-sm font-semibold">
            Тамагочи
          </span>
        </div>
        <h1 className="text-2xl md:text-3xl font-bold text-on-surface">С возвращением</h1>
        <p className="text-sm text-on-surface-variant mt-1">
          Войдите, чтобы продолжить приключение со своим питомцем
        </p>
      </div>
      <LoginForm />
    </AuthLayout>
  )
}

export default Login
