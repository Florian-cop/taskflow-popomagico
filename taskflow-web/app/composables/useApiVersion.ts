export type ApiVersion = 'v1' | 'v2'

const STORAGE_KEY = 'taskflow.apiVersion'

export function useApiVersion() {
  const version = useState<ApiVersion>('api-version', () => {
    if (import.meta.client) {
      const stored = window.localStorage.getItem(STORAGE_KEY)
      if (stored === 'v1' || stored === 'v2') return stored
    }
    return 'v1'
  })

  function setVersion(next: ApiVersion) {
    version.value = next
    if (import.meta.client) {
      window.localStorage.setItem(STORAGE_KEY, next)
    }
  }

  return { version, setVersion }
}
