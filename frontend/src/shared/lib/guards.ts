import { type FetchBaseQueryError } from '@reduxjs/toolkit/query'
import type { ApiError } from '../model/types'

export function isFetchBaseQueryError(error: unknown): error is FetchBaseQueryError {
  return typeof error === 'object' && error !== null && 'status' in error
}

export function isApiError<TCode extends string = string>(
  error: unknown,
): error is ApiError<TCode> {
  return typeof error === 'object' && error !== null && 'code' in error && 'message' in error
}
