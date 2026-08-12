import { ROUTES_PATHS } from '@/shared/config'
import { NavLink } from 'react-router-dom'
import { TodayQuestsList } from './TodayQuestsList'

export const TodayQuestsWidget = () => {
  return (
    <section className="bg-surface-lowest p-6 rounded-card shadow-level-1 flex flex-col gap-4">
      <div className="flex items-center justify-between gap-4 flex-wrap">
        <h4 className="text-on-surface font-bold text-2xl">Задания на сегодня</h4>
        <NavLink
          to={ROUTES_PATHS.TASKS}
          className="font-bold text-sm text-avito-blue-dark hover:text-avito-blue-dark/80 cursor-pointer transition-all tr"
        >
          Все задания
        </NavLink>
      </div>
      <TodayQuestsList />
    </section>
  )
}
