<script setup lang="ts">
const route = useRoute()
const { projects } = useProjects()
const { user, fullName, logout } = useAuth()
const toast = useToast()

const avatarText = computed(() => {
  if (!user.value) return '?'
  const source = user.value.firstName || user.value.email || '?'
  return (source.charAt(0) || '?').toUpperCase()
})

const breadcrumbs = computed(() => {
  const items: Array<{ label: string, to: string, icon: string }> = [
    { label: 'Projets', to: '/', icon: 'i-lucide-folder' }
  ]
  if (route.params.id) {
    const project = projects.value.find(p => p.id === route.params.id)
    if (project) {
      items.push({ label: project.name, to: `/projects/${project.id}`, icon: 'i-lucide-kanban' })
    }
  }
  if (route.path === '/audit') {
    items.push({ label: 'Audit', to: '/audit', icon: 'i-lucide-scroll-text' })
  }
  if (route.path.startsWith('/settings')) {
    items.push({ label: 'Paramètres', to: '/settings/notifications', icon: 'i-lucide-settings' })
  }
  return items
})

const userMenuItems = computed(() => [
  [
    {
      label: fullName.value || 'Utilisateur',
      slot: 'account',
      disabled: true
    }
  ],
  [
    { label: 'Préférences de notif.', icon: 'i-lucide-bell', to: '/settings/notifications' },
    { label: 'Journal d\'audit', icon: 'i-lucide-scroll-text', to: '/audit' }
  ],
  [
    { label: 'Se déconnecter', icon: 'i-lucide-log-out', onSelect: handleLogout }
  ]
])

async function handleLogout() {
  logout()
  toast.add({ title: 'Déconnecté', icon: 'i-lucide-check', color: 'info' })
  await navigateTo('/login')
}
</script>

<template>
  <UHeader>
    <template #left>
      <NuxtLink to="/" class="flex items-center gap-2">
        <UIcon name="i-lucide-check-square" class="size-6 text-primary" />
        <span class="text-lg font-bold">TaskFlow</span>
      </NuxtLink>
    </template>

    <UBreadcrumb :items="breadcrumbs" class="hidden md:flex" />

    <template #right>
      <UButton to="/audit" icon="i-lucide-scroll-text" color="neutral" variant="ghost" square aria-label="Audit" />
      <NotificationBell />
      <UColorModeButton />

      <UDropdownMenu v-if="user" :items="userMenuItems">
        <UButton color="neutral" variant="ghost" trailing-icon="i-lucide-chevron-down">
          <UAvatar
            :alt="fullName"
            size="xs"
            :text="avatarText"
          />
          <span class="hidden sm:inline max-w-[10rem] truncate">
            {{ fullName || user.email }}
          </span>
        </UButton>

        <template #account>
          <div class="text-left">
            <p class="text-xs text-dimmed">Connecté en tant que</p>
            <p class="truncate font-medium">{{ user.email }}</p>
          </div>
        </template>
      </UDropdownMenu>
    </template>
  </UHeader>
</template>
