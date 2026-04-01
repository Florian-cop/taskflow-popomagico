import type { Project } from '~/types'

export function useProjects() {
  const projects = useState<Project[]>('projects', () => [])

  function createProject(data: { name: string, description: string }) {
    const now = new Date().toISOString()
    const project: Project = {
      id: crypto.randomUUID(),
      name: data.name,
      description: data.description,
      members: [{ userId: 'user-1', role: 'owner', joinedAt: now }],
      createdAt: now,
      updatedAt: now
    }
    projects.value.push(project)
    return project
  }

  function getProject(id: string) {
    return projects.value.find(p => p.id === id)
  }

  function addMember(projectId: string, userId: string) {
    const project = getProject(projectId)
    if (!project) return
    if (project.members.some(m => m.userId === userId)) return
    project.members.push({
      userId,
      role: 'member',
      joinedAt: new Date().toISOString()
    })
  }

  return {
    projects: readonly(projects),
    createProject,
    getProject,
    addMember
  }
}
