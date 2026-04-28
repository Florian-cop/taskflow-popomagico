/**
 * useApi — $fetch instance avec JWT automatique + gestion 401.
 *
 * Tous les composables métier (useProjects, useTasks, useNotifications, useAudit)
 * passent par ici. Ajouter un header global = modifier un seul fichier.
 *
 * Versioning : si l'utilisateur a basculé sur v2 (cf. useApiVersion), les routes
 * implémentées en v2 (projects, tasks) sont reroutées et l'enveloppe {data, meta}
 * est déballée pour que les composables ne voient aucune différence.
 */
const V2_ROUTES = /^\/?(projects|tasks)(\/|$|\?)/

export function useApi() {
  const config = useRuntimeConfig()
  const { token, logout } = useAuth()
  const { version } = useApiVersion()

  const baseV1 = config.public.apiBase as string
  const baseV2 = baseV1.replace(/\/v1(\/?$)/, '/v2$1')

  return $fetch.create({
    baseURL: baseV1,
    onRequest({ request, options }) {
      if (token.value) {
        const headers = new Headers(options.headers)
        headers.set('Authorization', `Bearer ${token.value}`)
        options.headers = headers
      }
      if (version.value === 'v2' && V2_ROUTES.test(String(request))) {
        options.baseURL = baseV2
      }
    },
    onResponse({ response }) {
      const body = response._data
      if (
        body
        && typeof body === 'object'
        && !Array.isArray(body)
        && 'data' in body
        && 'meta' in body
        && body.meta
        && typeof body.meta === 'object'
        && 'apiVersion' in body.meta
      ) {
        response._data = body.data
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
