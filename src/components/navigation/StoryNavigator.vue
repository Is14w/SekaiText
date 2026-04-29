<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { useStoryStore } from '../../stores/story'
import { useEditorStore } from '../../stores/editor'
import { useSettingsStore } from '../../stores/settings'
import { api } from '../../api/client'
import { useToast } from '../../composables/useToast'
import { useDebugLog } from '../../composables/useDebugLog'
import { useDownloadFloat } from '../../composables/useDownloadFloat'

const story = useStoryStore()
const editor = useEditorStore()
const settingsStore = useSettingsStore()

const displayIndices = computed(() => {
  if (settingsStore.settings.indexOrder === 'desc') {
    return [...story.indices].reverse()
  }
  return story.indices
})
const refreshing = ref(false)
const toast = useToast()
const debug = useDebugLog()
const dlFloat = useDownloadFloat()
const fileInput = ref<HTMLInputElement | null>(null)

debug.log('StoryNavigator mounted, fetching types...')
onMounted(() => {
  const t0 = performance.now()
  const doFetch = async () => {
    for (let i = 0; i < 5; i++) {
      try {
        await story.fetchTypes()
        if (story.storyTypes.length > 0) {
          debug.log(`故事类型已加载 (${story.storyTypes.length}个) 耗时:${Math.round(performance.now()-t0)}ms`)
          return
        }
      } catch (e: any) {
        debug.log(`获取类型失败尝试 ${i+1}: ${e.message}`, 'warn')
      }
      await new Promise(r => setTimeout(r, 1000))
    }
    debug.log('无法加载故事类型', 'error')
  }
  doFetch()
})



watch(() => story.selectedType, async (type) => {
  debug.log(`selectedType changed: "${type}"`)
  if (!type) return
  story.selectedSort = ''
  story.selectedIndex = ''
  story.selectedChapter = -1
  story.sorts = []
  story.indices = []
  story.chapters = []
  await story.fetchSorts(type)
  if (story.sorts.length === 0) {
    story.fetchIndex(type, '')
  }
})

watch(() => story.selectedSort, (sort) => {
  debug.log(`selectedSort changed: "${sort}"`)
  if (!sort || !story.selectedType) return
  story.selectedIndex = ''
  story.selectedChapter = -1
  story.indices = []
  story.chapters = []
  story.fetchIndex(story.selectedType, sort)
})

watch(() => story.selectedIndex, (idx) => {
  debug.log(`selectedIndex changed: "${idx}"`)
  if (!idx || !story.selectedType) return
  story.selectedChapter = -1
  story.chapters = []
  if (idx !== '-') {
    story.fetchChapters(story.selectedType, story.selectedSort, idx)
  }
})

async function handleRefresh() {
  refreshing.value = true
  const taskId = dlFloat.add('刷新元数据')
  dlFloat.start(taskId)
  try {
    await api.update()
    let done = false
    while (!done) {
      await new Promise(r => setTimeout(r, 500))
      const progress = await api.updateProgress()
      if (progress.total > 0) {
        dlFloat.progress(taskId, progress.current, progress.total, Math.round((progress.current / progress.total) * 100))
      }
      done = progress.done
    }
    dlFloat.done(taskId, '元数据已刷新')
  } catch (e: any) {
    dlFloat.fail(taskId, e.message || '刷新失败')
  } finally {
    refreshing.value = false
  }
}

async function handleLocalLoad() {
  const file = fileInput.value?.files?.[0]
  if (!file) return
  try {
    const text = await file.text()
    await story.loadStoryLocal(text)
    if (story.sourceTalks.length > 0) {
      debug.log(`本地故事载入成功 ${story.sourceTalks.length}行`)
      editor.setSourceTalks(story.sourceTalks)
      const dstTalks = await api.translationCreate({
        sourceTalks: story.sourceTalks,
        jp: false,
      })
      editor.setTalks(dstTalks, dstTalks, [])
      editor.majorClue = null
      toast.show(`已载入 ${story.sourceTalks.length} 行`, 'success')
    }
  } catch (e: any) {
    debug.log('本地载入失败: ' + e.message, 'error')
    toast.show('本地载入失败: ' + (e.message || '未知错误'), 'error')
  }
  if (fileInput.value) fileInput.value.value = ''
}
async function handleLoad() {
  debug.log(`载入按钮点击 loading=${story.loading} selectedType="${story.selectedType}"`)
  const taskId = dlFloat.add('载入故事')
  dlFloat.start(taskId)
  try {
    await story.loadStory()
    if (story.sourceTalks.length > 0) {
      debug.log(`故事载入成功 ${story.sourceTalks.length}行`)
      editor.setSourceTalks(story.sourceTalks)
      const dstTalks = await api.translationCreate({
        sourceTalks: story.sourceTalks,
        jp: false,
      })
      editor.setTalks(dstTalks, dstTalks, [])
      editor.majorClue = null
      dlFloat.done(taskId, `已载入 ${story.sourceTalks.length} 行`)
    } else {
      debug.log('故事载入返回0行', 'warn')
      dlFloat.done(taskId, '载入完成（0行）')
    }
  } catch (e: any) {
    debug.log('载入失败: ' + e.message, 'error')
    dlFloat.fail(taskId, e.message || '载入失败')
  }
}
</script>

<template>
  <div class="flex items-center gap-3 flex-wrap">
    <select v-model="story.selectedType" class="px-2 py-1 rounded border border-[var(--color-border)] bg-[var(--color-surface)] text-sm">
      <option value="" disabled>故事类型</option>
      <option v-for="t in story.storyTypes" :key="t" :value="t">{{ t }}</option>
    </select>

    <select v-if="story.sorts?.length" v-model="story.selectedSort" class="px-2 py-1 rounded border border-[var(--color-border)] bg-[var(--color-surface)] text-sm">
      <option value="" disabled>排序</option>
      <option v-for="s in story.sorts" :key="s.value" :value="s.value">{{ s.label }}</option>
    </select>

    <select v-model="story.selectedIndex" class="px-2 py-1 rounded border border-[var(--color-border)] bg-[var(--color-surface)] text-sm">
      <option value="" disabled>索引</option>
      <option v-for="i in displayIndices" :key="i.value" :value="i.value" v-text="i.label" />
    </select>

    <select v-model="story.selectedChapter" class="px-2 py-1 rounded border border-[var(--color-border)] bg-[var(--color-surface)] text-sm">
      <option :value="-1" disabled>章节</option>
      <option v-for="c in story.chapters" :key="c.number" :value="c.number">{{ c.label }}</option>
    </select>

    <span class="px-2 py-1 text-xs text-[var(--color-text-secondary)]">harukineo</span>

    <button
      @click="handleRefresh"
      class="px-3 py-1 rounded text-white text-sm transition-all hover:brightness-110 disabled:opacity-50"
      style="background-color: var(--color-primary)"
      :disabled="refreshing"
    >
      {{ refreshing ? '刷新中...' : '刷新' }}
    </button>

    <button
      @click="handleLoad"
      class="px-3 py-1 rounded text-white text-sm transition-all hover:brightness-110 disabled:opacity-50"
      style="background-color: var(--color-secondary)"
      :disabled="story.loading || !story.selectedType || story.selectedChapter < 0"
    >
      {{ story.loading ? '载入中...' : '载入' }}
    </button>

    <input
      ref="fileInput"
      type="file"
      accept=".json"
      class="hidden"
      @change="handleLocalLoad"
    />
    <button
      @click="fileInput?.click()"
      class="px-3 py-1 rounded text-sm transition-colors border border-[var(--color-border)] text-[var(--color-text-secondary)] hover:text-[var(--color-primary)]"
      :disabled="story.loading"
    >
      本地
    </button>
  </div>
</template>
