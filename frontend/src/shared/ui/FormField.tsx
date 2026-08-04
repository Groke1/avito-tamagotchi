import type { FC, ReactNode } from 'react'
import { Input } from './input'

interface FormFieldProps {
  label: string
  id: string
  type: string
  placeholder: string
  leftIcon: ReactNode
  error?: string
}

export const FormField: FC<FormFieldProps> = ({
  label,
  id,
  type,
  placeholder,
  leftIcon,
  error,
}) => {
  return (
    <div className="flex flex-col gap-2">
      <label htmlFor={id} className="block text-sm font-semibold text-on-surface ml-1">
        {label}
      </label>
      <Input id={id} type={type} placeholder={placeholder} leftIcon={leftIcon} />
      {error && <p className="text-xs text-avito-red ml-1">{error}</p>}
    </div>
  )
}
