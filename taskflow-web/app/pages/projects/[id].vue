<script setup lang="ts">
import type { Project, RealtimeEvent } from '~/types'

const route = useRoute()
const projectId = route.params.id as string

const { getProject, fetchProjects, refreshProject } = useProjects()
const { fetchTasksByProject } = useTasks()

const project = computed(() => getProject(projectId))

onMounted(async () => {
  await Promise.all([
    fetchProjects(),
    fetchTasksByProject(projectId)
  ])
})

function onRealtime(event: RealtimeEvent) {
  // Les events métier transportent un projectId ; on ne réagit qu'à ceux du projet courant.
  if (event.projectId !== projectId) return

  if (event.type.startsWith('task.')) {
    fetchTasksByProject(projectId)
  }
  if (event.type === 'member.added') {
    refreshProject(projectId)
  }
}

const { connected } = useRealtime(() => projectId, onRealtime)
</script>

<template>
  <UContainer class="py-8">
    <div v-if="!project" class="py-16">
      <UEmpty
        icon="i-lucide-folder-x"
        title="Project not found"
        description="This project does not exist or has been deleted."
      >
        <template #actions>
          <UButton
            label="Back to Projects"
            icon="i-lucide-arrow-left"
            to="/"
          />
        </template>
      </UEmpty>
    </div>

    <template v-else>
      <div class="mb-6">
        <UButton
          label="Back to Projects"
          icon="i-lucide-arrow-left"
          variant="link"
          color="neutral"
          to="/"
        />
      </div>

      <div class="mb-6 flex items-start justify-between gap-4">
        <div>
          <h1 class="text-2xl font-bold">{{ project.name }}</h1>
          <p v-if="project.description" class="mt-1 text-sm text-muted">
            {{ project.description }}
          </p>
        </div>
        <UBadge
          :color="connected ? 'success' : 'neutral'"
          :variant="connected ? 'subtle' : 'outline'"
          :icon="connected ? 'i-lucide-radio' : 'i-lucide-radio-tower'"
        >
          {{ connected ? 'Temps réel' : 'Déconnecté' }}
        </UBadge>
      </div>

      <ProjectMembers :project="project as Project" />

      <KanbanBoard :project-id="projectId" />
    </template>
  </UContainer>
</template>
