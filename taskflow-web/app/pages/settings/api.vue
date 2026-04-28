<script setup lang="ts">
import type { ApiVersion } from '~/composables/useApiVersion'

const { version, setVersion } = useApiVersion()
const config = useRuntimeConfig()
const { token } = useAuth()
const toast = useToast()

const baseV1 = (config.public.apiBase as string)
const baseV2 = baseV1.replace(/\/v1(\/?$)/, '/v2$1')

const sample = ref<string>('')
const sampleVersion = ref<ApiVersion | null>(null)
const loadingSample = ref(false)

async function fetchSample(target: ApiVersion) {
  loadingSample.value = true
  sampleVersion.value = target
  try {
    const url = `${target === 'v2' ? baseV2 : baseV1}/projects`
    const raw = await $fetch<unknown>(url, {
      headers: token.value ? { Authorization: `Bearer ${token.value}` } : {}
    })
    sample.value = JSON.stringify(raw, null, 2)
  } catch (err) {
    sample.value = `Erreur : ${(err as Error).message}`
  } finally {
    loadingSample.value = false
  }
}

function selectVersion(next: ApiVersion) {
  setVersion(next)
  toast.add({
    title: `API ${next} activée`,
    description: next === 'v2'
      ? 'Les appels projects/tasks utilisent maintenant /api/v2 avec enveloppe {data, meta}.'
      : 'Retour à /api/v1, format brut.',
    icon: 'i-lucide-check',
    color: 'success'
  })
}
</script>

<template>
  <UContainer class="py-8">
    <div class="mb-8 max-w-2xl">
      <h1 class="text-2xl font-bold">Version d'API</h1>
      <p class="mt-1 text-sm text-muted">
        Bascule la couche présentation entre <code>/api/v1</code> (réponse brute) et
        <code>/api/v2</code> (enveloppe <code>{ data, meta }</code>). Les services métier
        et le domaine sont rigoureusement identiques — seuls les adaptateurs entrants changent.
      </p>
    </div>

    <UCard class="max-w-2xl">
      <div class="flex flex-col gap-4">
        <div class="flex items-center justify-between gap-4">
          <div>
            <p class="text-sm font-semibold">Version active</p>
            <p class="text-xs text-muted">
              Affecte uniquement <code>/projects</code> et <code>/tasks</code> (seules routes implémentées en v2).
            </p>
          </div>
          <UBadge :color="version === 'v2' ? 'primary' : 'neutral'" variant="soft" size="lg">
            {{ version }}
          </UBadge>
        </div>

        <URadioGroup
          :model-value="version"
          :items="[
            { value: 'v1', label: 'v1 — réponse brute', description: 'Array ou objet JSON sans enveloppe.' },
            { value: 'v2', label: 'v2 — enveloppée', description: '{ data, meta: { apiVersion, generatedAt, count } }.' }
          ]"
          @update:model-value="selectVersion($event as ApiVersion)"
        />
      </div>
    </UCard>

    <UCard class="mt-6 max-w-2xl">
      <template #header>
        <div class="flex items-center justify-between gap-2">
          <div>
            <p class="text-sm font-semibold">Inspecter le format</p>
            <p class="text-xs text-muted">
              Appelle <code>GET /projects</code> sur la version choisie et affiche la réponse brute.
            </p>
          </div>
          <div class="flex gap-2">
            <UButton
              size="xs"
              color="neutral"
              variant="soft"
              :loading="loadingSample && sampleVersion === 'v1'"
              @click="fetchSample('v1')"
            >
              Voir v1
            </UButton>
            <UButton
              size="xs"
              :loading="loadingSample && sampleVersion === 'v2'"
              @click="fetchSample('v2')"
            >
              Voir v2
            </UButton>
          </div>
        </div>
      </template>

      <pre
        v-if="sample"
        class="max-h-96 overflow-auto rounded bg-elevated p-3 text-xs"
      >{{ sample }}</pre>
      <p v-else class="text-xs text-muted">
        Aucune réponse capturée. Clique sur « Voir v1 » ou « Voir v2 » pour comparer.
      </p>
    </UCard>

    <p class="mt-4 max-w-2xl text-xs text-dimmed">
      Détail technique : <code>useApi()</code> redirige automatiquement vers
      <code>/api/v2</code> et déballe <code>{ data, meta }</code> côté client, ce qui permet
      à toute l'app de continuer à fonctionner sans modification des composables. Les routes
      auth, audit et notifications restent en v1 (non portées).
    </p>
  </UContainer>
</template>
