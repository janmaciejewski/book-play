<template>
  <div class="min-h-screen bg-gray-50">
    <header class="bg-white shadow-sm">
      <nav class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div class="flex h-16 justify-between">
          <div class="flex">
            <NuxtLink to="/" class="flex flex-shrink-0 items-center">
              <span class="text-2xl font-bold text-primary-600">BookPlay</span>
            </NuxtLink>
            <div class="hidden sm:ml-6 sm:flex sm:space-x-8">
              <NuxtLink
                to="/facilities"
                class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-900 hover:text-primary-600"
                active-class="border-b-2 border-primary-500 text-gray-900"
              >
                Obiekty
              </NuxtLink>
              <NuxtLink
                v-if="authStore.isAuthenticated"
                to="/reservations"
                class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-900 hover:text-primary-600"
                active-class="border-b-2 border-primary-500 text-gray-900"
              >
                Moje rezerwacje
              </NuxtLink>
              <NuxtLink
                v-if="authStore.isAuthenticated"
                to="/teams"
                class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-900 hover:text-primary-600"
                active-class="border-b-2 border-primary-500 text-gray-900"
              >
                Drużyny
              </NuxtLink>
              <NuxtLink
                v-if="authStore.user?.role === 'ADMIN'"
                to="/users"
                class="inline-flex items-center px-1 pt-1 text-sm font-medium text-gray-900 hover:text-primary-600"
                active-class="border-b-2 border-primary-500 text-gray-900"
              >
                Użytkownicy
              </NuxtLink>
            </div>
          </div>
          <div class="flex items-center">
            <template v-if="authStore.isAuthenticated">
              <div class="hidden sm:flex sm:items-center sm:space-x-4">
                <span class="text-sm text-gray-700">{{ authStore.user?.email }}</span>
                <button
                  @click="handleLogout"
                  class="rounded-md bg-white px-3 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50"
                >
                  Wyloguj
                </button>
              </div>
            </template>
            <template v-else>
              <NuxtLink
                to="/login"
                class="rounded-md bg-primary-600 px-3 py-2 text-sm font-medium text-white hover:bg-primary-700"
              >
                Zaloguj
              </NuxtLink>
            </template>
          </div>
        </div>
      </nav>
    </header>

    <main>
      <slot />
    </main>

    <footer class="bg-white border-t border-gray-200 mt-auto">
      <div class="mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8">
        <p class="text-center text-sm text-gray-500">
          &copy; {{ new Date().getFullYear() }} BookPlay. All rights reserved.
        </p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
const authStore = useAuthStore()

const handleLogout = async () => {
  await authStore.logout()
  navigateTo('/login')
}
</script>