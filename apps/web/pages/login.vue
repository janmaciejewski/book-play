<template>
  <div class="flex min-h-[calc(100vh-4rem)] items-center justify-center bg-gray-50 dark:bg-gray-900 py-12 px-4 sm:px-6 lg:px-8">
    <div class="w-full max-w-md space-y-8">
      <div>
        <h2 class="mt-6 text-center text-3xl font-bold tracking-tight text-gray-900 dark:text-white">
          Zaloguj się na swoje konto
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600 dark:text-gray-400">
          Lub
          <NuxtLink to="/register" class="font-medium text-primary-600 hover:text-primary-500">
            utwórz nowe konto
          </NuxtLink>
        </p>
      </div>
      
      <form class="mt-8 space-y-6" @submit.prevent="handleLogin">
        <div class="-space-y-px rounded-md shadow-sm">
          <div>
            <label for="email" class="sr-only">Adres email</label>
            <input
              id="email"
              v-model="form.email"
              name="email"
              type="email"
              autocomplete="email"
              required
              class="relative block w-full rounded-t-md border-0 py-1.5 bg-white dark:bg-gray-800 text-gray-900 dark:text-white ring-1 ring-inset ring-gray-300 dark:ring-gray-600 placeholder:text-gray-400 dark:placeholder:text-gray-500 focus:z-10 focus:ring-2 focus:ring-inset focus:ring-primary-600 sm:text-sm sm:leading-6"
              placeholder="Adres email"
            />
          </div>
          <div>
            <label for="password" class="sr-only">Hasło</label>
            <input
              id="password"
              v-model="form.password"
              name="password"
              type="password"
              autocomplete="current-password"
              required
              class="relative block w-full rounded-b-md border-0 py-1.5 bg-white dark:bg-gray-800 text-gray-900 dark:text-white ring-1 ring-inset ring-gray-300 dark:ring-gray-600 placeholder:text-gray-400 dark:placeholder:text-gray-500 focus:z-10 focus:ring-2 focus:ring-inset focus:ring-primary-600 sm:text-sm sm:leading-6"
              placeholder="Hasło"
            />
          </div>
        </div>

        <div class="flex justify-end">
          <a href="#" @click.prevent="showForgotPassword = true" class="text-sm font-medium text-primary-600 hover:text-primary-500">
            Zapomniałem hasła
          </a>
        </div>

        <!-- Forgot Password Modal -->
        <div v-if="showForgotPassword" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div class="bg-white dark:bg-gray-800 rounded-lg p-6 w-full max-w-md">
            <h2 class="text-xl font-bold text-gray-900 dark:text-white mb-4">Resetowanie hasła</h2>
            <p v-if="forgotPasswordStep === 'email'" class="text-sm text-gray-600 dark:text-gray-400 mb-4">
              Podaj swój adres email, a wyślemy Ci kod weryfikacyjny.
            </p>
            <p v-else class="text-sm text-gray-600 dark:text-gray-400 mb-4">
              Wpisz kod z emaila i ustaw nowe hasło.
            </p>
            <form @submit.prevent="handleForgotPassword" class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Email</label>
                <input
                  v-model="forgotPasswordEmail"
                  type="email"
                  required
                  :disabled="forgotPasswordStep === 'reset'"
                  class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500 disabled:opacity-50"
                  placeholder="twoj@email.com"
                />
              </div>
              <template v-if="forgotPasswordStep === 'reset'">
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Kod weryfikacyjny</label>
                  <input
                    v-model="forgotPasswordCode"
                    type="text"
                    required
                    maxlength="6"
                    class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500"
                    placeholder="000000"
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Nowe hasło</label>
                  <input
                    v-model="forgotPasswordNewPassword"
                    type="password"
                    required
                    minlength="6"
                    class="w-full rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500"
                    placeholder="Minimum 6 znaków"
                  />
                </div>
              </template>
              <p v-if="forgotPasswordMessage" :class="['text-sm', forgotPasswordMessage.includes('zmienione') ? 'text-green-600 dark:text-green-400' : 'text-gray-600 dark:text-gray-400']">{{ forgotPasswordMessage }}</p>
              <div class="flex gap-2 justify-end">
                <button type="button" @click="showForgotPassword = false; forgotPasswordStep = 'email'; forgotPasswordCode = ''; forgotPasswordNewPassword = ''" class="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-md text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 text-sm">Anuluj</button>
                <button type="submit" :disabled="forgotPasswordLoading" class="px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 disabled:opacity-50 text-sm">
                  {{ forgotPasswordLoading ? 'Wysyłanie...' : forgotPasswordStep === 'email' ? 'Wyślij kod' : 'Zmień hasło' }}
                </button>
              </div>
            </form>
          </div>
        </div>

        <div v-if="error" class="rounded-md bg-red-50 dark:bg-red-900/50 p-4 border border-red-200 dark:border-red-800">
          <p class="text-sm text-red-700 dark:text-red-300">{{ error }}</p>
        </div>

        <div>
          <button
            type="submit"
            :disabled="loading"
            class="group relative flex w-full justify-center rounded-md bg-primary-600 px-3 py-2 text-sm font-semibold text-white hover:bg-primary-700 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="loading">Logowanie...</span>
            <span v-else>Zaloguj się</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'default',
})

const authStore = useAuthStore()
const router = useRouter()

const form = reactive({
  email: '',
  password: '',
})

const loading = ref(false)
const error = ref('')

const showForgotPassword = ref(false)
const forgotPasswordEmail = ref('')
const forgotPasswordLoading = ref(false)
const forgotPasswordMessage = ref('')
const forgotPasswordStep = ref<'email' | 'reset'>('email')
const forgotPasswordCode = ref('')
const forgotPasswordNewPassword = ref('')

const handleForgotPassword = async () => {
  forgotPasswordLoading.value = true
  forgotPasswordMessage.value = ''

  try {
    if (forgotPasswordStep.value === 'email') {
      // Step 1: Send OTP
      await $fetch('http://localhost:8080/api/v1/auth/send-otp', {
        method: 'POST',
        body: { email: forgotPasswordEmail.value }
      })
      forgotPasswordMessage.value = 'Kod weryfikacyjny został wysłany na Twój email.'
      forgotPasswordStep.value = 'reset'
    } else {
      // Step 2: Reset password with OTP
      await $fetch('http://localhost:8080/api/v1/auth/reset-password', {
        method: 'POST',
        body: {
          email: forgotPasswordEmail.value,
          code: forgotPasswordCode.value,
          password: forgotPasswordNewPassword.value
        }
      })
      forgotPasswordMessage.value = 'Hasło zostało zmienione! Możesz się teraz zalogować.'
      setTimeout(() => {
        showForgotPassword.value = false
        forgotPasswordStep.value = 'email'
        forgotPasswordCode.value = ''
        forgotPasswordNewPassword.value = ''
      }, 2000)
    }
  } catch (e: any) {
    forgotPasswordMessage.value = e.data?.error || 'Wystąpił błąd. Spróbuj ponownie.'
  } finally {
    forgotPasswordLoading.value = false
  }
}

const handleLogin = async () => {
  loading.value = true
  error.value = ''

  try {
    await authStore.login(form.email, form.password)
    router.push('/')
  } catch (e: any) {
    error.value = e.data?.error || 'Nie udało się zalogować. Sprawdź swoje dane.'
  } finally {
    loading.value = false
  }
}
</script>