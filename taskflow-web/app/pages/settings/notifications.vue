<script setup lang="ts">
import { NOTIFICATION_CHANNELS } from '~/types'

const { preferences, fetchPreferences, updatePreferences } = useNotifications()
const toast = useToast()

const local = reactive<Record<string, boolean>>({})
const saving = ref(false)

onMounted(async () => {
  await fetchPreferences()
  syncLocal()
})

watch(preferences, syncLocal, { deep: true })

function syncLocal() {
  for (const channel of NOTIFICATION_CHANNELS) {
    const current = preferences.value.enabled[channel.key]
    local[channel.key] = current === undefined ? true : current
  }
}

async function save() {
  saving.value = true
  try {
    await updatePreferences({ ...local })
    toast.add({ title: 'Préférences sauvegardées', icon: 'i-lucide-check', color: 'success' })
  } catch {
    toast.add({ title: 'Impossible de sauvegarder', icon: 'i-lucide-x', color: 'error' })
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <UContainer class="py-8">
    <div class="mb-8 max-w-2xl">
      <h1 class="text-2xl font-bold">Préférences de notifications</h1>
      <p class="mt-1 text-sm text-muted">
        Choisissez les canaux sur lesquels vous souhaitez recevoir les notifications de l'application.
      </p>
    </div>

    <UCard class="max-w-2xl">
      <ul class="divide-y divide-default">
        <li
          v-for="channel in NOTIFICATION_CHANNELS"
          :key="channel.key"
          class="flex items-start justify-between gap-4 py-4 first:pt-0 last:pb-0"
        >
          <div class="flex items-start gap-3">
            <UIcon :name="channel.icon" class="mt-0.5 size-5 text-muted" />
            <div>
              <p class="text-sm font-semibold">{{ channel.label }}</p>
              <p class="text-xs text-muted">{{ channel.description }}</p>
            </div>
          </div>
          <USwitch v-model="local[channel.key]" />
        </li>
      </ul>

      <template #footer>
        <div class="flex justify-end">
          <UButton label="Enregistrer" icon="i-lucide-save" :loading="saving" @click="save" />
        </div>
      </template>
    </UCard>

    <p class="mt-4 max-w-2xl text-xs text-dimmed">
      Le canal <strong>email</strong> est simulé par un log serveur dans cette livraison
      (raccourci assumé — voir ADR-007). Les notifications <strong>in-app</strong> apparaissent
      dans la cloche en haut à droite.
    </p>
  </UContainer>
</template>
