<template>
  <div class="px-4 md:px-6 lg:px-8 mx-auto max-w-7xl py-8">
    <AppLoading v-if="pending" />
    <AppError v-else-if="error" message="Nie udało się załadować obiektu.">
      <NuxtLink to="/facilities" class="underline">Wróć do listy</NuxtLink>
    </AppError>

    <!-- Facility Details -->
    <div v-else-if="facility">
      <!-- Header -->
      <div class="mb-8">
        <NuxtLink to="/facilities" class="text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300 text-sm">
          &larr; Wróć do obiektów
        </NuxtLink>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white mt-2">{{ facility.name }}</h1>
        <p class="text-gray-600 dark:text-gray-400 mt-1">{{ facility.address }}, {{ facility.city }}</p>
        <FacilityTypeBadge :type="facility.type" class="inline-block mt-2 px-3 py-1 rounded-full text-sm" />
        <div v-if="isFacilityClosed" class="mt-2 bg-red-50 dark:bg-red-900/30 border border-red-300 dark:border-red-700 rounded-lg p-2 text-sm text-red-800 dark:text-red-300">
          ⚠ Obiekt zamknięty{{ facility.closedUntil ? ' do ' + formatClosedDate(facility.closedUntil) : '' }}
        </div>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Facility Info -->
        <div class="lg:col-span-2 space-y-6">
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 border border-gray-200 dark:border-gray-700">
            <h2 class="text-lg font-semibold mb-4 text-gray-900 dark:text-white">Informacje</h2>
            <p v-if="facility.description" class="text-gray-600 dark:text-gray-400">{{ facility.description }}</p>
            <p v-else class="text-gray-400 dark:text-gray-500 italic">Brak opisu</p>
            <div class="mt-4 grid grid-cols-2 gap-4">
              <div>
                <p class="text-sm text-gray-500 dark:text-gray-400">Cena za godzinę</p>
                <p class="text-xl font-semibold text-primary-600 dark:text-primary-400">{{ facility.hourlyRate }} PLN</p>
              </div>
              <div>
                <p class="text-sm text-gray-500 dark:text-gray-400">Godziny otwarcia</p>
                <p class="text-lg text-gray-900 dark:text-white">08:00 - 22:00</p>
              </div>
            </div>
          </div>

          <!-- Calendar -->
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 border border-gray-200 dark:border-gray-700">
            <h2 class="text-lg font-semibold mb-4 text-gray-900 dark:text-white">Wybierz termin</h2>
            <p v-if="isFacilityClosed" class="text-gray-500 dark:text-gray-400 py-4 text-center">Ten obiekt jest obecnie zamknięty. Rezerwacje nie są dostępne.</p>
            <template v-else>
              <div class="mb-4">
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Data</label>
                <input type="date" v-model="selectedDate" :min="minDate" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500" />
              </div>
              <div v-if="selectedDate">
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Dostępne godziny</label>
                <AppLoading v-if="availabilityPending" class="py-4">
                  <template #default><div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500"></div></template>
                </AppLoading>
                <div v-else-if="availableSlots.length === 0" class="text-gray-500 dark:text-gray-400 py-4">Brak dostępnych terminów na ten dzień</div>
                <div v-else class="grid grid-cols-4 sm:grid-cols-6 gap-2">
                  <button v-for="slot in availableSlots" :key="slot" @click="selectedSlot = slot" :class="['px-3 py-2 rounded-md text-sm font-medium transition-colors', selectedSlot === slot ? 'bg-primary-600 text-white' : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600']">
                    {{ slot }}
                  </button>
                </div>
              </div>
            </template>
          </div>
        </div>

        <!-- Booking Summary -->
        <div class="lg:col-span-1">
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 sticky top-4 border border-gray-200 dark:border-gray-700">
            <h2 class="text-lg font-semibold mb-4 text-gray-900 dark:text-white">Twoja rezerwacja</h2>
            <div v-if="!selectedDate || !selectedSlot" class="text-gray-500 dark:text-gray-400 text-center py-8">
              Wybierz datę i godzinę
            </div>
            <div v-else class="space-y-4">
              <div class="flex justify-between">
                <span class="text-gray-500 dark:text-gray-400">Obiekt</span>
                <span class="font-medium text-gray-900 dark:text-white">{{ facility.name }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-500 dark:text-gray-400">Data</span>
                <span class="font-medium text-gray-900 dark:text-white">{{ formatDate(selectedDate) }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-500 dark:text-gray-400">Godzina</span>
                <span class="font-medium text-gray-900 dark:text-white">{{ selectedSlot }}</span>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Czas trwania</label>
                <select v-model="duration" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500">
                  <option value="1">1 godzina</option>
                  <option value="2">2 godziny</option>
                  <option value="3">3 godziny</option>
                </select>
              </div>
              <div v-if="facility?.bookingMode !== 'INDIVIDUAL'">
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  Drużyna
                  <span v-if="facility?.bookingMode === 'TEAM'" class="text-red-500">(wymagane)</span>
                  <span v-else class="text-gray-400">(opcjonalnie)</span>
                </label>
                <select v-model="selectedTeam" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500">
                  <option v-if="facility?.bookingMode !== 'TEAM'" value="">Indywidualnie</option>
                  <option v-for="team in teams" :key="team.id" :value="team.id">{{ team.name }}</option>
                </select>
              </div>
              <div v-else class="text-sm text-gray-500 dark:text-gray-400 mb-4 py-2 px-3 bg-gray-50 dark:bg-gray-700/50 rounded-lg">
                Ten obiekt przyjmuje tylko rezerwacje indywidualne.
              </div>
              <div v-if="facility?.requiresPrepayment" class="bg-red-50 dark:bg-red-900/30 border border-red-300 dark:border-red-700 rounded-lg p-3 text-sm text-red-800 dark:text-red-300">
                <p class="font-semibold mb-1">⚠ UWAGA! Ten obiekt wymaga przedpłaty!</p>
                <p class="mb-1">Koszt przedpłaty: <strong>{{ facility.prepaymentCost }} PLN</strong></p>
                <p class="mb-1">Konto bankowe: <strong>{{ facility.bankAccount }}</strong></p>
                <p>Tytuł przelewu: <strong>{{ facility.transferTitle }}</strong></p>
                <p class="mt-2 text-xs">Przedpłata powinna zostać wykonana najpóźniej 24 godziny przed rezerwacją.</p>
              </div>
              <hr />
              <div class="flex justify-between text-lg font-semibold text-gray-900 dark:text-white">
                <span>Do zapłaty</span>
                <span class="text-primary-600 dark:text-primary-400">{{ totalPrice }} PLN</span>
              </div>
              <button @click="bookReservation" :disabled="booking" class="w-full bg-primary-600 text-white py-3 rounded-lg font-medium hover:bg-primary-700 disabled:opacity-50">
                <span v-if="booking">Rezerwowanie...</span>
                <span v-else>Zarezerwuj</span>
              </button>
            </div>
            <div v-if="bookingError" class="mt-4 text-red-600 dark:text-red-400 text-sm">{{ bookingError }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ middleware: 'auth' })

interface Facility {
  id: string; name: string; description?: string; type: string; address: string; city: string
  hourlyRate: number | string; closedUntil?: string; requiresPrepayment?: boolean
  prepaymentCost?: number | string; bankAccount?: string; transferTitle?: string
  bookingMode?: string
}
interface FacilityResponse { data: Facility }
interface AvailabilityResponse { data: string[] }
interface Team { id: string; name: string; members?: any[] }
interface TeamsResponse { data: Team[] }

const authStore = useAuthStore()
const route = useRoute()
const facilityId = route.params.id as string

const { data: facilityResponse, pending, error } = await useFetch<FacilityResponse>(`http://localhost:8080/api/v1/facilities/${facilityId}`)
const facility = computed(() => facilityResponse.value?.data)

interface MyTeamResponse { data: Team[] }
const { data: myTeamsResponse } = await useFetch<MyTeamResponse>('http://localhost:8080/api/v1/my-teams', {
  headers: { Authorization: `Bearer ${useCookie('token').value || ''}` }
})
const teams = computed(() => {
  const raw = myTeamsResponse.value?.data || []
  const currentUserId = authStore.user?.id
  return raw.filter((t: any) => {
    const member = t.members?.find((m: any) => m.user_id === currentUserId)
    return member?.role === 'CAPTAIN'
  })
})

const selectedDate = ref('')
const selectedSlot = ref('')
const duration = ref('1')
const selectedTeam = ref('')
const booking = ref(false)
const bookingError = ref('')

const minDate = computed(() => new Date().toISOString().split('T')[0])

const availabilityPending = ref(false)
const availableSlots = ref<string[]>([])

watch(selectedDate, async (newDate) => {
  if (!newDate) { availableSlots.value = []; return }
  selectedSlot.value = ''
  availabilityPending.value = true
  try {
    const response = await $fetch<AvailabilityResponse>(`http://localhost:8080/api/v1/facilities/${facilityId}/availability?date=${newDate}`)
    availableSlots.value = response.data || []
  } catch (err) {
    console.error('Failed to fetch availability:', err)
    availableSlots.value = []
  } finally { availabilityPending.value = false }
})

const isFacilityClosed = computed(() => {
  if (!facility.value?.closedUntil) return false
  return new Date(facility.value.closedUntil) > new Date()
})

function formatClosedDate(dateStr: string): string {
  return new Date(dateStr).toLocaleDateString('pl-PL', { year: 'numeric', month: 'long', day: 'numeric' })
}

const totalPrice = computed(() => {
  const rawRate = facility.value?.hourlyRate || 0
  const rate = typeof rawRate === 'string' ? parseFloat(rawRate) : rawRate
  const hours = parseInt(duration.value) || 1
  return (rate * hours).toFixed(2)
})

const formatDate = (dateStr: string): string => {
  const date = new Date(dateStr)
  return date.toLocaleDateString('pl-PL', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })
}

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
    if (selectedTeam.value) body.team_id = selectedTeam.value
    await $fetch('http://localhost:8080/api/v1/reservations', {
      method: 'POST',
      headers: { Authorization: `Bearer ${useCookie('token').value || ''}` },
      body
    })
    navigateTo('/reservations')
  } catch (err: any) {
    console.error('Failed to book:', err)
    bookingError.value = err.data?.error || 'Nie udało się dokonać rezerwacji'
  } finally { booking.value = false }
}
</script>