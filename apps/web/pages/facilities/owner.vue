<template>
  <div class="px-4 md:px-6 lg:px-8 mx-auto max-w-7xl py-8">
    <div class="mb-8 flex items-start justify-between">
      <div>
        <NuxtLink to="/" class="text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300 text-sm">&larr; Powrót do panelu</NuxtLink>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white mt-1">{{ isAdmin ? 'Obiekty' : 'Moje obiekty' }}</h1>
        <p class="mt-2 text-gray-600 dark:text-gray-400">{{ isAdmin ? 'Zarządzaj wszystkimi obiektami w systemie' : 'Zarządzaj swoimi obiektami sportowymi' }}</p>
      </div>
      <button v-if="isAdmin" @click="showCreateModal = true" class="bg-primary-600 text-white px-4 py-2 rounded-md hover:bg-primary-700 transition-colors">+ Dodaj obiekt</button>
    </div>

    <AppLoading v-if="pending" />
    <AppError v-else-if="error" message="Nie udało się załadować obiektów." />
    <div v-else-if="facilities.length === 0" class="text-center py-12 bg-white dark:bg-gray-800 rounded-lg shadow border border-gray-200 dark:border-gray-700">
      <h3 class="text-lg font-medium text-gray-900 dark:text-white">Brak obiektów</h3>
      <p class="mt-1 text-gray-500 dark:text-gray-400">{{ isAdmin ? 'W systemie nie ma jeszcze żadnych obiektów.' : 'Nie masz jeszcze żadnych obiektów.' }}</p>
      <button v-if="isAdmin" @click="showCreateModal = true" class="mt-4 bg-primary-600 text-white px-4 py-2 rounded-md hover:bg-primary-700 transition-colors">+ Dodaj obiekt</button>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="fac in facilities" :key="fac.id" class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 hover:shadow-lg transition-shadow cursor-pointer border border-gray-200 dark:border-gray-700" @click="openEdit(fac)">
        <div class="flex items-start justify-between mb-3">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">{{ fac.name }}</h3>
          <span v-if="fac.closedUntil" class="px-2 py-0.5 text-xs font-medium rounded bg-red-100 dark:bg-red-900/30 text-red-800 dark:text-red-300">Zamknięty</span>
        </div>
        <p class="text-sm text-gray-500 dark:text-gray-400 mb-2">{{ getTypeName(fac.type) }} — {{ fac.address }}, {{ fac.city }}</p>
        <div v-if="isAdmin && fac.ownerEmail" class="text-xs text-gray-500 dark:text-gray-400 mb-2">Właściciel: {{ fac.ownerEmail }}</div>
        <div class="flex justify-between items-center">
          <span class="text-lg font-bold text-primary-600 dark:text-primary-400">{{ fac.hourlyRate }} PLN/h</span>
          <span v-if="!isAdmin" class="text-sm text-gray-400 dark:text-gray-500">Kliknij, aby edytować</span>
        </div>
        <div v-if="fac.closedUntil" class="mt-2 text-xs text-red-600 dark:text-red-400">Zamknięty do: {{ formatDate(fac.closedUntil) }}</div>
      </div>
    </div>

    <!-- Create Modal (Admin only) -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click.self="showCreateModal = false">
      <div class="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-md">
        <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-4">Dodaj obiekt</h2>
        <div class="space-y-3">
          <div><label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Nazwa</label><input v-model="createForm.name" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500" placeholder="Nazwa obiektu" /></div>
          <div><label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Adres</label><input v-model="createForm.address" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500" placeholder="Ulica i numer" /></div>
          <div><label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Miasto</label><input v-model="createForm.city" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500" placeholder="Miasto" /></div>
          <div><label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Typ</label>
            <select v-model="createForm.type" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500">
              <option v-for="opt in typeOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
            </select>
          </div>
          <div><label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Właściciel (email)</label><input v-model="createForm.ownerEmail" type="email" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500" placeholder="email@example.com" /></div>
        </div>
        <div class="flex gap-2 justify-end mt-6">
          <button @click="showCreateModal = false" class="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700">Anuluj</button>
          <button @click="createFacility" :disabled="!createForm.name || !createForm.address || !createForm.city || !createForm.ownerEmail || creating" class="px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 disabled:opacity-50">{{ creating ? 'Tworzenie...' : 'Dodaj obiekt' }}</button>
        </div>
      </div>
    </div>

    <!-- Edit Modal -->
    <div v-if="editing" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click.self="editing = null">
      <div class="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-lg max-h-[90vh] overflow-y-auto">
        <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-4">Edytuj: {{ editing.name }}</h2>
        <div class="space-y-3 mb-6">
          <div><label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Nazwa</label>
            <input v-if="isAdmin" v-model="editForm.name" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500" />
            <input v-else :value="editing.name" disabled class="w-full rounded-md bg-gray-100 dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-500 dark:text-gray-400 shadow-sm" />
          </div>
          <div><label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Opis</label><textarea v-model="editForm.description" rows="2" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500"></textarea></div>
          <div class="grid grid-cols-2 gap-3">
            <div><label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Adres</label>
              <input v-if="isAdmin" v-model="editForm.address" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500" />
              <input v-else :value="editing.address" disabled class="w-full rounded-md bg-gray-100 dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-500 dark:text-gray-400 shadow-sm" />
            </div>
            <div><label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Miasto</label>
              <input v-if="isAdmin" v-model="editForm.city" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500" />
              <input v-else :value="editing.city" disabled class="w-full rounded-md bg-gray-100 dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-500 dark:text-gray-400 shadow-sm" />
            </div>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div><label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Typ</label>
              <select v-if="isAdmin" v-model="editForm.type" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500">
                <option v-for="opt in typeOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
              </select>
              <input v-else :value="getTypeName(editing.type)" disabled class="w-full rounded-md bg-gray-100 dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-500 dark:text-gray-400 shadow-sm" />
            </div>
            <div><label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Cena za godzinę (PLN)</label><input v-model.number="editForm.hourlyRate" type="number" step="0.01" min="0" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500" /></div>
          </div>
          <div><label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Tryb rezerwacji</label>
            <select v-model="editForm.bookingMode" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500">
              <option value="BOTH">Indywidualnie i drużynowo</option>
              <option value="INDIVIDUAL">Tylko indywidualnie</option>
              <option value="TEAM">Tylko drużynowo</option>
            </select>
          </div>
          <button @click="saveInfo" :disabled="saving" class="px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 disabled:opacity-50 text-sm">{{ saving ? 'Zapisywanie...' : 'Zapisz zmiany' }}</button>
        </div>

        <hr class="border-gray-200 dark:border-gray-700 my-4" />

        <div class="mb-6">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3">Godziny otwarcia</h3>
          <div v-for="day in dayNames" :key="day.index" class="flex items-center gap-3 mb-2">
            <span class="w-32 text-sm text-gray-700 dark:text-gray-300">{{ day.label }}</span>
            <input v-model="slotsForm[day.index].open" type="time" class="form-input w-32 rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500" />
            <span class="text-gray-500 dark:text-gray-400">–</span>
            <input v-model="slotsForm[day.index].close" type="time" class="form-input w-32 rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500" />
          </div>
          <button @click="saveSlots" :disabled="savingSlots" class="mt-2 px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 disabled:opacity-50 text-sm">{{ savingSlots ? 'Zapisywanie...' : 'Zapisz godziny' }}</button>
        </div>

        <hr class="border-gray-200 dark:border-gray-700 my-4" />

        <div>
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3">Status obiektu</h3>
          <div v-if="editing.closedUntil" class="mb-3">
            <p class="text-sm text-red-600 dark:text-red-400">Zamknięty do {{ formatDate(editing.closedUntil) }}</p>
            <button @click="reopenFacility" :disabled="savingClose" class="mt-2 px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 disabled:opacity-50 text-sm">Otwórz ponownie</button>
          </div>
          <div v-else>
            <div class="flex items-end gap-3">
              <div><label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Zamknij do (data)</label><input v-model="closeUntil" type="date" class="w-48 rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500" /></div>
              <button @click="closeFacility" :disabled="!closeUntil || savingClose" class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 disabled:opacity-50 text-sm">{{ savingClose ? '...' : 'Zamknij obiekt' }}</button>
            </div>
          </div>
        </div>

        <hr class="border-gray-200 dark:border-gray-700 my-4" />

        <div>
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3">Przedpłata</h3>
          <div class="space-y-3">
            <label class="flex items-center gap-2"><input v-model="editForm.requiresPrepayment" type="checkbox" class="rounded" /><span class="text-sm text-gray-700 dark:text-gray-300">Wymagaj przedpłaty</span></label>
            <div><label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Koszt przedpłaty (PLN)</label><input v-model.number="editForm.prepaymentCost" type="number" step="0.01" min="0" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500" :disabled="!editForm.requiresPrepayment" /></div>
            <div><label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Konto bankowe</label><input v-model="editForm.bankAccount" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500" placeholder="Numer konta" :disabled="!editForm.requiresPrepayment" /></div>
            <div><label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Tytuł przelewu</label><input v-model="editForm.transferTitle" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500" placeholder="Np. Rezerwacja obiektu" :disabled="!editForm.requiresPrepayment" /></div>
            <button @click="saveInfo" :disabled="saving" class="mt-2 px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 disabled:opacity-50 text-sm">{{ saving ? 'Zapisywanie...' : 'Zapisz przedpłatę' }}</button>
          </div>
        </div>

        <div class="flex justify-end mt-6 pt-4 border-t border-gray-200 dark:border-gray-700">
          <button @click="editing = null" class="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700">Zamknij</button>
        </div>
      </div>
    </div>

    <AppToast />
  </div>
</template>

<script setup lang="ts">
definePageMeta({ middleware: 'auth' })

const apiBase = useRuntimeConfig().public.apiBase
const authStore = useAuthStore()
const { getTypeName, typeOptions } = useFacilityTypes()
const { showToast } = useToast()

const isAdmin = computed(() => authStore.user?.role === 'ADMIN')
const isFacilityOwner = computed(() => authStore.user?.role === 'FACILITY_OWNER')

if (!isAdmin.value && !isFacilityOwner.value) {
  await navigateTo('/')
}

interface Slot { id?: string; facility_id?: string; day_of_week: number; open_time: string; close_time: string }
interface Facility {
  id: string; name: string; description?: string; type: string; address: string; city: string
  hourlyRate: number; closedUntil?: string; slots?: Slot[]; ownerEmail?: string
  bookingMode?: string; requiresPrepayment?: boolean; prepaymentCost?: number
  bankAccount?: string; transferTitle?: string
}
interface ApiResponse { data: Facility[] }

const fetchUrl = computed(() => isAdmin.value ? `${apiBase}/facilities` : `${apiBase}/facilities/mine`)

const { data: response, pending, error, refresh } = await useFetch<ApiResponse>(fetchUrl, {
  headers: { Authorization: `Bearer ${useCookie('token').value || ''}` },
})

const facilities = computed(() =>
  (response.value?.data || []).map((f: any) => ({
    ...f,
    hourlyRate: typeof f.hourlyRate === 'string' ? parseFloat(f.hourlyRate) : (f.hourlyRate || 0),
    prepaymentCost: typeof f.prepaymentCost === 'string' ? parseFloat(f.prepaymentCost) : (f.prepaymentCost || 0),
  })),
)

const showCreateModal = ref(false)
const creating = ref(false)
const createForm = reactive({ name: '', address: '', city: '', type: 'FOOTBALL', ownerEmail: '' })

async function createFacility() {
  creating.value = true
  try {
    const r = await fetch(`${apiBase}/facilities`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${useCookie('token').value || ''}` },
      body: JSON.stringify({ name: createForm.name, address: createForm.address, city: createForm.city, type: createForm.type, owner_email: createForm.ownerEmail }),
    })
    if (r.ok) {
      showToast('Obiekt utworzony', 'success')
      showCreateModal.value = false
      Object.assign(createForm, { name: '', address: '', city: '', type: 'FOOTBALL', ownerEmail: '' })
      refresh()
    } else {
      const e = await r.json()
      showToast(e.error || 'Błąd', 'error')
    }
  } catch { showToast('Błąd połączenia', 'error') }
  creating.value = false
}

const editing = ref<Facility | null>(null)
const saving = ref(false)
const savingSlots = ref(false)
const savingClose = ref(false)

const editForm = reactive({
  name: '', description: '', address: '', city: '', type: '', hourlyRate: 0,
  bookingMode: 'BOTH', requiresPrepayment: false, prepaymentCost: 0,
  bankAccount: '', transferTitle: '',
})

const dayNames = [
  { index: 0, label: 'Niedziela' }, { index: 1, label: 'Poniedziałek' },
  { index: 2, label: 'Wtorek' }, { index: 3, label: 'Środa' },
  { index: 4, label: 'Czwartek' }, { index: 5, label: 'Piątek' }, { index: 6, label: 'Sobota' },
]

const slotsForm = reactive<Record<number, { open: string; close: string }>>({
  0: { open: '08:00', close: '20:00' }, 1: { open: '06:00', close: '22:00' },
  2: { open: '06:00', close: '22:00' }, 3: { open: '06:00', close: '22:00' },
  4: { open: '06:00', close: '22:00' }, 5: { open: '06:00', close: '22:00' },
  6: { open: '08:00', close: '20:00' },
})

const closeUntil = ref('')

function openEdit(fac: any) {
  editing.value = fac
  editForm.name = fac.name; editForm.description = fac.description || ''
  editForm.address = fac.address; editForm.city = fac.city; editForm.type = fac.type
  editForm.hourlyRate = fac.hourlyRate; editForm.bookingMode = fac.bookingMode || 'BOTH'
  editForm.requiresPrepayment = fac.requiresPrepayment || false
  editForm.prepaymentCost = typeof fac.prepaymentCost === 'string' ? parseFloat(fac.prepaymentCost) : (fac.prepaymentCost || 0)
  editForm.bankAccount = fac.bankAccount || ''; editForm.transferTitle = fac.transferTitle || ''
  for (let i = 0; i < 7; i++) slotsForm[i] = { open: '08:00', close: '20:00' }
  if (fac.slots) for (const s of fac.slots) slotsForm[s.day_of_week] = { open: s.open_time, close: s.close_time }
  closeUntil.value = ''
}

async function saveInfo() {
  if (!editing.value) return
  saving.value = true
  try {
    const r = await fetch(`${apiBase}/facilities/${editing.value.id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${useCookie('token').value || ''}` },
      body: JSON.stringify({
        ...(isAdmin.value ? { name: editForm.name, address: editForm.address, city: editForm.city, type: editForm.type } : {}),
        description: editForm.description || null,
        booking_mode: editForm.bookingMode,
        hourly_rate: editForm.hourlyRate,
        requires_prepayment: editForm.requiresPrepayment,
        prepayment_cost: editForm.requiresPrepayment ? parseFloat(String(editForm.prepaymentCost)) || 0 : null,
        bank_account: editForm.requiresPrepayment ? editForm.bankAccount : null,
        transfer_title: editForm.requiresPrepayment ? editForm.transferTitle : null,
      }),
    })
    if (r.ok) { showToast('Zapisano', 'success'); refresh() }
    else { const e = await r.json(); showToast(e.error || 'Błąd', 'error') }
  } catch { showToast('Błąd połączenia', 'error') }
  saving.value = false
}

async function saveSlots() {
  if (!editing.value) return
  savingSlots.value = true
  const slots = dayNames.map(d => ({ day_of_week: d.index, open_time: slotsForm[d.index].open, close_time: slotsForm[d.index].close }))
  try {
    const r = await fetch(`${apiBase}/facilities/${editing.value.id}/slots`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${useCookie('token').value || ''}` },
      body: JSON.stringify({ slots }),
    })
    if (r.ok) { showToast('Godziny zapisane', 'success'); refresh() }
    else { const e = await r.json(); showToast(e.error || 'Błąd', 'error') }
  } catch { showToast('Błąd połączenia', 'error') }
  savingSlots.value = false
}

async function closeFacility() {
  if (!editing.value || !closeUntil.value) return
  savingClose.value = true
  try {
    const r = await fetch(`${apiBase}/facilities/${editing.value.id}/close`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${useCookie('token').value || ''}` },
      body: JSON.stringify({ closed_until: closeUntil.value }),
    })
    if (r.ok) { showToast('Obiekt zamknięty', 'success'); refresh(); editing.value = null }
    else { const e = await r.json(); showToast(e.error || 'Błąd', 'error') }
  } catch { showToast('Błąd połączenia', 'error') }
  savingClose.value = false
}

async function reopenFacility() {
  if (!editing.value) return
  savingClose.value = true
  try {
    const r = await fetch(`${apiBase}/facilities/${editing.value.id}/close`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${useCookie('token').value || ''}` },
      body: JSON.stringify({ closed_until: null }),
    })
    if (r.ok) { showToast('Obiekt otwarty', 'success'); refresh(); editing.value = null }
    else { const e = await r.json(); showToast(e.error || 'Błąd', 'error') }
  } catch { showToast('Błąd połączenia', 'error') }
  savingClose.value = false
}

function formatDate(date: string): string {
  if (!date) return ''
  return new Date(date).toLocaleDateString('pl-PL', { year: 'numeric', month: 'long', day: 'numeric' })
}
</script>