<template>
  <div>
    <!-- Loading -->
    <div v-if="pending" class="flex justify-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="bg-red-50 text-red-700 p-4 rounded-lg">
      Nie udało się załadować obiektu.
      <NuxtLink to="/facilities" class="underline">Wróć do listy</NuxtLink>
    </div>

    <!-- Facility Details -->
    <div v-else-if="facility">
      <!-- Header -->
      <div class="mb-8">
        <NuxtLink to="/facilities" class="text-indigo-600 hover:text-indigo-700 text-sm">
          &larr; Wróć do obiektów
        </NuxtLink>
        <h1 class="text-3xl font-bold text-gray-900 mt-2">{{ facility.name }}</h1>
        <p class="text-gray-600 mt-1">{{ facility.address }}, {{ facility.city }}</p>
        <span class="inline-block mt-2 px-3 py-1 bg-indigo-100 text-indigo-800 rounded-full text-sm">
          {{ getTypeLabel(facility.type) }}
        </span>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Facility Info -->
        <div class="lg:col-span-2 space-y-6">
          <div class="bg-white rounded-lg shadow p-6">
            <h2 class="text-lg font-semibold mb-4">Informacje</h2>
            <p v-if="facility.description" class="text-gray-600">{{ facility.description }}</p>
            <p v-else class="text-gray-400 italic">Brak opisu</p>
            
            <div class="mt-4 grid grid-cols-2 gap-4">
              <div>
                <p class="text-sm text-gray-500">Cena za godzinę</p>
                <p class="text-xl font-semibold text-indigo-600">{{ facility.hourlyRate }} PLN</p>
              </div>
              <div>
                <p class="text-sm text-gray-500">Godziny otwarcia</p>
                <p class="text-lg">08:00 - 22:00</p>
              </div>
            </div>
          </div>

          <!-- Calendar -->
          <div class="bg-white rounded-lg shadow p-6">
            <h2 class="text-lg font-semibold mb-4">Wybierz termin</h2>
            
            <!-- Date Selection -->
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-2">Data</label>
              <input
                type="date"
                v-model="selectedDate"
                :min="minDate"
                class="w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
              />
            </div>

            <!-- Time Slots -->
            <div v-if="selectedDate">
              <label class="block text-sm font-medium text-gray-700 mb-2">Dostępne godziny</label>
              <div v-if="availabilityPending" class="flex justify-center py-4">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
              </div>
              <div v-else-if="availableSlots.length === 0" class="text-gray-500 py-4">
                Brak dostępnych terminów na ten dzień
              </div>
              <div v-else class="grid grid-cols-4 sm:grid-cols-6 gap-2">
                <button
                  v-for="slot in availableSlots"
                  :key="slot"
                  @click="selectedSlot = slot"
                  :class="[
                    'px-3 py-2 rounded-md text-sm font-medium transition-colors',
                    selectedSlot === slot
                      ? 'bg-indigo-600 text-white'
                      : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                  ]"
                >
                  {{ slot }}
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Booking Summary -->
        <div class="lg:col-span-1">
          <div class="bg-white rounded-lg shadow p-6 sticky top-4">
            <h2 class="text-lg font-semibold mb-4">Twoja rezerwacja</h2>
            
            <div v-if="!selectedDate || !selectedSlot" class="text-gray-500 text-center py-8">
              Wybierz datę i godzinę
            </div>
            
            <div v-else class="space-y-4">
              <div class="flex justify-between">
                <span class="text-gray-500">Obiekt</span>
                <span class="font-medium">{{ facility.name }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-500">Data</span>
                <span class="font-medium">{{ formatDate(selectedDate) }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-500">Godzina</span>
                <span class="font-medium">{{ selectedSlot }}</span>
              </div>
              
              <!-- Duration Selection -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Czas trwania</label>
                <select
                  v-model="duration"
                  class="w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
                >
                  <option value="1">1 godzina</option>
                  <option value="2">2 godziny</option>
                  <option value="3">3 godziny</option>
                </select>
              </div>

              <!-- Team Selection (optional) -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Drużyna (opcjonalnie)</label>
                <select
                  v-model="selectedTeam"
                  class="w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
                >
                  <option value="">Indywidualnie</option>
                  <option v-for="team in teams" :key="team.id" :value="team.id">
                    {{ team.name }}
                  </option>
                </select>
              </div>

              <hr />

              <div class="flex justify-between text-lg font-semibold">
                <span>Do zapłaty</span>
                <span class="text-indigo-600">{{ totalPrice }} PLN</span>
              </div>

              <button
                @click="bookReservation"
                :disabled="booking"
                class="w-full bg-indigo-600 text-white py-3 rounded-lg font-medium hover:bg-indigo-700 disabled:opacity-50"
              >
                <span v-if="booking">Rezerwowanie...</span>
                <span v-else>Zarezerwuj</span>
              </button>
            </div>

            <div v-if="bookingError" class="mt-4 text-red-600 text-sm">
              {{ bookingError }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth'
})

interface Facility {
  id: string
  name: string
  description?: string
  type: string
  address: string
  city: string
  hourlyRate: number | string // API returns camelCase
}

interface FacilityResponse {
  data: Facility
}

interface AvailabilityResponse {
  data: string[]
}

interface Team {
  id: string
  name: string
}

interface TeamsResponse {
  data: Team[]
}

const route = useRoute()
const facilityId = route.params.id as string

// Fetch facility
const { data: facilityResponse, pending, error } = await useFetch<FacilityResponse>(`http://localhost:8080/api/v1/facilities/${facilityId}`)
const facility = computed(() => facilityResponse.value?.data)

// Fetch teams for selection
const { data: teamsResponse } = await useFetch<TeamsResponse>('http://localhost:8080/api/v1/teams', {
  headers: {
    Authorization: `Bearer ${useCookie('token').value || ''}`
  }
})
const teams = computed(() => teamsResponse.value?.data || [])

// Booking state
const selectedDate = ref('')
const selectedSlot = ref('')
const duration = ref('1')
const selectedTeam = ref('')
const booking = ref(false)
const bookingError = ref('')

// Min date is today
const minDate = computed(() => {
  const today = new Date()
  return today.toISOString().split('T')[0]
})

// Fetch availability when date changes
const availabilityPending = ref(false)
const availableSlots = ref<string[]>([])

watch(selectedDate, async (newDate) => {
  if (!newDate) {
    availableSlots.value = []
    return
  }
  
  selectedSlot.value = ''
  availabilityPending.value = true
  
  try {
    const response = await $fetch<AvailabilityResponse>(
      `http://localhost:8080/api/v1/facilities/${facilityId}/availability?date=${newDate}`
    )
    availableSlots.value = response.data || []
  } catch (err) {
    console.error('Failed to fetch availability:', err)
    availableSlots.value = []
  } finally {
    availabilityPending.value = false
  }
})

// Calculate total price
const totalPrice = computed(() => {
  const rawRate = facility.value?.hourlyRate || 0
  const rate = typeof rawRate === 'string' ? parseFloat(rawRate) : rawRate
  const hours = parseInt(duration.value) || 1
  return (rate * hours).toFixed(2)
})

// Format date for display
const formatDate = (dateStr: string): string => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('pl-PL', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

// Get type label
const getTypeLabel = (type: string): string => {
  const labels: Record<string, string> = {
    'FOOTBALL': 'Piłka nożna',
    'BASKETBALL': 'Koszykówka',
    'VOLLEYBALL': 'Siatkówka',
    'TENNIS': 'Tenis',
    'PADEL': 'Padel',
    'SWIMMING': 'Basen',
  }
  return labels[type] || type
}

// Book reservation
const bookReservation = async () => {
  if (!selectedDate.value || !selectedSlot.value) return
  
  booking.value = true
  bookingError.value = ''
  
  try {
    const startHour = parseInt(selectedSlot.value.split(':')[0])
    const endHour = startHour + parseInt(duration.value)
    const endTime = `${endHour.toString().padStart(2, '0')}:00`
    
    const body: Record<string, unknown> = {
      facility_id: facilityId,
      date: selectedDate.value,
      start_time: selectedSlot.value,
      end_time: endTime,
      total_price: totalPrice.value
    }
    
    if (selectedTeam.value) {
      body.team_id = selectedTeam.value
    }
    
    await $fetch('http://localhost:8080/api/v1/reservations', {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${useCookie('token').value || ''}`
      },
      body
    })
    
    // Redirect to reservations page
    navigateTo('/reservations')
  } catch (err: any) {
    console.error('Failed to book:', err)
    bookingError.value = err.data?.error || 'Nie udało się dokonać rezerwacji'
  } finally {
    booking.value = false
  }
}
</script>