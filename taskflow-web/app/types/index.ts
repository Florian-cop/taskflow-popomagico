export type TaskStatus = 'todo' | 'in_progress' | 'done'

export interface Task {
  id: string
  title: string
  description: string
  status: TaskStatus
  assigneeId: string | null
  projectId: string
  createdAt: string
  updatedAt: string
}

export interface Member {
  userId: string
  role: 'owner' | 'member'
  joinedAt: string
}

export interface Project {
  id: string
  name: string
  description: string
  members: Member[]
  createdAt: string
  updatedAt: string
}

export const STATUS_LABELS: Record<TaskStatus, string> = {
  todo: 'To Do',
  in_progress: 'In Progress',
  done: 'Done'
}

export const STATUS_ICONS: Record<TaskStatus, string> = {
  todo: 'i-lucide-circle',
  in_progress: 'i-lucide-loader',
  done: 'i-lucide-check-circle'
}

export const STATUS_COLORS: Record<TaskStatus, string> = {
  todo: 'neutral',
  in_progress: 'info',
  done: 'success'
}

export const STATUSES: TaskStatus[] = ['todo', 'in_progress', 'done']

export const VALID_TRANSITIONS: Record<TaskStatus, TaskStatus | null> = {
  todo: 'in_progress',
  in_progress: 'done',
  done: null
}
