<script setup lang="ts">
import { STATUS_LABELS, STATUS_ICONS, STATUS_COLORS } from '~/types'
import type { Task, TaskStatus } from '~/types'

defineProps<{
  status: TaskStatus
  tasks: Task[]
}>()

const emit = defineEmits<{
  moveTask: [taskId: string]
}>()
</script>

<template>
  <div class="flex w-80 shrink-0 flex-col rounded-lg border border-default bg-default/50">
    <div class="flex items-center gap-2 border-b border-default px-4 py-3">
      <UIcon :name="STATUS_ICONS[status]" :class="`text-${STATUS_COLORS[status]}`" />
      <span class="text-sm font-semibold">{{ STATUS_LABELS[status] }}</span>
      <UBadge :label="String(tasks.length)" size="xs" variant="subtle" color="neutral" class="ml-auto" />
    </div>

    <div class="flex flex-col gap-3 overflow-y-auto p-3" style="max-height: calc(100vh - 280px);">
      <UEmpty
        v-if="tasks.length === 0"
        icon="i-lucide-inbox"
        title="No tasks"
        class="py-6"
      />
      <TaskCard
        v-for="task in tasks"
        :key="task.id"
        :task="task"
        @move="emit('moveTask', $event)"
      />
    </div>
  </div>
</template>
