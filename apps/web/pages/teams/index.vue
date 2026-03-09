<template>
  <div>
    <div class="mb-8 flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Drużyny</h1>
        <p class="mt-2 text-gray-600">Twórz i zarządzaj swoimi drużynami sportowymi</p>
      </div>
      <button
        @click="showCreateModal = true"
        class="bg-indigo-600 text-white py-2 px-4 rounded-md hover:bg-indigo-700 transition-colors"
      >
        Utwórz drużynę
      </button>
    </div>

    <!-- Teams Grid -->
    <div v-if="pending" class="flex justify-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
    </div>

    <div v-else-if="error" class="bg-red-50 text-red-700 p-4 rounded-lg">
      Nie udało się załadować drużyn. Spróbuj ponownie.
    </div>

    <div v-else-if="teams.length === 0" class="text-center py-12 bg-white rounded-lg shadow">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">Brak drużyn</h3>
      <p class="mt-1 text-sm text-gray-500">Utwórz drużynę, aby zacząć grać z przyjaciółmi.</p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div
        v-for="team in teams"
        :key="team.id"
        class="bg-white rounded-lg shadow p-6 hover:shadow-lg transition-shadow cursor-pointer"
        @click="viewTeam(team.id)"
      >
        <div class="flex items-center justify-between mb-4">
          <div class="w-12 h-12 bg-indigo-100 rounded-full flex items-center justify-center">
            <span class="text-xl font-bold text-indigo-600">{{ team.name.charAt(0) }}</span>
          </div>
          <span
            :class="[
              'px-2 py-1 rounded text-xs font-medium',
              team.userRole === 'CAPTAIN' ? 'bg-indigo-100 text-indigo-800' : 'bg-gray-100 text-gray-800'
            ]"
          >
            {{ team.userRole || 'MEMBER' }}
          </span>
        </div>
        <h3 class="text-lg font-semibold text-gray-900 mb-1">{{ team.name }}</h3>
        <p class="text-sm text-gray-600 mb-4">{{ team.description || 'Brak opisu' }}</p>
        <div class="flex items-center text-sm text-gray-500">
          <svg class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
          </svg>
          {{ team.memberCount || 1 }} członków
        </div>
      </div>
    </div>

    <!-- Create Team Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 w-full max-w-md">
        <h2 class="text-xl font-bold text-gray-900 mb-4">Utwórz nową drużynę</h2>
        <form @submit.prevent="createTeam">
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-1">Nazwa drużyny</label>
            <input
              v-model="newTeam.name"
              type="text"
              required
              class="w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
              placeholder="Wprowadź nazwę drużyny"
            />
          </div>
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-1">Opis</label>
            <textarea
              v-model="newTeam.description"
              rows="3"
              class="w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
              placeholder="Wprowadź opis drużyny"
            ></textarea>
          </div>
          <div class="flex gap-2 justify-end">
            <button
              type="button"
              @click="showCreateModal = false"
              class="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
            >
              Anuluj
            </button>
            <button
              type="submit"
              class="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700"
            >
              Utwórz drużynę
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth'
})

const showCreateModal = ref(false)
const newTeam = reactive({
  name: '',
  description: ''
})

interface Team {
  id: string
  name: string
  description?: string
  captainId: string
  userRole?: string
  memberCount?: number
}

interface TeamsResponse {
  data: Team[]
}

const { data: response, pending, error, refresh } = await useFetch<TeamsResponse>('http://localhost:8080/api/v1/teams', {
  headers: {
    Authorization: `Bearer ${useCookie('token').value || ''}`
  }
})

const teams = computed(() => response.value?.data || [])

const createTeam = async () => {
  try {
    const response = await fetch('http://localhost:8080/api/v1/teams', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${useCookie('token').value || ''}`
      },
      body: JSON.stringify(newTeam)
    })
    
    if (response.ok) {
      showCreateModal.value = false
      newTeam.name = ''
      newTeam.description = ''
      refresh()
    }
  } catch (err) {
    console.error('Failed to create team:', err)
  }
}

const viewTeam = (id: string) => {
  navigateTo(`/teams/${id}`)
}
</script>