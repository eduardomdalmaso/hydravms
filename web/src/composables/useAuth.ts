import { ref, computed } from 'vue'
import type { LoginCredentials } from '../types/auth'

export function useAuth() {
  const isAuthenticated = ref(localStorage.getItem('hydra_auth') === 'true')
  const isLoading = ref(false)
  const errorMessage = ref<string | null>(null)
  const username = ref(localStorage.getItem('hydra_user') || 'admin')
  const storedRole = ref(localStorage.getItem('hydra_role') || '')

  const userRole = computed<'admin' | 'supervisor' | 'operator'>(() => {
    const roleLower = storedRole.value.toLowerCase()
    if (roleLower === 'admin' || roleLower === 'superadmin' || roleLower === 'administrator') return 'admin'
    if (roleLower.includes('super')) return 'supervisor'
    if (roleLower === 'operator') return 'operator'

    const userLower = username.value.toLowerCase()
    if (userLower.includes('admin')) return 'admin'
    if (userLower.includes('super')) return 'supervisor'
    return 'operator'
  })

  const isAdmin = computed(() => userRole.value === 'admin')

  const handleLogin = async (creds: LoginCredentials) => {
    isLoading.value = true
    errorMessage.value = null

    try {
      if (!creds.username || !creds.password) {
        throw new Error('Preencha o usuário e a senha.')
      }

      const apiBase = import.meta.env.VITE_API_URL || 'http://localhost:8083'
      const res = await fetch(`${apiBase}/api/v1/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username: creds.username, password: creds.password })
      })

      if (!res.ok) {
        const errData = await res.json().catch(() => ({}))
        throw new Error(errData.message || 'Usuário ou senha incorretos.')
      }

      const data = await res.json()
      username.value = data.user?.username || creds.username
      const role = data.user?.role || (username.value.toLowerCase().includes('admin') ? 'admin' : 'operator')
      storedRole.value = role
      isAuthenticated.value = true
      localStorage.setItem('hydra_auth', 'true')
      localStorage.setItem('hydra_token', data.token)
      localStorage.setItem('hydra_user', username.value)
      localStorage.setItem('hydra_role', role)
    } catch (err: any) {
      errorMessage.value = err.message || 'Credenciais inválidas.'
    } finally {
      isLoading.value = false
    }
  }

  const handleLogout = () => {
    isAuthenticated.value = false
    localStorage.removeItem('hydra_auth')
    localStorage.removeItem('hydra_token')
    localStorage.removeItem('hydra_user')
    localStorage.removeItem('hydra_role')
  }

  return {
    isAuthenticated,
    isLoading,
    errorMessage,
    username,
    userRole,
    isAdmin,
    handleLogin,
    handleLogout
  }
}
