<script setup lang="ts">
const { items, unreadCount, fetchAll, markAsRead, markAllAsRead } = useNotifications()
const { isAuthenticated } = useAuth()
const open = ref(false)

let poller: ReturnType<typeof setInterval> | null = null

onMounted(() => {
  if (!isAuthenticated.value) return
  fetchAll()
  // Polling 15s pour capturer les notifs envoyées par d'autres users.
  // Alternative plus propre : s'abonner au bus WS global — ADR futur.
  poller = setInterval(() => fetchAll(), 15_000)
})

onBeforeUnmount(() => {
  if (poller) clearInterval(poller)
})

function formatDate(iso: string): string {
  return new Date(iso).toLocaleString('fr-FR', {
    day: '2-digit', month: '2-digit', hour: '2-digit', minute: '2-digit'
  })
}
</script>

<template>
  <UPopover v-model:open="open">
    <UChip :text="unreadCount" :show="unreadCount > 0" color="error" size="sm">
      <UButton
        icon="i-lucide-bell"
        color="neutral"
        variant="ghost"
        square
        aria-label="Notifications"
      />
    </UChip>

    <template #content>
      <div class="w-80 max-w-[90vw]">
        <div class="flex items-center justify-between border-b border-default px-4 py-3">
          <span class="text-sm font-semibold">Notifications</span>
          <UButton
            v-if="unreadCount > 0"
            label="Tout lire"
            size="xs"
            variant="link"
            @click="markAllAsRead()"
          />
        </div>

        <div class="max-h-96 overflow-y-auto">
          <div v-if="items.length === 0" class="px-4 py-8">
            <UEmpty icon="i-lucide-bell-off" title="Aucune notification" />
          </div>
          <ul v-else class="divide-y divide-default">
            <li
              v-for="n in items"
              :key="n.id"
              class="px-4 py-3 text-sm hover:bg-elevated cursor-pointer"
              :class="{ 'bg-primary/5': !n.readAt }"
              @click="!n.readAt && markAsRead(n.id)"
            >
              <div class="flex items-start gap-3">
                <UIcon
                  :name="n.type === 'task.assigned' ? 'i-lucide-user-check' : 'i-lucide-arrow-right-circle'"
                  class="mt-0.5 size-4 shrink-0 text-primary"
                />
                <div class="flex-1 min-w-0">
                  <p class="font-medium truncate">{{ n.title }}</p>
                  <p class="mt-0.5 text-xs text-muted line-clamp-2">{{ n.body }}</p>
                  <p class="mt-1 text-[11px] text-dimmed">{{ formatDate(n.createdAt) }}</p>
                </div>
                <span v-if="!n.readAt" class="mt-1 size-2 shrink-0 rounded-full bg-primary" />
              </div>
            </li>
          </ul>
        </div>

        <div class="border-t border-default px-4 py-2">
          <UButton
            to="/settings/notifications"
            label="Préférences"
            icon="i-lucide-settings"
            size="xs"
            variant="link"
            block
            @click="open = false"
          />
        </div>
      </div>
    </template>
  </UPopover>
</template>
