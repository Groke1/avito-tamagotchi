export const usePetVitalityBars = ({
  happiness,
  satiety,
  xp,
  totalLevelXp,
}: {
  satiety: number
  happiness: number
  totalLevelXp: number
  xp: number
}) => {
  const satietyClamped = Math.min(Math.max(satiety, 0), 100)
  const happinessClamped = Math.min(Math.max(happiness, 0), 100)

  const xpPercentage = totalLevelXp > 0 ? Math.min(Math.round((xp / totalLevelXp) * 100), 100) : 0

  return { satietyClamped, happinessClamped, xpPercentage }
}
