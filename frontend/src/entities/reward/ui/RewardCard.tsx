import { cn, formatExpirationDate } from '@/shared/lib/utils'
import { Button } from '@/shared/ui'
import { Check, Clock, Copy, Gift } from 'lucide-react'
import type { UserReward } from '../model/types'
import { useRewardCard } from '../model/useRewardCard'

interface RewardCardProps {
  reward: UserReward
  onRedeem?: (promoCode: string) => void
  isRedeeming?: boolean
}

export const RewardCard = ({ reward, onRedeem, isRedeeming = false }: RewardCardProps) => {
  const {
    hasCopied,
    description,
    name,
    isActive,
    isExpired,
    isRedeemed,
    expires_at,
    promo_code,
    handleCopyCode,
  } = useRewardCard(reward)

  return (
    <div
      className={cn(
        'p-5 rounded-card border bg-surface-lowest transition-all flex flex-col justify-between gap-4 shadow-level-1 relative',
        isActive && 'border-surface-highest',
        isRedeemed && 'border-surface-highest opacity-75 grayscale-[0.2]',
        isExpired && 'border-rose-200 bg-rose-50/20 opacity-70',
      )}
    >
      <div className="flex items-start justify-between gap-2">
        <div className="w-12 h-12 rounded-full flex items-center justify-center shrink-0 bg-avito-blue/10 text-avito-blue-dark">
          <Gift className="size-6 stroke-2" />
        </div>
        <div>
          {isActive && (
            <span className="px-3 py-1 rounded-full text-xs font-bold bg-avito-green/20 text-avito-green-dark border border-avito-green/30">
              Активен
            </span>
          )}
          {isRedeemed && (
            <span className="px-3 py-1 rounded-full text-xs font-bold bg-surface-high text-on-surface-variant border border-surface-highest">
              Использован
            </span>
          )}
          {isExpired && (
            <span className="px-3 py-1 rounded-full text-xs font-bold bg-rose-100 text-rose-800 border border-rose-200">
              Истёк
            </span>
          )}
        </div>
      </div>
      <div className="space-y-1">
        <h4 className="font-extrabold text-base text-on-surface leading-tight">{name}</h4>
        <p className="text-xs text-on-surface-variant leading-relaxed line-clamp-2">
          {description}
        </p>
      </div>
      <div className="p-2.5 rounded-xl bg-surface-high/70 border border-dashed border-outline-variant/60 flex items-center justify-between gap-2">
        <div className="flex flex-col">
          <span className="text-[10px] uppercase font-bold tracking-wider text-on-surface-variant/70">
            Промокод
          </span>
          <span className="font-mono font-extrabold text-sm text-on-surface tracking-widest select-all">
            {promo_code}
          </span>
        </div>
        <button
          type="button"
          onClick={handleCopyCode}
          disabled={hasCopied}
          className={cn(
            'p-2 rounded-lg bg-surface-lowest text-on-surface-variant transition-all border border-surface-highest',
            !hasCopied &&
              'hover:bg-avito-blue/10 hover:text-avito-blue active:scale-95 cursor-pointer',
            hasCopied && 'cursor-default bg-emerald-50/80 border-emerald-200 text-emerald-600',
          )}
          title={hasCopied ? 'Скопировано' : 'Скопировать промокод'}
        >
          {hasCopied ? (
            <Check className="size-4 text-emerald-600 stroke-[2.5]" />
          ) : (
            <Copy className="size-4" />
          )}
        </button>
      </div>
      <div className="flex items-center justify-between text-xs text-on-surface-variant">
        <span className="inline-flex items-center gap-1.5 font-medium">
          <Clock className="size-3.5" />
          {isRedeemed ? 'Использован' : `Срок: ${formatExpirationDate(expires_at)}`}
        </span>
      </div>
      {isActive ? (
        <Button
          size="sm"
          onClick={() => onRedeem?.(promo_code)}
          disabled={isRedeeming}
          isLoading={isRedeeming}
          className="w-full"
        >
          Применить
        </Button>
      ) : isRedeemed ? (
        <div className="w-full py-2 px-4 rounded-xl bg-surface-high text-on-surface-variant font-bold text-xs flex items-center justify-center gap-1.5 border border-surface-highest">
          <Check className="size-4 text-emerald-600" />
          Использовано
        </div>
      ) : (
        <div className="w-full py-2 px-4 rounded-xl bg-rose-100/70 text-rose-800 font-bold text-xs flex items-center justify-center border border-rose-200/60">
          Срок действия истёк
        </div>
      )}
    </div>
  )
}
