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
