<script setup lang="ts">
import type { Project } from '~/types'
import { useProjects } from '~/composables/useProjects'

const { projects, loading, fetchProjects, createProject } = useProjects()
const toast = useToast()

const showCreateModal = ref(false)
const newProjectName = ref('')
const newProjectDescription = ref('')

onMounted(() => {
  fetchProjects()
})

async function handleCreateProject() {
  if (!newProjectName.value.trim()) return

  try {
    await createProject({
      name: newProjectName.value.trim(),
      description: newProjectDescription.value.trim()
    })

    toast.add({ title: 'Project created', icon: 'i-lucide-check', color: 'success' })
    newProjectName.value = ''
    newProjectDescription.value = ''
    showCreateModal.value = false
  } catch {
    toast.add({ title: 'Failed to create project', icon: 'i-lucide-x', color: 'error' })
  }
}
</script>

<template>
  <UContainer class="py-8">
    <div class="mb-8 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold">Projects</h1>
        <p class="mt-1 text-sm text-muted">Manage your Kanban projects</p>
      </div>
      <UButton
        label="New Project"
        icon="i-lucide-plus"
        @click="showCreateModal = true"
      />
    </div>

    <div v-if="loading" class="py-16 text-center">
      <UIcon name="i-lucide-loader" class="animate-spin size-8 text-muted" />
    </div>

    <div v-else-if="projects.length === 0" class="py-16">
      <UEmpty
        icon="i-lucide-folder-open"
        title="No projects yet"
        description="Create your first project to get started."
      >
        <template #actions>
          <UButton
            label="Create Project"
            icon="i-lucide-plus"
            @click="showCreateModal = true"
          />
        </template>
      </UEmpty>
    </div>

    <div v-else class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
      <ProjectCard
        v-for="project in projects"
        :key="project.id"
        :project="project as Project"
      />
    </div>

    <UModal v-model:open="showCreateModal">
      <template #header>
        <h3 class="text-lg font-semibold">New Project</h3>
      </template>

      <template #body>
        <div class="space-y-4">
          <UFormField label="Name" required>
            <UInput
              v-model="newProjectName"
              placeholder="Project name"
              autofocus
              class="w-full"
            />
          </UFormField>
          <UFormField label="Description">
            <UTextarea
              v-model="newProjectDescription"
              placeholder="Brief description"
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
            :disabled="!newProjectName.trim()"
            @click="handleCreateProject"
          />
        </div>
      </template>
    </UModal>
  </UContainer>
</template>
