export interface User {
  id: string
  email: string
  username: string
  hasPet?: boolean
  createdAt?: string
}

export interface UserResponse {
  user_id: string
  username: string
  email: string
  coins: number
}

export interface AuthTokens {
  access_token: string
  refresh_token: string
}

export interface AuthResponse {
  user: User
  accessToken: string
  refreshToken: string
}

export interface LoginDto {
  email: string
  password: string
}

export interface RegisterDto {
  username: string
  email: string
  password: string
}

export interface UserState {
  user: User | null
  accessToken: string | null
  isAuthenticated: boolean
  isInitialized: boolean
}
