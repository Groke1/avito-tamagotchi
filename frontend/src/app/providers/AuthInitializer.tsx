import {
  authApi,
  getStoredRefreshToken,
  logout,
  setAccessToken,
  setInitialized,
  setStoredRefreshToken,
  setUser,
} from '@/entities/user'
import { type FC, type PropsWithChildren, useEffect } from 'react'
import { useAppDispatch, useAppSelector } from '../store/hooks'

let initAuthPromise: Promise<void> | null = null

export const AuthInitializer: FC<PropsWithChildren> = ({ children }) => {
  const dispatch = useAppDispatch()
  const isInitialized = useAppSelector((state) => state.user.isInitialized)

  useEffect(() => {
    if (initAuthPromise) return

    const initAuth = async () => {
      const refreshToken = getStoredRefreshToken()

      if (!refreshToken) {
        dispatch(setInitialized(true))
        return
      }

      try {
        const { access_token, refresh_token } = await dispatch(
          authApi.endpoints.refreshToken.initiate(refreshToken),
        ).unwrap()

        dispatch(setAccessToken(access_token))

        if (refresh_token) {
          setStoredRefreshToken(refresh_token)
        }

        const profile = await dispatch(authApi.endpoints.getProfile.initiate()).unwrap()
        dispatch(setUser(profile))
      } catch {
        dispatch(logout())
      } finally {
        dispatch(setInitialized(true))
      }
    }

    initAuthPromise = initAuth()
  }, [dispatch])


  if (!isInitialized) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-surface-bg">
        <div className="flex flex-col items-center gap-3">
          <div className="w-10 h-10 border-4 border-avito-green border-t-transparent rounded-full animate-spin" />
          <span className="text-sm font-semibold text-on-surface-variant">Загрузка...</span>
        </div>
      </div>
    )
  }

  return <>{children}</>
}
