<script setup lang="ts">
import type { Project } from '~/types'

defineProps<{
  project: Project
}>()
</script>

<template>
  <NuxtLink :to="`/projects/${project.id}`">
    <UCard class="h-full transition-shadow hover:shadow-lg cursor-pointer">
      <div class="space-y-3">
        <div class="flex items-start justify-between">
          <h3 class="text-base font-semibold">{{ project.name }}</h3>
          <UBadge size="xs" variant="subtle" color="neutral">
            {{ project.members.length }} member{{ project.members.length !== 1 ? 's' : '' }}
          </UBadge>
        </div>
        <p v-if="project.description" class="text-sm text-muted line-clamp-2">
          {{ project.description }}
        </p>
        <p class="text-xs text-dimmed">
          Created {{ new Date(project.createdAt).toLocaleDateString('fr-FR') }}
        </p>
      </div>
    </UCard>
  </NuxtLink>
</template>
