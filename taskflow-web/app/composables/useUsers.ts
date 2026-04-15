import type { User } from '~/types'

export function useUsers() {
  const api = useApi()

  async function lookupByEmail(email: string): Promise<User> {
    return api<User>('/users/by-email', { query: { email } })
  }

  async function searchByEmail(query: string, limit = 10): Promise<User[]> {
    if (query.trim().length < 2) return []
    return api<User[]>('/users', { query: { search: query, limit } })
  }

  return { lookupByEmail, searchByEmail }
}
