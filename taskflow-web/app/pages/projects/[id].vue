<script setup lang="ts">
import { useProjects } from '~/composables/useProjects'
import { useTasks } from '~/composables/useTasks'

const route = useRoute()
const projectId = route.params.id as string

const { getProject, fetchProjects } = useProjects()
const { fetchTasksByProject } = useTasks()

const project = computed(() => getProject(projectId))

onMounted(async () => {
  await fetchProjects()
  await fetchTasksByProject(projectId)
})
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

      <div class="mb-6">
        <h1 class="text-2xl font-bold">{{ project.name }}</h1>
        <p v-if="project.description" class="mt-1 text-sm text-muted">
          {{ project.description }}
        </p>
      </div>

      <KanbanBoard :project-id="projectId" />
    </template>
  </UContainer>
</template>
