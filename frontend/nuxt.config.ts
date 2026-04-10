// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  ssr: false,
  telemetry: false, // why is this true by default
  devtools: { enabled: true }
})
