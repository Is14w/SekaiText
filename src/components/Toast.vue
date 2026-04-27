<script setup lang="ts">
import { useToast } from '../composables/useToast'

const { toasts } = useToast()
</script>

<template>
  <div class="fixed top-4 right-4 z-[9999] flex flex-col gap-2 pointer-events-none">
    <div
      v-for="t in toasts"
      :key="t.id"
      class="pointer-events-auto px-4 py-2.5 rounded-lg shadow-lg border text-sm transition-all duration-300 animate-in"
      :class="{
        'bg-green-50 border-green-200 text-green-800 dark:bg-green-900/30 dark:border-green-800 dark:text-green-300': t.type === 'success',
        'bg-red-50 border-red-200 text-red-800 dark:bg-red-900/30 dark:border-red-800 dark:text-red-300': t.type === 'error',
        'bg-blue-50 border-blue-200 text-blue-800 dark:bg-blue-900/30 dark:border-blue-800 dark:text-blue-300': t.type === 'info',
      }"
    >
      {{ t.message }}
    </div>
  </div>
</template>

<style scoped>
@keyframes toast-in {
  from { opacity: 0; transform: translateX(20px); }
  to { opacity: 1; transform: translateX(0); }
}
.animate-in {
  animation: toast-in 0.2s ease-out;
}
</style>
