import type { DailyStatResponse } from '@/entities/user'
import { render, screen } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { DailyReportCards } from './DailyReportCards'

describe('DailyReportCards Component', () => {
  const mockData: DailyStatResponse = {
    user_id: 'user-1',
    streak: 7,
    pet: {
      daily_gained_xp: 150,
    },
    tasks: [
      {
        id: 't-1',
        title: 'Задание 1',
        reward_coins: 20,
        reward_xp: 10,
        finished_desc: 'Завершено',
        updated_at: '2026-08-10',
      },
      {
        id: 't-2',
        title: 'Задание 2',
        reward_coins: 30,
        reward_xp: 15,
        finished_desc: 'Завершено',
        updated_at: '2026-08-10',
      },
    ],
    rewards: [
      {
        promo_code: 'AVITO100',
        name: 'Награда 1',
        description: 'Скидка 100р',
        finished_desc: 'Завершено',
        created_time: '2026-08-10',
        status: 'active',
      },
    ],
  }

  it('renders stats correctly on report cards', () => {
    render(<DailyReportCards data={mockData} />)

    expect(screen.getByText('Заработано')).toBeInTheDocument()
    expect(screen.getByText('+150 XP')).toBeInTheDocument()

    expect(screen.getByText('Выполнено заданий')).toBeInTheDocument()
    expect(screen.getByText('+50 монет')).toBeInTheDocument()

    expect(screen.getByText('Получено наград')).toBeInTheDocument()
    expect(screen.getByText('1')).toBeInTheDocument()

    expect(screen.getByText('Серия дней')).toBeInTheDocument()
    expect(screen.getByText('7 дней')).toBeInTheDocument()
  })
})
