import { Award, FastForward, type LucideIcon, Star, Zap } from 'lucide-react'

export const taskIconMap: Record<string, LucideIcon> = {
  'Быстрый ответ покупателю': Zap,
  'Лояльный продавец': Star,
  'Первая продажа месяца': FastForward,
}

export const DefaultTaskIcon: LucideIcon = Award
