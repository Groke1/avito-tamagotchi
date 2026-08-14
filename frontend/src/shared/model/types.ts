export interface ApiError<TCode extends string = string> {
  code: TCode
  message: string
  retry_after?: number
}

export type AuthErrorCode =
  | 'INVALID_CREDENTIALS'
  | 'USER_ALREADY_EXISTS'
  | 'INVALID_REFRESH_TOKEN'
  | 'UNAUTHORIZED'
  | 'VALIDATION_ERROR'

export type PetErrorCode =
  | 'PET_NOT_FOUND'
  | 'PET_ALREADY_EXISTS'
  | 'PET_ACTION_UNAVAILABLE'
  | 'INSUFFICIENT_COINS'
  | 'TRIP_GENERATION_ERROR'

export type TaskErrorCode = 'TASK_NOT_FOUND' | 'TASK_ALREADY_COMPLETED' | 'TASK_EXPIRED'

export type ApiErrorCode = AuthErrorCode | PetErrorCode | TaskErrorCode
