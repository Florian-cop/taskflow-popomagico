import type { AuditLog } from '~/types'

export interface AuditQuery {
  aggregateType?: string
  aggregateId?: string
  userId?: string
  limit?: number
}

export function useAudit() {
  const entries = useState<AuditLog[]>('audit-entries', () => [])
  const loading = useState<boolean>('audit-loading', () => false)
  const api = useApi()

  async function fetchLogs(query: AuditQuery = {}) {
    loading.value = true
    try {
      const params = new URLSearchParams()
      if (query.aggregateType) params.set('aggregateType', query.aggregateType)
      if (query.aggregateId) params.set('aggregateId', query.aggregateId)
      if (query.userId) params.set('userId', query.userId)
      if (query.limit) params.set('limit', String(query.limit))
      const qs = params.toString()
      const data = await api<AuditLog[]>(`/audit/logs${qs ? '?' + qs : ''}`)
      entries.value = data ?? []
    } catch (error) {
      console.error('[TaskFlow] Failed to fetch audit logs:', error)
    } finally {
      loading.value = false
    }
  }

  return {
    entries: readonly(entries),
    loading: readonly(loading),
    fetchLogs
  }
}
