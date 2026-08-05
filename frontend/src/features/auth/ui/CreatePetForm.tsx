import { Button, FormField } from '@/shared/ui'
import { PawPrint } from 'lucide-react'
import { useCreatePetForm } from '../model/useCreatePetForm'

export const CreatePetForm = () => {
  const { register, errors, isLoading, handleSubmit } = useCreatePetForm()

  return (
    <div className="bg-surface-lowest rounded-3xl shadow-level-1 p-6 md:p-8 w-full border-2 border-white relative overflow-hidden z-10">
      <form onSubmit={handleSubmit} className="flex flex-col gap-5 mt-1" noValidate>
        <FormField
          label="Как назовем питомца?"
          id="petName"
          type="text"
          placeholder="например, Барсик"
          disabled={isLoading}
          leftIcon={<PawPrint className="w-5 h-5" />}
          error={errors.name?.message}
          {...register('name')}
        />
        <Button type="submit" isLoading={isLoading} className="mt-2">
          {isLoading ? 'Создание...' : 'Начать игру'}
        </Button>
      </form>
    </div>
  )
}
