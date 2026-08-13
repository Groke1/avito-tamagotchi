import { formatTaskDate } from '@/shared/lib/utils'
import { useState } from 'react'
import { toast } from 'sonner'
import type { UserReward } from './types'

export const useRewardCard = (reward: UserReward) => {
  const [hasCopied, setHasCopied] = useState(false)
  const { promo_code, name, description, status, expires_at } = reward

  const isActive = status === 'active'
  const isRedeemed = status === 'redeemed'
  const isExpired = status === 'expired'
  const isSupportedCopy = navigator.clipboard && window.isSecureContext

  const formattedExpiresAt = formatTaskDate({ dateStr: expires_at, hasMinutes: true })

  const handleCopyCode = async () => {
    if (isSupportedCopy) {
      try {
        await navigator.clipboard.writeText(promo_code)
        setHasCopied(true)
        toast.success('Промокод скопирован!')
        setTimeout(() => setHasCopied(false), 2000)
      } catch {
        toast.error('Не удалось скопировать промокод')
      }
    }
  }

  return {
    hasCopied,
    isActive,
    isRedeemed,
    isExpired,
    isSupportedCopy,
    name,
    description,
    expires_at,
    formattedExpiresAt,
    promo_code,
    handleCopyCode,
  }
}

