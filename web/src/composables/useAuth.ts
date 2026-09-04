import { ref, computed } from 'vue'
import type { LoginCredentials } from '../types/auth'

export function useAuth() {
  const isAuthenticated = ref(localStorage.getItem('hydra_auth') === 'true')
  const isLoading = ref(false)
  const errorMessage = ref<string | null>(null)
  const username = ref(localStorage.getItem('hydra_user') || 'admin')

  const userRole = computed<'admin' | 'supervisor' | 'operator'>(() => {
    if (username.value.toLowerCase() === 'admin') return 'admin'
    if (username.value.toLowerCase().includes('super')) return 'supervisor'
    return 'operator'
  })

  const isAdmin = computed(() => userRole.value === 'admin')

  const handleLogin = async (creds: LoginCredentials) => {
    isLoading.value = true
    errorMessage.value = null

    try {
      if (!creds.username || !creds.password) {
        throw new Error('Preencha o usuario e a senha.')
      }

      username.value = creds.username
      isAuthenticated.value = true
      localStorage.setItem('hydra_auth', 'true')
      localStorage.setItem('hydra_user', creds.username)
    } catch (err: any) {
      errorMessage.value = err.message || 'Credenciais invalidas.'
    } finally {
      isLoading.value = false
    }
  }

  const handleLogout = () => {
    isAuthenticated.value = false
    localStorage.removeItem('hydra_auth')
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
