<script setup lang="ts">
import { useDownloadFloat } from '../composables/useDownloadFloat'

export interface DownloadTask {
  id: number
  name: string
  status: 'pending' | 'downloading' | 'done' | 'error'
  result?: string
  error?: string
  read?: number
  total?: number
  percent?: number
}

const { tasks } = useDownloadFloat()
</script>

<template>
  <div class="fixed top-4 right-4 z-[9998] flex flex-col gap-2 pointer-events-none">
    <div
      v-for="t in tasks"
      :key="t.id"
      class="pointer-events-auto px-4 py-3 rounded-lg shadow-lg border text-sm min-w-64 transition-all duration-300 animate-in"
      :class="{
        'bg-blue-50 border-blue-200 text-blue-800 dark:bg-blue-900 dark:border-blue-700 dark:text-blue-200': t.status === 'downloading' || t.status === 'pending',
        'bg-green-50 border-green-200 text-green-800 dark:bg-green-900 dark:border-green-700 dark:text-green-200': t.status === 'done',
        'bg-red-50 border-red-200 text-red-800 dark:bg-red-900 dark:border-red-700 dark:text-red-200': t.status === 'error',
      }"
    >
      <div class="flex items-center gap-2">
        <span v-if="t.status === 'downloading' || t.status === 'pending'" class="animate-spin text-base">⏳</span>
        <span v-else-if="t.status === 'done'" class="text-base">✅</span>
        <span v-else-if="t.status === 'error'" class="text-base">❌</span>
        <div class="flex-1 min-w-0">
          <div class="font-medium truncate">{{ t.name }}</div>
          <div v-if="t.status === 'downloading'" class="mt-1">
            <div class="flex justify-between text-xs opacity-75 mb-0.5">
              <span>正在下载...</span>
              <span v-if="t.percent !== undefined">{{ t.percent }}%</span>
            </div>
            <div class="w-full h-1 bg-black/10 dark:bg-white/10 rounded-full overflow-hidden">
              <div
                class="h-full rounded-full transition-all duration-300"
                style="background-color: var(--color-primary)"
                :style="{ width: (t.percent || 0) + '%' }"
              />
            </div>
          </div>
          <div v-else-if="t.status === 'done' && t.result" class="text-xs opacity-75 mt-0.5 truncate">{{ t.result }}</div>
          <div v-else-if="t.status === 'error' && t.error" class="text-xs opacity-75 mt-0.5">{{ t.error }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
@keyframes float-in {
  from { opacity: 0; transform: translateX(20px); }
  to { opacity: 1; transform: translateX(0); }
}
.animate-in {
  animation: float-in 0.2s ease-out;
}
</style>
