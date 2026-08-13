import { cn } from '@/shared/lib/utils'
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
    isSupportedCopy,
    formattedExpiresAt,
    promo_code,
    handleCopyCode,
  } = useRewardCard(reward)

  return (
    <div
      className={cn(
        'p-3.5 sm:p-5 rounded-card border bg-surface-lowest transition-all flex flex-col justify-between gap-3 sm:gap-4 shadow-level-1 relative min-w-0',
        isActive && 'border-surface-highest',
        isRedeemed && 'border-surface-highest opacity-75 grayscale-[0.2]',
        isExpired && 'border-rose-200 bg-rose-50/20 opacity-70',
      )}
    >
      <div className="flex items-start justify-between gap-2">
        <div className="w-10 h-10 sm:w-12 sm:h-12 rounded-full flex items-center justify-center shrink-0 bg-avito-blue/10 text-avito-blue-dark">
          <Gift className="size-5 sm:size-6 stroke-2" />
        </div>
        <div>
          {isActive && (
            <span className="px-2.5 sm:px-3 py-0.5 sm:py-1 rounded-full text-[11px] sm:text-xs font-bold bg-avito-green/20 text-avito-green-dark border border-avito-green/30">
              Активен
            </span>
          )}
          {isRedeemed && (
            <span className="px-2.5 sm:px-3 py-0.5 sm:py-1 rounded-full text-[11px] sm:text-xs font-bold bg-surface-high text-on-surface-variant border border-surface-highest">
              Использован
            </span>
          )}
          {isExpired && (
            <span className="px-2.5 sm:px-3 py-0.5 sm:py-1 rounded-full text-[11px] sm:text-xs font-bold bg-rose-100 text-rose-800 border border-rose-200">
              Истёк
            </span>
          )}
        </div>
      </div>
      <div className="space-y-1 min-w-0">
        <h4 className="font-extrabold text-sm sm:text-base text-on-surface leading-tight wrap-break-word">
          {name}
        </h4>
        <p className="text-xs text-on-surface-variant leading-relaxed line-clamp-2 wrap-break-word">
          {description}
        </p>
      </div>
      <div className="p-2 sm:p-2.5 rounded-xl bg-surface-high/70 border border-dashed border-outline-variant/60 flex items-center justify-between gap-2 min-w-0">
        <div className="flex flex-col min-w-0 flex-1">
          <span className="text-[9px] sm:text-[10px] uppercase font-bold tracking-wider text-on-surface-variant/70 truncate">
            Промокод
          </span>
          <span
            className={cn(
              'font-mono font-extrabold text-xs sm:text-sm text-on-surface tracking-wider sm:tracking-widest select-all truncate',
              (isRedeemed || isExpired) && 'select-none',
            )}
          >
            {promo_code}
          </span>
        </div>
        {isSupportedCopy && (
          <button
            type="button"
            onClick={handleCopyCode}
            disabled={hasCopied || isRedeemed || isExpired}
            className={cn(
              'p-1.5 sm:p-2 rounded-lg bg-surface-lowest text-on-surface-variant transition-all border border-surface-highest shrink-0',
              !hasCopied &&
                'hover:bg-avito-blue/10 hover:text-avito-blue active:scale-95 cursor-pointer',
              hasCopied && 'cursor-default bg-emerald-50/80 border-emerald-200 text-emerald-600',
              (isRedeemed || isExpired) && 'pointer-events-none opacity-50',
            )}
            title={hasCopied ? 'Скопировано' : 'Скопировать промокод'}
          >
            {hasCopied ? (
              <Check className="size-3.5 sm:size-4 text-emerald-600 stroke-[2.5]" />
            ) : (
              <Copy className="size-3.5 sm:size-4" />
            )}
          </button>
        )}
      </div>
      <div className="flex items-center justify-between text-[11px] sm:text-xs text-on-surface-variant min-w-0">
        <span className="inline-flex items-center gap-1.5 font-medium min-w-0 truncate">
          <Clock className="size-3.5 shrink-0" />
          <span className="truncate" title={formattedExpiresAt}>
            {isRedeemed ? 'Использован' : `Срок до: ${formattedExpiresAt}`}
          </span>
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
