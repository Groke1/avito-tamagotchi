import * as petEntity from '@/entities/pet'
import * as sharedModel from '@/shared/model'
import { act, renderHook } from '@testing-library/react'
import { toast } from 'sonner'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { useHandleStroke } from './useHandleStroke'

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
    useStrokePetMutation: vi.fn(),
  }
})

vi.mock('@/shared/model', () => ({
  useAppDispatch: vi.fn(),
  useAppSelector: vi.fn(),
}))

describe('useHandleStroke', () => {
  const mockDispatch = vi.fn()
  const mockStrokePetTrigger = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(sharedModel.useAppDispatch).mockReturnValue(mockDispatch)
    vi.mocked(petEntity.useStrokePetMutation).mockReturnValue([
      mockStrokePetTrigger,
      { isLoading: false },
    ] as unknown as ReturnType<typeof petEntity.useStrokePetMutation>)
  })

  it('shows info toast and prevents API call if pet happiness is 100', async () => {
    vi.mocked(sharedModel.useAppSelector).mockReturnValue({
      id: 1,
      name: 'Кот',
      satiety: 80,
      happiness: 100,
      level: 1,
      xp: 0,
      next_level_xp: 100,
    })

    const { result } = renderHook(() => useHandleStroke())

    await act(async () => {
      await result.current.handleStroke()
    })

    expect(toast.info).toHaveBeenCalledWith('Питомец уже максимально счастлив! (100/100)')
    expect(mockStrokePetTrigger).not.toHaveBeenCalled()
  })

  it('successfully strokes pet and dispatches setPet when happiness is below 100', async () => {
    const updatedPet = {
      id: 1,
      name: 'Кот',
      satiety: 80,
      happiness: 83,
      level: 1,
      xp: 10,
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

    mockStrokePetTrigger.mockReturnValue({
      unwrap: vi.fn().mockResolvedValue(updatedPet),
    })

    const { result } = renderHook(() => useHandleStroke())

    await act(async () => {
      await result.current.handleStroke()
    })

    expect(mockStrokePetTrigger).toHaveBeenCalled()
    expect(mockDispatch).toHaveBeenCalledWith(petEntity.setPet(updatedPet))
    expect(toast.success).toHaveBeenCalledWith('Вы погладили питомца! 🖐️ (+3 счастья, +10 XP)')
  })

  it('handles cooldown retry_after error correctly', async () => {
    vi.mocked(sharedModel.useAppSelector).mockReturnValue({
      id: 1,
      name: 'Кот',
      satiety: 80,
      happiness: 50,
      level: 1,
      xp: 0,
      next_level_xp: 100,
    })

    const cooldownError = {
      status: 400,
      data: {
        code: 'STROKE_COOLDOWN',
        message: 'Питомец устал от ласок',
        retry_after: 30,
      },
    }

    mockStrokePetTrigger.mockReturnValue({
      unwrap: vi.fn().mockRejectedValue(cooldownError),
    })

    const { result } = renderHook(() => useHandleStroke())

    await act(async () => {
      await result.current.handleStroke()
    })

    expect(toast.info).toHaveBeenCalledWith(
      'Питомец устал от ласок. Попробуйте через 30 секунд',
    )
  })
})
