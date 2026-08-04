import { Button, Checkbox, FormField } from '@/shared/ui'
import { Lock, Mail, PawPrint, User } from 'lucide-react'
import { Link } from 'react-router-dom'

export const RegisterForm = () => {
  return (
    <div className="bg-surface-lowest rounded-3xl shadow-level-1 p-6 md:p-8 w-full max-w-md border-2 border-white relative overflow-hidden">
      <div className="absolute top-0 left-0 w-full h-2 bg-linear-to-r from-avito-green to-avito-blue" />
      <form className="flex flex-col gap-5 mt-1" onSubmit={(e) => e.preventDefault()}>
        <FormField
          label="Имя пользователя"
          id="username"
          type="text"
          placeholder="Ваш никнейм"
          leftIcon={<User size="20" />}
        />
        <FormField
          label="Email"
          id="email"
          type="email"
          placeholder="mail@example.com"
          leftIcon={<Mail size="20" />}
        />
        <FormField
          label="Пароль"
          id="password"
          type="password"
          placeholder="••••••••"
          leftIcon={<Lock size="20" />}
        />
        <div className="flex items-start bg-surface-bright p-3 rounded-xl border border-surface-container">
          <div className="flex items-center h-5 mt-0.5">
            <Checkbox id="terms" />
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
        <Button type="submit">
          <span>Создать аккаунт</span>
          <PawPrint size="20" />
        </Button>
      </form>
      <div className="mt-6 text-center text-sm text-on-surface-variant">
        Уже есть аккаунт?{' '}
        <Link to="/login" className="text-avito-blue font-semibold hover:underline">
          Войти
        </Link>
      </div>
    </div>
  )
}
