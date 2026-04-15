<template>
  <div class="min-h-screen bg-default flex items-center justify-center px-4 py-12">
    <div class="w-full max-w-md">
      <div class="mb-8 flex flex-col items-center gap-2">
        <UIcon name="i-lucide-check-square" class="size-10 text-primary" />
        <span class="text-2xl font-bold">TaskFlow</span>
        <p class="text-sm text-muted">Kanban projet · événementiel</p>
      </div>
      <slot />
    </div>
  </div>
</template>
