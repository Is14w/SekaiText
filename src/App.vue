<template>
  <router-view />
  <Toast />
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useSettingsStore } from './stores/settings'
import { useDebugLog } from './composables/useDebugLog'
import Toast from './components/Toast.vue'

const settings = useSettingsStore()
const { enabled } = useDebugLog()

onMounted(async () => {
  try {
    await settings.fetchSettings()
  } catch {
    // backend not available, use defaults
  }
  enabled.value = settings.settings.debugEnabled
})
</script>
