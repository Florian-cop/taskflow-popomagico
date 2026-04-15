<script setup lang="ts">
const { entries, loading, fetchLogs } = useAudit()

const filters = reactive({
  aggregateType: '',
  aggregateId: '',
  userId: '',
  limit: 100
})

const typeOptions = [
  { label: 'Tous les types', value: '' },
  { label: 'Projets', value: 'project' },
  { label: 'Tâches', value: 'task' },
  { label: 'Membres', value: 'member' },
  { label: 'Utilisateurs', value: 'user' }
]

onMounted(() => fetchLogs({ limit: 100 }))

function applyFilters() {
  fetchLogs({
    aggregateType: filters.aggregateType || undefined,
    aggregateId: filters.aggregateId || undefined,
    userId: filters.userId || undefined,
    limit: filters.limit
  })
}

function resetFilters() {
  filters.aggregateType = ''
  filters.aggregateId = ''
  filters.userId = ''
  filters.limit = 100
  fetchLogs({ limit: 100 })
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleString('fr-FR', {
    day: '2-digit', month: '2-digit', year: 'numeric',
    hour: '2-digit', minute: '2-digit', second: '2-digit'
  })
}

function iconFor(eventName: string): string {
  if (eventName.startsWith('task.created')) return 'i-lucide-plus-circle'
  if (eventName.startsWith('task.moved')) return 'i-lucide-arrow-right-circle'
  if (eventName.startsWith('task.assigned')) return 'i-lucide-user-check'
  if (eventName.startsWith('project.created')) return 'i-lucide-folder-plus'
  if (eventName.startsWith('member.added')) return 'i-lucide-user-plus'
  if (eventName.startsWith('user.created')) return 'i-lucide-user'
  return 'i-lucide-activity'
}
</script>

<template>
  <UContainer class="py-8">
    <div class="mb-8">
      <h1 class="text-2xl font-bold">Journal d'audit</h1>
      <p class="mt-1 text-sm text-muted">
        Historique de toutes les actions d'écriture sur les projets et les tâches.
        Source : table <code>audit_logs</code>, alimentée par le handler d'events du bounded context <code>audit</code>.
      </p>
    </div>

    <UCard class="mb-6">
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <UFormField label="Type d'agrégat">
          <USelect v-model="filters.aggregateType" :items="typeOptions" class="w-full" />
        </UFormField>
        <UFormField label="ID d'agrégat">
          <UInput v-model="filters.aggregateId" placeholder="uuid…" class="w-full" />
        </UFormField>
        <UFormField label="Auteur (userId)">
          <UInput v-model="filters.userId" placeholder="uuid…" class="w-full" />
        </UFormField>
        <UFormField label="Limite">
          <UInputNumber v-model="filters.limit" :min="1" :max="500" class="w-full" />
        </UFormField>
      </div>

      <template #footer>
        <div class="flex justify-end gap-2">
          <UButton label="Réinitialiser" color="neutral" variant="outline" @click="resetFilters" />
          <UButton label="Filtrer" icon="i-lucide-search" @click="applyFilters" />
        </div>
      </template>
    </UCard>

    <div v-if="loading" class="py-16 text-center">
      <UIcon name="i-lucide-loader" class="animate-spin size-8 text-muted" />
    </div>

    <UCard v-else-if="entries.length === 0">
      <UEmpty icon="i-lucide-inbox" title="Aucun événement" description="Aucun événement ne correspond aux filtres." />
    </UCard>

    <ol v-else class="relative border-s border-default ml-3 space-y-6 py-2">
      <li v-for="entry in entries" :key="entry.id" class="ms-6">
        <span class="absolute -start-3 flex size-6 items-center justify-center rounded-full bg-primary/10 ring-4 ring-default">
          <UIcon :name="iconFor(entry.eventName)" class="size-3.5 text-primary" />
        </span>
        <UCard>
          <div class="flex flex-wrap items-center gap-2">
            <UBadge color="primary" variant="subtle">{{ entry.eventName }}</UBadge>
            <UBadge color="neutral" variant="outline">{{ entry.aggregateType }}</UBadge>
            <span class="ml-auto text-xs text-dimmed">{{ formatDate(entry.occurredAt) }}</span>
          </div>
          <dl class="mt-3 grid grid-cols-1 gap-2 text-xs sm:grid-cols-2">
            <div>
              <dt class="text-dimmed">Aggregate ID</dt>
              <dd class="font-mono break-all">{{ entry.aggregateId }}</dd>
            </div>
            <div>
              <dt class="text-dimmed">Auteur</dt>
              <dd class="font-mono break-all">{{ entry.userId || '— (anonyme)' }}</dd>
            </div>
          </dl>
          <details v-if="entry.payload && entry.payload !== '{}'" class="mt-3">
            <summary class="cursor-pointer text-xs text-muted hover:text-default">Payload</summary>
            <pre class="mt-2 overflow-x-auto rounded bg-elevated p-3 text-xs">{{ entry.payload }}</pre>
          </details>
        </UCard>
      </li>
    </ol>
  </UContainer>
</template>
