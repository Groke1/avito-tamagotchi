import { type ReactNode, forwardRef } from 'react'
import { Input, type InputProps } from './input'

interface FormFieldProps extends InputProps {
  label: string
  id: string
  leftIcon?: ReactNode
  error?: string
}

export const FormField = forwardRef<HTMLInputElement, FormFieldProps>(
  ({ label, id, error, ...props }, ref) => {
    return (
      <div className="flex flex-col gap-2">
        <label
          htmlFor={id}
          className="block text-sm font-semibold text-on-surface ml-1 select-none cursor-pointer"
        >
          {label}
        </label>
        <Input id={id} ref={ref} error={error} {...props} />
        {error && <p className="text-xs text-avito-red ml-1 font-medium">{error}</p>}
      </div>
    )
  },
)
