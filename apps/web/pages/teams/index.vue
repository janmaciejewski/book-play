<template>
  <div class="px-4 md:px-6 lg:px-8 mx-auto max-w-7xl py-8">
    <div class="mb-8 flex items-center justify-between flex-wrap gap-4">
      <div>
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
          {{ showMyTeams ? 'Moje drużyny' : 'Wszystkie drużyny' }}
        </h1>
        <p class="mt-2 text-gray-600 dark:text-gray-400">
          {{ showMyTeams ? 'Zarządzaj drużynami, do których należysz' : 'Przeglądaj drużyny i aplikuj do tych z otwartą rekrutacją' }}
        </p>
      </div>
      <div class="flex gap-2">
        <button @click="showMyTeams = !showMyTeams"
          :class="['py-2 px-4 rounded-md transition-colors text-sm font-medium', showMyTeams ? 'bg-primary-600 text-white' : 'border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700']">
          {{ showMyTeams ? '← Wszystkie drużyny' : 'Moje drużyny' }}
        </button>
        <button @click="showCreateModal = true" class="bg-primary-600 text-white py-2 px-4 rounded-md hover:bg-primary-700 transition-colors">Utwórz drużynę</button>
      </div>
    </div>

    <div v-if="pending" class="flex justify-center py-12"><div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500"></div></div>

    <div v-else-if="error" class="bg-red-50 dark:bg-red-900/50 text-red-700 dark:text-red-300 p-4 rounded-lg border border-red-200 dark:border-red-800">Nie udało się załadować drużyn.</div>

    <div v-else-if="teams.length === 0" class="text-center py-12 bg-white dark:bg-gray-800 rounded-lg shadow border border-gray-200 dark:border-gray-700">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" /></svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900 dark:text-white">{{ showMyTeams ? 'Nie należysz do żadnej drużyny' : 'Brak drużyn' }}</h3>
      <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">{{ showMyTeams ? 'Utwórz drużynę lub dołącz do istniejącej.' : 'Utwórz drużynę, aby zacząć grać z przyjaciółmi.' }}</p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="team in teams" :key="team.id"
        class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 hover:shadow-lg transition-shadow cursor-pointer border border-gray-200 dark:border-gray-700 relative"
        @click="viewTeam(team.id)">
        <div class="flex items-center justify-between mb-4">
          <div v-if="team.logo" class="w-12 h-12 rounded-full overflow-hidden bg-gray-100 dark:bg-gray-700">
            <img :src="getLogoUrl(team.logo)" :alt="team.name" class="w-full h-full object-cover" />
          </div>
          <div v-else class="w-12 h-12 bg-primary-100 dark:bg-primary-900/30 rounded-full flex items-center justify-center">
            <span class="text-xl font-bold text-primary-600 dark:text-primary-400">{{ team.name.charAt(0) }}</span>
          </div>
          <div class="flex flex-col items-end gap-1">
            <span :class="['px-2 py-0.5 rounded text-xs font-medium', team.userRole === 'CAPTAIN' ? 'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-300' : team.userRole === 'ADMIN' ? 'bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-300' : team.userRole ? 'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300' : 'bg-transparent text-transparent']">
              {{ getRoleLabel(team.userRole || '') || '&nbsp;' }}
            </span>
            <span v-if="team.recruitment_open" class="px-2 py-0.5 rounded text-xs font-medium bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300">Rekrutacja otwarta</span>
          </div>
        </div>
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-1">{{ team.name }}</h3>
        <p class="text-sm text-gray-600 dark:text-gray-400 mb-4 line-clamp-2">{{ team.description || 'Brak opisu' }}</p>
        <div class="flex items-center text-sm text-gray-500 dark:text-gray-400">
          <svg class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" /></svg>
          {{ (team.members || []).length || 1 }} członków
        </div>
        <div v-if="showAllTeams && team.recruitment_open && !isMember(team)" class="mt-3" @click.stop>
          <button @click="openApplyModal(team)" class="w-full px-3 py-1.5 text-xs font-medium bg-green-600 text-white rounded-md hover:bg-green-700">Aplikuj</button>
        </div>
      </div>
    </div>

    <div v-if="showApplyModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-md">
        <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-4">Aplikuj do {{ applyingTeam?.name }}</h2>
        <div class="mb-4"><label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Twoje imię</label><input :value="authStore.user?.first_name + ' ' + authStore.user?.last_name" disabled class="w-full rounded-md bg-gray-100 dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-500 dark:text-gray-400 shadow-sm" /></div>
        <div class="mb-4"><label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Dlaczego chcesz dołączyć?</label><textarea v-model="applyMessage" rows="4" required class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500" placeholder="Napisz kilka słów o sobie i dlaczego chcesz dołączyć..."></textarea></div>
        <div class="flex gap-2 justify-end"><button @click="showApplyModal = false; applyMessage = ''" class="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700">Anuluj</button><button @click="submitApplication" :disabled="!applyMessage.trim() || applying" class="px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 disabled:opacity-50">{{ applying ? 'Wysyłanie...' : 'Wyślij aplikację' }}</button></div>
      </div>
    </div>

    <div v-if="showCreateModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-md">
        <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-4">Utwórz nową drużynę</h2>
        <form @submit.prevent="createTeam">
          <div class="mb-4"><label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Nazwa drużyny</label><input v-model="newTeam.name" type="text" required class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500" placeholder="Wprowadź nazwę drużyny" /></div>
          <div class="mb-4"><label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Opis</label><textarea v-model="newTeam.description" rows="3" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500" placeholder="Wprowadź opis drużyny"></textarea></div>
          <div class="flex gap-2 justify-end"><button type="button" @click="showCreateModal = false" class="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700">Anuluj</button><button type="submit" class="px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700">Utwórz drużynę</button></div>
        </form>
      </div>
    </div>

    <div v-if="toast.message" :class="['fixed bottom-4 right-4 px-4 py-3 rounded-lg shadow-lg z-50 text-sm', toast.type === 'success' ? 'bg-green-500 text-white' : 'bg-red-500 text-white']">{{ toast.message }}</div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({ middleware: 'auth' })
const apiBase = useRuntimeConfig().public.apiBase
const authStore = useAuthStore()
const showMyTeams = ref(false)
const showCreateModal = ref(false)
const newTeam = reactive({ name: '', description: '' })

interface User { id: string; email: string; first_name: string; last_name: string }
interface Member { id: string; role: string; user?: User }
interface Team { id: string; name: string; description?: string; logo?: string; captain_id?: string; userRole?: string; recruitment_open?: boolean; members?: Member[] }
interface TeamsResponse { data: Team[] }

const fetchUrl = computed(() => showMyTeams.value ? `${apiBase}/my-teams` : `${apiBase}/teams`)
const headers = { Authorization: `Bearer ${useCookie('token').value || ''}` }
const { data: response, pending, error, refresh } = await useFetch<TeamsResponse>(fetchUrl, { headers })

const teams = computed(() => {
  const raw = response.value?.data || []
  if (authStore.user?.id) { return raw.map(t => { const my = t.members?.find(m => m.user?.id === authStore.user!.id); return { ...t, userRole: my?.role } }) }
  return raw
})
const showAllTeams = computed(() => !showMyTeams.value)
function isMember(team: Team): boolean { return !!team.userRole }
function getRoleLabel(role: string): string { const l: Record<string,string>={CAPTAIN:'Kapitan',ADMIN:'Administrator',MEMBER:'Członek'}; return l[role]||role }
watch(showMyTeams, () => refresh())

const createTeam = async () => {
  try {
    const r = await fetch(`${apiBase}/teams`, { method:'POST', headers:{'Content-Type':'application/json', Authorization: `Bearer ${useCookie('token').value || ''}`}, body: JSON.stringify(newTeam) })
    if (r.ok) { showCreateModal.value=false; newTeam.name=''; newTeam.description=''; showToast('Drużyna utworzona','success'); refresh() }
  } catch { }
}

const showApplyModal = ref(false); const applyingTeam = ref<Team|null>(null); const applyMessage = ref(''); const applying = ref(false)
function openApplyModal(team: Team) { applyingTeam.value=team; applyMessage.value=''; showApplyModal.value=true }
async function submitApplication() {
  if (!applyingTeam.value||!applyMessage.value.trim()) return; applying.value=true
  try {
    const r = await fetch(`${apiBase}/teams/${applyingTeam.value.id}/apply`, { method:'POST', headers:{'Content-Type':'application/json',Authorization:`Bearer ${useCookie('token').value||''}`}, body: JSON.stringify({message:applyMessage.value}) })
    if (r.ok) { showToast('Aplikacja wysłana','success'); showApplyModal.value=false; applyMessage.value='' }
    else { const e=await r.json(); showToast(e.error||'Błąd','error') }
  } catch { showToast('Błąd połączenia','error') }; applying.value=false
}

function getLogoUrl(path: string): string { if (path.startsWith('http')) return path; return `http://localhost:8080${path}` }
const viewTeam = (id: string) => navigateTo(`/teams/${id}`)
const toast = reactive({ message:'', type:'' })
function showToast(m:string,t:string){toast.message=m;toast.type=t;setTimeout(()=>{toast.message=''},3000)}
</script>