import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  const fontSize = ref(18)
  const editorMode = ref<0 | 1 | 2>(0)
  const showFlashback = ref(true)
  const syncScroll = ref(true)
  const showDiff = ref(false)
  const saveN = ref(true)
  const isDark = ref(false)

  function toggleTheme() {
    isDark.value = !isDark.value
    document.documentElement.classList.toggle('dark', isDark.value)
  }

  function setEditorMode(mode: 0 | 1 | 2) {
    editorMode.value = mode
  }

  return {
    fontSize,
    editorMode,
    showFlashback,
    syncScroll,
    showDiff,
    saveN,
    isDark,
    toggleTheme,
    setEditorMode,
  }
})
