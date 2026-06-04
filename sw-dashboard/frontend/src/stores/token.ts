import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useTokenStore = defineStore('token', () => {
  const token = ref<string | null>(localStorage.getItem('dashboard_token'))

  const isAuthenticated = computed(() => token.value !== null && token.value !== '')

  function setToken(newToken: string) {
    token.value = newToken
    localStorage.setItem('dashboard_token', newToken)
  }

  function clearToken() {
    token.value = null
    localStorage.removeItem('dashboard_token')
  }

  return { token, isAuthenticated, setToken, clearToken }
})
