/**
 * useApi — $fetch instance avec JWT automatique + gestion 401.
 *
 * Tous les composables métier (useProjects, useTasks, useNotifications, useAudit)
 * passent par ici. Ajouter un header global = modifier un seul fichier.
 */
export function useApi() {
  const config = useRuntimeConfig()
  const { token, logout } = useAuth()

  return $fetch.create({
    baseURL: config.public.apiBase as string,
    onRequest({ options }) {
      if (token.value) {
        const headers = new Headers(options.headers)
        headers.set('Authorization', `Bearer ${token.value}`)
        options.headers = headers
      }
    },
    onResponseError({ response }) {
      if (response.status === 401) {
        logout()
        if (import.meta.client) {
          navigateTo('/login')
        }
      }
    }
  })
}
