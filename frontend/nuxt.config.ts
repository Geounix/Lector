export default defineNuxtConfig({
  app: {
    head: {
      title: 'Lector Comics',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'Self-hosted comic reader' }
      ]
    }
  },

  css: ['~/assets/css/main.css'],

  modules: ['@pinia/nuxt'],

  runtimeConfig: {
    public: {
      apiBase: process.env.API_URL || 'http://localhost:3000/api/v1'
    }
  },

  devtools: { enabled: true }
})