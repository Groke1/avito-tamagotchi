import * as petEntity from '@/entities/pet'
import * as sharedModel from '@/shared/model'
import { act, renderHook } from '@testing-library/react'
import { toast } from 'sonner'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { useHandleFeed } from './useHandleFeed'

vi.mock('sonner', () => ({
  toast: {
    info: vi.fn(),
    success: vi.fn(),
    error: vi.fn(),
  },
}))

vi.mock('@/entities/pet', async (importOriginal) => {
  const actual = await importOriginal<typeof petEntity>()
  return {
    ...actual,
    useFeedPetMutation: vi.fn(),
  }
})

vi.mock('@/shared/model', () => ({
  useAppDispatch: vi.fn(),
  useAppSelector: vi.fn(),
}))

describe('useHandleFeed', () => {
  const mockDispatch = vi.fn()
  const mockFeedPetTrigger = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(sharedModel.useAppDispatch).mockReturnValue(mockDispatch)
    vi.mocked(petEntity.useFeedPetMutation).mockReturnValue([
      mockFeedPetTrigger,
      { isLoading: false },
    ] as unknown as ReturnType<typeof petEntity.useFeedPetMutation>)
  })

  it('shows info toast and prevents API call if pet satiety is 100', async () => {
    vi.mocked(sharedModel.useAppSelector).mockReturnValue({
      id: 1,
      name: 'Кот',
      satiety: 100,
      happiness: 80,
      level: 1,
      xp: 0,
      next_level_xp: 100,
    })

    const { result } = renderHook(() => useHandleFeed())

    await act(async () => {
      await result.current.handleFeed()
    })

    expect(toast.info).toHaveBeenCalledWith('Питомец полностью сыт! (100/100)')
    expect(mockFeedPetTrigger).not.toHaveBeenCalled()
  })

  it('successfully feeds pet and dispatches setPet when satiety is below 100', async () => {
    const updatedPet = {
      id: 1,
      name: 'Кот',
      satiety: 85,
      happiness: 80,
      level: 1,
      xp: 4,
      next_level_xp: 100,
    }

    vi.mocked(sharedModel.useAppSelector).mockReturnValue({
      id: 1,
      name: 'Кот',
      satiety: 80,
      happiness: 80,
      level: 1,
      xp: 0,
      next_level_xp: 100,
    })

    mockFeedPetTrigger.mockReturnValue({
      unwrap: vi.fn().mockResolvedValue(updatedPet),
    })

    const { result } = renderHook(() => useHandleFeed())

    await act(async () => {
      await result.current.handleFeed()
    })

    expect(mockFeedPetTrigger).toHaveBeenCalled()
    expect(mockDispatch).toHaveBeenCalledWith(petEntity.setPet(updatedPet))
    expect(toast.success).toHaveBeenCalledWith('Вы покормили питомца! 🍎 (+5 сытости, +4 XP)')
  })

  it('handles cooldown retry_after error correctly', async () => {
    vi.mocked(sharedModel.useAppSelector).mockReturnValue({
      id: 1,
      name: 'Кот',
      satiety: 50,
      happiness: 80,
      level: 1,
      xp: 0,
      next_level_xp: 100,
    })

    const cooldownError = {
      status: 400,
      data: {
        code: 'FEED_COOLDOWN',
        message: 'Питомец ещё не проголодался',
        retry_after: 45,
      },
    }

    mockFeedPetTrigger.mockReturnValue({
      unwrap: vi.fn().mockRejectedValue(cooldownError),
    })

    const { result } = renderHook(() => useHandleFeed())

    await act(async () => {
      await result.current.handleFeed()
    })

    expect(toast.info).toHaveBeenCalledWith(
      'Питомец ещё не проголодался. Попробуйте через 45 секунд',
    )
  })

  it('handles generic ApiError without retry_after', async () => {
    vi.mocked(sharedModel.useAppSelector).mockReturnValue({
      id: 1,
      name: 'Кот',
      satiety: 50,
      happiness: 80,
      level: 1,
      xp: 0,
      next_level_xp: 100,
    })

    const genericApiError = {
      status: 400,
      data: {
        code: 'BAD_REQUEST',
        message: 'Недостаточно средств',
      },
    }

    mockFeedPetTrigger.mockReturnValue({
      unwrap: vi.fn().mockRejectedValue(genericApiError),
    })

    const { result } = renderHook(() => useHandleFeed())

    await act(async () => {
      await result.current.handleFeed()
    })

    expect(toast.error).toHaveBeenCalledWith('Недостаточно средств')
  })

  it('handles unknown non-API error gracefully', async () => {
    vi.mocked(sharedModel.useAppSelector).mockReturnValue({
      id: 1,
      name: 'Кот',
      satiety: 50,
      happiness: 80,
      level: 1,
      xp: 0,
      next_level_xp: 100,
    })

    mockFeedPetTrigger.mockReturnValue({
      unwrap: vi.fn().mockRejectedValue(new Error('Network Crash')),
    })

    const { result } = renderHook(() => useHandleFeed())

    await act(async () => {
      await result.current.handleFeed()
    })

    expect(toast.error).toHaveBeenCalledWith('Это действие пока недоступно')
  })
})
