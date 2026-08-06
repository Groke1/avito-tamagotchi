import { Button } from '@/shared/ui'

export const PetActionButtons = () => {
  return (
    <div className="grid grid-cols-2 gap-3 mt-4">
      <Button
        variant="outline"
        size="sm"
        className="hover:bg-avito-green/10 hover:border-avito-green/60 hover:text-avito-green-dark transition-all duration-200"
      >
        🍎 Покормить
      </Button>
      <Button
        variant="outline"
        size="sm"
        className="hover:bg-avito-blue/10 hover:border-avito-blue/60 hover:text-avito-blue-dark transition-all duration-200"
      >
        🖐️ Погладить
      </Button>
    </div>
  )
}
