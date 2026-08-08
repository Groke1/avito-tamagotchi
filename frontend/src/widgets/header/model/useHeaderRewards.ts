import { useGetRewardsQuery } from '@/entities/reward'
import { useGetProfileQuery } from '@/entities/user'
import { useAppSelector } from '@/shared/model'

export const useHeaderRewards = () => {
  useGetProfileQuery()
  const user = useAppSelector((state) => state.user.user)
  const { data: rewardsData } = useGetRewardsQuery()

  const totalCount = rewardsData?.items.length ?? 0
  const activeCount = rewardsData?.items.filter((r) => r.status === 'active').length ?? 0

  return {
    coins: user?.coins ?? 0,
    totalCount,
    activeCount,
  }
}
