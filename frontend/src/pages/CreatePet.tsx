import { CreatePetForm } from '@/features/auth'
import { LazyImage } from '@/shared/ui'

export const CreatePet = () => {
  return (
    <>
      <div className="flex flex-col items-center mb-6 text-center max-w-md">
        <div className="flex items-center gap-1.5 mb-2">
          <span className="text-2xl font-black tracking-tight text-avito-green">Авито</span>
          <span className="bg-avito-blue/15 text-avito-blue px-2.5 py-0.5 rounded-lg text-sm font-semibold">
            Тамагочи
          </span>
        </div>
        <h1 className="text-2xl md:text-3xl font-bold text-on-surface">Назовите вашего питомца</h1>
        <p className="text-sm text-on-surface-variant mt-1.5 px-2">
          Этот питомец будет сопровождать вас в ваших приключениях на Авито.
        </p>
        <div className="relative my-4 flex justify-center items-center w-48 h-48 md:w-56 md:h-56 z-10">
          <LazyImage
            src="/avito-kot.png"
            alt="Питомец Авито Тамагочи"
            className="object-cover hover:scale-105 transition-all"
          />
        </div>
        <div className="flex items-center gap-2.5 mt-2 w-44">
          <div className="flex-1 h-1.5 bg-surface-high rounded-full overflow-hidden">
            <div className="h-full bg-avito-green w-full rounded-full transition-all duration-300" />
          </div>
          <span className="text-xs font-semibold text-on-surface">Шаг 2 из 2</span>
        </div>
      </div>
      <CreatePetForm />
    </>
  )
}
