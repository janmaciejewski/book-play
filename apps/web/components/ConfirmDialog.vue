<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-sm">
      <h2 class="text-lg font-bold text-gray-900 dark:text-white mb-2">{{ title }}</h2>
      <p class="text-gray-600 dark:text-gray-400 mb-4">{{ message }}</p>
      <div class="flex gap-2 justify-end">
        <button
          @click="$emit('cancel')"
          class="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700"
        >
          {{ cancelText }}
        </button>
        <button
          @click="$emit('confirm')"
          :disabled="loading"
          class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 disabled:opacity-50"
        >
          {{ loading ? (loadingText || `${confirmText}...`) : confirmText }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  title: string
  message: string
  confirmText?: string
  cancelText?: string
  loadingText?: string
  loading?: boolean
}>(), {
  confirmText: 'Potwierdź',
  cancelText: 'Anuluj',
  loading: false,
})

defineEmits<{
  confirm: []
  cancel: []
}>()
</script>