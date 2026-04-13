<template>
  <div class="min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors duration-300 flex flex-col">
    <header class="bg-white dark:bg-gray-800 shadow-sm border-b border-gray-200 dark:border-gray-700">
      <nav class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div class="flex h-16 justify-between">
          <div class="flex">
            <NuxtLink to="/" class="flex flex-shrink-0 items-center">
              <span class="text-2xl font-bold text-primary-600">BookPlay</span>
            </NuxtLink>
            <div class="hidden sm:ml-6 sm:flex sm:space-x-8">
                <NuxtLink
                  to="/facilities"
                  class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-900 dark:text-gray-100 hover:text-primary-600 dark:hover:text-primary-400"
                  active-class="border-b-2 border-primary-500 text-gray-900 dark:text-white"
              >
                Obiekty
              </NuxtLink>
                <NuxtLink
                  v-if="authStore.isAuthenticated"
                  to="/reservations"
                  class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-900 dark:text-gray-100 hover:text-primary-600 dark:hover:text-primary-400"
                  active-class="border-b-2 border-primary-500 text-gray-900 dark:text-white"
              >
                Moje rezerwacje
              </NuxtLink>
                <NuxtLink
                  v-if="authStore.isAuthenticated"
                  to="/teams"
                  class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-900 dark:text-gray-100 hover:text-primary-600 dark:hover:text-primary-400"
                  active-class="border-b-2 border-primary-500 text-gray-900 dark:text-white"
              >
                Drużyny
              </NuxtLink>
                <NuxtLink
                  v-if="authStore.user?.role === 'ADMIN'"
                  to="/users"
                  class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-900 dark:text-gray-100 hover:text-primary-600 dark:hover:text-primary-400"
                  active-class="border-b-2 border-primary-500 text-gray-900 dark:text-white"
              >
                Użytkownicy
              </NuxtLink>
            </div>
          </div>
          <div class="flex items-center gap-3">
            <!-- Przełącznik motywu jasny/ciemny -->
            <button
              @click="themeStore.toggleTheme"
              class="p-2 rounded-lg bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
              :title="themeStore.isDark ? 'Włącz jasny motyw' : 'Włącz ciemny motyw'"
            >
              <span v-if="themeStore.isDark" class="text-xl">☀️</span>
              <span v-else class="text-xl">🌙</span>
            </button>

            <template v-if="authStore.isAuthenticated">
              <div class="hidden sm:flex sm:items-center sm:space-x-4">
                <span class="text-sm text-gray-700 dark:text-gray-300">{{ authStore.user?.email }}</span>
                <button
                  @click="handleLogout"
                  class="rounded-md bg-white dark:bg-gray-700 px-3 py-2 text-sm font-medium text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-600"
                >
                  Wyloguj
                </button>
              </div>
            </template>
            <template v-else>
              <NuxtLink
                to="/register"
                class="rounded-md border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-700 px-3 py-2 text-sm font-medium text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-600 transition-colors"
              >
                Zarejestruj
              </NuxtLink>
              <NuxtLink
                to="/login"
                class="rounded-md bg-primary-600 px-3 py-2 text-sm font-medium text-white hover:bg-primary-700 transition-colors"
              >
                Zaloguj
              </NuxtLink>
            </template>
          </div>
        </div>
      </nav>
    </header>

    <main class="flex-grow">
      <slot />
    </main>

    <footer class="bg-white dark:bg-gray-800 border-t border-gray-200 dark:border-gray-700 mt-auto">
      <div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
        <p class="text-center text-sm text-gray-500 dark:text-gray-400">
          &copy; {{ new Date().getFullYear() }} BookPlay. All rights reserved.
        </p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
const authStore = useAuthStore()
const themeStore = useThemeStore()

onMounted(() => {
  themeStore.loadFromLocalStorage()
})

const handleLogout = async () => {
  await authStore.logout()
  navigateTo('/login')
}
</script>