import { useGetLeaderboardQuery } from '@/entities/leaderboard'
import { cn } from '@/shared/lib/utils'
import { useState } from 'react'
import { CurrentUserinfo } from './CurrentUserinfo'
import { LeaderboardTop } from './LeaderboardTop'
import { LeadersList } from './LeadersList'

const LIMIT_OPTIONS = [10, 20, 50]

export const LeaderboardWidget = () => {
  const [limit, setLimit] = useState<number>(10)
  const { data, isLoading } = useGetLeaderboardQuery({ limit })

  return (
    <div className="space-y-8 mt-6">
      <LeaderboardTop leaders={data?.items} />
      <div className="space-y-6">
        <CurrentUserinfo user={data?.current_user} isLoading={isLoading} />
        <div className="bg-surface-lowest rounded-card p-6 sm:p-8 shadow-level-1 overflow-hidden">
          <div className="flex items-center justify-between flex-wrap gap-4 mb-2">
            <h2 className="text-on-surface text-2xl font-bold">Топ Игроков</h2>
            <div className="flex items-center bg-surface-low p-1 rounded-xl border border-surface-highest">
              {LIMIT_OPTIONS.map((option) => (
                <button
                  key={option}
                  onClick={() => setLimit(option)}
                  className={cn(
                    'px-3.5 py-1.5 rounded-lg text-xs font-bold transition-all duration-200 cursor-pointer',
                    limit === option
                      ? 'bg-surface-lowest text-on-surface shadow-xs'
                      : 'text-on-surface-variant hover:text-on-surface',
                  )}
                >
                  Топ {option}
                </button>
              ))}
            </div>
          </div>
          <LeadersList
            leaders={data?.items}
            currentUserRank={data?.current_user?.rank}
            currentUserName={data?.current_user?.user_name}
            isLoading={isLoading}
          />
        </div>
      </div>
    </div>
  )
}
