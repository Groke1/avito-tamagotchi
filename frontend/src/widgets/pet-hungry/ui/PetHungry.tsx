import { useAppSelector } from '@/shared/model'
import { TriangleAlert } from 'lucide-react'

export const PetHungry = () => {
  const pet = useAppSelector((state) => state.pet.pet)

  if (!pet || pet.satiety > 20) return null

  return (
    <div className="flex items-center gap-2 bg-avito-red/10 text-avito-red font-semibold text-sm rounded-lg p-2">
      <TriangleAlert className="size-5" />
      Питомец голоден! Покормите его до конца дня
    </div>
  )
}
