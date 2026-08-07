import { useAppDispatch, useAppSelector } from '@/app/store/hooks'
import { logout, useGetProfileQuery } from '@/entities/user'
import { Button } from '@/shared/ui'
import { Coins, LogOut, User as UserIcon } from 'lucide-react'

export const Header = () => {
  useGetProfileQuery()
  const user = useAppSelector((state) => state.user.user)
  const dispatch = useAppDispatch()

  const handleLogout = () => {
    dispatch(logout())
  }

  return (
    <header className="flex items-center justify-between pb-6 border-b border-surface-high">
      <div className="max-w-sm">
        <h2 className="text-2xl font-bold text-on-surface truncate">
          Привет, {user?.username || 'Пользователь'}! 👋
        </h2>
        <p className="text-sm text-on-surface-variant">С возвращением в Авито Тамагочи</p>
      </div>
      <div className="flex items-center gap-3">
        <div className="flex items-center gap-2 px-4 py-2 bg-avito-yellow/15 text-on-surface rounded-xl border border-avito-yellow/30 font-bold text-sm">
          <Coins className="size-5 text-avito-yellow shrink-0" />
          <span>{user?.coins ?? 0} монет</span>
        </div>
        <div className="flex items-center gap-2 px-3.5 py-2 bg-surface-low rounded-xl text-on-surface font-semibold text-sm">
          <UserIcon className="size-4 text-on-surface-variant" />
          <span className="truncate max-w-32">{user?.username}</span>
        </div>
        <Button variant="ghost" size="sm" className="group" onClick={handleLogout} title="Выйти">
          <LogOut className="size-4 text-on-surface-variant group-hover:text-avito-red transition-colors" />
        </Button>
      </div>
    </header>
  )
}
