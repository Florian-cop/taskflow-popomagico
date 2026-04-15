import type { Task, TaskStatus } from '~/types'
import { VALID_TRANSITIONS } from '~/types'

export function useTasks() {
  const tasks = useState<Task[]>('tasks', () => [])
  const loading = useState<boolean>('tasks-loading', () => false)
  const api = useApi()

  async function fetchTasksByProject(projectId: string) {
    loading.value = true
    try {
      const data = await api<Task[]>(`/projects/${projectId}/tasks`)
      tasks.value = [
        ...tasks.value.filter(t => t.projectId !== projectId),
        ...(data ?? [])
      ]
    } catch (error) {
      console.error('[TaskFlow] Failed to fetch tasks:', error)
    } finally {
      loading.value = false
    }
  }

  function getTasksByProject(projectId: string) {
    return computed(() => tasks.value.filter(t => t.projectId === projectId))
  }

  function getTasksByStatus(projectId: string, status: TaskStatus) {
    return computed(() =>
      tasks.value.filter(t => t.projectId === projectId && t.status === status)
    )
  }

  async function createTask(data: { title: string, description: string, projectId: string }) {
    const task = await api<Task>(`/projects/${data.projectId}/tasks`, {
      method: 'POST',
      body: { title: data.title, description: data.description }
    })
    tasks.value.push(task)
    return task
  }

  async function moveTask(taskId: string) {
    const task = tasks.value.find(t => t.id === taskId)
    if (!task) return
    const next = VALID_TRANSITIONS[task.status]
    if (!next) return

    const updated = await api<Task>(`/tasks/${taskId}/move`, {
      method: 'PUT',
      body: { status: next }
    })
    const idx = tasks.value.findIndex(t => t.id === taskId)
    if (idx !== -1) tasks.value[idx] = updated
    return updated
  }

  return {
    tasks: readonly(tasks),
    loading: readonly(loading),
    fetchTasksByProject,
    getTasksByProject,
    getTasksByStatus,
    createTask,
    moveTask
  }
}
