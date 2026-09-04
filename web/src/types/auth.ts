export interface UserTenant {
  id: string
  name: string
  slug: string
  role: 'admin' | 'operator' | 'viewer'
}

export interface UserProfile {
  id: string
  username: string
  full_name: string
  tenants: UserTenant[]
}

export interface LoginCredentials {
  username: string
  password: string
  remember_me: boolean
}

export interface AuthTokens {
  access_token: string
  refresh_token: string
  expires_in: number
}

export type AuthStep = 'credentials' | 'two_factor' | 'tenant_select'
