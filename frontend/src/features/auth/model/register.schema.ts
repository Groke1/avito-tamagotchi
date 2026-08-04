import { z } from 'zod'

export const registerSchema = z.object({
  username: z
    .string()
    .min(1, 'Имя пользователя обязательно')
    .min(3, 'Имя пользователя должно быть от 3 символов')
    .max(20, 'Имя пользователя не должно превышать 20 символов')
    .regex(/^[a-zA-Z0-9_]+$/, 'Используйте только латинские буквы, цифры и символ подчеркивания'),
  email: z.email('Введите корректный E-mail').min(1, 'Email обязателен'),
  password: z
    .string()
    .min(1, 'Пароль обязателен')
    .min(8, 'Пароль должен содержать минимум 8 символов'),
  terms: z.boolean().refine((val) => val === true, {
    message: 'Необходимо согласиться с условиями',
  }),
})

export type RegisterFormData = z.infer<typeof registerSchema>
