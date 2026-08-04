import { cn } from '@/shared/lib/utils'
import { CheckIcon } from 'lucide-react'
import { Checkbox as CheckboxPrimitive } from 'radix-ui'
import * as React from 'react'

function Checkbox({ className, ...props }: React.ComponentProps<typeof CheckboxPrimitive.Root>) {
  return (
    <CheckboxPrimitive.Root
      data-slot="checkbox"
      className={cn(
        'peer size-5 shrink-0 rounded-md border border-outline bg-surface-lowest shadow-xs transition-all outline-none focus-visible:ring-2 focus-visible:ring-avito-green/20 disabled:cursor-not-allowed disabled:opacity-50 data-[state=checked]:bg-avito-green data-[state=checked]:border-avito-green data-[state=checked]:text-white cursor-pointer',
        className,
      )}
      {...props}
    >
      <CheckboxPrimitive.Indicator
        data-slot="checkbox-indicator"
        className="grid place-content-center text-current transition-none"
      >
        <CheckIcon className="size-3.5 stroke-3" />
      </CheckboxPrimitive.Indicator>
    </CheckboxPrimitive.Root>
  )
}

export { Checkbox }
