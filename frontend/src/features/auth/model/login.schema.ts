import { z } from 'zod'

export const loginSchema = z.object({
  email: z.email('Введите корректный E-mail').min(1, 'Email обязателен'),
  password: z
    .string()
    .min(1, 'Пароль обязателен')
    .min(8, 'Пароль должен содержать минимум 8 символов'),
})

export type LoginFormData = z.infer<typeof loginSchema>
