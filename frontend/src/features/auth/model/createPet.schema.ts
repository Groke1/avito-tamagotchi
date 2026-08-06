import { z } from 'zod'

export const createPetSchema = z.object({
  name: z
    .string()
    .min(1, 'Имя питомца обязательно')
    .max(25, 'Имя питомца не должно превышать 25 символов')
    .trim(),
})

export type CreatePetFormData = z.infer<typeof createPetSchema>
