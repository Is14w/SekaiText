<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft } from 'lucide-vue-next'
import { useStoryStore } from '../stores/story'
import { useSettingsStore } from '../stores/settings'
import { api } from '../api/client'
import { useToast } from '../composables/useToast'
import { useDownloadFloat } from '../composables/useDownloadFloat'

const router = useRouter()
const story = useStoryStore()
const settings = useSettingsStore()
const toast = useToast()
const dlFloat = useDownloadFloat()

const outputDir = ref(settings.settings.jsonDownloadDir || '')

const displayIndices = computed(() => {
  if (settings.settings.indexOrder === 'desc') {
    return [...story.indices].reverse()
  }
  return story.indices
})

onMounted(async () => {
  await story.fetchTypes()
})

watch(() => story.selectedType, async (type) => {
  if (!type) return
  story.selectedSort = ''
  story.selectedIndex = ''
  story.selectedChapter = -1
  story.sorts = []
  story.indices = []
  story.chapters = []
  await story.fetchSorts(type)
  if (story.sorts.length === 0) story.fetchIndex(type, '')
})

watch(() => story.selectedSort, async (sort) => {
  if (!sort || !story.selectedType) return
  story.selectedIndex = ''
  story.selectedChapter = -1
  story.indices = []
  story.chapters = []
  await story.fetchIndex(story.selectedType, sort)
})

watch(() => story.selectedIndex, async (idx) => {
  if (!idx || !story.selectedType) return
  story.selectedChapter = -1
  story.chapters = []
  if (idx !== '-') {
    await story.fetchChapters(story.selectedType, story.selectedSort, idx)
  }
})

async function handleDownload() {
  if (!story.selectedType || !story.selectedIndex || story.selectedChapter < 0 || !outputDir.value) {
    toast.show('请填写所有字段', 'warn')
    return
  }
  const taskId = dlFloat.add(story.selectedIndex + ' ch' + story.selectedChapter)
  dlFloat.start(taskId)
  try {
    const { taskId: tid } = await api.downloadJson({
      storyType: story.selectedType,
      sort: story.selectedSort,
      index: story.selectedIndex,
      chapter: story.selectedChapter,
      source: 'haruki',
      outputDir: outputDir.value,
    })
    let done = false
    while (!done) {
      await new Promise(r => setTimeout(r, 300))
      const prog = await api.downloadProgress(tid)
      if (prog.total > 0) {
        dlFloat.progress(taskId, prog.read, prog.total, Math.round((prog.read / prog.total) * 100))
      }
      if (prog.status === 'done') {
        dlFloat.done(taskId, prog.filePath || '')
        done = true
      } else if (prog.status === 'error') {
        dlFloat.fail(taskId, prog.error || '下载失败')
        done = true
      }
    }
  } catch (e: any) {
    dlFloat.fail(taskId, e.message || '下载失败')
  }
}
</script>

<template>
  <div class="min-h-screen bg-[var(--color-bg)]">
    <header class="border-b border-[var(--color-border)] bg-[var(--color-surface)] px-6 py-3 flex items-center justify-between">
      <div class="flex items-center gap-4">
        <button
          @click="router.push('/')"
          class="flex items-center gap-1.5 text-sm text-[var(--color-text-secondary)] hover:text-[var(--color-text)] transition-colors"
        >
          <ArrowLeft :size="18" />
          返回编辑器
        </button>
        <span class="text-sm font-medium">JSON 下载</span>
      </div>
    </header>

    <main class="max-w-2xl mx-auto p-6">
      <div class="bg-[var(--color-surface)] border border-[var(--color-border)] rounded-xl p-6 space-y-5">
        <h2 class="text-base font-semibold">选择故事</h2>

        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-xs text-[var(--color-text-secondary)] mb-1">故事类型</label>
            <select v-model="story.selectedType" class="w-full px-2 py-1.5 rounded border border-[var(--color-border)] bg-[var(--color-bg)] text-sm">
              <option value="" disabled>选择类型</option>
              <option v-for="t in story.storyTypes" :key="t" :value="t">{{ t }}</option>
            </select>
          </div>

          <div v-if="story.sorts.length">
            <label class="block text-xs text-[var(--color-text-secondary)] mb-1">排序</label>
            <select v-model="story.selectedSort" class="w-full px-2 py-1.5 rounded border border-[var(--color-border)] bg-[var(--color-bg)] text-sm">
              <option value="" disabled>选择排序</option>
              <option v-for="s in story.sorts" :key="s.value" :value="s.value">{{ s.label }}</option>
            </select>
          </div>

          <div>
            <label class="block text-xs text-[var(--color-text-secondary)] mb-1">索引</label>
            <select v-model="story.selectedIndex" class="w-full px-2 py-1.5 rounded border border-[var(--color-border)] bg-[var(--color-bg)] text-sm">
              <option value="" disabled>选择索引</option>
              <option v-for="i in displayIndices" :key="i.value" :value="i.value">{{ i.label }}</option>
            </select>
          </div>

          <div>
            <label class="block text-xs text-[var(--color-text-secondary)] mb-1">章节</label>
            <select v-model="story.selectedChapter" class="w-full px-2 py-1.5 rounded border border-[var(--color-border)] bg-[var(--color-bg)] text-sm">
              <option :value="-1" disabled>选择章节</option>
              <option v-for="c in story.chapters" :key="c.number" :value="c.number">{{ c.label }}</option>
            </select>
          </div>
        </div>

        <div class="border-t border-[var(--color-border)]" />

        <h2 class="text-base font-semibold">下载选项</h2>

        <div>
          <label class="block text-xs text-[var(--color-text-secondary)] mb-1">输出目录</label>
          <div class="flex gap-2">
            <input
              v-model="outputDir"
              type="text"
              placeholder="输入保存目录路径..."
              class="flex-1 px-2 py-1.5 rounded border border-[var(--color-border)] bg-[var(--color-bg)] text-sm"
            />
            <button
              @click="handleDownload"
              class="px-6 py-1.5 rounded-lg text-sm text-white transition-opacity hover:opacity-90"
              style="background-color: var(--color-primary)"
            >
              下载
            </button>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>
