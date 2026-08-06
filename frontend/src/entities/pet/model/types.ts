export interface Pet {
  id: string
  name: string
  level: number
  xp: number
  next_level_xp: number
  satiety: number
  happiness: number
}

export type PetResponse = Pet

export interface PetDto {
  name: string
}

export interface PetState {
  pet: Pet | null
  hasPet: boolean | null
}
