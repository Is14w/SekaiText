<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { api } from '../../api/client'
import { useEditorStore } from '../../stores/editor'

const editor = useEditorStore()
const speakers = ref<{ japanese: string; chinese: string }[]>([])
const loading = ref(false)

const emit = defineEmits<{
  close: []
  save: []
}>()

onMounted(async () => {
  if (editor.talks.length === 0) return
  loading.value = true
  try {
    const result = await api.speakerCount({
      talks: editor.talks,
      sourceTalks: editor.sourceTalks,
    })
    speakers.value = result.speakers.map(s => ({
      japanese: s.japanese,
      chinese: s.chinese || s.japanese,
    }))
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/40" @click.self="emit('close')">
    <div class="bg-[var(--color-surface)] rounded-xl shadow-xl border border-[var(--color-border)] w-full max-w-lg p-6">
      <h2 class="text-lg font-semibold mb-4">批量修改说话人</h2>

      <div v-if="loading" class="text-center py-8 text-sm text-[var(--color-text-secondary)]">加载中...</div>

      <div v-else class="max-h-80 overflow-y-auto">
        <table class="w-full text-sm">
          <thead>
            <tr class="border-b border-[var(--color-border)]">
              <th class="text-left py-2 font-medium">日文/原翻译</th>
              <th class="text-left py-2 font-medium">翻译</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(s, idx) in speakers" :key="idx" class="border-b border-[var(--color-border)]">
              <td class="py-1.5 text-[var(--color-text-secondary)]">{{ s.japanese }}</td>
              <td class="py-1.5">
                <input
                  v-model="s.chinese"
                  class="w-full px-2 py-1 rounded border border-[var(--color-border)] bg-[var(--color-bg)] text-sm"
                />
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="flex justify-end gap-2 mt-4">
        <button
          class="px-4 py-1.5 rounded text-sm border border-[var(--color-border)] hover:bg-black/5 dark:hover:bg-white/10"
          @click="emit('close')"
        >
          取消
        </button>
        <button
          class="px-4 py-1.5 rounded text-sm text-white"
          style="background-color: var(--color-primary)"
          @click="emit('save')"
        >
          应用
        </button>
      </div>
    </div>
  </div>
</template>
