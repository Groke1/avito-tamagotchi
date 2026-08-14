export interface Pet {
  id: number
  name: string
  level: number
  xp: number
  next_level_xp: number
  satiety: number
  happiness: number
}

export type PetResponse = Pet
export interface PetTicketResponse {
  ticket: string
}

export interface TripReward {
  id: string
  promo_code: string
  name: string
  description: string
  status: string
  expires_at: string
  earned_reason: string
  redeemed_at: string | null
}

export interface TripResult {
  story: string
  coins: number
  xp: number
  reward?: TripReward | null
}

export interface PetDto {
  name: string
}

export interface PetState {
  pet: Pet | null
  latestTrip: TripResult | null
  isInitialized: boolean
}

export interface PetTripResponse {
  story: string
  xp: number
  coins: number
}
