import { cn } from '@/shared/lib/utils'
import * as React from 'react'

export interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  leftIcon?: React.ReactNode
}

const Input = React.forwardRef<HTMLInputElement, InputProps>(
  ({ className, type, leftIcon, ...props }, ref) => {
    return (
      <div className="group relative w-full flex items-center">
        {leftIcon && (
          <div className="absolute left-3.5 inset-y-0 flex items-center justify-center text-on-surface-variant/60 pointer-events-none transition-colors group-focus-within:text-avito-green">
            {leftIcon}
          </div>
        )}
        <input
          type={type}
          className={cn(
            'w-full bg-surface-container-low rounded-xl px-4 py-3 text-sm text-on-surface border border-surface-container placeholder:text-on-surface-variant/50 focus:outline-none focus:ring-1 focus:ring-avito-green transition-all disabled:opacity-50 disabled:cursor-not-allowed',
            leftIcon && 'pl-11',
            className,
          )}
          autoComplete="on"
          ref={ref}
          {...props}
        />
      </div>
    )
  },
)

Input.displayName = 'Input'

export { Input }
