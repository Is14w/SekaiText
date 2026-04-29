<script setup lang="ts">
import { CheckCircle, XCircle, Info, AlertTriangle } from 'lucide-vue-next'
import { useToast } from '../composables/useToast'

const { toasts } = useToast()

const iconMap: Record<string, typeof CheckCircle> = {
  success: CheckCircle,
  error: XCircle,
  info: Info,
  warn: AlertTriangle,
}
</script>

<template>
  <div class="fixed top-4 right-4 z-[9999] flex flex-col gap-2 pointer-events-none">
    <div
      v-for="t in toasts"
      :key="t.id"
      class="pointer-events-auto px-4 py-2.5 rounded-lg shadow-lg border text-sm transition-all duration-300 animate-in flex items-center gap-2"
      :class="{
        'bg-green-50 border-green-200 text-green-800 dark:bg-green-900 dark:border-green-700 dark:text-green-200': t.type === 'success',
        'bg-red-50 border-red-200 text-red-800 dark:bg-red-900 dark:border-red-700 dark:text-red-200': t.type === 'error',
        'bg-blue-50 border-blue-200 text-blue-800 dark:bg-blue-900 dark:border-blue-700 dark:text-blue-200': t.type === 'info',
        'bg-yellow-50 border-yellow-200 text-yellow-800 dark:bg-yellow-900 dark:border-yellow-700 dark:text-yellow-200': t.type === 'warn',
      }"
    >
      <component :is="iconMap[t.type]" :size="16" />
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
