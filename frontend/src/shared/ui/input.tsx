import { cn } from '@/shared/lib/utils'
import { Eye, EyeOff } from 'lucide-react'
import { type ReactNode, useState } from 'react'

export interface InputProps extends React.ComponentProps<'input'> {
  leftIcon?: ReactNode
  error?: string
}

function Input({ className, type = 'text', leftIcon, error, ...props }: InputProps) {
  const [showPassword, setShowPassword] = useState(false)
  const isPassword = type === 'password'

  const currentType = isPassword ? (showPassword ? 'text' : 'password') : type

  return (
    <div className="relative w-full flex items-center">
      <input
        type={currentType}
        data-slot="input"
        className={cn(
          'peer w-full bg-surface-container-low rounded-xl px-4 py-3 text-sm text-on-surface border border-surface-container placeholder:text-on-surface-variant/50 focus:outline-none focus:ring-1 focus:ring-avito-green transition-all disabled:opacity-50 disabled:cursor-not-allowed',
          leftIcon && 'pl-11',
          isPassword && 'pr-11',
          error && 'ring-1 ring-avito-red focus:ring-avito-red',
          className,
        )}
        {...props}
      />
      {leftIcon && (
        <div
          className={cn(
            'absolute left-3.5 inset-y-0 flex items-center justify-center text-on-surface-variant/60 pointer-events-none transition-colors peer-focus:text-avito-green',
            error && 'peer-focus:text-avito-red text-avito-red',
          )}
        >
          {leftIcon}
        </div>
      )}
      {isPassword && (
        <button
          type="button"
          onClick={() => setShowPassword((prev) => !prev)}
          tabIndex={-1}
          aria-label={showPassword ? 'Скрыть пароль' : 'Показать пароль'}
          className="absolute right-3.5 inset-y-0 flex items-center justify-center text-on-surface-variant/60 hover:text-on-surface focus:outline-none cursor-pointer transition-colors"
        >
          {showPassword ? <EyeOff size="20" /> : <Eye size="20" />}
        </button>
      )}
    </div>
  )
}

export { Input }
