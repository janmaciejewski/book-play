<template>
<<<<<<< Updated upstream
  <div>
    <div v-if="pending" class="flex justify-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
    </div>

    <div v-else-if="error" class="bg-red-50 text-red-700 p-4 rounded-lg">
      Nie udało się załadować drużyny. Spróbuj ponownie.
    </div>

    <div v-else-if="team">
      <!-- Team Header -->
      <div class="bg-white rounded-lg shadow p-6 mb-6">
        <div class="flex justify-between items-start">
          <div>
            <h1 class="text-3xl font-bold text-gray-900">{{ team.name }}</h1>
            <p v-if="team.description" class="mt-2 text-gray-600">{{ team.description }}</p>
          </div>
          <NuxtLink
            to="/teams"
            class="text-indigo-600 hover:text-indigo-800"
          >
            &larr; Powrót do drużyn
          </NuxtLink>
        </div>
      </div>

      <!-- Team Members -->
      <div class="bg-white rounded-lg shadow">
        <div class="px-6 py-4 border-b border-gray-200">
          <h2 class="text-xl font-semibold text-gray-900">Członkowie drużyny</h2>
        </div>
        <div class="divide-y divide-gray-200">
          <div
            v-for="member in team.members"
            :key="member.id"
            class="px-6 py-4 flex items-center justify-between"
          >
            <div class="flex items-center">
              <div class="w-10 h-10 rounded-full bg-indigo-100 flex items-center justify-center">
                <span class="text-indigo-600 font-medium">
                  {{ member.user?.firstName?.[0] }}{{ member.user?.lastName?.[0] }}
                </span>
              </div>
              <div class="ml-4">
                <p class="text-sm font-medium text-gray-900">
                  {{ member.user?.firstName }} {{ member.user?.lastName }}
                </p>
                <p class="text-sm text-gray-500">{{ member.user?.email }}</p>
              </div>
            </div>
            <span
              :class="[
                'px-2 py-1 text-xs font-medium rounded',
                member.role === 'CAPTAIN' ? 'bg-yellow-100 text-yellow-800' :
                member.role === 'ADMIN' ? 'bg-blue-100 text-blue-800' :
                'bg-gray-100 text-gray-800'
              ]"
            >
              {{ getRoleLabel(member.role) }}
            </span>
          </div>
        </div>
        <div v-if="!team.members || team.members.length === 0" class="px-6 py-12 text-center text-gray-500">
          Brak członków w drużynie
=======
  <div class="px-4 md:px-6 lg:px-8 mx-auto max-w-7xl py-8">
    <div v-if="pending" class="flex justify-center py-12"><div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500"></div></div>
    <div v-else-if="error" class="bg-red-50 dark:bg-red-900/50 text-red-700 dark:text-red-300 p-4 rounded-lg border border-red-200 dark:border-red-800">Nie udało się załadować drużyny.</div>
    <div v-else-if="team">
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 mb-6 border border-gray-200 dark:border-gray-700">
        <div class="flex flex-col md:flex-row justify-between items-start gap-4">
          <div class="flex items-center gap-4">
            <div class="relative flex-shrink-0">
              <div v-if="team.logo" class="w-20 h-20 rounded-lg overflow-hidden bg-gray-100 dark:bg-gray-700"><img :src="getLogoUrl(team.logo)" :alt="team.name" class="w-full h-full object-cover" /></div>
              <div v-else class="w-20 h-20 rounded-lg bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center"><span class="text-3xl font-bold text-primary-600 dark:text-primary-400">{{ team.name.charAt(0) }}</span></div>
              <button v-if="canManage" @click="triggerLogoUpload" class="absolute -bottom-1 -right-1 bg-primary-600 text-white rounded-full p-1 hover:bg-primary-700" title="Zmień logo">
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" /></svg>
              </button>
              <input ref="logoInput" type="file" accept="image/jpeg,image/png,image/webp,image/svg+xml" class="hidden" @change="handleLogoUpload" />
            </div>
            <div>
              <div class="flex items-center gap-2 flex-wrap">
                <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ team.name }}</h1>
                <button v-if="canManage" @click="toggleEditInfo" class="text-gray-400 hover:text-primary-600 dark:hover:text-primary-400" title="Edytuj informacje">
                  <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>
                </button>
              </div>
              <div v-if="editingInfo" class="mt-3 p-4 bg-gray-50 dark:bg-gray-700/50 rounded-lg border border-gray-200 dark:border-gray-600">
                <div class="mb-3"><label class="block text-sm font-medium mb-1">Nazwa</label><input v-model="editForm.name" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500" /></div>
                <div class="mb-3"><label class="block text-sm font-medium mb-1">Opis</label><textarea v-model="editForm.description" rows="3" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500"></textarea></div>
                <div class="flex gap-2"><button @click="saveInfo" :disabled="saving" class="px-3 py-1.5 bg-primary-600 text-white rounded-md hover:bg-primary-700 text-sm disabled:opacity-50">{{ saving?'Zapisywanie...':'Zapisz' }}</button><button @click="editingInfo=false" class="px-3 py-1.5 border border-gray-300 dark:border-gray-600 rounded-md text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 text-sm">Anuluj</button></div>
              </div>
              <p v-if="team.description" class="mt-2 text-gray-600 dark:text-gray-400">{{ team.description }}</p>
              <p v-else class="mt-2 text-gray-400 dark:text-gray-500 italic">Brak opisu</p>
            </div>
          </div>
          <NuxtLink to="/teams" class="text-primary-600 dark:text-primary-400 hover:text-primary-800 dark:hover:text-primary-300 flex-shrink-0">&larr; Powrót do drużyn</NuxtLink>
        </div>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div class="lg:col-span-2">
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow border border-gray-200 dark:border-gray-700">
            <div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex justify-between items-center">
              <h2 class="text-xl font-semibold text-gray-900 dark:text-white">Członkowie drużyny<span class="text-sm font-normal text-gray-500 dark:text-gray-400 ml-2">({{ members.length }})</span></h2>
              <button v-if="canManage" @click="showAddMemberModal=true" class="bg-primary-600 text-white px-3 py-1.5 rounded-md hover:bg-primary-700 text-sm">+ Dodaj członka</button>
            </div>
            <div v-if="members.length===0" class="px-6 py-12 text-center text-gray-500 dark:text-gray-400">Brak członków</div>
            <div v-else class="divide-y divide-gray-200 dark:divide-gray-700">
              <div v-for="member in members" :key="member.id" class="px-6 py-4 flex items-center justify-between hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
                <div class="flex items-center">
                  <div class="w-10 h-10 rounded-full bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center overflow-hidden"><span class="text-primary-600 dark:text-primary-400 font-medium">{{ member.user?.firstName?.[0] }}{{ member.user?.lastName?.[0] }}</span></div>
                  <div class="ml-4">
                    <p class="text-sm font-medium text-gray-900 dark:text-white">{{ member.user?.firstName }} {{ member.user?.lastName }}</p>
                    <NuxtLink :to="`/users/${member.user?.id}`" class="text-sm text-primary-600 dark:text-primary-400 hover:underline">{{ member.user?.email }}</NuxtLink>
                  </div>
                </div>
                <div class="flex items-center gap-3">
                  <div v-if="canManage && member.role!=='CAPTAIN'"><select :value="member.role" @change="changeRole(member,($event.target as HTMLSelectElement).value)" class="text-xs font-medium rounded px-2 py-1 border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:border-primary-500 focus:ring-primary-500"><option value="ADMIN">Administrator</option><option value="MEMBER">Członek</option><option value="CAPTAIN">Kapitan (przekaż)</option></select></div>
                  <span v-else :class="['px-2 py-1 text-xs font-medium rounded',member.role==='CAPTAIN'?'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-300':member.role==='ADMIN'?'bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-300':'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300']">{{ getRoleLabel(member.role) }}</span>
                  <button v-if="canManage && member.role!=='CAPTAIN'" @click="confirmRemove(member)" class="text-red-500 hover:text-red-700 dark:hover:text-red-400 p-1" title="Usuń">
                    <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="space-y-6">
          <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 border border-gray-200 dark:border-gray-700">
            <h3 class="font-semibold text-gray-900 dark:text-white mb-3">Informacje</h3>
            <div class="space-y-3 text-sm">
              <div><span class="text-gray-500 dark:text-gray-400">Twój status:</span><span :class="['ml-2 px-2 py-0.5 text-xs font-medium rounded',myMembership?.role==='CAPTAIN'?'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-300':myMembership?.role==='ADMIN'?'bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-300':'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300']">{{ getRoleLabel(myMembership?.role||'') }}</span></div>
              <div><span class="text-gray-500 dark:text-gray-400">Liczba członków:</span><span class="ml-2 text-gray-900 dark:text-white font-medium">{{ members.length }}</span></div>
              <div v-if="team.created_at"><span class="text-gray-500 dark:text-gray-400">Utworzono:</span><span class="ml-2 text-gray-900 dark:text-white">{{ formatDate(team.created_at) }}</span></div>
            </div>
          </div>
          <div v-if="myMembership" class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 border border-gray-200 dark:border-gray-700">
            <button @click="leaveTeam" class="w-full px-4 py-2 border border-red-300 dark:border-red-600 text-red-700 dark:text-red-400 rounded-md hover:bg-red-50 dark:hover:bg-red-900/20 text-sm">{{ myMembership.role==='CAPTAIN'?'Opuść drużynę (przekaż kapitana)':'Opuść drużynę' }}</button>
          </div>
          <div v-if="canManage" class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 border border-gray-200 dark:border-gray-700">
            <h3 class="font-semibold text-gray-900 dark:text-white mb-3">Rekrutacja</h3>
            <div class="flex items-center justify-between mb-3"><span class="text-sm text-gray-600 dark:text-gray-400">Otwarta</span><button @click="toggleRecruitment" :disabled="togglingRecruitment" :class="['relative inline-flex h-6 w-11 items-center rounded-full transition-colors',team.recruitment_open?'bg-green-600':'bg-gray-300 dark:bg-gray-600']"><span :class="['inline-block h-4 w-4 transform rounded-full bg-white transition-transform',team.recruitment_open?'translate-x-6':'translate-x-1']" /></button></div>
            <button @click="loadApplications" class="w-full px-3 py-1.5 text-xs font-medium border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 rounded-md hover:bg-gray-50 dark:hover:bg-gray-700 mb-3">Pokaż aplikacje</button>
            <div v-if="applications.length>0" class="space-y-3 max-h-60 overflow-y-auto">
              <div v-for="app in applications" :key="app.id" class="p-3 bg-gray-50 dark:bg-gray-700/50 rounded-lg border border-gray-200 dark:border-gray-600">
                <div class="flex items-center justify-between mb-1"><p class="text-sm font-medium text-gray-900 dark:text-white">{{ app.user?.first_name }} {{ app.user?.last_name }}</p><span :class="['px-2 py-0.5 text-xs rounded',app.status==='PENDING'?'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-300':app.status==='ACCEPTED'?'bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300':'bg-red-100 dark:bg-red-900/30 text-red-800 dark:text-red-300']">{{ app.status==='PENDING'?'Oczekuje':app.status==='ACCEPTED'?'Zaakceptowany':'Odrzucony' }}</span></div>
                <p class="text-xs text-gray-500 dark:text-gray-400 break-words">{{ app.message }}</p>
                <div v-if="app.status==='PENDING'" class="flex gap-2 mt-2"><button @click="handleApplication(app.id,'ACCEPTED')" class="px-2 py-1 text-xs bg-green-600 text-white rounded hover:bg-green-700">Akceptuj</button><button @click="handleApplication(app.id,'REJECTED')" class="px-2 py-1 text-xs bg-red-600 text-white rounded hover:bg-red-700">Odrzuć</button></div>
              </div>
            </div>
            <p v-else-if="showApplications" class="text-sm text-gray-500 dark:text-gray-400 text-center">Brak aplikacji</p>
          </div>
>>>>>>> Stashed changes
        </div>
      </div>

      <!-- Add Member Modal -->
      <div v-if="showAddMemberModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-md">
          <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-4">Dodaj członka</h2>
          <div class="mb-4"><label class="block text-sm font-medium mb-1">Szukaj użytkownika (email)</label>
            <div class="relative"><input v-model="searchQuery" type="text" placeholder="Email..." class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500" @input="searchUsers" />
              <div v-if="searchResults.length>0" class="absolute z-10 mt-1 w-full bg-white dark:bg-gray-700 border border-gray-200 dark:border-gray-600 rounded-md shadow-lg max-h-48 overflow-y-auto">
                <div v-for="user in searchResults" :key="user.id" @click="selectUser(user)" class="px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 cursor-pointer flex items-center justify-between"><div><p class="text-sm font-medium">{{ user.first_name }} {{ user.last_name }}</p><p class="text-xs text-gray-500">{{ user.email }}</p></div><svg class="w-5 h-5 text-primary-600" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" /></svg></div>
              </div>
            </div>
          </div>
          <div v-if="selectedUser" class="mb-4 p-3 bg-gray-50 dark:bg-gray-700/50 rounded-lg border border-gray-200 dark:border-gray-600"><p class="text-sm font-medium">{{ selectedUser.first_name }} {{ selectedUser.last_name }}</p><p class="text-xs text-gray-500">{{ selectedUser.email }}</p></div>
          <div class="mb-4"><label class="block text-sm font-medium mb-1">Rola</label><select v-model="addMemberRole" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500"><option value="MEMBER">Członek</option><option value="ADMIN">Administrator</option><option value="CAPTAIN">Kapitan</option></select></div>
          <div class="flex gap-2 justify-end"><button @click="showAddMemberModal=false;selectedUser=null;searchQuery=''" class="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700">Anuluj</button><button @click="addMember" :disabled="!selectedUser||addingMember" class="px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 disabled:opacity-50">{{ addingMember?'Dodawanie...':'Dodaj' }}</button></div>
        </div>
      </div>

      <!-- Confirm Remove -->
      <div v-if="showRemoveConfirm" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
        <div class="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-sm"><h2 class="text-lg font-bold mb-2">Usuń członka</h2><p class="text-gray-600 dark:text-gray-400 mb-4">Usunąć <strong>{{ memberToRemove?.user?.firstName }} {{ memberToRemove?.user?.lastName }}</strong>?</p>
          <div class="flex gap-2 justify-end"><button @click="showRemoveConfirm=false;memberToRemove=null" class="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700">Anuluj</button><button @click="removeMember" :disabled="removingMember" class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 disabled:opacity-50">{{ removingMember?'Usuwanie...':'Usuń' }}</button></div>
        </div>
      </div>

      <div v-if="toast.message" :class="['fixed bottom-4 right-4 px-4 py-3 rounded-lg shadow-lg z-50 text-sm',toast.type==='success'?'bg-green-500 text-white':'bg-red-500 text-white']">{{ toast.message }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({middleware:'auth'})
const apiBase = useRuntimeConfig().public.apiBase

interface User { id:string;first_name:string;last_name:string;email:string }
interface Member { id:string;role:string;user?:User }
interface Team { id:string;name:string;description?:string;logo?:string;captain_id:string;recruitment_open?:boolean;members?:Member[];created_at?:string }
interface TeamResponse { data:Team }
interface SearchResponse { data:User[] }
interface Application { id:string;user_id:string;message:string;status:string;created_at:string;user?:User }
interface ApplicationsResponse { data:Application[] }

const route=useRoute(); const teamId=route.params.id as string; const authStore=useAuthStore(); const currentUserId=computed(()=>authStore.user?.id)
const { data:response, pending, error, refresh }=await useFetch<TeamResponse>(`${apiBase}/teams/${teamId}`,{headers:{Authorization:`Bearer ${useCookie('token').value||''}`}})
const team=computed(()=>response.value?.data); const members=computed(()=>team.value?.members||[])
const myMembership=computed(()=>{if(!currentUserId.value||!members.value)return null;return members.value.find(m=>m.user?.id===currentUserId.value)})
const canManage=computed(()=>{const r=myMembership.value?.role;return r==='CAPTAIN'||r==='ADMIN'})

const editingInfo=ref(false); const saving=ref(false); const editForm=reactive({name:'',description:''})
function toggleEditInfo(){if(team.value){editForm.name=team.value.name;editForm.description=team.value.description||''}editingInfo.value=!editingInfo.value}
async function saveInfo(){saving.value=true;try{const r=await fetch(`${apiBase}/teams/${teamId}`,{method:'PUT',headers:{'Content-Type':'application/json',Authorization:`Bearer ${useCookie('token').value||''}`},body:JSON.stringify({name:editForm.name,description:editForm.description||null})});if(r.ok){editingInfo.value=false;showToast('Zaktualizowano','success');refresh()}else{const e=await r.json();showToast(e.error||'Błąd','error')}}catch{showToast('Błąd połączenia','error')}saving.value=false}

const logoInput=ref<HTMLInputElement|null>(null)
function triggerLogoUpload(){logoInput.value?.click()}
async function handleLogoUpload(e:Event){const input=e.target as HTMLInputElement;if(!input.files?.length)return;const file=input.files[0];if(!['image/jpeg','image/png','image/webp','image/svg+xml'].includes(file.type)){showToast('Nieobsługiwany format','error');return};if(file.size>5*1024*1024){showToast('Zbyt duży plik (max 5MB)','error');return};const fd=new FormData();fd.append('logo',file);try{const r=await fetch(`${apiBase}/teams/${teamId}/logo`,{method:'POST',headers:{Authorization:`Bearer ${useCookie('token').value||''}`},body:fd});if(r.ok){showToast('Logo zaktualizowane','success');refresh()}else{const e=await r.json();showToast(e.error||'Błąd','error')}}catch{showToast('Błąd połączenia','error')}input.value=''}
function getLogoUrl(path:string):string{if(path.startsWith('http'))return path;return`http://localhost:8080${path}`}

const showAddMemberModal=ref(false);const searchQuery=ref('');const searchResults=ref<User[]>([]);const selectedUser=ref<User|null>(null);const addMemberRole=ref('MEMBER');const addingMember=ref(false);let searchTimeout:ReturnType<typeof setTimeout>|null=null
async function searchUsers(){if(searchTimeout)clearTimeout(searchTimeout);if(searchQuery.value.length<2){searchResults.value=[];return};searchTimeout=setTimeout(async()=>{try{const r=await fetch(`${apiBase}/teams/${teamId}/search-users?q=${encodeURIComponent(searchQuery.value)}`,{headers:{Authorization:`Bearer ${useCookie('token').value||''}`}});if(r.ok){const d:SearchResponse=await r.json();searchResults.value=d.data||[]}}catch{}},400)}
function selectUser(user:User){selectedUser.value=user;searchQuery.value=user.email;searchResults.value=[]}
async function addMember(){if(!selectedUser.value)return;addingMember.value=true;try{const r=await fetch(`${apiBase}/teams/${teamId}/members`,{method:'POST',headers:{'Content-Type':'application/json',Authorization:`Bearer ${useCookie('token').value||''}`},body:JSON.stringify({email:selectedUser.value.email,role:addMemberRole.value})});if(r.ok){showToast('Dodano','success');showAddMemberModal.value=false;selectedUser.value=null;searchQuery.value='';addMemberRole.value='MEMBER';refresh()}else{const e=await r.json();showToast(e.error||'Błąd','error')}}catch{showToast('Błąd połączenia','error')}addingMember.value=false}

const showRemoveConfirm=ref(false);const memberToRemove=ref<Member|null>(null);const removingMember=ref(false)
function confirmRemove(member:Member){memberToRemove.value=member;showRemoveConfirm.value=true}
async function removeMember(){if(!memberToRemove.value)return;removingMember.value=true;try{const r=await fetch(`${apiBase}/teams/${teamId}/members/${memberToRemove.value.id}`,{method:'DELETE',headers:{Authorization:`Bearer ${useCookie('token').value||''}`}});if(r.ok){showToast('Usunięto','success');showRemoveConfirm.value=false;memberToRemove.value=null;refresh()}else{const e=await r.json();showToast(e.error||'Błąd','error')}}catch{showToast('Błąd połączenia','error')}removingMember.value=false}

async function changeRole(member:Member,newRole:string){if(newRole==='CAPTAIN'&&!confirm(`Przekazać kapitana?`))return;try{const r=await fetch(`${apiBase}/teams/${teamId}/members/${member.id}/role`,{method:'PUT',headers:{'Content-Type':'application/json',Authorization:`Bearer ${useCookie('token').value||''}`},body:JSON.stringify({role:newRole})});if(r.ok){showToast('Rola zmieniona','success');refresh()}else{const e=await r.json();showToast(e.error||'Błąd','error')}}catch{showToast('Błąd połączenia','error')}}
async function leaveTeam(){if(!myMembership.value)return;if(!confirm('Opuścić drużynę?'))return;try{const r=await fetch(`${apiBase}/teams/${teamId}/members/${myMembership.value.id}`,{method:'DELETE',headers:{Authorization:`Bearer ${useCookie('token').value||''}`}});if(r.ok){showToast('Opuściłeś drużynę','success');navigateTo('/teams')}else{const e=await r.json();showToast(e.error||'Błąd','error')}}catch{showToast('Błąd połączenia','error')}}

const togglingRecruitment=ref(false);const showApplications=ref(false);const applications=ref<Application[]>([])
async function toggleRecruitment(){togglingRecruitment.value=true;try{const ns=!team.value?.recruitment_open;const r=await fetch(`${apiBase}/teams/${teamId}/recruitment`,{method:'PUT',headers:{'Content-Type':'application/json',Authorization:`Bearer ${useCookie('token').value||''}`},body:JSON.stringify({open:ns})});if(r.ok){showToast(ns?'Otwarta':'Zamknięta','success');refresh()}}catch{showToast('Błąd','error')}togglingRecruitment.value=false}
async function loadApplications(){showApplications.value=true;try{const r=await fetch(`${apiBase}/teams/${teamId}/applications`,{headers:{Authorization:`Bearer ${useCookie('token').value||''}`}});if(r.ok){const d:ApplicationsResponse=await r.json();applications.value=d.data||[]}}catch{}}
async function handleApplication(appId:string,status:string){try{const r=await fetch(`${apiBase}/teams/${teamId}/applications/${appId}`,{method:'PUT',headers:{'Content-Type':'application/json',Authorization:`Bearer ${useCookie('token').value||''}`},body:JSON.stringify({status})});if(r.ok){showToast(status==='ACCEPTED'?'Zaakceptowano':'Odrzucono','success');loadApplications();refresh()}else{const e=await r.json();showToast(e.error||'Błąd','error')}}catch{showToast('Błąd połączenia','error')}}

const toast=reactive({message:'',type:''})
function showToast(m:string,t:string){toast.message=m;toast.type=t;setTimeout(()=>{toast.message=''},3000)}
function getRoleLabel(r:string):string{const l:Record<string,string>={CAPTAIN:'Kapitan',ADMIN:'Administrator',MEMBER:'Członek'};return l[r]||r}
function formatDate(d:string):string{if(!d)return'';return new Date(d).toLocaleDateString('pl-PL',{year:'numeric',month:'long',day:'numeric'})}
</script>