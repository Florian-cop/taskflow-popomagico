import type { Project } from '~/types'

export function useProjects() {
  const projects = useState<Project[]>('projects', () => [])
  const loading = useState<boolean>('projects-loading', () => false)
  const api = useApi()

  async function fetchProjects() {
    loading.value = true
    try {
      const data = await api<Project[]>('/projects')
      projects.value = data ?? []
    } catch (error) {
      console.error('[TaskFlow] Failed to fetch projects:', error)
    } finally {
      loading.value = false
    }
  }

  async function createProject(data: { name: string, description: string }) {
    const project = await api<Project>('/projects', {
      method: 'POST',
      body: { name: data.name, description: data.description }
    })
    projects.value.push(project)
    return project
  }

  function getProject(id: string) {
    return projects.value.find(p => p.id === id)
  }

  async function refreshProject(id: string) {
    const fresh = await api<Project>(`/projects/${id}`)
    const idx = projects.value.findIndex(p => p.id === id)
    if (idx !== -1) {
      projects.value[idx] = fresh
    } else {
      projects.value.push(fresh)
    }
    return fresh
  }

  async function addMember(projectId: string, userId: string) {
    const project = await api<Project>(`/projects/${projectId}/members`, {
      method: 'POST',
      body: { userId }
    })
    const idx = projects.value.findIndex(p => p.id === projectId)
    if (idx !== -1) projects.value[idx] = project
    return project
  }

  return {
    projects: readonly(projects),
    loading: readonly(loading),
    fetchProjects,
    createProject,
    getProject,
    refreshProject,
    addMember
  }
}
