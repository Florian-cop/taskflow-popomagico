<script setup lang="ts">
import type { Project, User } from '~/types'

const props = defineProps<{
  project: Project
}>()

const { addMember } = useProjects()
const { searchByEmail } = useUsers()
const toast = useToast()

type UserItem = {
  id: string
  email: string
  label: string
  description: string
  avatar: { text: string }
  disabled?: boolean
}

const showAddModal = ref(false)
const searchTerm = ref('')
const selected = ref<UserItem | undefined>(undefined)
const items = ref<UserItem[]>([])
const searching = ref(false)
const submitting = ref(false)

const memberIds = computed(() => new Set(props.project.members.map(m => m.userId)))

function toItem(u: User): UserItem {
  const name = `${u.firstName ?? ''} ${u.lastName ?? ''}`.trim()
  const avatarSource = u.firstName || u.email || '?'
  return {
    id: u.id,
    email: u.email,
    label: u.email,
    description: name || '—',
    avatar: { text: (avatarSource.charAt(0) || '?').toUpperCase() }
  }
}

let debounceTimer: ReturnType<typeof setTimeout> | null = null

watch(searchTerm, (query) => {
  if (debounceTimer) clearTimeout(debounceTimer)
  const q = (query ?? '').trim()

  if (q.length < 2) {
    items.value = []
    searching.value = false
    return
  }

  searching.value = true
  debounceTimer = setTimeout(async () => {
    try {
      const users = await searchByEmail(q)
      items.value = users.map((u) => {
        const item = toItem(u)
        if (memberIds.value.has(u.id)) {
          return { ...item, description: `${item.description} · déjà membre`, disabled: true }
        }
        return item
      })
    } catch (err) {
      console.error('[TaskFlow] user search failed:', err)
      items.value = []
    } finally {
      searching.value = false
    }
  }, 250)
})

function resetForm() {
  searchTerm.value = ''
  selected.value = undefined
  items.value = []
}

function initials(userId: string, role: string): string {
  if (role === 'owner') return 'O'
  return (userId.charAt(0) || '?').toUpperCase()
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString('fr-FR', {
    day: '2-digit', month: '2-digit', year: 'numeric'
  })
}

async function handleAdd() {
  if (!selected.value) return

  if (memberIds.value.has(selected.value.id)) {
    toast.add({
      title: 'Déjà membre',
      description: `${selected.value.email} fait déjà partie du projet.`,
      icon: 'i-lucide-info',
      color: 'warning'
    })
    return
  }

  submitting.value = true
  try {
    await addMember(props.project.id, selected.value.id)
    const label = selected.value.description !== '—' ? selected.value.description : selected.value.email
    toast.add({
      title: 'Membre ajouté',
      description: `${label} a rejoint le projet.`,
      icon: 'i-lucide-check',
      color: 'success'
    })
    resetForm()
    showAddModal.value = false
  } catch (err: unknown) {
    const status = (err as { statusCode?: number, response?: { status?: number } })?.statusCode
      ?? (err as { response?: { status?: number } })?.response?.status

    if (status === 409) {
      toast.add({ title: 'Déjà membre', icon: 'i-lucide-info', color: 'warning' })
    } else {
      toast.add({
        title: 'Impossible d\'ajouter le membre',
        description: err instanceof Error ? err.message : undefined,
        icon: 'i-lucide-x',
        color: 'error'
      })
    }
  } finally {
    submitting.value = false
  }
}

watch(showAddModal, (open) => {
  if (!open) resetForm()
})
</script>

<template>
  <UCard class="mb-6">
    <template #header>
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2">
          <UIcon name="i-lucide-users" class="size-5 text-muted" />
          <h2 class="text-base font-semibold">Membres</h2>
          <UBadge size="xs" variant="subtle" color="neutral">
            {{ project.members.length }}
          </UBadge>
        </div>
        <UButton
          label="Ajouter un membre"
          icon="i-lucide-user-plus"
          size="sm"
          @click="showAddModal = true"
        />
      </div>
    </template>

    <ul class="divide-y divide-default">
      <li
        v-for="member in project.members"
        :key="member.userId"
        class="flex items-center justify-between py-3 first:pt-0 last:pb-0"
      >
        <div class="flex items-center gap-3 min-w-0">
          <UAvatar
            :text="initials(member.userId, member.role)"
            size="sm"
            :color="member.role === 'owner' ? 'primary' : 'neutral'"
          />
          <div class="min-w-0">
            <p class="text-sm font-mono truncate">{{ member.userId }}</p>
            <p class="text-xs text-dimmed">
              Membre depuis le {{ formatDate(member.joinedAt) }}
            </p>
          </div>
        </div>
        <UBadge
          :color="member.role === 'owner' ? 'primary' : 'neutral'"
          variant="subtle"
          size="xs"
        >
          {{ member.role === 'owner' ? 'Owner' : 'Member' }}
        </UBadge>
      </li>
    </ul>

    <UModal v-model:open="showAddModal">
      <template #header>
        <h3 class="text-lg font-semibold">Ajouter un membre</h3>
      </template>

      <template #body>
        <div class="space-y-4">
          <UFormField
            label="Rechercher un utilisateur"
            hint="Saisissez au moins 2 caractères de l'email."
            required
          >
            <UInputMenu
              v-model="selected"
              v-model:search-term="searchTerm"
              :items="items"
              :loading="searching"
              by="id"
              placeholder="email@exemple.com"
              icon="i-lucide-search"
              :ignore-filter="true"
              class="w-full"
            >
              <template #empty>
                <span v-if="searchTerm.trim().length < 2" class="text-xs text-muted">
                  Tapez au moins 2 caractères…
                </span>
                <span v-else-if="searching" class="text-xs text-muted">
                  Recherche en cours…
                </span>
                <span v-else class="text-xs text-muted">
                  Aucun utilisateur trouvé.
                </span>
              </template>
            </UInputMenu>
          </UFormField>

          <div v-if="selected" class="flex items-start gap-3 rounded-md border border-default p-3">
            <UAvatar :text="selected.avatar.text" size="sm" color="primary" />
            <div class="min-w-0">
              <p class="text-sm font-semibold">
                {{ selected.description !== '—' ? selected.description : selected.email }}
              </p>
              <p class="text-xs text-muted truncate">{{ selected.email }}</p>
              <p class="text-[11px] text-dimmed font-mono mt-1 truncate">{{ selected.id }}</p>
            </div>
          </div>
        </div>
      </template>

      <template #footer>
        <div class="flex justify-end gap-3">
          <UButton
            label="Annuler"
            color="neutral"
            variant="outline"
            :disabled="submitting"
            @click="showAddModal = false"
          />
          <UButton
            label="Ajouter"
            icon="i-lucide-user-plus"
            :loading="submitting"
            :disabled="!selected"
            @click="handleAdd"
          />
        </div>
      </template>
    </UModal>
  </UCard>
</template>
