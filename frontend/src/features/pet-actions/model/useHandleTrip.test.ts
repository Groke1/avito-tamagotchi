import * as petEntity from '@/entities/pet'
import * as sharedModel from '@/shared/model'
import { act, renderHook } from '@testing-library/react'
import { toast } from 'sonner'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { useHandleTrip } from './useHandleTrip'

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
    useTripPetMutation: vi.fn(),
  }
})

vi.mock('@/shared/model', () => ({
  useAppSelector: vi.fn(),
}))

describe('useHandleTrip', () => {
  const mockTripPetTrigger = vi.fn()

  const mockPet = {
    id: 1,
    name: 'Кот',
    satiety: 80,
    happiness: 80,
    level: 1,
    xp: 0,
    next_level_xp: 100,
  }

  const mockUser = {
    user_id: '123',
    username: 'test',
    email: 'test@example.com',
    coins: 150,
    streak: 3,
  }

  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(petEntity.useTripPetMutation).mockReturnValue([
      mockTripPetTrigger,
      { isLoading: false },
    ] as unknown as ReturnType<typeof petEntity.useTripPetMutation>)
  })

  it('shows error toast if pet does not exist', async () => {
    vi.mocked(sharedModel.useAppSelector).mockImplementation((selector) =>
      selector({
        pet: { pet: null, latestTrip: null, isInitialized: true },
        user: { user: mockUser, accessToken: null, isAuthenticated: true, isInitialized: true },
      } as unknown as Parameters<typeof selector>[0]),
    )

    const { result } = renderHook(() => useHandleTrip())

    await act(async () => {
      await result.current.handleTrip()
    })

    expect(toast.error).toHaveBeenCalledWith('Сначала создайте питомца')
    expect(mockTripPetTrigger).not.toHaveBeenCalled()
  })

  it('shows error toast if user has less than 100 coins', async () => {
    vi.mocked(sharedModel.useAppSelector).mockImplementation((selector) =>
      selector({
        pet: { pet: mockPet, latestTrip: null, isInitialized: true },
        user: {
          user: { ...mockUser, coins: 50 },
          accessToken: null,
          isAuthenticated: true,
          isInitialized: true,
        },
      } as unknown as Parameters<typeof selector>[0]),
    )

    const { result } = renderHook(() => useHandleTrip())

    await act(async () => {
      await result.current.handleTrip()
    })

    expect(toast.error).toHaveBeenCalledWith('Недостаточно монет для путешествия (нужно 100 монет)')
    expect(mockTripPetTrigger).not.toHaveBeenCalled()
  })

  it('successfully starts trip when pet exists and coins >= 100', async () => {
    vi.mocked(sharedModel.useAppSelector).mockImplementation((selector) =>
      selector({
        pet: { pet: mockPet, latestTrip: null, isInitialized: true },
        user: { user: mockUser, accessToken: null, isAuthenticated: true, isInitialized: true },
      } as unknown as Parameters<typeof selector>[0]),
    )

    mockTripPetTrigger.mockReturnValue({
      unwrap: vi.fn().mockResolvedValue(undefined),
    })

    const { result } = renderHook(() => useHandleTrip())

    await act(async () => {
      await result.current.handleTrip()
    })

    expect(mockTripPetTrigger).toHaveBeenCalledWith(1)
    expect(toast.success).toHaveBeenCalledWith(
      'Питомец отправился в путешествие! 🧭',
      expect.objectContaining({
        description: 'Когда путешествие завершится, мы покажем историю и награду.',
      }),
    )
  })

  it('handles retry_after cooldown error gracefully', async () => {
    vi.mocked(sharedModel.useAppSelector).mockImplementation((selector) =>
      selector({
        pet: { pet: mockPet, latestTrip: null, isInitialized: true },
        user: { user: mockUser, accessToken: null, isAuthenticated: true, isInitialized: true },
      } as unknown as Parameters<typeof selector>[0]),
    )

    const cooldownError = {
      status: 409,
      data: {
        code: 'PET_ACTION_UNAVAILABLE',
        message: 'Питомец уже в путешествии',
        retry_after: 30,
      },
    }

    mockTripPetTrigger.mockReturnValue({
      unwrap: vi.fn().mockRejectedValue(cooldownError),
    })

    const { result } = renderHook(() => useHandleTrip())

    await act(async () => {
      await result.current.handleTrip()
    })

    expect(toast.info).toHaveBeenCalledWith('Питомец уже в путешествии. Попробуйте через 30 секунд')
  })

  it('handles generic API error', async () => {
    vi.mocked(sharedModel.useAppSelector).mockImplementation((selector) =>
      selector({
        pet: { pet: mockPet, latestTrip: null, isInitialized: true },
        user: { user: mockUser, accessToken: null, isAuthenticated: true, isInitialized: true },
      } as unknown as Parameters<typeof selector>[0]),
    )

    const apiError = {
      status: 409,
      data: {
        code: 'INSUFFICIENT_COINS',
        message: 'Недостаточно монет',
      },
    }

    mockTripPetTrigger.mockReturnValue({
      unwrap: vi.fn().mockRejectedValue(apiError),
    })

    const { result } = renderHook(() => useHandleTrip())

    await act(async () => {
      await result.current.handleTrip()
    })

    expect(toast.error).toHaveBeenCalledWith('Недостаточно монет')
  })
})
