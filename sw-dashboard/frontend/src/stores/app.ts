import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  const sidebarCollapsed = ref(false)
  const theme = ref<'dark' | 'light'>('dark')

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  function setCollapsed(val: boolean) {
    sidebarCollapsed.value = val
  }

  return { sidebarCollapsed, theme, toggleSidebar, setCollapsed }
})
