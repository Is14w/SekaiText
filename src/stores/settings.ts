import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '../api/client'
import type { Settings } from '../types/api'

export const useSettingsStore = defineStore('settings', () => {
  const settings = ref<Settings>({
    fontSize: 18,
    downloadSource: 'haruki',
    saveN: true,
    saveVoice: false,
    disableSSL: false,
    debugEnabled: false,
    jsonDownloadDir: './downloads/json',

    indexOrder: 'asc',
  })
  const loading = ref(false)

  async function fetchSettings() {
    loading.value = true
    try {
      settings.value = await api.getSettings()
    } finally {
      loading.value = false
    }
  }

  async function saveSettings() {
    await api.putSettings(settings.value)
  }

  return { settings, loading, fetchSettings, saveSettings }
})
