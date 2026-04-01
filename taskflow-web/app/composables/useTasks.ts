import type { Task, TaskStatus } from '~/types'
import { VALID_TRANSITIONS } from '~/types'

export function useTasks() {
  const tasks = useState<Task[]>('tasks', () => [])

  function getTasksByProject(projectId: string) {
    return computed(() => tasks.value.filter(t => t.projectId === projectId))
  }

  function getTasksByStatus(projectId: string, status: TaskStatus) {
    return computed(() =>
      tasks.value.filter(t => t.projectId === projectId && t.status === status)
    )
  }

  function createTask(data: { title: string, description: string, projectId: string }) {
    const now = new Date().toISOString()
    const task: Task = {
      id: crypto.randomUUID(),
      title: data.title,
      description: data.description,
      status: 'todo',
      assigneeId: null,
      projectId: data.projectId,
      createdAt: now,
      updatedAt: now
    }
    tasks.value.push(task)
    return task
  }

  function moveTask(taskId: string) {
    const task = tasks.value.find(t => t.id === taskId)
    if (!task) return
    const next = VALID_TRANSITIONS[task.status]
    if (!next) return
    task.status = next
    task.updatedAt = new Date().toISOString()
  }

  return {
    tasks: readonly(tasks),
    getTasksByProject,
    getTasksByStatus,
    createTask,
    moveTask
  }
}
