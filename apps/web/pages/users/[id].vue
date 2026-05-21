<template>
  <div class="px-4 md:px-6 lg:px-8 mx-auto max-w-7xl py-8">
    <div v-if="pending" class="flex justify-center py-12"><div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500"></div></div>
    <div v-else-if="error" class="bg-red-50 dark:bg-red-900/50 text-red-700 dark:text-red-300 p-4 rounded-lg">Nie udało się załadować profilu.</div>
    <div v-else-if="user">
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow border border-gray-200 dark:border-gray-700 mb-6">
        <div class="p-6 border-b border-gray-200 dark:border-gray-700">
          <div class="flex items-center gap-4">
            <div class="relative flex-shrink-0">
              <img v-if="user.avatar" :src="user.avatar" class="w-16 h-16 rounded-full object-cover" />
              <div v-else class="w-16 h-16 rounded-full bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center"><span class="text-2xl font-bold text-primary-600 dark:text-primary-400">{{ user.first_name?.[0] }}{{ user.last_name?.[0] }}</span></div>
              <button v-if="isOwnProfile" @click="triggerAvatarUpload" class="absolute -bottom-1 -right-1 bg-primary-600 text-white rounded-full p-1 hover:bg-primary-700" title="Zmień awatar"><svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" /></svg></button>
              <input ref="avatarInput" type="file" accept="image/jpeg,image/png,image/webp,image/svg+xml" class="hidden" @change="handleAvatarUpload" />
            </div>
            <div>
              <h1 class="text-2xl font-bold text-gray-900 dark:text-white">{{ user.first_name }} {{ user.last_name }}</h1>
              <p class="text-gray-500 dark:text-gray-400">{{ user.email }}</p>
              <span :class="['mt-1 inline-block px-2 py-0.5 text-xs font-medium rounded',user.role==='ADMIN'?'bg-purple-100 dark:bg-purple-900/30 text-purple-800 dark:text-purple-300':user.role==='FACILITY_OWNER'?'bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300':'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300']">{{ getRoleLabel(user.role) }}</span>
            </div>
            <div class="ml-auto"><NuxtLink to="/" class="text-primary-600 dark:text-primary-400 hover:underline text-sm">&larr; Powrót</NuxtLink></div>
          </div>
        </div>
        <div class="p-6">
          <!-- FACILITY_OWNER and ADMIN see only basic details -->
          <template v-if="user.role === 'FACILITY_OWNER' || user.role === 'ADMIN'">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div><h3 class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase mb-1">Email</h3><p class="text-gray-900 dark:text-white">{{ user.email }}</p></div>
              <div><h3 class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase mb-1">Telefon</h3><p class="text-gray-900 dark:text-white">{{ user.phone||'—' }}</p></div>
            </div>
          </template>
          <!-- Other roles see full details -->
          <template v-else>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div><h3 class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase mb-1">Email</h3><p class="text-gray-900 dark:text-white">{{ user.email }}</p></div>
              <div><h3 class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase mb-1">Telefon</h3><p class="text-gray-900 dark:text-white">{{ user.phone||'—' }}</p></div>
              <div><h3 class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase mb-1">Miasto</h3><p class="text-gray-900 dark:text-white">{{ user.city||'—' }}</p></div>
              <div><h3 class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase mb-1">Kraj</h3><p class="text-gray-900 dark:text-white">{{ user.country||'—' }}</p></div>
              <div><h3 class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase mb-1">Wiek</h3><p class="text-gray-900 dark:text-white">{{ user.age||'—' }}</p></div>
              <div><h3 class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase mb-1">Preferowana pozycja</h3><p class="text-gray-900 dark:text-white">{{ user.preferred_position||'—' }}</p></div>
              <div><h3 class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase mb-1">Aktywny</h3><p><span :class="user.is_active?'text-green-600 dark:text-green-400':'text-red-600 dark:text-red-400'">{{ user.is_active?'Tak':'Nie' }}</span></p></div>
            </div>
            <div v-if="user.bio" class="mt-6"><h3 class="text-sm font-medium text-gray-500 dark:text-gray-400 uppercase mb-1">O mnie</h3><p class="text-gray-900 dark:text-white whitespace-pre-wrap">{{ user.bio }}</p></div>
          </template>
          <div v-if="isOwnProfile" class="mt-6 pt-4 border-t border-gray-200 dark:border-gray-700"><button @click="showEditModal=true" class="px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 text-sm">Edytuj profil</button></div>
        </div>
      </div>
    </div>

    <div v-if="showEditModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-lg max-h-[90vh] overflow-y-auto">
        <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-4">Edytuj profil</h2>
        <form @submit.prevent="saveProfile" class="space-y-4">
          <div><label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Telefon</label><input v-model="editProfile.phone" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500" /></div>
          <template v-if="user?.role !== 'FACILITY_OWNER' && user?.role !== 'ADMIN'">
            <div><label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Miasto</label><input v-model="editProfile.city" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500" /></div>
            <div><label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Kraj</label><input v-model="editProfile.country" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500" /></div>
            <div><label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Wiek</label><input v-model.number="editProfile.age" type="number" min="0" max="120" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500" /></div>
            <div><label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Preferowana pozycja</label><input v-model="editProfile.preferred_position" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500" /></div>
            <div><label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">O mnie</label><textarea v-model="editProfile.bio" rows="4" class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500" placeholder="Kilka słów o sobie..."></textarea></div>
          </template>
          <div class="flex gap-2 justify-end"><button type="button" @click="showEditModal=false" class="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700">Anuluj</button><button type="submit" :disabled="saving" class="px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 disabled:opacity-50">{{ saving?'Zapisywanie...':'Zapisz' }}</button></div>
        </form>
      </div>
    </div>
    <div v-if="toast.message" :class="['fixed bottom-4 right-4 px-4 py-3 rounded-lg shadow-lg z-50 text-sm',toast.type==='success'?'bg-green-500 text-white':'bg-red-500 text-white']">{{ toast.message }}</div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({middleware:'auth'})
const apiBase=useRuntimeConfig().public.apiBase; const route=useRoute(); const userId=route.params.id as string; const authStore=useAuthStore()
const isOwnProfile=computed(()=>authStore.user?.id===userId)
interface UserData { id:string;email:string;first_name:string;last_name:string;phone?:string;avatar?:string;role:string;bio?:string;city?:string;country?:string;position?:string;preferred_position?:string;age?:number;is_active:boolean }
interface UserResponse { data:UserData }
const {data:response,pending,error,refresh}=await useFetch<UserResponse>(`${apiBase}/users/${userId}`,{headers:{Authorization:`Bearer ${useCookie('token').value||''}`}})
const user=computed(()=>response.value?.data)
const showEditModal=ref(false); const saving=ref(false); const avatarInput=ref<HTMLInputElement|null>(null)
function triggerAvatarUpload() { avatarInput.value?.click() }
const editProfile=reactive({phone:'',city:'',country:'',age:null as number|null,position:'',preferred_position:'',bio:''})
watch(showEditModal,(s)=>{if(s&&user.value){editProfile.phone=user.value.phone||'';editProfile.city=user.value.city||'';editProfile.country=user.value.country||'';editProfile.age=user.value.age||null;editProfile.position=user.value.position||'';editProfile.preferred_position=user.value.preferred_position||'';editProfile.bio=user.value.bio||''}})
async function saveProfile(){saving.value=true;try{const r=await fetch(`${apiBase}/users/${userId}`,{method:'PUT',headers:{'Content-Type':'application/json',Authorization:`Bearer ${useCookie('token').value||''}`},body:JSON.stringify({phone:editProfile.phone||null,city:editProfile.city||null,country:editProfile.country||null,age:editProfile.age||null,position:editProfile.position||null,preferred_position:editProfile.preferred_position||null,bio:editProfile.bio||null})});if(r.ok){showToast('Zaktualizowano','success');showEditModal.value=false;refresh()}else{const e=await r.json();showToast(e.error||'Błąd','error')}}catch{showToast('Błąd połączenia','error')}saving.value=false}
const toast=reactive({message:'',type:''})
function showToast(m:string,t:string){toast.message=m;toast.type=t;setTimeout(()=>{toast.message=''},3000)}
function getRoleLabel(r:string):string{const l:Record<string,string>={ADMIN:'Administrator',FACILITY_OWNER:'Właściciel obiektu',PLAYER:'Gracz'};return l[r]||r}

async function handleAvatarUpload(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  const formData = new FormData()
  formData.append('avatar', file)
  try {
    const baseUrl = useRuntimeConfig().public.apiBase
    const token = useCookie('token').value || ''
    const r = await fetch(`${baseUrl}/users/${userId}/avatar`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: formData
    })
    if (r.ok) {
      showToast('Awatar zaktualizowany', 'success')
      refresh()
    } else {
      const e = await r.json()
      showToast(e.error || 'Błąd', 'error')
    }
  } catch {
    showToast('Błąd połączenia', 'error')
  }
}
</script>