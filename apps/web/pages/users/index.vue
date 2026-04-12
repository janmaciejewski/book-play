<template>
  <div class="px-4 md:px-6 lg:px-8 mx-auto max-w-7xl py-8">
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Użytkownicy</h1>
      <p class="mt-2 text-gray-600 dark:text-gray-400">Zarządzaj użytkownikami systemu</p>
    </div>

    <!-- Users Table -->
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow border border-gray-200 dark:border-gray-700">
      <div v-if="pending" class="flex justify-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500"></div>
      </div>

      <div v-else-if="error" class="p-4 text-red-700 dark:text-red-300">
        Nie udało się załadować użytkowników.
      </div>

      <table v-else class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
        <thead class="bg-gray-50 dark:bg-gray-700">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Użytkownik</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Email</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Rola</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Status</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-300 uppercase tracking-wider">Akcje</th>
          </tr>
        </thead>
        <tbody class="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
          <tr v-for="user in users" :key="user.id">
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="flex items-center">
                <div class="w-10 h-10 rounded-full bg-primary-100 dark:bg-primary-900/30 flex items-center justify-center">
                  <span class="text-primary-600 dark:text-primary-400 font-medium">
                    {{ user.first_name?.[0] }}{{ user.last_name?.[0] }}
                  </span>
                </div>
                <div class="ml-4">
                  <div class="text-sm font-medium text-gray-900 dark:text-white">{{ user.first_name }} {{ user.last_name }}</div>
                  <div v-if="user.phone" class="text-sm text-gray-500 dark:text-gray-400">{{ user.phone }}</div>
                </div>
              </div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-300">{{ user.email }}</td>
            <td class="px-6 py-4 whitespace-nowrap">
              <select
                v-model="user.role"
                @change="updateRole(user.id, user.role)"
                class="text-sm rounded-md bg-white dark:bg-gray-700 border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white shadow-sm focus:border-primary-500 focus:ring-primary-500"
              >
                <option value="ADMIN">Administrator</option>
                <option value="FACILITY_OWNER">Właściciel obiektu</option>
                <option value="PLAYER">Gracz</option>
              </select>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span
                :class="[
                  'px-2 py-1 text-xs font-medium rounded',
                  user.is_active ? 'bg-green-100 dark:bg-green-900/30 text-green-800 dark:text-green-300' : 'bg-red-100 dark:bg-red-900/30 text-red-800 dark:text-red-300'
                ]"
              >
                {{ user.is_active ? 'Aktywny' : 'Nieaktywny' }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm">
              <button
                @click="deleteUser(user.id)"
                class="text-red-600 dark:text-red-400 hover:text-red-900 dark:hover:text-red-300"
                :disabled="user.id === authStore.user?.id"
                :class="{ 'opacity-50 cursor-not-allowed': user.id === authStore.user?.id }"
              >
                Usuń
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
definePageMeta({
  middleware: 'auth'
})

const authStore = useAuthStore()

// Check if user is admin
if (authStore.user?.role !== 'ADMIN') {
  navigateTo('/')
}

interface User {
  id: string
  email: string
  first_name: string
  last_name: string
  phone?: string
  role: string
  is_active: boolean
}

interface UsersResponse {
  data: User[]
}

const { data: response, pending, error, refresh } = await useFetch<UsersResponse>('http://localhost:8080/api/v1/users', {
  headers: {
    Authorization: `Bearer ${useCookie('token').value || ''}`
  }
})

const users = computed(() => response.value?.data || [])

const updateRole = async (userId: string, newRole: string) => {
  try {
    await $fetch(`http://localhost:8080/api/v1/users/${userId}/role`, {
      method: 'PUT',
      headers: {
        Authorization: `Bearer ${useCookie('token').value || ''}`
      },
      body: { role: newRole }
    })
  } catch (err) {
    console.error('Failed to update role:', err)
    refresh()
  }
}

const deleteUser = async (userId: string) => {
  if (!confirm('Czy na pewno chcesz usunąć tego użytkownika?')) return
  
  try {
    await $fetch(`http://localhost:8080/api/v1/users/${userId}`, {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${useCookie('token').value || ''}`
      }
    })
    refresh()
  } catch (err) {
    console.error('Failed to delete user:', err)
  }
}
</script>