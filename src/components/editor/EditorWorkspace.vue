<script setup lang="ts">
import { computed } from 'vue'
import { useAppStore } from '../../stores/app'
import { useScrollSync } from '../../composables/useScrollSync'
import SourcePanel from './SourcePanel.vue'
import DestPanel from './DestPanel.vue'

const app = useAppStore()
const { sourceRef, destRef, onSourceScroll, onDestScroll } = useScrollSync()

function handleScroll(fn: () => void) {
  if (app.syncScroll) fn()
}
</script>

<template>
  <div class="flex gap-4 h-full">
    <!-- Source Panel -->
    <div class="flex-1 min-w-0 flex flex-col">
      <div
        ref="sourceRef"
        class="flex-1 overflow-y-auto border border-[var(--color-border)] rounded-lg bg-[var(--color-surface)]"
        @scroll="handleScroll(onSourceScroll)"
      >
        <SourcePanel />
      </div>
    </div>

    <!-- Dest Panel -->
    <div class="flex-1 min-w-0 flex flex-col">
      <div
        ref="destRef"
        class="flex-1 overflow-y-auto border border-[var(--color-border)] rounded-lg bg-[var(--color-surface)]"
        @scroll="handleScroll(onDestScroll)"
      >
        <DestPanel />
      </div>
    </div>
  </div>
</template>
