<template>
  <router-view />
  <Toast />
  <DownloadFloat />
</template>

<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useSettingsStore } from './stores/settings'
import { useDebugLog } from './composables/useDebugLog'
import Toast from './components/Toast.vue'
import DownloadFloat from './components/DownloadFloat.vue'

const settings = useSettingsStore()
const { enabled } = useDebugLog()

function applyFontSize(size: number) {
  document.documentElement.style.setProperty('--editor-font-size', size + 'px')
}

watch(() => settings.settings.fontSize, applyFontSize, { immediate: true })

onMounted(async () => {
  try {
    await settings.fetchSettings()
  } catch {
    // backend not available, use defaults
  }
  enabled.value = settings.settings.debugEnabled
  applyFontSize(settings.settings.fontSize)
})
</script>
