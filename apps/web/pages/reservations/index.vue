<template>
  <div class="px-4 md:px-6 lg:px-8 mx-auto max-w-7xl py-8">
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Moje rezerwacje</h1>
      <p class="mt-2 text-gray-600 dark:text-gray-400">Przeglądaj i zarządzaj swoimi rezerwacjami</p>
    </div>

    <!-- Reservations List -->
    <div v-if="pending" class="flex justify-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500"></div>
    </div>

    <div v-else-if="error" class="bg-red-50 dark:bg-red-900/50 text-red-700 dark:text-red-300 p-4 rounded-lg border border-red-200 dark:border-red-800">
      Nie udało się załadować rezerwacji. Spróbuj ponownie.
    </div>

    <div v-else-if="reservations.length === 0" class="text-center py-12 bg-white dark:bg-gray-800 rounded-lg shadow border border-gray-200 dark:border-gray-700">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-white">Brak rezerwacji</h3>
      <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">Rozpocznij od zarezerwowania obiektu.</p>
      <div class="mt-6">
        <NuxtLink to="/facilities" class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700">
          Przeglądaj obiekty
        </NuxtLink>
      </div>
    </div>

    <div v-else class="space-y-4">
      <div
        v-for="reservation in reservations"
        :key="reservation.id"
        class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 border border-gray-200 dark:border-gray-700"
      >
        <div class="flex items-center justify-between">
          <div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ reservation.facility?.name || 'Facility' }}</h3>
            <p class="text-sm text-gray-600 dark:text-gray-400">{{ reservation.facility?.address }}, {{ reservation.facility?.city }}</p>
          </div>
          <span
            :class="[
              'px-3 py-1 rounded-full text-sm font-medium',
              reservation.status === 'CONFIRMED' ? 'bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300' :
              reservation.status === 'PENDING' ? 'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-300' :
              reservation.status === 'CANCELLED' ? 'bg-red-100 dark:bg-red-900/30 text-red-800 dark:text-red-300' :
              'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300'
            ]"
          >
            {{ reservation.status }}
          </span>
        </div>
        <div class="mt-4 grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
          <div>
            <p class="text-gray-500 dark:text-gray-400">Data</p>
            <p class="font-medium text-gray-900 dark:text-white">{{ formatDate(reservation.date) }}</p>
          </div>
          <div>
            <p class="text-gray-500 dark:text-gray-400">Godzina</p>
            <p class="font-medium text-gray-900 dark:text-white">{{ reservation.startTime }} - {{ reservation.endTime }}</p>
          </div>
          <div>
            <p class="text-gray-500 dark:text-gray-400">Cena całkowita</p>
            <p class="font-medium text-gray-900 dark:text-white">{{ reservation.totalPrice }} PLN</p>
          </div>
          <div>
            <p class="text-gray-500 dark:text-gray-400">Drużyna</p>
            <p class="font-medium text-gray-900 dark:text-white">{{ reservation.team?.name || 'Indywidualnie' }}</p>
          </div>
        </div>
        <div v-if="reservation.status === 'PENDING' || reservation.status === 'CONFIRMED'" class="mt-4 flex gap-2">
          <button 
            @click="cancelReservation(reservation.id)"
            :disabled="cancellingId === reservation.id"
            class="text-red-600 dark:text-red-400 hover:text-red-700 dark:hover:text-red-300 text-sm font-medium disabled:opacity-50"
          >
            {{ cancellingId === reservation.id ? 'Anulowanie...' : 'Anuluj rezerwację' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth'
})

interface Reservation {
  id: string
  facilityId: string
  userId: string
  teamId?: string
  date: string
  startTime: string
  endTime: string
  status: 'PENDING' | 'CONFIRMED' | 'CANCELLED' | 'COMPLETED'
  totalPrice: number
  notes?: string
  facility?: {
    name: string
    address: string
    city: string
  }
  team?: {
    name: string
  }
}

interface ReservationsResponse {
  data: Reservation[]
}

const { data: response, pending, error } = await useFetch<ReservationsResponse>('http://localhost:8080/api/v1/reservations', {
  headers: {
    Authorization: `Bearer ${useCookie('token').value || ''}`
  }
})

const reservations = computed(() => response.value?.data || [])

const cancellingId = ref<string | null>(null)

const cancelReservation = async (id: string) => {
  if (!confirm('Czy na pewno chcesz anulować tę rezerwację?')) return
  
  cancellingId.value = id
  try {
    await $fetch(`http://localhost:8080/api/v1/reservations/${id}/cancel`, {
      method: 'PUT',
      headers: {
        Authorization: `Bearer ${useCookie('token').value || ''}`
      }
    })
    // Update local state
    const index = reservations.value.findIndex(r => r.id === id)
    if (index !== -1) {
      reservations.value[index].status = 'CANCELLED'
    }
  } catch (err) {
    console.error('Failed to cancel reservation:', err)
    alert('Nie udało się anulować rezerwacji')
  } finally {
    cancellingId.value = null
  }
}

const formatDate = (date: string): string => {
  return new Date(date).toLocaleDateString('pl-PL', {
    weekday: 'short',
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}
</script>