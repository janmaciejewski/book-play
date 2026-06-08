<template>
  <div class="relative flex-shrink-0" :class="sizeClass">
    <img
      v-if="avatar"
      :src="avatar"
      class="w-full h-full rounded-full object-cover"
    />
    <div
      v-else
      :class="['rounded-full bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center', textClass]"
    >
      <span class="text-primary-600 dark:text-primary-400 font-medium">{{ initials }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    firstName?: string
    lastName?: string
    avatar?: string | null
    size?: 'sm' | 'md' | 'lg'
  }>(),
  {
    size: 'md',
  },
)

const sizeClass = computed(() => {
  switch (props.size) {
    case 'sm':
      return 'w-10 h-10'
    case 'lg':
      return 'w-20 h-20'
    default:
      return 'w-16 h-16'
  }
})

const textClass = computed(() => {
  switch (props.size) {
    case 'lg':
      return 'text-3xl'
    case 'sm':
      return 'text-lg'
    default:
      return 'text-2xl'
  }
})

const initials = computed(() => {
  const first = (props.firstName || '')[0] || ''
  const last = (props.lastName || '')[0] || ''
  return first + last
})
</script>