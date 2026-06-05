<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>

<script setup lang="ts">
useHead({
  titleTemplate: (title) => title ? `${title} | BookPlay` : 'BookPlay',
})

// Inicjalizuje stan autoryzacji z ciasteczek
const authStore = useAuthStore()
onMounted(() => {
  authStore.initialize()
})
</script>
