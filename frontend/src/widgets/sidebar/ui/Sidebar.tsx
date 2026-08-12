import { useHandleFeed } from '@/features/pet-actions'
import { useAppSelector } from '@/shared/model'
import { Button, LazyImage } from '@/shared/ui'
import { Utensils } from 'lucide-react'
import { SidebarNavList } from './SidebarNavList'

export const Sidebar = () => {
  const pet = useAppSelector((state) => state.pet.pet)
  const { handleFeed, isFeedLoading } = useHandleFeed()

  return (
    <aside className="hidden lg:flex flex-col gap-6 bg-surface-lowest px-4 py-6 rounded-r-section shadow-level-1 h-screen sticky top-0 overflow-x-hidden">
      <h1 className="flex items-center gap-1.5 px-1 mb-2">
        <span className="text-4xl font-black tracking-tight text-avito-green">Авито</span>
        <span className="bg-avito-blue/15 text-avito-blue px-2.5 py-0.5 rounded-lg text-sm font-semibold">
          Тамагочи
        </span>
      </h1>
      <div className="flex items-center gap-3 p-2.5 bg-surface-low rounded-xl">
        <div className="size-10 rounded-full flex items-center justify-center shrink-0">
          <LazyImage src="/avito-kot.png" alt="Питомец Авито Тамагочи" className="object-contain" />
        </div>
        <div className="min-w-0 flex-1">
          <h4 className="font-semibold text-sm text-on-surface truncate">{pet?.name}</h4>
          <p className="text-xs text-on-surface-variant font-medium truncate">
            Level {pet?.level} • XP {pet?.xp}
          </p>
        </div>
      </div>
      <div className="flex-1">
        <SidebarNavList />
      </div>
      <Button
        variant="avitoGreen"
        className="text-sm"
        onClick={handleFeed}
        disabled={isFeedLoading}
        isLoading={isFeedLoading}
      >
        <Utensils className="size-5 shrink-0" />
        <span>Покормить питомца</span>
      </Button>
    </aside>
  )
}
