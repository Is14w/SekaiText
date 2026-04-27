<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../../api/client'
import { useEditorStore } from '../../stores/editor'

const editor = useEditorStore()
const speakers = ref<{ japanese: string; chinese: string; count: number }[]>([])
const loading = ref(false)

const emit = defineEmits<{
  close: []
}>()

onMounted(async () => {
  if (editor.talks.length === 0) return
  loading.value = true
  try {
    const result = await api.speakerCount({
      talks: editor.talks,
      sourceTalks: editor.sourceTalks,
    })
    speakers.value = result.speakers
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40" @click.self="emit('close')">
    <div class="bg-[var(--color-surface)] rounded-xl shadow-xl border border-[var(--color-border)] w-full max-w-md p-6">
      <h2 class="text-lg font-semibold mb-4">说话人统计</h2>

      <div v-if="loading" class="text-center py-8 text-sm text-[var(--color-text-secondary)]">加载中...</div>

      <div v-else class="max-h-80 overflow-y-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-[var(--color-border)]">
              <th class="text-left py-2 font-medium">说话人</th>
              <th class="text-right py-2 font-medium">台词数</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="s in speakers" :key="s.japanese" class="border-b border-[var(--color-border)]">
              <td class="py-1.5">{{ s.japanese }}</td>
              <td class="text-right py-1.5">{{ s.count }}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="flex justify-end mt-4">
        <button
          class="px-4 py-1.5 rounded text-sm border border-[var(--color-border)] hover:bg-black/5 dark:hover:bg-white/10"
          @click="emit('close')"
        >
          关闭
        </button>
      </div>
    </div>
  </div>
</template>
