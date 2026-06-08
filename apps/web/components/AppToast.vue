<template>
  <div v-if="toast.message" :class="toastClass">
    {{ toast.message }}
  </div>
</template>

<script setup lang="ts">
const { toast, toastClass } = useToast()
</script>