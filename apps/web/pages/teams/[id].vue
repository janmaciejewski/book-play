<template>
  <div class="px-4 md:px-6 lg:px-8 mx-auto max-w-7xl py-8">
    <div v-if="pending" class="flex justify-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500"></div>
    </div>

    <div v-else-if="error" class="bg-red-50 dark:bg-red-900/50 text-red-700 dark:text-red-300 p-4 rounded-lg border border-red-200 dark:border-red-800">
      Nie udało się załadować drużyny. Spróbuj ponownie.
    </div>

    <div v-else-if="team">
      <!-- Team Header -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6 mb-6 border border-gray-200 dark:border-gray-700">
        <div class="flex justify-between items-start">
          <div>
            <h1 class="text-3xl font-bold text-gray-900 dark:text-white">{{ team.name }}</h1>
            <p v-if="team.description" class="mt-2 text-gray-600 dark:text-gray-400">{{ team.description }}</p>
          </div>
          <NuxtLink
            to="/teams"
            class="text-primary-600 dark:text-primary-400 hover:text-primary-800 dark:hover:text-primary-300"
          >
            &larr; Powrót do drużyn
          </NuxtLink>
        </div>
      </div>

      <!-- Team Members -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow border border-gray-200 dark:border-gray-700">
        <div class="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
          <h2 class="text-xl font-semibold text-gray-900 dark:text-white">Członkowie drużyny</h2>
        </div>
        <div class="divide-y divide-gray-200 dark:divide-gray-700">
          <div
            v-for="member in team.members"
            :key="member.id"
            class="px-6 py-4 flex items-center justify-between"
          >
            <div class="flex items-center">
              <div class="w-10 h-10 rounded-full bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center">
                <span class="text-primary-600 dark:text-primary-400 font-medium">
                  {{ member.user?.firstName?.[0] }}{{ member.user?.lastName?.[0] }}
                </span>
              </div>
              <div class="ml-4">
                <p class="text-sm font-medium text-gray-900 dark:text-white">
                  {{ member.user?.firstName }} {{ member.user?.lastName }}
                </p>
                <p class="text-sm text-gray-500 dark:text-gray-400">{{ member.user?.email }}</p>
              </div>
            </div>
            <span
              :class="[
                'px-2 py-1 text-xs font-medium rounded',
                member.role === 'CAPTAIN' ? 'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-300' :
                member.role === 'ADMIN' ? 'bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-300' :
                'bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300'
              ]"
            >
              {{ getRoleLabel(member.role) }}
            </span>
          </div>
        </div>
        <div v-if="!team.members || team.members.length === 0" class="px-6 py-12 text-center text-gray-500 dark:text-gray-400">
          Brak członków w drużynie
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth'
})

interface Member {
  id: string
  role: string
  user?: {
    firstName: string
    lastName: string
    email: string
  }
}

interface Team {
  id: string
  name: string
  description?: string
  members?: Member[]
}

interface TeamResponse {
  data: Team
}

const route = useRoute()
const teamId = route.params.id

const { data: response, pending, error } = await useFetch<TeamResponse>(`http://localhost:8080/api/v1/teams/${teamId}`, {
  headers: {
    Authorization: `Bearer ${useCookie('token').value || ''}`
  }
})

const team = computed(() => response.value?.data)

const getRoleLabel = (role: string): string => {
  const labels: Record<string, string> = {
    CAPTAIN: 'Kapitan',
    ADMIN: 'Administrator',
    MEMBER: 'Członek'
  }
  return labels[role] || role
}
</script>