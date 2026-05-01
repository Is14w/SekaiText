import { defineStore } from 'pinia'
import { useLocalStorage, usePreferredDark } from '@vueuse/core'
import { computed, ref, watch } from 'vue'

export type ThemeMode = 'system' | 'light' | 'dark'

const THEME_STORAGE_KEY = 'sekaitext-theme-mode'

export const useAppStore = defineStore('app', () => {
  const fontSize = ref(18)
  const editorMode = ref<0 | 1 | 2>(0)
  const showFlashback = ref(true)
  const syncScroll = ref(true)
  const showDiff = ref(false)
  const saveN = ref(true)
  const themeMode = useLocalStorage<ThemeMode>(THEME_STORAGE_KEY, 'system')
  const isSystemDark = usePreferredDark()
  const isDark = computed(() => themeMode.value === 'dark' || (themeMode.value === 'system' && isSystemDark.value))

  function applyTheme(dark: boolean) {
    document.documentElement.classList.toggle('dark', dark)
  }

  watch(isDark, applyTheme, { immediate: true })

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
    themeMode,
    isDark,
    setEditorMode,
  }
})
