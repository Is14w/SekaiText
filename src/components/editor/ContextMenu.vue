<script setup lang="ts">
const props = defineProps<{
  x: number
  y: number
}>()

const emit = defineEmits<{
  select: [brackets: string]
  close: []
}>()

const bracketOptions = [
  { label: '替换为「」', value: '「」' },
  { label: '替换为『』', value: '『』' },
  { label: '替换为（）', value: '（）' },
  { label: '替换为“”', value: '“”' },
  { label: '替换为‘’', value: '‘’' },
]
</script>

<template>
  <div
    class="fixed z-50 bg-[var(--color-surface)] border border-[var(--color-border)] rounded-lg shadow-xl py-1 min-w-[160px]"
    :style="{ left: x + 'px', top: y + 'px' }"
    @click.self="emit('close')"
  >
    <button
      v-for="opt in bracketOptions"
      :key="opt.value"
      class="w-full text-left px-3 py-1.5 text-sm hover:text-[var(--color-primary)] transition-colors"
      @click="emit('select', opt.value)"
    >
      {{ opt.label }}
    </button>
  </div>
</template>
