<script setup lang="ts">
import { useProjects } from '~/composables/useProjects'

const route = useRoute()

const breadcrumbs = computed(() => {
  const items = [{ label: 'Projects', to: '/', icon: 'i-lucide-folder' }]
  if (route.params.id) {
    const { getProject } = useProjects()
    const project = getProject(route.params.id as string)
    if (project) {
      items.push({ label: project.name, to: `/projects/${project.id}`, icon: 'i-lucide-kanban' })
    }
  }
  return items
})
</script>

<template>
  <UHeader>
    <template #left>
      <NuxtLink to="/" class="flex items-center gap-2">
        <UIcon name="i-lucide-check-square" class="size-6 text-primary" />
        <span class="text-lg font-bold">TaskFlow</span>
      </NuxtLink>
    </template>

    <template #center>
      <UBreadcrumb :items="breadcrumbs" />
    </template>

    <template #right>
      <UBadge variant="subtle" color="neutral" icon="i-lucide-user">
        user-1
      </UBadge>
      <UColorModeButton />
    </template>
  </UHeader>
</template>
