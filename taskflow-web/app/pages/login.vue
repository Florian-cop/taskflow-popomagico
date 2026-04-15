<script setup lang="ts">
definePageMeta({ layout: 'auth' })

const { login, register } = useAuth()
const toast = useToast()
const route = useRoute()

type Tab = 'login' | 'register'
const tab = ref<Tab>('login')

const loginForm = reactive({ email: '', password: '' })
const registerForm = reactive({ email: '', password: '', firstName: '', lastName: '' })
const loading = ref(false)

const tabs = [
  { label: 'Connexion', value: 'login' as const },
  { label: 'Inscription', value: 'register' as const }
]

async function handleLogin() {
  if (!loginForm.email || !loginForm.password) return
  loading.value = true
  try {
    await login({ email: loginForm.email, password: loginForm.password })
    toast.add({ title: 'Connecté', icon: 'i-lucide-check', color: 'success' })
    const redirect = (route.query.redirect as string) || '/'
    await navigateTo(redirect)
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : 'Identifiants invalides'
    toast.add({ title: 'Échec de connexion', description: msg, icon: 'i-lucide-x', color: 'error' })
  } finally {
    loading.value = false
  }
}

async function handleRegister() {
  if (!registerForm.email || !registerForm.password) return
  loading.value = true
  try {
    await register({ ...registerForm })
    toast.add({ title: 'Compte créé', icon: 'i-lucide-check', color: 'success' })
    await navigateTo('/')
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : 'Impossible de créer le compte'
    toast.add({ title: 'Échec d\'inscription', description: msg, icon: 'i-lucide-x', color: 'error' })
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <UCard>
    <template #header>
      <UTabs v-model="tab" :items="tabs" :content="false" variant="link" />
    </template>

    <form v-if="tab === 'login'" class="space-y-4" @submit.prevent="handleLogin">
      <UFormField label="Email" required>
        <UInput
          v-model="loginForm.email"
          type="email"
          placeholder="vous@exemple.com"
          autocomplete="email"
          autofocus
          class="w-full"
        />
      </UFormField>
      <UFormField label="Mot de passe" required>
        <UInput
          v-model="loginForm.password"
          type="password"
          placeholder="••••••••"
          autocomplete="current-password"
          class="w-full"
        />
      </UFormField>
      <UButton
        type="submit"
        label="Se connecter"
        icon="i-lucide-log-in"
        :loading="loading"
        :disabled="!loginForm.email || !loginForm.password"
        block
      />
    </form>

    <form v-else class="space-y-4" @submit.prevent="handleRegister">
      <div class="grid grid-cols-2 gap-3">
        <UFormField label="Prénom">
          <UInput v-model="registerForm.firstName" placeholder="Ada" class="w-full" />
        </UFormField>
        <UFormField label="Nom">
          <UInput v-model="registerForm.lastName" placeholder="Lovelace" class="w-full" />
        </UFormField>
      </div>
      <UFormField label="Email" required>
        <UInput
          v-model="registerForm.email"
          type="email"
          placeholder="vous@exemple.com"
          autocomplete="email"
          class="w-full"
        />
      </UFormField>
      <UFormField label="Mot de passe" required hint="Minimum 8 caractères">
        <UInput
          v-model="registerForm.password"
          type="password"
          placeholder="••••••••"
          autocomplete="new-password"
          class="w-full"
        />
      </UFormField>
      <UButton
        type="submit"
        label="Créer le compte"
        icon="i-lucide-user-plus"
        :loading="loading"
        :disabled="!registerForm.email || registerForm.password.length < 8"
        block
      />
    </form>
  </UCard>
</template>
