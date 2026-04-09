<template>
  <NuxtLayout>
    <NuxtPage />
  </NuxtLayout>
</template>

<script setup lang="ts">
useHead({
  titleTemplate: (title) => title ? `${title} | BookPlay` : 'BookPlay',
})

// Initialize auth state from cookies
const authStore = useAuthStore()
onMounted(() => {
  authStore.initialize()
})
</script>
