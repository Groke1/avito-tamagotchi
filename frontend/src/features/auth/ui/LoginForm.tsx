import { ROUTES_PATHS } from '@/app/router/config'
import { type LoginFormData, loginSchema } from '@/features/auth/model/login.schema'
import { Button, FormField } from '@/shared/ui'
import { zodResolver } from '@hookform/resolvers/zod'
import { ArrowRight, Lock, Mail } from 'lucide-react'
import { type SubmitHandler, useForm } from 'react-hook-form'
import { Link } from 'react-router-dom'

export const LoginForm = () => {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>({ resolver: zodResolver(loginSchema), mode: 'onTouched' })

  const onSubmit: SubmitHandler<LoginFormData> = (data) => console.log(data)

  return (
    <div className="bg-surface-lowest rounded-3xl shadow-level-1 p-6 md:p-8 w-full border-2 border-white relative overflow-hidden z-10">
      <div className="absolute top-0 left-0 w-full h-2 bg-linear-to-r from-avito-green to-avito-blue" />
      <form className="flex flex-col gap-5 mt-1" onSubmit={handleSubmit(onSubmit)} noValidate>
        <FormField
          label="Email"
          id="email"
          type="email"
          placeholder="mail@example.com"
          error={errors.email?.message}
          leftIcon={<Mail size="20" />}
          {...register('email')}
        />
        <FormField
          label="Пароль"
          id="password"
          type="password"
          placeholder="••••••••"
          leftIcon={<Lock size="20" />}
          {...register('password')}
          error={errors.password?.message}
        />
        <Button type="submit">
          <span>Войти</span>
          <ArrowRight size="20" />
        </Button>
      </form>
      <div className="mt-6 text-center text-sm text-on-surface-variant">
        Ещё нет аккаунта?{' '}
        <Link to={ROUTES_PATHS.REGISTER} className="text-avito-blue font-semibold hover:underline">
          Зарегистрироваться
        </Link>
      </div>
    </div>
  )
}
