import type { Notification, NotificationPreferences } from '~/types'

export function useNotifications() {
  const items = useState<Notification[]>('notifications', () => [])
  const loading = useState<boolean>('notifications-loading', () => false)
  const preferences = useState<NotificationPreferences>('notification-preferences', () => ({ enabled: {} }))
  const api = useApi()

  const unreadCount = computed(() => items.value.filter(n => !n.readAt).length)

  async function fetchAll() {
    loading.value = true
    try {
      const data = await api<Notification[]>('/notifications')
      items.value = data ?? []
    } catch (error) {
      console.error('[TaskFlow] Failed to fetch notifications:', error)
    } finally {
      loading.value = false
    }
  }

  async function markAsRead(id: string) {
    const updated = await api<Notification>(`/notifications/${id}/read`, { method: 'PATCH' })
    const idx = items.value.findIndex(n => n.id === id)
    if (idx !== -1) items.value[idx] = updated
    return updated
  }

  async function markAllAsRead() {
    await Promise.all(items.value.filter(n => !n.readAt).map(n => markAsRead(n.id)))
  }

  async function fetchPreferences() {
    try {
      const data = await api<NotificationPreferences>('/notifications/preferences')
      preferences.value = data ?? { enabled: {} }
    } catch (error) {
      console.error('[TaskFlow] Failed to fetch preferences:', error)
    }
  }

  async function updatePreferences(enabled: Record<string, boolean>) {
    const data = await api<NotificationPreferences>('/notifications/preferences', {
      method: 'PUT',
      body: { enabled }
    })
    preferences.value = data
    return data
  }

  return {
    items: readonly(items),
    loading: readonly(loading),
    preferences: readonly(preferences),
    unreadCount,
    fetchAll,
    markAsRead,
    markAllAsRead,
    fetchPreferences,
    updatePreferences
  }
}
