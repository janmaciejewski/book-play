<template>
  <div>
    <!-- ===== DASHBOARD: FACILITY_OWNER ===== -->
    <template v-if="authStore.user?.role === 'FACILITY_OWNER'">
      <div class="px-4 md:px-6 lg:px-8 mx-auto max-w-7xl py-8">
        <div class="mb-8">
          <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Panel właściciela</h1>
          <p class="mt-2 text-gray-600 dark:text-gray-400">Witaj, {{ authStore.user.first_name }}! Zarządzaj swoimi obiektami i rezerwacjami.</p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <NuxtLink to="/facilities/owner" class="bg-white dark:bg-gray-800 rounded-lg shadow border border-gray-200 dark:border-gray-700 p-6 hover:shadow-lg hover:scale-[1.02] transition-all">
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 rounded-lg bg-green-100 dark:bg-green-900/30 flex items-center justify-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-green-600 dark:text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                </svg>
              </div>
              <div>
                <h2 class="text-xl font-semibold text-gray-900 dark:text-white">Moje obiekty</h2>
                <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">Przeglądaj i zarządzaj swoimi obiektami sportowymi</p>
              </div>
            </div>
          </NuxtLink>

          <NuxtLink to="/facilities/my-reservations" class="bg-white dark:bg-gray-800 rounded-lg shadow border border-gray-200 dark:border-gray-700 p-6 hover:shadow-lg hover:scale-[1.02] transition-all">
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 rounded-lg bg-blue-100 dark:bg-blue-900/30 flex items-center justify-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-blue-600 dark:text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              </div>
              <div>
                <h2 class="text-xl font-semibold text-gray-900 dark:text-white">Rezerwacje</h2>
                <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">Zobacz i zarządzaj rezerwacjami swoich obiektów</p>
              </div>
            </div>
          </NuxtLink>
        </div>
      </div>
    </template>

    <!-- ===== DASHBOARD: ADMIN ===== -->
    <template v-else-if="authStore.user?.role === 'ADMIN'">
      <div class="px-4 md:px-6 lg:px-8 mx-auto max-w-7xl py-8">
        <div class="mb-8">
          <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Panel administratora</h1>
          <p class="mt-2 text-gray-600 dark:text-gray-400">Witaj, {{ authStore.user.first_name }}! Zarządzaj platformą BookPlay.</p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <NuxtLink to="/facilities/owner" class="bg-white dark:bg-gray-800 rounded-lg shadow border border-gray-200 dark:border-gray-700 p-6 hover:shadow-lg hover:scale-[1.02] transition-all">
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 rounded-lg bg-green-100 dark:bg-green-900/30 flex items-center justify-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-green-600 dark:text-green-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                </svg>
              </div>
              <div>
                <h2 class="text-xl font-semibold text-gray-900 dark:text-white">Obiekty</h2>
                <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">Przeglądaj wszystkie obiekty sportowe</p>
              </div>
            </div>
          </NuxtLink>

          <NuxtLink to="/reservations" class="bg-white dark:bg-gray-800 rounded-lg shadow border border-gray-200 dark:border-gray-700 p-6 hover:shadow-lg hover:scale-[1.02] transition-all">
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 rounded-lg bg-blue-100 dark:bg-blue-900/30 flex items-center justify-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-blue-600 dark:text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              </div>
              <div>
                <h2 class="text-xl font-semibold text-gray-900 dark:text-white">Rezerwacje</h2>
                <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">Zarządzaj wszystkimi rezerwacjami</p>
              </div>
            </div>
          </NuxtLink>

          <NuxtLink to="/users" class="bg-white dark:bg-gray-800 rounded-lg shadow border border-gray-200 dark:border-gray-700 p-6 hover:shadow-lg hover:scale-[1.02] transition-all">
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 rounded-lg bg-purple-100 dark:bg-purple-900/30 flex items-center justify-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-purple-600 dark:text-purple-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
                </svg>
              </div>
              <div>
                <h2 class="text-xl font-semibold text-gray-900 dark:text-white">Użytkownicy</h2>
                <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">Zarządzaj użytkownikami i rolami</p>
              </div>
            </div>
          </NuxtLink>

          <NuxtLink to="/teams" class="bg-white dark:bg-gray-800 rounded-lg shadow border border-gray-200 dark:border-gray-700 p-6 hover:shadow-lg hover:scale-[1.02] transition-all">
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 rounded-lg bg-yellow-100 dark:bg-yellow-900/30 flex items-center justify-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 text-yellow-600 dark:text-yellow-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                </svg>
              </div>
              <div>
                <h2 class="text-xl font-semibold text-gray-900 dark:text-white">Drużyny</h2>
                <p class="text-sm text-gray-500 dark:text-gray-400 mt-1">Przeglądaj wszystkie drużyny</p>
              </div>
            </div>
          </NuxtLink>
        </div>
      </div>
    </template>

    <!-- ===== DEFAULT: Regular user homepage with slideshow ===== -->
    <template v-else>
      <!-- Hero Section -->
      <section class="relative overflow-hidden h-[600px] flex items-center justify-center">
        <!-- Background Slideshow -->
        <div 
          v-for="(image, index) in heroImages" 
          :key="index"
          class="absolute inset-0 transition-opacity duration-1000 ease-in-out"
          :class="{
            'opacity-100': currentSlide === index,
            'opacity-0': currentSlide !== index
          }"
        >
          <img 
            :src="image" 
            alt="Sport background"
            class="w-full h-full object-cover object-center"
          />
        </div>

        <!-- Dark Overlay for text readability -->
        <div class="absolute inset-0 bg-black/60 z-10"></div>

        <!-- Hero Content -->
        <div class="relative z-20 mx-auto max-w-7xl px-4 py-16 sm:px-6 lg:px-8">
          <div class="text-center">
            <h1 class="text-4xl font-bold tracking-tight text-white sm:text-5xl md:text-6xl">
              Zarezerwuj swój <span class="text-primary-400">mecz</span>
            </h1>
            <p class="mx-auto mt-3 max-w-md text-base text-gray-200 sm:text-lg md:mt-5 md:max-w-3xl md:text-xl">
              Rezerwuj obiekty sportowe, zarządzaj drużynami i śledź swoje mecze z łatwością.
            </p>
            <div class="mx-auto mt-5 max-w-md sm:flex sm:justify-center md:mt-8 gap-3">
              <NuxtLink
                to="/facilities"
                class="rounded-md shadow bg-primary-600 px-6 py-2.5 text-sm font-medium text-white hover:bg-primary-700 transition-colors whitespace-nowrap"
              >
                Przeglądaj obiekty
              </NuxtLink>
              <NuxtLink
                :to="authStore.isAuthenticated ? '/reservations' : '/register'"
                class="rounded-md shadow bg-white/10 backdrop-blur-sm text-white px-6 py-2.5 text-sm font-medium hover:bg-white/20 transition-colors whitespace-nowrap"
              >
                {{ authStore.isAuthenticated ? 'Sprawdź swoje rezerwacje' : 'Rozpocznij' }}
              </NuxtLink>
            </div>
          </div>
        </div>
      </section>

      <!-- Features Section -->
      <section class="bg-gray-50 dark:bg-gray-800 py-16">
        <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div class="text-center">
            <h2 class="text-3xl font-bold text-gray-900 dark:text-white">Wszystko, czego potrzebujesz do gry</h2>
            <p class="mt-4 text-lg text-gray-500 dark:text-gray-400">
              Zarządzaj swoimi aktywnościami sportowymi z jednej wygodnej platformy.
            </p>
          </div>

          <div class="mt-12 grid gap-8 md:grid-cols-3">
            <!-- Feature 1 -->
            <div class="rounded-lg bg-white dark:bg-gray-700 p-6 shadow border border-gray-200 dark:border-gray-600">
              <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-md bg-primary-500 text-white">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              </div>
              <h3 class="text-lg font-semibold text-gray-900 dark:text-white">Łatwa rezerwacja</h3>
              <p class="mt-2 text-gray-500 dark:text-gray-400">
                Rezerwuj korty, boiska i obiekty w kilku kliknięciach. Dostępność w czasie rzeczywistym gwarantuje brak konfliktów.
              </p>
            </div>

            <!-- Feature 2 -->
            <div class="rounded-lg bg-white dark:bg-gray-700 p-6 shadow border border-gray-200 dark:border-gray-600">
              <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-md bg-primary-500 text-white">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                </svg>
              </div>
              <h3 class="text-lg font-semibold text-gray-900 dark:text-white">Zarządzanie drużyną</h3>
              <p class="mt-2 text-gray-500 dark:text-gray-400">
                Twórz drużyny, zapraszaj członków i koordynujcie mecze razem. Utrzymujcie wszystkich w kontakcie.
              </p>
            </div>

            <!-- Feature 3 -->
            <div class="rounded-lg bg-white dark:bg-gray-700 p-6 shadow border border-gray-200 dark:border-gray-600">
              <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-md bg-primary-500 text-white">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                </svg>
              </div>
              <h3 class="text-lg font-semibold text-gray-900 dark:text-white">Śledzenie aktywności</h3>
              <p class="mt-2 text-gray-500 dark:text-gray-400">
                Śledź swoje rezerwacje i historię meczów. Przeglądaj nadchodzące i przeszłe aktywności.
              </p>
            </div>
          </div>
        </div>
      </section>
    </template>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  layout: 'default',
})

const authStore = useAuthStore()

const heroImages = [
  'https://images.unsplash.com/photo-1461896836934-ffe607ba8211?ixlib=rb-4.0.3&auto=format&fit=crop&w=1920&q=80',
  'https://images.unsplash.com/photo-1546519638-68e109498ffc?ixlib=rb-4.0.3&auto=format&fit=crop&w=1920&q=80',
  'https://images.unsplash.com/photo-1579952363873-27f3bade9f55?ixlib=rb-4.0.3&auto=format&fit=crop&w=1920&q=80',
  'https://images.unsplash.com/photo-1622279457486-62dcc4a431d6?ixlib=rb-4.0.3&auto=format&fit=crop&w=1920&q=80',
  'https://images.unsplash.com/photo-1517649763962-0c623066013b?ixlib=rb-4.0.3&auto=format&fit=crop&w=1920&q=80'
]

const currentSlide = ref(0)

onMounted(() => {
  setInterval(() => {
    currentSlide.value = (currentSlide.value + 1) % heroImages.length
  }, 3000)
})
</script>