import { ROUTES_PATHS } from '@/shared/config'
import { cn } from '@/shared/lib/utils'
import { LazyImage } from '@/shared/ui'
import { Award, CalendarCheck, CheckSquare, Gift, Home } from 'lucide-react'
import { NavLink } from 'react-router-dom'

const NAV_LIST = [
  {
    label: 'Главная',
    icon: Home,
    link: ROUTES_PATHS.DASHBOARD,
  },
  {
    label: 'Задания',
    icon: CheckSquare,
    link: ROUTES_PATHS.TASKS,
  },
  {
    label: 'Награды',
    icon: Gift,
    link: ROUTES_PATHS.REWARDS,
  },
  {
    label: 'Лидерборд',
    icon: Award,
    link: ROUTES_PATHS.LEADERBOARD,
  },
  {
    label: 'Отчёт',
    icon: CalendarCheck,
    link: ROUTES_PATHS.DAILY_REPORT,
  },
] as const

export const MobileNav = () => {
  return (
    <>
      <header className="lg:hidden sticky top-0 z-40 bg-surface-lowest/90 backdrop-blur-md border-b border-surface-highest px-4 py-2.5 flex items-center justify-between shadow-xs">
        <div className="flex items-center gap-2">
          <div className="size-8 rounded-full flex items-center justify-center shrink-0">
            <LazyImage src="/avito-kot.png" alt="Питомец" className="object-contain" />
          </div>
          <h1 className="flex items-center gap-1.5 px-1">
            <span className="text-2xl font-black tracking-tight text-avito-green">Авито</span>
            <span className="bg-avito-blue/15 text-avito-blue px-2.5 py-0.5 rounded-lg text-sm font-semibold">
              Тамагочи
            </span>
          </h1>
        </div>
      </header>
      <nav className="lg:hidden fixed bottom-0 left-0 right-0 z-50 bg-surface-lowest/95 backdrop-blur-md border-t border-surface-highest px-2 py-2.5 flex items-center justify-around shadow-level-2">
        {NAV_LIST.map(({ icon: Icon, label, link }) => (
          <NavLink
            key={label}
            to={link}
            className={({ isActive }) =>
              cn(
                'flex flex-col items-center gap-1 px-3 py-1.5 rounded-xl font-semibold text-[10px] transition-all  min-w-14 text-center',
                isActive ? 'bg-avito-blue-container text-[#003A5C]' : 'text-on-surface-variant ',
              )
            }
          >
            <Icon className="size-5 shrink-0" />
            <span className="truncate max-w-14">{label}</span>
          </NavLink>
        ))}
      </nav>
    </>
  )
}
