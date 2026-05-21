<template>
  <div class="page-container">
    <div class="page-header">
      <NuxtLink to="/" class="text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300 text-sm">&larr; Powrót do panelu</NuxtLink>
      <h1 class="page-title mt-1">Rezerwacje na moje obiekty</h1>
      <p class="page-subtitle">
        Zarządzaj rezerwacjami dla swoich obiektów sportowych
      </p>
    </div>

    <div v-if="pending" class="loading-spinner">
      <div class="spinner"></div>
    </div>

    <div
      v-else-if="error"
      class="error-box"
    >
      Nie udało się załadować rezerwacji. Spróbuj ponownie.
    </div>

    <div
      v-else-if="reservations.length === 0"
      class="empty-state"
    >
      <svg
        class="empty-icon"
        fill="none"
        viewBox="0 0 24 24"
        stroke="currentColor"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
        />
      </svg>
      <h3 class="empty-title">Brak rezerwacji</h3>
      <p class="empty-text">
        Nie masz jeszcze żadnych rezerwacji na swoje obiekty.
      </p>
    </div>

    <div
      v-else
      class="data-table-wrapper"
    >
      <table class="data-table">
        <thead>
          <tr>
            <th>Obiekt</th>
            <th>Użytkownik</th>
            <th>Data</th>
            <th>Godzina</th>
            <th>Status</th>
            <th>Akcje</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="r in reservations"
            :key="r.id"
          >
            <td>
              <div class="facility-name">{{ r.facility?.name }}</div>
              <div class="facility-address">
                {{ r.facility?.address }}, {{ r.facility?.city }}
              </div>
            </td>
            <td>
              <div class="user-name">
                {{ r.user?.first_name }} {{ r.user?.last_name }}
              </div>
              <div class="user-email">{{ r.user?.email }}</div>
            </td>
            <td class="text-gray-900 dark:text-white text-sm">{{ formatDate(r.date) }}</td>
            <td class="text-gray-900 dark:text-white text-sm">{{ r.startTime }} – {{ r.endTime }}</td>
            <td>
              <span :class="statusBadgeClass(r.status)">
                {{ statusLabel(r.status) }}
              </span>
            </td>
            <td>
              <div
                v-if="r.status === 'PENDING'"
                class="action-buttons"
              >
                <button
                  class="btn-confirm"
                  :disabled="actionId === r.id"
                  @click="handleAction(r.id, 'CONFIRMED')"
                >
                  {{ actionId === r.id && actionType === 'CONFIRMED' ? '...' : 'Potwierdź' }}
                </button>
                <button
                  class="btn-decline"
                  :disabled="actionId === r.id"
                  @click="handleAction(r.id, 'CANCELLED')"
                >
                  {{ actionId === r.id && actionType === 'CANCELLED' ? '...' : 'Odrzuć' }}
                </button>
              </div>
              <span v-else class="text-gray-400 dark:text-gray-500 text-sm">{{ statusLabel(r.status) }}</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div
      v-if="toast.message"
      :class="toastClass"
    >
      {{ toast.message }}
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ middleware: 'auth' })

const apiBase = useRuntimeConfig().public.apiBase
const authStore = useAuthStore()

if (authStore.user?.role !== 'FACILITY_OWNER' && authStore.user?.role !== 'ADMIN') {
  await navigateTo('/')
}

interface Facility {
  name: string
  address: string
  city: string
}

interface User {
  first_name: string
  last_name: string
  email: string
}

interface Reservation {
  id: string
  facilityId: string
  date: string
  startTime: string
  endTime: string
  status: string
  totalPrice: number
  facility?: Facility
  user?: User
}

interface ApiResponse {
  data: Reservation[]
}

const { data: response, pending, error, refresh } = await useFetch<ApiResponse>(
  `${apiBase}/facilities/my/reservations`,
  {
    headers: { Authorization: `Bearer ${useCookie('token').value || ''}` },
  },
)

const reservations = computed(() => response.value?.data || [])

const actionId = ref<string | null>(null)
const actionType = ref<string>('')

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
  } catch {
    showToast('Błąd połączenia', 'error')
  }
  actionId.value = null
  actionType.value = ''
}

function statusBadgeClass(status: string): string {
  const map: Record<string, string> = {
    PENDING: 'badge-yellow',
    CONFIRMED: 'badge-green',
    CANCELLED: 'badge-red',
    COMPLETED: 'badge-gray',
  }
  return `status-badge ${map[status] || 'badge-gray'}`
}

function statusLabel(status: string): string {
  const map: Record<string, string> = {
    PENDING: 'Oczekująca',
    CONFIRMED: 'Potwierdzona',
    CANCELLED: 'Anulowana',
    COMPLETED: 'Zakończona',
  }
  return map[status] || status
}

function formatDate(date: string): string {
  if (!date) return ''
  return new Date(date).toLocaleDateString('pl-PL', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}

const toast = reactive({ message: '', type: '' })
const toastClass = computed(() =>
  toast.type === 'success'
    ? 'toast toast-success'
    : 'toast toast-error',
)

function showToast(m: string, t: string) {
  toast.message = m
  toast.type = t
  setTimeout(() => { toast.message = '' }, 3000)
}
</script>