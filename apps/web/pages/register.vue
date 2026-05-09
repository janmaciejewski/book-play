<template>
  <div class="flex min-h-[calc(100vh-4rem)] items-center justify-center bg-gray-50 dark:bg-gray-900 py-12 px-4 sm:px-6 lg:px-8">
    <div class="w-full max-w-md space-y-8">
      <div>
        <h2 class="mt-6 text-center text-3xl font-bold tracking-tight text-gray-900 dark:text-white">
        <h2 class="mt-6 text-center text-3xl font-bold tracking-tight text-gray-900">
          Utwórz swoje konto
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600 dark:text-gray-400">
          Lub
          <NuxtLink to="/login" class="font-medium text-primary-600 hover:text-primary-500">
            zaloguj się na istniejące konto
          </NuxtLink>
        </p>
        <h2 class="mt-6 text-center text-3xl font-bold tracking-tight text-gray-900 dark:text-white">Utwórz swoje konto</h2>
        <p class="mt-2 text-center text-sm text-gray-600 dark:text-gray-400">Lub <NuxtLink to="/login" class="font-medium text-primary-600 hover:text-primary-500">zaloguj się na istniejące konto</NuxtLink></p>
      </div>
      <form class="mt-8 space-y-6" @submit.prevent="handleRegister">
        <div class="space-y-4 rounded-md shadow-sm">
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label for="first_name" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Imię</label>
              <input
                id="first_name"
                v-model="form.first_name"
                name="first_name"
                type="text"
                required
              class="relative block w-full rounded-md border-0 py-1.5 bg-white dark:bg-gray-800 text-gray-900 dark:text-white ring-1 ring-inset ring-gray-300 dark:ring-gray-600 placeholder:text-gray-400 dark:placeholder:text-gray-500 focus:ring-2 focus:ring-inset focus:ring-primary-600 sm:text-sm sm:leading-6"
                placeholder="Imię"
              />
            </div>
            <div>
              <label for="last_name" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Nazwisko</label>
              <input
                id="last_name"
                v-model="form.last_name"
                name="last_name"
                type="text"
                required
              class="relative block w-full rounded-md border-0 py-1.5 bg-white dark:bg-gray-800 text-gray-900 dark:text-white ring-1 ring-inset ring-gray-300 dark:ring-gray-600 placeholder:text-gray-400 dark:placeholder:text-gray-500 focus:ring-2 focus:ring-inset focus:ring-primary-600 sm:text-sm sm:leading-6"
                placeholder="Nazwisko"
              />
            </div>
          </div>
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Adres email</label>
            <label for="email" class="block text-sm font-medium text-gray-700 mb-1">Adres email</label>
            <input
              id="email"
              v-model="form.email"
              name="email"
              type="email"
              autocomplete="email"
              required
              class="relative block w-full rounded-md border-0 py-1.5 bg-white dark:bg-gray-800 text-gray-900 dark:text-white ring-1 ring-inset ring-gray-300 dark:ring-gray-600 placeholder:text-gray-400 dark:placeholder:text-gray-500 focus:ring-2 focus:ring-inset focus:ring-primary-600 sm:text-sm sm:leading-6"
              placeholder="Adres email"
            />
            <label for="email" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Adres email</label>
            <div class="flex gap-2">
              <input id="email" v-model="form.email" type="email" autocomplete="email" required class="w-full rounded-md border-0 py-1.5 bg-white dark:bg-gray-800 text-gray-900 dark:text-white ring-1 ring-inset ring-gray-300 dark:ring-gray-600 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-primary-600 sm:text-sm sm:leading-6" placeholder="Adres email" />
              <button type="button" @click="sendOTP" :disabled="!form.email||otpSent||sendingOTP" class="px-3 py-1.5 text-sm font-medium rounded-md border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50 whitespace-nowrap">{{ otpSent?'Wysłano ✓':sendingOTP?'Wysyłanie...':'Wyślij kod' }}</button>
            </div>
          </div>
          <div v-if="otpSent">
            <label for="otp" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Kod weryfikacyjny</label>
            <div class="flex gap-2">
              <input id="otp" v-model="form.otp" type="text" maxlength="6" class="w-full rounded-md border-0 py-1.5 bg-white dark:bg-gray-800 text-gray-900 dark:text-white ring-1 ring-inset ring-gray-300 dark:ring-gray-600 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-primary-600 sm:text-sm sm:leading-6" placeholder="6-cyfrowy kod" />
              <button type="button" @click="verifyOTP" :disabled="!form.otp||form.otp.length!==6||verifyingOTP||otpVerified" class="px-3 py-1.5 text-sm font-medium rounded-md border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50 whitespace-nowrap">{{ otpVerified?'✓':verifyingOTP?'Sprawdzanie...':'Zweryfikuj' }}</button>
            </div>
          </div>
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Hasło</label>
            <label for="password" class="block text-sm font-medium text-gray-700 mb-1">Hasło</label>
            <input
              id="password"
              v-model="form.password"
              name="password"
              type="password"
              autocomplete="new-password"
              required
              minlength="8"
              class="relative block w-full rounded-md border-0 py-1.5 bg-white dark:bg-gray-800 text-gray-900 dark:text-white ring-1 ring-inset ring-gray-300 dark:ring-gray-600 placeholder:text-gray-400 dark:placeholder:text-gray-500 focus:ring-2 focus:ring-inset focus:ring-primary-600 sm:text-sm sm:leading-6"
              placeholder="Hasło (min. 8 znaków)"
            />
          </div>
        </div>

        <div v-if="error" class="rounded-md bg-red-50 dark:bg-red-900/50 p-4 border border-red-200 dark:border-red-800">
          <p class="text-sm text-red-700 dark:text-red-300">{{ error }}</p>
        </div>

            <label for="password" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Hasło</label>
            <input id="password" v-model="form.password" type="password" autocomplete="new-password" required minlength="8" class="w-full rounded-md border-0 py-1.5 bg-white dark:bg-gray-800 text-gray-900 dark:text-white ring-1 ring-inset ring-gray-300 dark:ring-gray-600 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-primary-600 sm:text-sm sm:leading-6" placeholder="Hasło (min. 8 znaków)" />
          </div>
        </div>
        <div v-if="success" class="rounded-md bg-green-50 dark:bg-green-900/50 p-4 border border-green-200 dark:border-green-800"><p class="text-sm text-green-700 dark:text-green-300">{{ success }}</p></div>
        <div v-if="error" class="rounded-md bg-red-50 dark:bg-red-900/50 p-4 border border-red-200 dark:border-red-800"><p class="text-sm text-red-700 dark:text-red-300">{{ error }}</p></div>
        <div>
          <button type="submit" :disabled="loading||!otpVerified" class="group relative flex w-full justify-center rounded-md bg-primary-600 px-3 py-2 text-sm font-semibold text-white hover:bg-primary-700 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 disabled:opacity-50 disabled:cursor-not-allowed">
            <span v-if="loading">Tworzenie konta...</span>
            <span v-else-if="!otpSent">Wyślij kod, aby kontynuować</span>
            <span v-else-if="!otpVerified">Zweryfikuj kod, aby kontynuować</span>
            <span v-else>Utwórz konto</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({layout:'default'})
const apiBase=useRuntimeConfig().public.apiBase; const authStore=useAuthStore(); const router=useRouter()
const form=reactive({email:'',password:'',first_name:'',last_name:'',otp:''})
const loading=ref(false); const error=ref(''); const success=ref(''); const otpSent=ref(false); const otpVerified=ref(false); const sendingOTP=ref(false); const verifyingOTP=ref(false)

async function sendOTP(){sendingOTP.value=true;error.value='';success.value='';try{const r=await fetch(`${apiBase}/auth/send-otp`,{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({email:form.email})});const d=await r.json();if(r.ok){otpSent.value=true;success.value='Kod wysłany na Twój email.'}else{error.value=d.error||'Błąd'}}catch{error.value='Błąd połączenia'}sendingOTP.value=false}
async function verifyOTP(){verifyingOTP.value=true;error.value='';success.value='';try{const r=await fetch(`${apiBase}/auth/verify-otp`,{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({email:form.email,code:form.otp})});const d=await r.json();if(r.ok){otpVerified.value=true;success.value='Email zweryfikowany!'}else{error.value=d.error||'Nieprawidłowy kod'}}catch{error.value='Błąd połączenia'}verifyingOTP.value=false}
const handleRegister=async()=>{if(!otpVerified.value){error.value='Zweryfikuj email';return};loading.value=true;error.value='';try{await authStore.register({email:form.email,password:form.password,first_name:form.first_name,last_name:form.last_name});router.push('/login')}catch(e:any){error.value=e.data?.error||'Nie udało się utworzyć konta.'}finally{loading.value=false}}
</script>