import { ROUTES_PATHS } from '@/app/router/paths'
import { cn } from '@/shared/lib/utils'
import { Award, CalendarCheck, CheckSquare, Home, ShoppingBag } from 'lucide-react'
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
    label: 'Магазин наград',
    icon: ShoppingBag,
    link: ROUTES_PATHS.SHOP,
  },
  {
    label: 'Лидерборд',
    icon: Award,
    link: ROUTES_PATHS.LEADERBOARD,
  },
  {
    label: 'Ежедневный отчёт',
    icon: CalendarCheck,
    link: ROUTES_PATHS.DAILY_REPORT,
  },
] as const

export const SidebarNavList = () => {
  return (
    <nav>
      <ul className="flex flex-col gap-2">
        {NAV_LIST.map(({ icon: Icon, label, link }) => (
          <li key={label}>
            <NavLink
              to={link}
              className={({ isActive }) =>
                cn(
                  'flex items-center gap-4 p-4 rounded-xl font-semibold text-sm transition-all duration-200',
                  isActive
                    ? 'bg-avito-blue-container text-[#003A5C]'
                    : 'text-on-surface-variant hover:text-on-surface hover:bg-surface-high',
                )
              }
            >
              <Icon className="size-5 shrink-0" />
              <span className="truncate">{label}</span>
            </NavLink>
          </li>
        ))}
      </ul>
    </nav>
  )
}
