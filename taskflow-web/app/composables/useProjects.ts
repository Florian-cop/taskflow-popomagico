import type { Project } from '~/types'

export function useProjects() {
  const projects = useState<Project[]>('projects', () => [])
  const loading = useState<boolean>('projects-loading', () => false)

  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase as string

  async function fetchProjects() {
    loading.value = true
    try {
      const data = await $fetch<Project[]>(`${apiBase}/projects`)
      projects.value = data ?? []
    } catch (error) {
      console.error('[TaskFlow] Failed to fetch projects:', error)
    } finally {
      loading.value = false
    }
  }

  async function createProject(data: { name: string, description: string }) {
    try {
      const project = await $fetch<Project>(`${apiBase}/projects`, {
        method: 'POST',
        headers: { 'X-User-Id': 'user-1' },
        body: { name: data.name, description: data.description }
      })

      console.log('[TaskFlow] Project created:', project)

      projects.value.push(project)
      return project
    } catch (error) {
      console.error('[TaskFlow] Failed to create project:', error)
      throw error
    }
  }

  function getProject(id: string) {
    return projects.value.find(p => p.id === id)
  }

  async function addMember(projectId: string, userId: string) {
    try {
      const project = await $fetch<Project>(`${apiBase}/projects/${projectId}/members`, {
        method: 'POST',
        body: { userId }
      })
      const idx = projects.value.findIndex(p => p.id === projectId)
      if (idx !== -1) projects.value[idx] = project
    } catch (error) {
      console.error('[TaskFlow] Failed to add member:', error)
      throw error
    }
  }

  return {
    projects: readonly(projects),
    loading: readonly(loading),
    fetchProjects,
    createProject,
    getProject,
    addMember
  }
}
