const WEEK_DAYS = ['Пн', 'Вт', 'Ср', 'Чт', 'Пт', 'Сб', 'Вс']

export const useStreakDays = (streak: number) => {
  const todayIndex = (new Date().getDay() + 6) % 7
  const streakInCycle = streak === 0 ? 0 : streak % 7 === 0 ? 7 : streak % 7
  const startDayIndex = (todayIndex - (streakInCycle > 0 ? streakInCycle - 1 : 0) + 7) % 7

  const daysCycle = Array.from({ length: 7 }, (_, stepIndex) => {
    const dayOfWeekIndex = (startDayIndex + stepIndex) % 7
    const dayStep = stepIndex + 1
    const isDone = dayStep <= streakInCycle
    const isNextGoal = dayStep === streakInCycle + 1
    const isFuture = dayStep > streakInCycle + 1
    const isBonus = dayStep === 7
    const isToday = dayOfWeekIndex === todayIndex

    return {
      label: WEEK_DAYS[dayOfWeekIndex],
      dayStep,
      isDone,
      isNextGoal,
      isFuture,
      isBonus,
      isToday,
    }
  })

  return { daysCycle }
}
