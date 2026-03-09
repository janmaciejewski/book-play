// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2024-04-03',
  devtools: { enabled: true },
  
  modules: [
    '@pinia/nuxt',
    '@nuxtjs/tailwindcss',
  ],

  runtimeConfig: {
    // Server-side only
    apiSecret: process.env.API_SECRET || '',
    
    // Public runtime config
    public: {
      apiBase: process.env.API_BASE_URL || 'http://localhost:8080/api/v1',
      appName: 'BookPlay',
    },
  },

  nitro: {
    preset: 'node-server',
  },

  app: {
    head: {
      title: 'BookPlay - Sports Facility Booking',
      meta: [
        { name: 'description', content: 'BookPlay - Reserve sports facilities, manage teams, and track your games' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { charset: 'utf-8' },
      ],
      link: [
        { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' },
      ],
    },
  },

  css: ['~/assets/css/main.css'],

  tailwindcss: {
    cssPath: '~/assets/css/tailwind.css',
    configPath: 'tailwind.config.ts',
  },

  pinia: {
    storesDirs: ['./stores/**'],
  },

  typescript: {
    strict: true,
  },
})