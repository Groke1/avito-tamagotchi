import { useState } from 'react'
import { toast } from 'sonner'
import type { UserReward } from './types'

export const useRewardCard = (reward: UserReward) => {
  const [hasCopied, setHasCopied] = useState(false)
  const { promo_code, name, description, status, expires_at } = reward

  const isActive = status === 'active'
  const isRedeemed = status === 'redeemed'
  const isExpired = status === 'expired'

  const handleCopyCode = async () => {
    try {
      await navigator.clipboard.writeText(promo_code)
      setHasCopied(true)
      toast.success('Промокод скопирован!')
      setTimeout(() => setHasCopied(false), 2000)
    } catch {
      toast.error('Не удалось скопировать промокод')
    }
  }

  return {
    hasCopied,
    isActive,
    isRedeemed,
    isExpired,
    name,
    description,
    expires_at,
    promo_code,
    handleCopyCode,
  }
}
