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

export interface User {
  id: string
  email: string
  firstName: string
  lastName: string
  createdAt: string
  updatedAt: string
}

export interface AuthResponse {
  token: string
  user: User
}

export interface Notification {
  id: string
  type: string
  title: string
  body: string
  readAt: string | null
  createdAt: string
}

export interface NotificationPreferences {
  enabled: Record<string, boolean>
}

export interface AuditLog {
  id: string
  userId: string
  eventName: string
  aggregateType: string
  aggregateId: string
  payload: string
  occurredAt: string
}

export interface RealtimeEvent {
  type: string
  aggregateId: string
  projectId: string
  occurredAt: string
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

export const NOTIFICATION_CHANNELS = [
  { key: 'email', label: 'Email', icon: 'i-lucide-mail', description: 'Recevoir un email à chaque notification' },
  { key: 'in_app', label: 'In-app', icon: 'i-lucide-bell', description: 'Afficher les notifications dans l\'application' }
] as const
