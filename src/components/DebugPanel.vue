<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useDebugLog } from '../composables/useDebugLog'
import { useSettingsStore } from '../stores/settings'

const { logs, enabled, clear } = useDebugLog()
const settings = useSettingsStore()

onMounted(() => {
  enabled.value = settings.settings.debugEnabled
})

watch(enabled, (v) => {
  settings.settings.debugEnabled = v
})
</script>

<template>
  <div
    v-show="enabled"
    class="fixed bottom-0 right-0 w-[420px] h-60 bg-[var(--color-surface)] border border-[var(--color-border)] rounded-tl-xl shadow-2xl flex flex-col z-[9998] text-xs font-mono"
  >
    <div class="flex items-center justify-between px-3 py-1.5 border-b border-[var(--color-border)] bg-gray-50 dark:bg-gray-900 rounded-tl-xl">
      <span class="font-semibold text-[var(--color-text-secondary)]">Debug Log</span>
      <button @click="clear" class="text-[var(--color-text-secondary)] hover:text-[var(--color-text)] px-1.5 py-0.5 rounded hover:bg-gray-200 dark:hover:bg-gray-700">清空</button>
    </div>
    <div class="flex-1 overflow-y-auto p-2 space-y-0.5">
      <div v-if="logs.length === 0" class="text-[var(--color-text-secondary)] text-center pt-8">暂无日志</div>
      <div v-for="(entry, i) in logs" :key="i" class="flex gap-2 leading-5">
        <span class="text-[var(--color-text-secondary)] flex-shrink-0 w-16">{{ entry.ts }}</span>
        <span :class="{
          'text-green-600 dark:text-green-400': entry.type === 'info',
          'text-yellow-600 dark:text-yellow-400': entry.type === 'warn',
          'text-red-600 dark:text-red-400': entry.type === 'error',
        }">{{ entry.msg }}</span>
      </div>
    </div>
  </div>
</template>
