import { useRedeemRewardMutation } from '@/entities/reward'
import { useState } from 'react'
import { toast } from 'sonner'

export const useRedeemReward = () => {
  const [redeemReward] = useRedeemRewardMutation()
  const [redeemingPromoCode, setRedeemingPromoCode] = useState<string | null>(null)

  const handleRedeem = async (promoCode: string) => {
    try {
      setRedeemingPromoCode(promoCode)
      await redeemReward({ promo_code: promoCode }).unwrap()
      toast.success('Промокод успешно применен!')
    } catch {
      toast.error('Не удалось применить промокод')
    } finally {
      setRedeemingPromoCode(null)
    }
  }

  return { redeemingPromoCode, handleRedeem }
}
