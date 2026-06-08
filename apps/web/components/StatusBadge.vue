<template>
  <span :class="classes">
    {{ label }}
  </span>
</template>

<script setup lang="ts">
const props = withDefaults(defineProps<{
  status: string
  size?: 'sm' | 'md'
}>(), {
  size: 'md',
})

const statusMap: Record<string, { label: string; class: string }> = {
  PENDING: {
    label: 'Oczekująca',
    class: 'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-300',
  },
  CONFIRMED: {
    label: 'Potwierdzona',
    class: 'bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300',
  },
  CANCELLED: {
    label: 'Anulowana',
    class: 'bg-red-100 dark:bg-red-900/30 text-red-800 dark:text-red-300',
  },
  CAPTAIN: {
    label: 'Kapitan',
    class: 'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-300',
  },
  MEMBER: {
    label: 'Członek',
    class: 'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300',
  },
  COMPLETED: {
    label: 'Zakończona',
    class: 'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300',
  },
  ACCEPTED: {
    label: 'Zaakceptowany',
    class: 'bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300',
  },
  REJECTED: {
    label: 'Odrzucony',
    class: 'bg-red-100 dark:bg-red-900/30 text-red-800 dark:text-red-300',
  },
}

const sizeClass = computed(() =>
  props.size === 'sm'
    ? 'px-2 py-0.5 text-xs font-medium rounded'
    : 'px-3 py-1 rounded-full text-sm font-medium'
)

const entry = computed(() => statusMap[props.status] || { label: props.status, class: 'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300' })

const label = computed(() => entry.value.label)
const classes = computed(() => [sizeClass.value, entry.value.class].join(' '))
</script>