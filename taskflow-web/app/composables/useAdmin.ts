export interface ChannelStatus {
  name: string
  failing: boolean
  toggleable: boolean
}

export interface FailedNotification {
  id: string
  notificationId: string
  userId: string
  channel: string
  type: string
  title: string
  body: string
  error: string
  retryCount: number
  status: 'pending' | 'retried'
  occurredAt: string
  lastRetriedAt: string | null
}

export function useAdmin() {
  const api = useApi()

  const channels = useState<ChannelStatus[]>('admin-channels', () => [])
  const failed = useState<FailedNotification[]>('admin-failed', () => [])
  const loading = useState<boolean>('admin-loading', () => false)

  async function fetchChannels() {
    const data = await api<ChannelStatus[]>('/admin/notifications/channels')
    channels.value = data ?? []
  }

  async function setChannelFailing(name: string, failing: boolean) {
    await api(`/admin/notifications/channels/${name}`, {
      method: 'PUT',
      body: { failing }
    })
    const idx = channels.value.findIndex(c => c.name === name)
    if (idx !== -1 && channels.value[idx]) {
      channels.value[idx] = { ...channels.value[idx], failing }
    }
  }

  async function fetchFailed() {
    loading.value = true
    try {
      const data = await api<FailedNotification[]>('/admin/notifications/failed', {
        query: { limit: 100 }
      })
      failed.value = data ?? []
    } finally {
      loading.value = false
    }
  }

  async function retryFailed(id: string) {
    await api(`/admin/notifications/failed/${id}/retry`, { method: 'POST' })
    failed.value = failed.value.filter(f => f.id !== id)
  }

  return {
    channels: readonly(channels),
    failed: readonly(failed),
    loading: readonly(loading),
    fetchChannels,
    setChannelFailing,
    fetchFailed,
    retryFailed
  }
}
