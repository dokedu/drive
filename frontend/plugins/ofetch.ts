import { ofetch } from 'ofetch'
import { useAuthStore } from '#imports'

export default defineNuxtPlugin((_nuxtApp) => {
  globalThis.$fetch = ofetch.create({
    onRequest({ request, options }) {
      const authStore = useAuthStore()
      if (authStore.userToken) {
        options.headers = { Authorization: `${authStore.userToken}` }
      }
    },
    onResponse({ response }) {
      if (response.status === 401) {
        useAuthStore().userToken = ''
        useAuthStore().user = {}
        navigateTo('/login')
      }
    },
  })
})