<template>
  <div class="page-container">
    <div class="section-header mb-8">
      <div>
        <NuxtLink to="/" class="text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300 text-sm">&larr; Powrót do panelu</NuxtLink>
        <h1 class="page-title mt-1">
          {{ isAdmin ? 'Obiekty' : 'Moje obiekty' }}
        </h1>
        <p class="page-subtitle">
          {{ isAdmin ? 'Zarządzaj wszystkimi obiektami w systemie' : 'Zarządzaj swoimi obiektami sportowymi' }}
        </p>
      </div>
      <button v-if="isAdmin" @click="showCreateModal = true" class="btn-primary">
        + Dodaj obiekt
      </button>
    </div>

    <div v-if="pending" class="loading-spinner">
      <div class="spinner"></div>
    </div>

    <div v-else-if="error" class="error-box">
      Nie udało się załadować obiektów.
    </div>

    <div v-else-if="facilities.length === 0" class="empty-state">
      <h3 class="empty-title">Brak obiektów</h3>
      <p class="empty-text">
        {{ isAdmin ? 'W systemie nie ma jeszcze żadnych obiektów.' : 'Nie masz jeszcze żadnych obiektów.' }}
      </p>
      <button v-if="isAdmin" @click="showCreateModal = true" class="btn-primary mt-4">
        + Dodaj obiekt
      </button>
    </div>

    <div v-else class="card-grid">
      <div
        v-for="fac in facilities"
        :key="fac.id"
        class="card-hover"
        @click="openEdit(fac)"
      >
        <div class="flex items-start justify-between mb-3">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ fac.name }}</h3>
          <span v-if="fac.closedUntil" class="badge-red">Zamknięty</span>
        </div>
        <p class="text-sm text-gray-500 dark:text-gray-400 mb-2">
          {{ typeLabel(fac.type) }} — {{ fac.address }}, {{ fac.city }}
        </p>
        <div v-if="isAdmin && fac.ownerEmail" class="text-xs text-gray-500 dark:text-gray-400 mb-2">
          Właściciel: {{ fac.ownerEmail }}
        </div>
        <div class="flex justify-between items-center">
          <span class="text-lg font-bold text-primary-600 dark:text-primary-400">
            {{ fac.hourlyRate }} PLN/h
          </span>
          <span v-if="!isAdmin" class="text-sm text-gray-400 dark:text-gray-500">Kliknij, aby edytować</span>
        </div>
        <div v-if="fac.closedUntil" class="mt-2 text-xs text-red-600 dark:text-red-400">
          Zamknięty do: {{ formatDate(fac.closedUntil) }}
        </div>
      </div>
    </div>

    <!-- Create Modal (Admin only) -->
    <div v-if="showCreateModal" class="modal-overlay" @click.self="showCreateModal = false">
      <div class="modal-content">
        <h2 class="modal-title">Dodaj obiekt</h2>
        <div class="space-y-3">
          <div>
            <label class="form-label">Nazwa</label>
            <input v-model="createForm.name" class="form-input" placeholder="Nazwa obiektu" />
          </div>
          <div>
            <label class="form-label">Adres</label>
            <input v-model="createForm.address" class="form-input" placeholder="Ulica i numer" />
          </div>
          <div>
            <label class="form-label">Miasto</label>
            <input v-model="createForm.city" class="form-input" placeholder="Miasto" />
          </div>
          <div>
            <label class="form-label">Typ</label>
            <select v-model="createForm.type" class="form-select">
              <option value="FOOTBALL">Piłka nożna</option>
              <option value="BASKETBALL">Koszykówka</option>
              <option value="TENNIS">Tenis</option>
              <option value="VOLLEYBALL">Siatkówka</option>
              <option value="SWIMMING">Pływanie</option>
              <option value="OTHER">Inne</option>
            </select>
          </div>
          <div>
            <label class="form-label">Właściciel (email)</label>
            <input v-model="createForm.ownerEmail" type="email" class="form-input" placeholder="email@example.com" />
          </div>
        </div>
        <div class="modal-actions mt-6">
          <button @click="showCreateModal = false" class="btn-secondary">Anuluj</button>
          <button
            @click="createFacility"
            :disabled="!createForm.name || !createForm.address || !createForm.city || !createForm.ownerEmail || creating"
            class="btn-primary"
          >
            {{ creating ? 'Tworzenie...' : 'Dodaj obiekt' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Edit Modal (Facility Owner only) -->
    <div v-if="editing" class="modal-overlay" @click.self="editing = null">
      <div class="modal-content-lg">
        <h2 class="modal-title">Edytuj: {{ editing.name }}</h2>

        <div class="space-y-3 mb-6">
          <div>
            <label class="form-label">Nazwa</label>
            <input v-if="isAdmin" v-model="editForm.name" class="form-input" />
            <input v-else :value="editing.name" disabled class="form-input-disabled" />
          </div>
          <div>
            <label class="form-label">Opis</label>
            <textarea v-model="editForm.description" rows="2" class="form-textarea"></textarea>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="form-label">Adres</label>
              <input v-if="isAdmin" v-model="editForm.address" class="form-input" />
              <input v-else :value="editing.address" disabled class="form-input-disabled" />
            </div>
            <div>
              <label class="form-label">Miasto</label>
              <input v-if="isAdmin" v-model="editForm.city" class="form-input" />
              <input v-else :value="editing.city" disabled class="form-input-disabled" />
            </div>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="form-label">Typ</label>
              <select v-if="isAdmin" v-model="editForm.type" class="form-select">
                <option value="FOOTBALL">Piłka nożna</option>
                <option value="BASKETBALL">Koszykówka</option>
                <option value="TENNIS">Tenis</option>
                <option value="VOLLEYBALL">Siatkówka</option>
                <option value="SWIMMING">Pływanie</option>
                <option value="OTHER">Inne</option>
              </select>
              <input v-else :value="typeLabel(editing.type)" disabled class="form-input-disabled" />
            </div>
            <div>
              <label class="form-label">Cena za godzinę (PLN)</label>
              <input v-model.number="editForm.hourlyRate" type="number" step="0.01" min="0" class="form-input" />
            </div>
          </div>
          <!-- Booking Mode -->
          <div>
            <label class="form-label">Tryb rezerwacji</label>
            <select v-model="editForm.bookingMode" class="form-select">
              <option value="BOTH">Indywidualnie i drużynowo</option>
              <option value="INDIVIDUAL">Tylko indywidualnie</option>
              <option value="TEAM">Tylko drużynowo</option>
            </select>
          </div>
          <button @click="saveInfo" :disabled="saving" class="btn-primary-sm">
            {{ saving ? 'Zapisywanie...' : 'Zapisz zmiany' }}
          </button>
        </div>

        <hr class="border-gray-200 dark:border-gray-700 my-4" />

        <div class="mb-6">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3">Godziny otwarcia</h3>
          <div v-for="day in dayNames" :key="day.index" class="flex items-center gap-3 mb-2">
            <span class="w-32 text-sm text-gray-700 dark:text-gray-300">{{ day.label }}</span>
            <input v-model="slotsForm[day.index].open" type="time" class="form-input w-32" />
            <span class="text-gray-500 dark:text-gray-400">–</span>
            <input v-model="slotsForm[day.index].close" type="time" class="form-input w-32" />
          </div>
          <button @click="saveSlots" :disabled="savingSlots" class="btn-primary-sm mt-2">
            {{ savingSlots ? 'Zapisywanie...' : 'Zapisz godziny' }}
          </button>
        </div>

        <hr class="border-gray-200 dark:border-gray-700 my-4" />

        <div>
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3">Status obiektu</h3>
          <div v-if="editing.closedUntil" class="mb-3">
            <p class="text-sm text-red-600 dark:text-red-400">
              Zamknięty do {{ formatDate(editing.closedUntil) }}
            </p>
            <button @click="reopenFacility" :disabled="savingClose" class="btn-primary-sm mt-2">
              Otwórz ponownie
            </button>
          </div>
          <div v-else>
            <div class="flex items-end gap-3">
              <div>
                <label class="form-label">Zamknij do (data)</label>
                <input v-model="closeUntil" type="date" class="form-input w-48" />
              </div>
              <button @click="closeFacility" :disabled="!closeUntil || savingClose" class="btn-danger">
                {{ savingClose ? '...' : 'Zamknij obiekt' }}
              </button>
            </div>
          </div>
        </div>

        <hr class="border-gray-200 dark:border-gray-700 my-4" />

        <div>
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3">Przedpłata</h3>
          <div class="space-y-3">
            <label class="flex items-center gap-2">
              <input v-model="editForm.requiresPrepayment" type="checkbox" class="rounded" />
              <span class="text-sm text-gray-700 dark:text-gray-300">Wymagaj przedpłaty</span>
            </label>
            <div>
              <label class="form-label">Koszt przedpłaty (PLN)</label>
              <input v-model.number="editForm.prepaymentCost" type="number" step="0.01" min="0" class="form-input" :disabled="!editForm.requiresPrepayment" />
            </div>
            <div>
              <label class="form-label">Konto bankowe</label>
              <input v-model="editForm.bankAccount" class="form-input" placeholder="Numer konta" :disabled="!editForm.requiresPrepayment" />
            </div>
            <div>
              <label class="form-label">Tytuł przelewu</label>
              <input v-model="editForm.transferTitle" class="form-input" placeholder="Np. Rezerwacja obiektu" :disabled="!editForm.requiresPrepayment" />
            </div>
            <button @click="saveInfo" :disabled="saving" class="btn-primary-sm mt-2">
              {{ saving ? 'Zapisywanie...' : 'Zapisz przedpłatę' }}
            </button>
          </div>
        </div>

        <div class="modal-actions mt-6 pt-4 border-t border-gray-200 dark:border-gray-700">
          <button @click="editing = null" class="btn-secondary">Zamknij</button>
        </div>
      </div>
    </div>

    <div v-if="toast.message" :class="toastClass">{{ toast.message }}</div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ middleware: 'auth' })

const apiBase = useRuntimeConfig().public.apiBase
const authStore = useAuthStore()

const isAdmin = computed(() => authStore.user?.role === 'ADMIN')
const isFacilityOwner = computed(() => authStore.user?.role === 'FACILITY_OWNER')

if (!isAdmin.value && !isFacilityOwner.value) {
  await navigateTo('/')
}

interface Slot {
  id?: string
  facility_id?: string
  day_of_week: number
  open_time: string
  close_time: string
}

interface Facility {
  id: string
  name: string
  description?: string
  type: string
  address: string
  city: string
  hourlyRate: number
  closedUntil?: string
  slots?: Slot[]
  ownerEmail?: string
}

interface ApiResponse {
  data: Facility[]
}

const fetchUrl = computed(() =>
  isAdmin.value ? `${apiBase}/facilities` : `${apiBase}/facilities/mine`,
)

const { data: response, pending, error, refresh } = await useFetch<ApiResponse>(fetchUrl, {
  headers: { Authorization: `Bearer ${useCookie('token').value || ''}` },
})

const facilities = computed(() =>
  (response.value?.data || []).map((f: any) => {
    const rate = typeof f.hourlyRate === 'string' ? parseFloat(f.hourlyRate) : f.hourlyRate
    const ppCost = typeof f.prepaymentCost === 'string' ? parseFloat(f.prepaymentCost) : (f.prepaymentCost || 0)
    return { ...f, hourlyRate: rate || 0, prepaymentCost: ppCost }
  }),
)

// Modal tworzenia obiektu (admin)
const showCreateModal = ref(false)
const creating = ref(false)
const createForm = reactive({
  name: '',
  address: '',
  city: '',
  type: 'FOOTBALL',
  ownerEmail: '',
})

async function createFacility() {
  creating.value = true
  try {
    const r = await fetch(`${apiBase}/facilities`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${useCookie('token').value || ''}`,
      },
      body: JSON.stringify({
        name: createForm.name,
        address: createForm.address,
        city: createForm.city,
        type: createForm.type,
        owner_email: createForm.ownerEmail,
      }),
    })
    if (r.ok) {
      showToast('Obiekt utworzony', 'success')
      showCreateModal.value = false
      createForm.name = ''
      createForm.address = ''
      createForm.city = ''
      createForm.type = 'FOOTBALL'
      createForm.ownerEmail = ''
      refresh()
    } else {
      const e = await r.json()
      showToast(e.error || 'Błąd', 'error')
    }
  } catch {
    showToast('Błąd połączenia', 'error')
  }
  creating.value = false
}

// Modal edycji obiektu (właściciel)
const editing = ref<Facility | null>(null)
const saving = ref(false)
const savingSlots = ref(false)
const savingClose = ref(false)

const editForm = reactive({
  name: '',
  description: '',
  address: '',
  city: '',
  type: '',
  hourlyRate: 0,
  requiresPrepayment: false,
  prepaymentCost: 0,
  bankAccount: '',
  transferTitle: '',
})

const dayNames = [
  { index: 0, label: 'Niedziela' },
  { index: 1, label: 'Poniedziałek' },
  { index: 2, label: 'Wtorek' },
  { index: 3, label: 'Środa' },
  { index: 4, label: 'Czwartek' },
  { index: 5, label: 'Piątek' },
  { index: 6, label: 'Sobota' },
]

const slotsForm = reactive<Record<number, { open: string; close: string }>>({
  0: { open: '08:00', close: '20:00' },
  1: { open: '06:00', close: '22:00' },
  2: { open: '06:00', close: '22:00' },
  3: { open: '06:00', close: '22:00' },
  4: { open: '06:00', close: '22:00' },
  5: { open: '06:00', close: '22:00' },
  6: { open: '08:00', close: '20:00' },
})

const closeUntil = ref('')

function openEdit(fac: any) {
  editing.value = fac
  editForm.name = fac.name
  editForm.description = fac.description || ''
  editForm.address = fac.address
  editForm.city = fac.city
  editForm.type = fac.type
  editForm.hourlyRate = fac.hourlyRate
  editForm.bookingMode = fac.bookingMode || 'BOTH'
  editForm.requiresPrepayment = fac.requiresPrepayment || false
  editForm.prepaymentCost = typeof fac.prepaymentCost === 'string' ? parseFloat(fac.prepaymentCost) : (fac.prepaymentCost || 0)
  editForm.bankAccount = fac.bankAccount || ''
  editForm.transferTitle = fac.transferTitle || ''

  for (let i = 0; i < 7; i++) {
    slotsForm[i] = { open: '08:00', close: '20:00' }
  }
  if (fac.slots) {
    for (const s of fac.slots) {
      slotsForm[s.day_of_week] = { open: s.open_time, close: s.close_time }
    }
  }

  closeUntil.value = ''
}

async function saveInfo() {
  if (!editing.value) return
  saving.value = true
  try {
    const r = await fetch(`${apiBase}/facilities/${editing.value.id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${useCookie('token').value || ''}`,
      },
      body: JSON.stringify({
        ...(isAdmin.value ? {
          name: editForm.name,
          address: editForm.address,
          city: editForm.city,
          type: editForm.type,
        } : {}),
        description: editForm.description || null,
        booking_mode: editForm.bookingMode,
        hourly_rate: editForm.hourlyRate,
        requires_prepayment: editForm.requiresPrepayment,
        prepayment_cost: editForm.requiresPrepayment ? parseFloat(String(editForm.prepaymentCost)) || 0 : null,
        bank_account: editForm.requiresPrepayment ? editForm.bankAccount : null,
        transfer_title: editForm.requiresPrepayment ? editForm.transferTitle : null,
      }),
    })
    if (r.ok) {
      showToast('Zapisano', 'success')
      refresh()
    } else {
      const e = await r.json()
      showToast(e.error || 'Błąd', 'error')
    }
  } catch {
    showToast('Błąd połączenia', 'error')
  }
  saving.value = false
}

async function saveSlots() {
  if (!editing.value) return
  savingSlots.value = true
  const slots = dayNames.map((d) => ({
    day_of_week: d.index,
    open_time: slotsForm[d.index].open,
    close_time: slotsForm[d.index].close,
  }))
  try {
    const r = await fetch(`${apiBase}/facilities/${editing.value.id}/slots`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${useCookie('token').value || ''}`,
      },
      body: JSON.stringify({ slots }),
    })
    if (r.ok) {
      showToast('Godziny zapisane', 'success')
      refresh()
    } else {
      const e = await r.json()
      showToast(e.error || 'Błąd', 'error')
    }
  } catch {
    showToast('Błąd połączenia', 'error')
  }
  savingSlots.value = false
}

async function closeFacility() {
  if (!editing.value || !closeUntil.value) return
  savingClose.value = true
  try {
    const r = await fetch(`${apiBase}/facilities/${editing.value.id}/close`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${useCookie('token').value || ''}`,
      },
      body: JSON.stringify({ closed_until: closeUntil.value }),
    })
    if (r.ok) {
      showToast('Obiekt zamknięty', 'success')
      refresh()
      editing.value = null
    } else {
      const e = await r.json()
      showToast(e.error || 'Błąd', 'error')
    }
  } catch {
    showToast('Błąd połączenia', 'error')
  }
  savingClose.value = false
}

async function reopenFacility() {
  if (!editing.value) return
  savingClose.value = true
  try {
    const r = await fetch(`${apiBase}/facilities/${editing.value.id}/close`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${useCookie('token').value || ''}`,
      },
      body: JSON.stringify({ closed_until: null }),
    })
    if (r.ok) {
      showToast('Obiekt otwarty', 'success')
      refresh()
      editing.value = null
    } else {
      const e = await r.json()
      showToast(e.error || 'Błąd', 'error')
    }
  } catch {
    showToast('Błąd połączenia', 'error')
  }
  savingClose.value = false
}

function typeLabel(type: string): string {
  const map: Record<string, string> = {
    FOOTBALL: 'Piłka nożna',
    BASKETBALL: 'Koszykówka',
    TENNIS: 'Tenis',
    VOLLEYBALL: 'Siatkówka',
    SWIMMING: 'Pływanie',
    OTHER: 'Inne',
  }
  return map[type] || type
}

function formatDate(date: string): string {
  if (!date) return ''
  return new Date(date).toLocaleDateString('pl-PL', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

const toast = reactive({ message: '', type: '' })
const toastClass = computed(() =>
  toast.type === 'success' ? 'toast toast-success' : 'toast toast-error',
)

function showToast(m: string, t: string) {
  toast.message = m
  toast.type = t
  setTimeout(() => { toast.message = '' }, 3000)
}
</script>