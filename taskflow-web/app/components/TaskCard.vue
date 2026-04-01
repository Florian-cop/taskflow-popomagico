<script setup lang="ts">
import { VALID_TRANSITIONS, STATUS_LABELS } from '~/types'
import type { Task } from '~/types'

const props = defineProps<{
  task: Task
}>()

const emit = defineEmits<{
  move: [taskId: string]
}>()

const nextStatus = computed(() => VALID_TRANSITIONS[props.task.status])
const nextLabel = computed(() => nextStatus.value ? STATUS_LABELS[nextStatus.value] : null)
</script>

<template>
  <UCard class="group">
    <div class="space-y-3">
      <div>
        <p class="text-sm font-semibold">{{ task.title }}</p>
        <p v-if="task.description" class="mt-1 text-xs text-muted line-clamp-2">
          {{ task.description }}
        </p>
      </div>

      <div class="flex items-center justify-between">
        <UBadge v-if="task.assigneeId" size="xs" variant="subtle" color="neutral" icon="i-lucide-user">
          {{ task.assigneeId }}
        </UBadge>
        <span v-else class="text-xs text-dimmed italic">Unassigned</span>

        <UButton
          v-if="nextStatus"
          size="xs"
          variant="soft"
          :label="`${nextLabel} →`"
          @click="emit('move', task.id)"
        />
      </div>
    </div>
  </UCard>
</template>
