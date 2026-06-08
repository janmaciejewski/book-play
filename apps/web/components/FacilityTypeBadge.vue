<template>
  <span :class="['text-xs font-medium px-2 py-1 rounded border', badgeClass]">
    {{ label }}
  </span>
</template>

<script setup lang="ts">
const props = defineProps<{ type: string }>()
const { getTypeName, getTypeBadgeClass } = useFacilityTypes()
const label = computed(() => getTypeName(props.type))
const badgeClass = computed(() => getTypeBadgeClass(props.type))
</script>