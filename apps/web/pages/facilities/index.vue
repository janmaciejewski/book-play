<template>
  <div class="px-4 md:px-6 lg:px-8 mx-auto max-w-7xl py-8">
    <div class="mb-8 flex items-start justify-between">
      <div>
        <NuxtLink v-if="authStore.user?.role === 'ADMIN' || authStore.user?.role === 'FACILITY_OWNER'" to="/" class="text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300 text-sm">&larr; Powrót do panelu</NuxtLink>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white" :class="{ 'mt-1': authStore.user?.role === 'ADMIN' || authStore.user?.role === 'FACILITY_OWNER' }">Obiekty</h1>
        <p class="mt-2 text-gray-600 dark:text-gray-400">Przeglądaj i rezerwuj obiekty sportowe w swojej okolicy</p>
      </div>
      <button
        @click="navigateTo('/facilities/report')"
        class="bg-primary-600 text-white px-4 py-2 rounded-md hover:bg-primary-700 transition-colors font-medium"
      >
        Zgłoś obiekt
      </button>
    </div>

    <!-- Search and Filters -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-4 mb-6">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Szukaj</label>
          <input v-model="search" type="text" placeholder="Szukaj obiektów..." class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white placeholder-gray-400 shadow-sm focus:border-primary-500 focus:ring-primary-500 px-3 py-2" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Typ</label>
          <select v-model="typeFilter" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500 px-3 py-2">
            <option value="">Wszystkie typy</option>
            <option v-for="opt in typeOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
          </select>
        </div>
        <div class="flex items-end">
          <button @click="searchFacilities" class="w-full bg-primary-600 text-white py-2 px-4 rounded-md hover:bg-primary-700 transition-colors">Szukaj</button>
        </div>
      </div>
    </div>

    <!-- Facilities Grid -->
    <AppLoading v-if="pending" />
    <AppError v-else-if="error" message="Nie udało się załadować obiektów. Spróbuj ponownie." />

    <div v-else-if="facilities.length === 0" class="text-center py-12">
      <p class="text-gray-500 dark:text-gray-400">Nie znaleziono obiektów. Spróbuj dostosować kryteria wyszukiwania.</p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div
        v-for="facility in facilities"
        :key="facility.id"
        class="bg-white dark:bg-gray-800 rounded-lg shadow overflow-hidden hover:shadow-xl hover:scale-[1.02] transition-all cursor-pointer border border-gray-200 dark:border-gray-700"
        @click="viewFacility(facility.id)"
      >
        <div class="h-48 bg-gray-200 relative overflow-hidden">
          <img
            v-if="facility.images && facility.images.length > 0"
            :src="facility.images[0]"
            :alt="facility.name"
            class="w-full h-full object-cover"
          />
          <img v-else :src="getTypeImage(facility.type)" :alt="facility.name" class="w-full h-full object-cover" />
        </div>
        <div class="p-4">
          <div class="flex items-center justify-between mb-2">
            <FacilityTypeBadge :type="facility.type" />
            <span class="text-lg font-bold text-gray-900 dark:text-white">{{ facility.hourlyRate }} PLN/h</span>
          </div>
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-1">{{ facility.name }}</h3>
          <p class="text-sm text-gray-600 dark:text-gray-400 mb-2">{{ facility.address }}, {{ facility.city }}</p>
          <div v-if="facility.amenities && facility.amenities.length > 0" class="flex flex-wrap gap-1">
            <span v-for="amenity in facility.amenities.slice(0, 3)" :key="amenity" class="text-xs bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 px-2 py-0.5 rounded">
              {{ amenity }}
            </span>
            <span v-if="facility.amenities.length > 3" class="text-xs text-gray-400 dark:text-gray-500">
              +{{ facility.amenities.length - 3 }} więcej
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const authStore = useAuthStore()
const { getTypeImage, typeOptions } = useFacilityTypes()
const search = ref('')
const typeFilter = ref('')
const cityFilter = ref('')

interface Facility {
  id: string
  name: string
  description?: string
  type: string
  address: string
  city: string
  hourlyRate: number
  amenities?: string[]
  images?: string[]
}

interface FacilitiesResponse { data: Facility[] }

const apiUrl = computed(() => {
  const params = new URLSearchParams()
  if (typeFilter.value) params.append('type', typeFilter.value)
  if (cityFilter.value) params.append('city', cityFilter.value)
  const queryString = params.toString()
  return `http://localhost:8080/api/v1/facilities${queryString ? '?' + queryString : ''}`
})

const { data: response, pending, error, refresh } = await useFetch<FacilitiesResponse>(apiUrl)

const facilities = computed(() => {
  const raw = response.value?.data || []
  if (!search.value.trim()) return raw
  const q = search.value.toLowerCase()
  return raw.filter(f => f.name.toLowerCase().includes(q))
})

const searchFacilities = () => refresh()

watch([typeFilter, cityFilter], () => refresh())

const viewFacility = (id: string) => navigateTo(`/facilities/${id}`)
</script>