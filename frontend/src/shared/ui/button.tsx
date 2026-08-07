import { cn } from '@/shared/lib/utils'
import { type VariantProps, cva } from 'class-variance-authority'
import { Loader2 } from 'lucide-react'
import * as React from 'react'

const buttonVariants = cva(
  'inline-flex shrink-0 items-center justify-center gap-2 rounded-2xl text-base font-semibold whitespace-nowrap transition-all outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 cursor-pointer active:scale-[0.98]',
  {
    variants: {
      variant: {
        default: 'bg-avito-blue text-white hover:bg-avito-blue/90 shadow-md',
        primary: 'bg-avito-blue text-white hover:bg-avito-blue/90 shadow-md',
        avitoGreen: 'bg-avito-green text-white hover:bg-avito-green/90 shadow-md',
        outline:
          'border-2 border-surface-container bg-surface-lowest text-on-surface hover:bg-surface-high hover:border-outline-variant transition-all duration-200',
        secondary: 'bg-surface-container text-on-surface hover:bg-surface-container-high',
        ghost: 'hover:bg-surface-container-low text-on-surface',
        destructive: 'bg-red-500 text-white hover:bg-red-600',
      },
      size: {
        default: 'py-3.5 px-6 text-base min-h-13',
        sm: 'py-2 px-4 text-sm rounded-xl',
        lg: 'py-4 px-8 text-lg rounded-2xl',
        icon: 'size-11 rounded-xl',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  },
)

export interface ButtonProps
  extends React.ComponentProps<'button'>, VariantProps<typeof buttonVariants> {
  isLoading?: boolean
}

function Button({
  className,
  variant = 'default',
  size = 'default',
  isLoading = false,
  children,
  disabled,
  ...props
}: ButtonProps) {
  return (
    <button
      data-slot="button"
      data-variant={variant}
      data-size={size}
      disabled={disabled || isLoading}
      className={cn(buttonVariants({ variant, size, className }))}
      {...props}
    >
      {isLoading ? <Loader2 className="animate-spin size-5" /> : children}
    </button>
  )
}

export { Button, buttonVariants }
