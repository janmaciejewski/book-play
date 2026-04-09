import { defineStore } from 'pinia'

interface User {
  id: string
  email: string
  first_name: string
  last_name: string
  role: string
}

interface AuthResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  token_type: string
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as User | null,
    token: null as string | null,
    refreshToken: null as string | null,
    loading: false,
  }),

  getters: {
    isAuthenticated: (state) => !!state.token && !!state.user,
    userRole: (state) => state.user?.role || null,
  },

  actions: {
    async register(data: { email: string; password: string; first_name: string; last_name: string }) {
      const config = useRuntimeConfig()
      
      const response = await $fetch(`${config.public.apiBase}/auth/register`, {
        method: 'POST',
        body: data,
      })
      
      return response
    },

    async login(email: string, password: string) {
      const config = useRuntimeConfig()
      
      const response = await $fetch<AuthResponse>(`${config.public.apiBase}/auth/login`, {
        method: 'POST',
        body: { email, password },
      })
      
      this.token = response.access_token
      this.refreshToken = response.refresh_token
      
      // Store in cookies for persistence (works with SSR)
      const tokenCookie = useCookie('token')
      const refreshTokenCookie = useCookie('refreshToken')
      tokenCookie.value = response.access_token
      refreshTokenCookie.value = response.refresh_token
      
      // Fetch user info
      await this.fetchUser()
    },

    async fetchUser() {
      if (!this.token) return
      
      const config = useRuntimeConfig()
      
      try {
        const response = await $fetch<{ 
          id: string
          email: string
          first_name: string
          last_name: string
          role: string 
        }>(`${config.public.apiBase}/auth/me`, {
          headers: {
            Authorization: `Bearer ${this.token}`,
          },
        })
        
        this.user = {
          id: response.id,
          email: response.email,
          first_name: response.first_name,
          last_name: response.last_name,
          role: response.role,
        }
      } catch (error) {
        this.logout()
        throw error
      }
    },

    async refresh() {
      if (!this.refreshToken) return
      
      const config = useRuntimeConfig()
      
      const response = await $fetch<AuthResponse>(`${config.public.apiBase}/auth/refresh`, {
        method: 'POST',
        body: { refresh_token: this.refreshToken },
      })
      
      this.token = response.access_token
      this.refreshToken = response.refresh_token
      
      const tokenCookie = useCookie('token')
      const refreshTokenCookie = useCookie('refreshToken')
      tokenCookie.value = response.access_token
      refreshTokenCookie.value = response.refresh_token
    },

    async logout() {
      const config = useRuntimeConfig()
      
      if (this.refreshToken) {
        try {
          await $fetch(`${config.public.apiBase}/auth/logout`, {
            method: 'POST',
          })
        } catch (error) {
          // Ignore logout errors
        }
      }
      
      this.token = null
      this.refreshToken = null
      this.user = null
      
      const tokenCookie = useCookie('token')
      const refreshTokenCookie = useCookie('refreshToken')
      tokenCookie.value = null
      refreshTokenCookie.value = null
    },

    initialize() {
      const tokenCookie = useCookie('token')
      const refreshTokenCookie = useCookie('refreshToken')
      
      if (tokenCookie.value) {
        this.token = tokenCookie.value as string
        this.refreshToken = refreshTokenCookie.value as string
        this.fetchUser()
      }
    },
  },
})