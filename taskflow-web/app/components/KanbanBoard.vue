<script setup lang="ts">
import { STATUSES } from '~/types'
import { useTasks } from '~/composables/useTasks'

const props = defineProps<{
  projectId: string
}>()

const { getTasksByStatus, createTask, moveTask } = useTasks()
const toast = useToast()

const showCreateModal = ref(false)
const newTaskTitle = ref('')
const newTaskDescription = ref('')

function handleCreateTask() {
  if (!newTaskTitle.value.trim()) return

  createTask({
    title: newTaskTitle.value.trim(),
    description: newTaskDescription.value.trim(),
    projectId: props.projectId
  })

  toast.add({ title: 'Task created', icon: 'i-lucide-check', color: 'success' })
  newTaskTitle.value = ''
  newTaskDescription.value = ''
  showCreateModal.value = false
}

function handleMoveTask(taskId: string) {
  moveTask(taskId)
  toast.add({ title: 'Task moved', icon: 'i-lucide-arrow-right', color: 'info' })
}
</script>

<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h2 class="text-lg font-bold">Board</h2>
      <UButton
        label="New Task"
        icon="i-lucide-plus"
        @click="showCreateModal = true"
      />
    </div>

    <div class="flex gap-6 overflow-x-auto pb-4">
      <KanbanColumn
        v-for="status in STATUSES"
        :key="status"
        :status="status"
        :tasks="getTasksByStatus(projectId, status).value"
        @move-task="handleMoveTask"
      />
    </div>

    <UModal v-model:open="showCreateModal">
      <template #header>
        <h3 class="text-lg font-semibold">New Task</h3>
      </template>

      <template #body>
        <div class="space-y-4">
          <UFormField label="Title" required>
            <UInput
              v-model="newTaskTitle"
              placeholder="Task title"
              autofocus
              class="w-full"
            />
          </UFormField>
          <UFormField label="Description">
            <UTextarea
              v-model="newTaskDescription"
              placeholder="Task description"
              :rows="3"
              class="w-full"
            />
          </UFormField>
        </div>
      </template>

      <template #footer>
        <div class="flex justify-end gap-3">
          <UButton
            label="Cancel"
            color="neutral"
            variant="outline"
            @click="showCreateModal = false"
          />
          <UButton
            label="Create"
            icon="i-lucide-plus"
            :disabled="!newTaskTitle.trim()"
            @click="handleCreateTask"
          />
        </div>
      </template>
    </UModal>
  </div>
</template>
