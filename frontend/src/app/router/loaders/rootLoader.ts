import { store } from '@/app/store'
import { petApi, setPet } from '@/entities/pet'
import {
  authApi,
  getStoredRefreshToken,
  logout,
  setAccessToken,
  setInitialized,
  setStoredRefreshToken,
  setUser,
} from '@/entities/user'

export const rootLoader = async () => {
  if (store.getState().user.isInitialized) return null

  const refreshToken = getStoredRefreshToken()

  if (!refreshToken) {
    store.dispatch(setInitialized(true))
    return null
  }

  try {
    const { access_token, refresh_token } = await store
      .dispatch(authApi.endpoints.refreshToken.initiate(refreshToken))
      .unwrap()

    store.dispatch(setAccessToken(access_token))

    if (refresh_token) {
      setStoredRefreshToken(refresh_token)
    }

    const profile = await store.dispatch(authApi.endpoints.getProfile.initiate()).unwrap()
    store.dispatch(setUser(profile))

    try {
      const pet = await store.dispatch(petApi.endpoints.getPet.initiate()).unwrap()
      store.dispatch(setPet(pet))
    } catch {
      store.dispatch(setPet(null))
    }
  } catch {
    store.dispatch(logout())
  } finally {
    store.dispatch(setInitialized(true))
  }

  return null
}
