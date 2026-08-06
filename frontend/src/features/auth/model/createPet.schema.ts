import { z } from 'zod'

export const createPetSchema = z.object({
  name: z
    .string()
    .min(1, 'Имя питомца обязательно')
    .max(30, 'Имя питомца не должно превышать 30 символов')
    .trim(),
})

export type CreatePetFormData = z.infer<typeof createPetSchema>
