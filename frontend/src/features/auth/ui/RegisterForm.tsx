import { ROUTES_PATHS } from '@/app/router/config'
import { Button, Checkbox, FormField } from '@/shared/ui'
import { Lock, Mail, PawPrint, User } from 'lucide-react'
import { Link } from 'react-router-dom'
import { useRegisterForm } from '../model/useRegisterForm'

export const RegisterForm = () => {
  const { termsChecked, errors, isLoading, register, setValue, handleSubmit } = useRegisterForm()

  return (
    <div className="bg-surface-lowest rounded-3xl shadow-level-1 p-6 md:p-8 w-full border-2 border-white relative overflow-hidden z-10">
      <div className="absolute top-0 left-0 w-full h-2 bg-linear-to-r from-avito-green to-avito-blue" />
      <form
        className="flex flex-col gap-5 mt-1"
        onSubmit={handleSubmit}
        noValidate
        autoComplete="off"
      >
        <FormField
          label="Имя пользователя"
          id="username"
          type="text"
          placeholder="Ваш никнейм"
          disabled={isLoading}
          leftIcon={<User size="20" />}
          error={errors.username?.message}
          {...register('username')}
        />
        <FormField
          label="Email"
          id="email"
          type="email"
          placeholder="mail@example.com"
          disabled={isLoading}
          leftIcon={<Mail size="20" />}
          error={errors.email?.message}
          autoComplete="email"
          {...register('email')}
        />
        <FormField
          label="Пароль"
          id="password"
          type="password"
          placeholder="••••••••"
          disabled={isLoading}
          leftIcon={<Lock size="20" />}
          error={errors.password?.message}
          autoComplete="new-password"
          {...register('password')}
        />
        <div className="flex flex-col gap-1">
          <div className="flex items-start bg-surface-bright p-3 rounded-xl border border-surface-container">
            <div className="flex items-center h-5 mt-0.5">
              <Checkbox
                id="terms"
                disabled={isLoading}
                checked={termsChecked}
                onCheckedChange={(checked) =>
                  setValue('terms', Boolean(checked), { shouldValidate: true })
                }
              />
            </div>
            <div className="ml-2.5 text-xs text-on-surface-variant leading-tight">
              <label htmlFor="terms" className="cursor-pointer select-none">
                Я согласен с{' '}
                <a href="#" className="text-avito-blue hover:underline font-semibold">
                  условиями использования
                </a>{' '}
                и политикой конфиденциальности.
              </label>
            </div>
          </div>
          {errors.terms?.message && (
            <span className="text-xs text-avito-red ml-1 font-medium">{errors.terms.message}</span>
          )}
        </div>
        <Button type="submit" isLoading={isLoading}>
          <span>Создать аккаунт</span>
          <PawPrint size="20" />
        </Button>
      </form>
      <div className="mt-6 text-center text-sm text-on-surface-variant">
        Уже есть аккаунт?{' '}
        <Link to={ROUTES_PATHS.LOGIN} className="text-avito-blue font-semibold hover:underline">
          Войти
        </Link>
      </div>
    </div>
  )
}
