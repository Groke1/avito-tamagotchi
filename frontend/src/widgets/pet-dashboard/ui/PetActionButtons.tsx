import { Button } from '@/shared/ui'
import { useHandleFeed } from '../model/useHandleFeed'
import { useHandleStroke } from '../model/useHandleStroke'

export const PetActionButtons = () => {
  const { handleFeed, isFeedLoading } = useHandleFeed()
  const { handleStroke, isStrokeLoading } = useHandleStroke()

  return (
    <div className="grid grid-cols-2 gap-3 mt-4">
      <Button
        variant="outline"
        size="sm"
        isLoading={isFeedLoading}
        disabled={isFeedLoading || isStrokeLoading}
        onClick={handleFeed}
        className="hover:bg-avito-green/10 hover:border-avito-green/60 hover:text-avito-green-dark transition-all duration-200"
      >
        🍎 Покормить (-5 монет)
      </Button>
      <Button
        variant="outline"
        size="sm"
        isLoading={isStrokeLoading}
        disabled={isFeedLoading || isStrokeLoading}
        onClick={handleStroke}
        className="hover:bg-avito-blue/10 hover:border-avito-blue/60 hover:text-avito-blue-dark transition-all duration-200"
      >
        🖐️ Погладить (-7 монет)
      </Button>
    </div>
  )
}
