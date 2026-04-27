<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useDebugLog } from '../composables/useDebugLog'

const router = useRouter()
const { logs, clear } = useDebugLog()
</script>

<template>
  <div class="min-h-screen bg-[var(--color-bg)] flex flex-col">
    <!-- Header -->
    <header data-tauri-drag-region class="border-b border-[var(--color-border)] bg-[var(--color-surface)] px-6 py-3 flex items-center justify-between">
      <div class="flex items-center gap-4">
        <button
          @click="router.push('/')"
          class="flex items-center gap-1.5 text-sm text-[var(--color-text-secondary)] hover:text-[var(--color-text)] transition-colors"
        >
          <span class="text-base leading-none">←</span>
          返回编辑器
        </button>
        <span class="text-sm font-medium">调试日志</span>
      </div>
      <button
        @click="clear"
        class="px-4 py-1.5 rounded-lg text-sm border border-[var(--color-border)] hover:bg-black/5 dark:hover:bg-white/10 transition-colors"
      >
        清空日志
      </button>
    </header>

    <!-- Content -->
    <main class="flex-1 p-6 overflow-y-auto">
      <div v-if="logs.length === 0" class="text-center text-[var(--color-text-secondary)] py-20">
        <p class="text-lg">暂无日志</p>
      </div>

      <div v-else class="space-y-1 font-mono text-sm">
        <div
          v-for="(entry, i) in logs"
          :key="i"
          class="flex gap-3 px-4 py-1.5 rounded hover:bg-[var(--color-surface)] transition-colors"
        >
          <span class="text-[var(--color-text-secondary)] flex-shrink-0 w-20">{{ entry.ts }}</span>
          <span :class="{
            'text-green-600 dark:text-green-400': entry.type === 'info',
            'text-yellow-600 dark:text-yellow-400': entry.type === 'warn',
            'text-red-600 dark:text-red-400': entry.type === 'error',
          }">{{ entry.msg }}</span>
        </div>
      </div>
    </main>
  </div>
</template>
