<template>
  <div>
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900">Obiekty</h1>
      <p class="mt-2 text-gray-600">Przeglądaj i rezerwuj obiekty sportowe w swojej okolicy</p>
    </div>

    <!-- Search and Filters -->
    <div class="bg-white rounded-lg shadow p-4 mb-6">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Szukaj</label>
          <input
            v-model="search"
            type="text"
            placeholder="Szukaj obiektów..."
            class="w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Typ</label>
          <select
            v-model="typeFilter"
            class="w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
          >
            <option value="">Wszystkie typy</option>
            <option value="FOOTBALL">Piłka nożna</option>
            <option value="BASKETBALL">Koszykówka</option>
            <option value="TENNIS">Tenis</option>
            <option value="VOLLEYBALL">Siatkówka</option>
            <option value="SWIMMING">Pływanie</option>
            <option value="OTHER">Inne</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Miasto</label>
          <input
            v-model="cityFilter"
            type="text"
            placeholder="Wprowadź miasto..."
            class="w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
          />
        </div>
        <div class="flex items-end">
          <button
            @click="searchFacilities"
            class="w-full bg-indigo-600 text-white py-2 px-4 rounded-md hover:bg-indigo-700 transition-colors"
          >
            Szukaj
          </button>
        </div>
      </div>
    </div>

    <!-- Facilities Grid -->
    <div v-if="pending" class="flex justify-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
    </div>

    <div v-else-if="error" class="bg-red-50 text-red-700 p-4 rounded-lg">
      Nie udało się załadować obiektów. Spróbuj ponownie.
    </div>

    <div v-else-if="facilities.length === 0" class="text-center py-12">
      <p class="text-gray-500">Nie znaleziono obiektów. Spróbuj dostosować kryteria wyszukiwania.</p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div
        v-for="facility in facilities"
        :key="facility.id"
        class="bg-white rounded-lg shadow overflow-hidden hover:shadow-lg transition-shadow cursor-pointer"
        @click="viewFacility(facility.id)"
      >
        <div class="h-48 bg-gray-200 flex items-center justify-center">
          <span class="text-6xl">{{ getTypeIcon(facility.type) }}</span>
        </div>
        <div class="p-4">
          <div class="flex items-center justify-between mb-2">
            <span class="text-xs font-medium text-indigo-600 bg-indigo-100 px-2 py-1 rounded">
              {{ facility.type }}
            </span>
            <span class="text-lg font-bold text-gray-900">{{ facility.hourlyRate }} PLN/h</span>
          </div>
          <h3 class="text-lg font-semibold text-gray-900 mb-1">{{ facility.name }}</h3>
          <p class="text-sm text-gray-600 mb-2">{{ facility.address }}, {{ facility.city }}</p>
          <div v-if="facility.amenities && facility.amenities.length > 0" class="flex flex-wrap gap-1">
            <span
              v-for="amenity in facility.amenities.slice(0, 3)"
              :key="amenity"
              class="text-xs bg-gray-100 text-gray-600 px-2 py-0.5 rounded"
            >
              {{ amenity }}
            </span>
            <span v-if="facility.amenities.length > 3" class="text-xs text-gray-400">
              +{{ facility.amenities.length - 3 }} więcej
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// No auth required for viewing facilities

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

interface FacilitiesResponse {
  data: Facility[]
}

// Build URL with query parameters
const apiUrl = computed(() => {
  const params = new URLSearchParams()
  if (typeFilter.value) params.append('type', typeFilter.value)
  if (cityFilter.value) params.append('city', cityFilter.value)
  const queryString = params.toString()
  return `http://localhost:8080/api/v1/facilities${queryString ? '?' + queryString : ''}`
})

const { data: response, pending, error, refresh } = await useFetch<FacilitiesResponse>(apiUrl)

const facilities = computed(() => response.value?.data || [])

const getTypeIcon = (type: string): string => {
  const icons: Record<string, string> = {
    FOOTBALL: '⚽',
    BASKETBALL: '🏀',
    TENNIS: '🎾',
    VOLLEYBALL: '🏐',
    SWIMMING: '🏊',
    OTHER: '🏟️'
  }
  return icons[type] || '🏟️'
}

const searchFacilities = () => {
  refresh()
}

// Watch for filter changes
watch([typeFilter, cityFilter], () => {
  refresh()
})

const viewFacility = (id: string) => {
  navigateTo(`/facilities/${id}`)
}
</script>
