import {
  CircleCheckIcon,
  InfoIcon,
  Loader2Icon,
  OctagonXIcon,
  TriangleAlertIcon,
} from 'lucide-react'
import { useTheme } from 'next-themes'
import { Toaster as Sonner, type ToasterProps } from 'sonner'

const Toaster = ({ ...props }: ToasterProps) => {
  const { theme = 'system' } = useTheme()

  return (
    <Sonner
      theme={theme as ToasterProps['theme']}
      position="bottom-right"
      className="toaster group"
      icons={{
        success: <CircleCheckIcon className="size-4 text-avito-green" />,
        info: <InfoIcon className="size-4 text-avito-blue" />,
        warning: <TriangleAlertIcon className="size-4 text-avito-yellow" />,
        error: <OctagonXIcon className="size-4 text-avito-red" />,
        loading: <Loader2Icon className="size-4 animate-spin text-avito-blue" />,
      }}
      toastOptions={{
        classNames: {
          toast:
            'group toast group-[.toaster]:bg-surface-lowest group-[.toaster]:text-on-surface group-[.toaster]:border-surface-container group-[.toaster]:shadow-level-2 group-[.toaster]:rounded-card p-4',
          title: 'group-[.toast]:text-on-surface group-[.toast]:font-bold text-sm',
          description:
            'group-[.toast]:text-on-surface-variant group-[.toast]:opacity-100 text-xs mt-1 font-medium',
          actionButton:
            'group-[.toast]:bg-avito-green group-[.toast]:text-white font-semibold rounded-lg px-3 py-1.5 text-xs',
          cancelButton: 'group-[.toast]:bg-surface-high group-[.toast]:text-on-surface',
          closeButton:
            'group-[.toast]:bg-surface-low group-[.toast]:text-on-surface group-[.toast]:border-surface-container group-[.toast]:hover:bg-surface-high group-[.toast]:top-3 group-[.toast]:right-3 group-[.toast]:left-auto group-[.toast]:transform-none',
        },
      }}
      style={
        {
          '--normal-bg': 'var(--color-surface-lowest)',
          '--normal-text': 'var(--color-on-surface)',
          '--normal-border': 'var(--color-surface-container)',
        } as React.CSSProperties
      }
      {...props}
    />
  )
}

export { Toaster }
