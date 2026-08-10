import type { DailyStatResponse } from '@/entities/user'
import { renderHook } from '@testing-library/react'
import { describe, expect, it } from 'vitest'
import { useDailyReportCards } from './useDailyReportCards'

describe('useDailyReportCards', () => {
  const mockData: DailyStatResponse = {
    user_id: 'user-1',
    streak: 5,
    pet: {
      daily_gained_xp: 50,
    },
    tasks: [
      {
        id: '1',
        title: 'Task 1',
        reward_coins: 10,
        reward_xp: 20,
        finished_desc: 'Done',
        updated_at: '2026-08-10',
      },
      {
        id: '2',
        title: 'Task 2',
        reward_coins: 15,
        reward_xp: 30,
        finished_desc: 'Done',
        updated_at: '2026-08-10',
      },
    ],
    rewards: [
      {
        promo_code: 'PROMO1',
        name: 'Reward 1',
        description: 'Desc 1',
        finished_desc: 'Finished',
        created_time: '2026-08-10',
        status: 'active',
      },
      {
        promo_code: 'PROMO2',
        name: 'Reward 2',
        description: 'Desc 2',
        finished_desc: 'Finished',
        created_time: '2026-08-10',
        status: 'redeemed',
      },
    ],
  }

  it('calculates metrics correctly', () => {
    const { result } = renderHook(() => useDailyReportCards(mockData))

    expect(result.current.completedTasksCount).toBe(2)
    expect(result.current.earnedRewardsCount).toBe(1)
    expect(result.current.totalCoinsEarned).toBe(25)
    expect(result.current.totalXp).toBe(50)
  })

  it('handles empty tasks and rewards gracefully', () => {
    const emptyData: DailyStatResponse = {
      user_id: 'user-2',
      streak: 0,
      pet: { daily_gained_xp: 0 },
      tasks: null,
      rewards: null,
    }

    const { result } = renderHook(() => useDailyReportCards(emptyData))

    expect(result.current.completedTasksCount).toBe(0)
    expect(result.current.earnedRewardsCount).toBe(0)
    expect(result.current.totalCoinsEarned).toBe(0)
    expect(result.current.totalXp).toBe(0)
  })
})
