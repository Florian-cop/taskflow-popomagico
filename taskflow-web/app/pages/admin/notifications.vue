<script setup lang="ts">
const { channels, failed, loading, fetchChannels, setChannelFailing, fetchFailed, retryFailed } = useAdmin()
const toast = useToast()

const retrying = ref<Set<string>>(new Set())

onMounted(async () => {
  await Promise.all([fetchChannels(), fetchFailed()])
})

async function toggleChannel(name: string, failing: boolean) {
  try {
    await setChannelFailing(name, failing)
    toast.add({
      title: failing ? `Panne simulée sur ${name}` : `Canal ${name} restauré`,
      icon: failing ? 'i-lucide-alert-triangle' : 'i-lucide-check',
      color: failing ? 'warning' : 'success'
    })
  } catch (err) {
    console.error(err)
    toast.add({ title: 'Échec du toggle', icon: 'i-lucide-x', color: 'error' })
  }
}

async function handleRetry(id: string) {
  retrying.value.add(id)
  try {
    await retryFailed(id)
    toast.add({ title: 'Notification rejouée', icon: 'i-lucide-check', color: 'success' })
  } catch (err: unknown) {
    toast.add({
      title: 'Échec du retry',
      description: err instanceof Error ? err.message : undefined,
      icon: 'i-lucide-x',
      color: 'error'
    })
    // Le repo a peut-être mis à jour le retryCount → on refresh la liste
    await fetchFailed()
  } finally {
    retrying.value.delete(id)
  }
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleString('fr-FR', {
    day: '2-digit', month: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit'
  })
}
</script>

<template>
  <UContainer class="py-8">
    <div class="mb-8">
      <h1 class="text-2xl font-bold">Administration — Notifications</h1>
      <p class="mt-1 text-sm text-muted">
        Démonstration de la résilience (chantier 1 Disruption #2). Simulez une panne
        d'un canal pour vérifier que les autres canaux et le reste du système continuent
        de fonctionner. Les messages échoués sont conservés pour retraitement manuel.
      </p>
    </div>

    <UCard class="mb-8">
      <template #header>
        <div class="flex items-center gap-2">
          <UIcon name="i-lucide-radio-tower" class="size-5 text-muted" />
          <h2 class="text-base font-semibold">Canaux</h2>
        </div>
      </template>

      <ul class="divide-y divide-default">
        <li
          v-for="channel in channels"
          :key="channel.name"
          class="flex items-center justify-between py-4 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-3">
            <UIcon
              :name="channel.failing ? 'i-lucide-alert-triangle' : 'i-lucide-check-circle'"
              :class="channel.failing ? 'text-warning size-5' : 'text-success size-5'"
            />
            <div>
              <p class="text-sm font-semibold">{{ channel.name }}</p>
              <p class="text-xs text-muted">
                {{ channel.failing ? 'En panne simulée — toute tentative d\'envoi échoue' : 'Opérationnel' }}
              </p>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <span class="text-xs text-muted">Simuler une panne</span>
            <USwitch
              :model-value="channel.failing"
              @update:model-value="(v: boolean) => toggleChannel(channel.name, v)"
            />
          </div>
        </li>
        <li v-if="channels.length === 0" class="py-4 text-center text-sm text-muted">
          Aucun canal toggleable.
        </li>
      </ul>

      <template #footer>
        <p class="text-xs text-dimmed">
          Note : seul le canal <code>email</code> est wrappé dans un <code>FaultInjectingChannel</code>
          pour la démo. <code>in_app</code> reste toujours opérationnel.
        </p>
      </template>
    </UCard>

    <UCard>
      <template #header>
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2">
            <UIcon name="i-lucide-inbox" class="size-5 text-muted" />
            <h2 class="text-base font-semibold">Messages échoués</h2>
            <UBadge size="xs" variant="subtle" color="neutral">
              {{ failed.length }}
            </UBadge>
          </div>
          <UButton
            label="Rafraîchir"
            icon="i-lucide-refresh-cw"
            size="sm"
            color="neutral"
            variant="outline"
            :loading="loading"
            @click="fetchFailed"
          />
        </div>
      </template>

      <div v-if="loading" class="py-12 text-center">
        <UIcon name="i-lucide-loader" class="animate-spin size-8 text-muted" />
      </div>

      <UEmpty
        v-else-if="failed.length === 0"
        icon="i-lucide-check"
        title="Aucun message en attente"
        description="Tous les envois sont passés, ou aucun n'a échoué."
      />

      <ul v-else class="divide-y divide-default">
        <li v-for="msg in failed" :key="msg.id" class="py-4 first:pt-0 last:pb-0">
          <div class="flex items-start justify-between gap-4">
            <div class="min-w-0 flex-1">
              <div class="flex items-center gap-2">
                <UBadge color="warning" variant="subtle" size="xs">{{ msg.channel }}</UBadge>
                <UBadge color="neutral" variant="outline" size="xs">{{ msg.type }}</UBadge>
                <span class="text-xs text-dimmed">{{ formatDate(msg.occurredAt) }}</span>
              </div>
              <p class="mt-2 text-sm font-semibold truncate">{{ msg.title }}</p>
              <p class="text-xs text-muted line-clamp-2">{{ msg.body }}</p>
              <p class="mt-2 text-xs text-error font-mono break-all">{{ msg.error }}</p>
              <p v-if="msg.retryCount > 0" class="mt-1 text-[11px] text-dimmed">
                Tentatives : {{ msg.retryCount }}
              </p>
            </div>
            <UButton
              label="Rejouer"
              icon="i-lucide-refresh-ccw"
              size="sm"
              :loading="retrying.has(msg.id)"
              @click="handleRetry(msg.id)"
            />
          </div>
        </li>
      </ul>
    </UCard>
  </UContainer>
</template>
