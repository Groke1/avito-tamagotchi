import { cn } from '@/shared/lib/utils'
import { type ComponentProps, type FC, useState } from 'react'
import { Skeleton } from './skeleton'

interface LazyImageProps extends ComponentProps<'img'> {
  skeletonClassName?: string
}

export const LazyImage: FC<LazyImageProps> = ({
  src,
  alt,
  className,
  skeletonClassName,
  onLoad,
  ...props
}) => {
  const [isLoaded, setIsLoaded] = useState(false)

  return (
    <div className="relative flex items-center justify-center h-full w-full">
      {!isLoaded && (
        <Skeleton
          className={cn(
            'absolute inset-0 rounded-full animate-pulse',
            skeletonClassName || className,
          )}
        />
      )}
      <img
        src={src}
        alt={alt}
        loading="lazy"
        onLoad={(e) => {
          setIsLoaded(true)
          onLoad?.(e)
        }}
        className={cn(
          'w-full h-full transition-opacity duration-300',
          isLoaded ? 'opacity-100' : 'opacity-0',
          className,
        )}
        {...props}
      />
    </div>
  )
}
