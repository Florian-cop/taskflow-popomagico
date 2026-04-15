import type { AuthResponse, User } from '~/types'

/**
 * useAuth — source de vérité pour l'identité utilisateur côté front.
 *
 * - Le token JWT est stocké dans un cookie (SSR-safe, survit au reload).
 * - L'objet user est stocké dans un cookie JSON séparé pour éviter de re-décoder le JWT.
 * - Les mutations passent par login/register/logout ; jamais d'écriture directe sur le cookie ailleurs.
 */
export function useAuth() {
  const token = useCookie<string | null>('taskflow_token', {
    default: () => null,
    maxAge: 60 * 60 * 24,
    sameSite: 'lax'
  })

  const user = useCookie<User | null>('taskflow_user', {
    default: () => null,
    maxAge: 60 * 60 * 24,
    sameSite: 'lax'
  })

  const isAuthenticated = computed(() => !!token.value)
  const fullName = computed(() =>
    user.value ? `${user.value.firstName} ${user.value.lastName}`.trim() || user.value.email : ''
  )

  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase as string

  async function register(payload: { email: string, password: string, firstName: string, lastName: string }) {
    const res = await $fetch<AuthResponse>(`${apiBase}/auth/register`, {
      method: 'POST',
      body: payload
    })
    token.value = res.token
    user.value = res.user
    return res
  }

  async function login(payload: { email: string, password: string }) {
    const res = await $fetch<AuthResponse>(`${apiBase}/auth/login`, {
      method: 'POST',
      body: payload
    })
    token.value = res.token
    user.value = res.user
    return res
  }

  function logout() {
    token.value = null
    user.value = null
  }

  return {
    token,
    user,
    isAuthenticated,
    fullName,
    register,
    login,
    logout
  }
}
