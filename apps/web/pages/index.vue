<template>
  <div>
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
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3-a2 0 11-4 0 2 0 014 0zM7 10a2 0 11-4 0 2 0 014 0z" />
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
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 0 00-2-2H5a2 0 00-2 2v6a2 0 002 2h2a2 0 002-2zm0 0V9a2 0 012-2h2a2 0 012 2v10m-6 0a2 0 002 2h2a2 0 002-2m0 0V5a2 0 012-2h2a2 0 012 2v14a2 0 01-2 2h-2a2 0 01-2-2z" />
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