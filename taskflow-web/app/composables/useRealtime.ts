import type { RealtimeEvent } from '~/types'

/**
 * useRealtime — WebSocket client scoped par projectId.
 *
 * - Construit l'URL ws(s):// à partir de l'apiBase (http → ws, https → wss).
 * - Passe le JWT via query param (les WebSocket navigateur ne supportent pas
 *   de header Authorization custom — c'est le raccourci assumé côté backend).
 * - Le composant fournit un handler ; nous nous occupons du cycle de vie.
 */
export function useRealtime(projectId: MaybeRefOrGetter<string>, onEvent: (event: RealtimeEvent) => void) {
  const config = useRuntimeConfig()
  const { token } = useAuth()

  const socket = ref<WebSocket | null>(null)
  const connected = ref(false)
  let reconnectTimer: ReturnType<typeof setTimeout> | null = null
  let manuallyClosed = false

  function buildUrl(pid: string, jwt: string): string {
    const base = config.public.apiBase as string
    const wsBase = base.replace(/^http/, 'ws')
    return `${wsBase}/projects/${pid}/ws?token=${encodeURIComponent(jwt)}`
  }

  function connect() {
    if (!import.meta.client) return
    const pid = toValue(projectId)
    if (!pid || !token.value) return

    manuallyClosed = false
    try {
      const ws = new WebSocket(buildUrl(pid, token.value))
      socket.value = ws

      ws.onopen = () => {
        connected.value = true
      }

      ws.onmessage = (msg) => {
        try {
          const parsed = JSON.parse(msg.data) as RealtimeEvent
          onEvent(parsed)
        } catch (err) {
          console.warn('[Realtime] message non JSON:', msg.data, err)
        }
      }

      ws.onclose = () => {
        connected.value = false
        socket.value = null
        if (!manuallyClosed) {
          reconnectTimer = setTimeout(connect, 2000)
        }
      }

      ws.onerror = (err) => {
        console.warn('[Realtime] erreur socket:', err)
      }
    } catch (err) {
      console.error('[Realtime] impossible d\'ouvrir la WS:', err)
    }
  }

  function disconnect() {
    manuallyClosed = true
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    socket.value?.close()
    socket.value = null
    connected.value = false
  }

  onMounted(connect)
  onBeforeUnmount(disconnect)

  return {
    connected: readonly(connected)
  }
}
