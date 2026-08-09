import { Button } from '@/shared/ui'
import { Coins, LogOut, User as UserIcon } from 'lucide-react'
import { useHeaderDashboard } from '../model/useHeaderDashboard'

export const HeaderDashboard = () => {
  const { user, coins, username, handleLogout, isLogoutLoading } = useHeaderDashboard()

  return (
    <header className="flex items-center justify-between pb-4 sm:pb-6 border-b border-surface-high flex-wrap gap-3 sm:gap-4">
      <div className="max-w-xs sm:max-w-sm">
        <h2 className="text-xl sm:text-2xl font-bold text-on-surface truncate">
          Привет, {username}! 👋
        </h2>
        <p className="text-xs sm:text-sm text-on-surface-variant">
          С возвращением в Авито Тамагочи
        </p>
      </div>
      <div className="flex items-center gap-2 sm:gap-3 flex-wrap">
        <div className="flex items-center gap-1.5 sm:gap-2 px-3 sm:px-4 py-1.5 sm:py-2 bg-avito-yellow/15 text-on-surface rounded-xl border border-avito-yellow/30 font-bold text-xs sm:text-sm">
          <Coins className="size-4 sm:size-5 text-avito-yellow shrink-0" />
          <span>{coins} монет</span>
        </div>
        <div className="flex items-center gap-1.5 sm:gap-2 px-3 sm:px-3.5 py-1.5 sm:py-2 bg-surface-low rounded-xl text-on-surface font-semibold text-xs sm:text-sm">
          <UserIcon className="size-3.5 sm:size-4 text-on-surface-variant shrink-0" />
          <span className="truncate max-w-24 sm:max-w-32">{user?.username}</span>
        </div>
        <Button
          variant="ghost"
          size="sm"
          className="group h-8 w-8 sm:h-9 sm:w-9 p-0"
          disabled={isLogoutLoading}
          onClick={handleLogout}
          title="Выйти"
        >
          <LogOut className="size-4 text-on-surface-variant group-hover:text-avito-red transition-colors" />
        </Button>
      </div>
    </header>
  )
}
