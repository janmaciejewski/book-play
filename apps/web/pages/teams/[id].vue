<template>
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