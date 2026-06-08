<template>
  <div class="px-4 md:px-6 lg:px-8 mx-auto max-w-7xl py-8">
    <div class="mb-8">
      <NuxtLink to="/" class="text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300 text-sm">&larr; Powrót do panelu</NuxtLink>
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white mt-1">Rezerwacje na moje obiekty</h1>
      <p class="mt-2 text-gray-600 dark:text-gray-400">Zarządzaj rezerwacjami dla swoich obiektów sportowych</p>
    </div>

    <AppLoading v-if="pending" />
    <AppError v-else-if="error" message="Nie udało się załadować rezerwacji. Spróbuj ponownie." />
    <AppEmpty v-else-if="reservations.length === 0" title="Brak rezerwacji" description="Nie masz jeszcze żadnych rezerwacji na swoje obiekty." icon="calendar" />

    <div v-else class="bg-white dark:bg-gray-800 rounded-lg shadow border border-gray-200 dark:border-gray-700 overflow-x-auto">
      <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
        <thead class="bg-gray-50 dark:bg-gray-700">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Obiekt</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Użytkownik</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Data</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Godzina</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Status</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Akcje</th>
          </tr>
        </thead>
        <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
          <tr v-for="r in reservations" :key="r.id">
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="font-medium text-gray-900 dark:text-white">{{ r.facility?.name }}</div>
              <div class="text-sm text-gray-500 dark:text-gray-400">{{ r.facility?.address }}, {{ r.facility?.city }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="font-medium text-gray-900 dark:text-white">{{ r.user?.first_name }} {{ r.user?.last_name }}</div>
              <div class="text-sm text-gray-500 dark:text-gray-400">{{ r.user?.email }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">{{ formatDate(r.date) }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-white">{{ r.startTime }} – {{ r.endTime }}</td>
            <td class="px-6 py-4 whitespace-nowrap">
              <StatusBadge :status="r.status" />
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div v-if="r.status === 'PENDING'" class="flex gap-2">
                <button :disabled="actionId === r.id" @click="handleAction(r.id, 'CONFIRMED')" class="px-3 py-1 text-xs font-medium bg-green-600 text-white rounded hover:bg-green-700 disabled:opacity-50">
                  {{ actionId === r.id && actionType === 'CONFIRMED' ? '...' : 'Potwierdź' }}
                </button>
                <button :disabled="actionId === r.id" @click="handleAction(r.id, 'CANCELLED')" class="px-3 py-1 text-xs font-medium bg-red-600 text-white rounded hover:bg-red-700 disabled:opacity-50">
                  {{ actionId === r.id && actionType === 'CANCELLED' ? '...' : 'Odrzuć' }}
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <AppToast />
  </div>
</template>

<script setup lang="ts">
definePageMeta({ middleware: 'auth' })

const apiBase = useRuntimeConfig().public.apiBase
const authStore = useAuthStore()
const { showToast } = useToast()

if (authStore.user?.role !== 'FACILITY_OWNER' && authStore.user?.role !== 'ADMIN') {
  await navigateTo('/')
}

interface Reservation {
  id: string; facilityId: string; date: string; startTime: string; endTime: string
  status: string; totalPrice: number; facility?: any; user?: any
}
interface ApiResponse { data: Reservation[] }

const { data: response, pending, error, refresh } = await useFetch<ApiResponse>(
  `${apiBase}/facilities/my/reservations`,
  { headers: { Authorization: `Bearer ${useCookie('token').value || ''}` } },
)

const reservations = computed(() => response.value?.data || [])

const actionId = ref<string | null>(null)
const actionType = ref('')

async function handleAction(id: string, status: string) {
  actionId.value = id
  actionType.value = status
  try {
    const r = await fetch(`${apiBase}/reservations/${id}/status`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${useCookie('token').value || ''}`,
      },
      body: JSON.stringify({ status }),
    })
    if (r.ok) {
      showToast(status === 'CONFIRMED' ? 'Potwierdzono' : 'Odrzucono', 'success')
      refresh()
    } else {
      const e = await r.json()
      showToast(e.error || 'Błąd', 'error')
    }
  } catch { showToast('Błąd połączenia', 'error') }
  actionId.value = null
  actionType.value = ''
}

function formatDate(date: string): string {
  if (!date) return ''
  return new Date(date).toLocaleDateString('pl-PL', { year: 'numeric', month: 'short', day: 'numeric' })
}
</script>