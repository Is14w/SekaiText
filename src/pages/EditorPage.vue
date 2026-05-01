<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAppStore } from '../stores/app'
import { useEditorStore } from '../stores/editor'
import { useStoryStore } from '../stores/story'
import { EditorModeLabel } from '../types/translation'
import type { SaveMetadata } from '../types/api'
import { useSettingsStore } from '../stores/settings'
import { useToast } from '../composables/useToast'
import { useFileDialog } from '../composables/useFileDialog'
import { useAutoSave } from '../composables/useAutoSave'
import { api } from '../api/client'
import { Pencil, Check, CircleDot, ChevronLeft, ChevronRight, Cog, Download, Bug } from 'lucide-vue-next'
import { getCurrentWindow } from '@tauri-apps/api/window'
import StoryNavigator from '../components/navigation/StoryNavigator.vue'
import EditorWorkspace from '../components/editor/EditorWorkspace.vue'
import SpeakerCountDialog from '../components/dialogs/SpeakerCountDialog.vue'
import SpeakerCheckDialog from '../components/dialogs/SpeakerCheckDialog.vue'

const app = useAppStore()
const editor = useEditorStore()
const story = useStoryStore()
const settings = useSettingsStore()
const toast = useToast()
const fileDialog = useFileDialog()
const autoSave = useAutoSave()

const isTauri = typeof window !== 'undefined' && !!(window as any).__TAURI_INTERNALS__

onMounted(async () => {
  autoSave.start()
  if (!isTauri) return
  try {
    const win = getCurrentWindow()
    await win.onCloseRequested(async (event) => {
      if (forceClose.value) return
      if (editor.hasUnsavedChanges) {
        event.preventDefault()
        await new Promise(r => setTimeout(r, 0))
        showCloseConfirm.value = true
      }
    })
  } catch (e: any) {
    tauriErr.value = `init: ${e.message || e}`
  }
})

async function handleCloseSave() {
  try {
    await handleSave()
    if (!editor.hasUnsavedChanges) {
      showCloseConfirm.value = false
      forceClose.value = true
      await closeWindow()
    }
  } catch { /* Save failed */ }
  showCloseConfirm.value = false
}

async function handleCloseDiscard() {
  showCloseConfirm.value = false
  forceClose.value = true
  await closeWindow()
}

async function closeWindow() {
  try { await getCurrentWindow().destroy() } catch {
    try { await getCurrentWindow().close() } catch {}
  }
}

function handleCloseCancel() { showCloseConfirm.value = false }

const showSpeakerCount = ref(false)
const tauriErr = ref('')
const showSpeakerCheck = ref(false)
const showCloseConfirm = ref(false)
const forceClose = ref(false)
const sidebarOpen = ref(true)

function setMode(key: number) {
  editor.switchMode(key as 0 | 1 | 2, settings.settings.preserveStoryOnModeSwitch)
  app.setEditorMode(key as 0 | 1 | 2)
  if (key === 2) app.showDiff = true
}

const modes = [ { key: 0, label: '翻译' }, { key: 1, label: '校对' }, { key: 2, label: '合意' } ]
const modeIcons: Record<number, typeof Pencil> = { 0: Pencil, 1: Check, 2: CircleDot }

async function handleOpen() {
  try {
    const result = await fileDialog.openTranslation()
    if (!result) return
    editor.setTalks(result.talks, result.talks, [])
    editor.currentFilePath = result.filePath || result.fileName || ''
    editor.markSaved()
    if (result.meta) {
      try {
        const m = result.meta
        story.selectedType = m.type
        story.selectedSort = m.sort || ''
        story.selectedIndex = m.index
        story.selectedChapter = m.chapter
        story.selectedSource = m.source
        await story.loadStory()
        if (story.sourceTalks.length > 0) {
          const aligned = await api.checkLines({ sourceTalks: story.sourceTalks, loadedTalks: result.talks })
          editor.setTalks(aligned, result.talks, [])
        }
      } catch { /* skip */ }
    }
    toast.show('已打开: ' + editor.currentFilePath, 'success')
  } catch (e: any) { toast.show('打开失败: ' + (e.message || '未知错误'), 'error') }
}

async function handleSave() {
  if (editor.talks.length === 0) return
  try {
    const modeLabel = EditorModeLabel[app.editorMode as 0 | 1 | 2]
    let defaultName = editor.currentFilePath
    if (!defaultName) {
      defaultName = '【' + modeLabel + '】' + (story.saveTitle || 'untitled')
      if (story.chapterTitle) defaultName += ' ' + story.chapterTitle
      defaultName += '.txt'
    }
    const meta: SaveMetadata | undefined = story.selectedType ? {
      type: story.selectedType, sort: story.selectedSort, index: story.selectedIndex,
      chapter: story.selectedChapter, source: story.selectedSource, scenarioId: story.scenarioId,
    } : undefined
    let saveTalks = editor.dstTalks
    if (app.editorMode === 2) { saveTalks = editor.dstTalks.filter(t => !t.checkmode) }
    const path = await fileDialog.saveTranslation(defaultName, saveTalks, app.saveN, meta)
    if (!path) return
    editor.currentFilePath = path
    editor.markSaved()
    api.recoveryClear().catch(() => {})
    toast.show('已保存', 'success')
  } catch (e: any) { toast.show('保存失败: ' + (e.message || '未知错误'), 'error') }
}

function handleClear() {
  if (editor.hasUnsavedChanges) { if (!confirm('有未保存的更改，确定清空吗？')) return }
  editor.clearAll()
  toast.show('已清空', 'info')
}

async function handleCompare() {
  if (editor.talks.length === 0) { toast.show('请先打开一份校对稿作为基准', 'warn'); return }
  try {
    const result = await fileDialog.openTranslation()
    if (!result) return
    const merged = await api.compare({ referTalks: editor.talks, checkTalks: result.talks, editorMode: 2 })
    editor.setTalks(merged.talks, merged.dstTalks, [])
    editor.currentFilePath = ''
    editor.markUnsaved()
    toast.show('已合并校对稿，请检查并解决差异', 'info')
  } catch (e: any) { toast.show('合意对比失败: ' + (e.message || '未知错误'), 'error') }
}

function handleConfirm() {
  if (editor.talks.length === 0) return
  if (!confirm('确认合意完成？所有差异将以当前内容为准合并。')) return
  const rowsToRemove: number[] = [], dstToRemove: number[] = [], idxHasProofread = new Set<number>()
  for (const talk of editor.talks) { if (talk.proofread) idxHasProofread.add(talk.idx) }
  for (let i = editor.talks.length - 1; i >= 0; i--) {
    const talk = editor.talks[i]
    if (talk.checkmode && idxHasProofread.has(talk.idx)) { rowsToRemove.push(i); dstToRemove.push(talk.dstidx) }
  }
  for (const i of rowsToRemove) editor.talks.splice(i, 1)
  dstToRemove.sort((a, b) => b - a)
  for (const d of dstToRemove) editor.dstTalks.splice(d, 1)
  for (const talk of editor.talks) {
    talk.checked = true; talk.checkmode = false; talk.proofread = null; talk.diff = undefined
    let removed = 0
    for (const d of dstToRemove) { if (talk.dstidx > d) removed++ }
    talk.dstidx -= removed
  }
  for (const talk of editor.dstTalks) { talk.checked = true }
  editor.markUnsaved()
  toast.show('合意已确认', 'success')
}

async function handleFullCheck() {
  if (editor.talks.length === 0) return
  let hasIssues = false, msgs: string[] = []
  for (const talk of editor.talks) {
    if (!talk.checked && talk.save) { hasIssues = true; msgs.push(`行 ${talk.idx}: ${talk.text.split('\n')[0]}`) }
  }
  if (hasIssues) { toast.show('发现 ' + msgs.length + ' 个问题', 'error') }
  else { toast.show('全文检查通过', 'success') }
}
</script>

<template>
  <div class="h-screen bg-[var(--color-bg)] flex flex-col">
    <div class="flex flex-1 min-h-0">
      <aside class="flex flex-col border-r border-[var(--color-border)] bg-[var(--color-surface)] flex-shrink-0 transition-all duration-200 overflow-hidden" :class="sidebarOpen ? 'w-36' : 'w-12'">
        <button @click="sidebarOpen = !sidebarOpen" class="flex items-center gap-2 h-10 px-3 text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] transition-colors flex-shrink-0">
          <ChevronLeft v-if="sidebarOpen" :size="18"/><ChevronRight v-else :size="18"/>
          <span v-if="sidebarOpen" class="text-xs font-medium">模式</span>
        </button>
        <div class="border-b border-[var(--color-border)]" />
        <div class="flex flex-col gap-0.5 p-1.5">
          <button v-for="m in modes" :key="m.key" @click="setMode(m.key)" class="flex items-center gap-2.5 h-9 px-2 rounded-lg transition-colors text-sm flex-shrink-0" :class="app.editorMode === m.key ? 'bg-[var(--color-primary)]/10 text-[var(--color-primary)] font-medium' : 'text-[var(--color-text-secondary)] hover:text-[var(--color-primary)]'">
            <component :is="modeIcons[m.key]" :size="18" /><span v-if="sidebarOpen" class="whitespace-nowrap">{{ m.label }}</span>
          </button>
        </div>
        <div class="flex-1" />
        <div class="border-t border-[var(--color-border)] p-1.5 space-y-0.5">
          <router-link to="/download" class="flex items-center gap-2.5 h-9 w-full px-2 rounded-lg transition-colors text-sm text-[var(--color-text-secondary)] hover:text-[var(--color-primary)]"><Download :size="18"/><span v-if="sidebarOpen" class="whitespace-nowrap">下载</span></router-link>
          <router-link v-if="settings.settings.debugEnabled" to="/debug" class="flex items-center gap-2.5 h-9 w-full px-2 rounded-lg transition-colors text-sm text-[var(--color-text-secondary)] hover:text-[var(--color-primary)]"><Bug :size="18"/><span v-if="sidebarOpen" class="whitespace-nowrap">调试</span></router-link>
          <router-link to="/settings" class="flex items-center gap-2.5 h-9 w-full px-2 rounded-lg transition-colors text-sm text-[var(--color-text-secondary)] hover:text-[var(--color-primary)]"><Cog :size="18"/><span v-if="sidebarOpen" class="whitespace-nowrap">设置</span></router-link>
        </div>
      </aside>
      <div class="flex-1 flex flex-col min-w-0">
        <header class="border-b border-[var(--color-border)] bg-[var(--color-surface)] px-4 py-2"><StoryNavigator/></header>
        <div class="border-b border-[var(--color-border)] bg-[var(--color-surface)] px-4 py-1.5">
          <div class="flex items-center gap-2 flex-wrap text-sm">
            <button @click="handleOpen" class="px-2.5 py-1 rounded text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] transition-colors">打开</button>
            <button @click="handleSave" class="px-2.5 py-1 rounded text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] transition-colors">保存</button>
            <button @click="handleClear" class="px-2.5 py-1 rounded text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] transition-colors">清空</button>
            <div class="w-px h-4 bg-[var(--color-border)]" />
            <label class="flex items-center gap-1 cursor-pointer text-[var(--color-text-secondary)] hover:text-[var(--color-primary)]"><input v-model="app.showFlashback" type="checkbox" class="accent-[var(--color-primary)] w-3 h-3" />闪回</label>
            <label class="flex items-center gap-1 cursor-pointer text-[var(--color-text-secondary)] hover:text-[var(--color-primary)]"><input v-model="app.syncScroll" type="checkbox" class="accent-[var(--color-primary)] w-3 h-3" />同步</label>
            <label class="flex items-center gap-1 cursor-pointer text-[var(--color-text-secondary)] hover:text-[var(--color-primary)]"><input v-model="app.showDiff" type="checkbox" class="accent-[var(--color-primary)] w-3 h-3" />差异</label>
            <div class="w-px h-4 bg-[var(--color-border)]" />
            <button @click="showSpeakerCheck = true" class="px-2.5 py-1 rounded text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] transition-colors">说话人</button>
            <button @click="handleFullCheck" class="px-2.5 py-1 rounded text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] transition-colors">检查</button>
            <button @click="showSpeakerCount = true" class="px-2.5 py-1 rounded text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] transition-colors">统计</button>
            <template v-if="app.editorMode === 2"><div class="w-px h-4 bg-[var(--color-border)]" /><button @click="handleCompare" class="px-2.5 py-1 rounded text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] transition-colors font-medium">对比</button><button @click="handleConfirm" class="px-2.5 py-1 rounded text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] transition-colors font-medium">确认</button></template>
          </div>
        </div>
        <main class="flex-1 min-h-0"><EditorWorkspace/></main>
      </div>
    </div>
    <SpeakerCountDialog v-if="showSpeakerCount" @close="showSpeakerCount = false"/>
    <SpeakerCheckDialog v-if="showSpeakerCheck" @close="showSpeakerCheck = false"/>
    <div v-if="showCloseConfirm" class="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div class="bg-[var(--color-surface)] border border-[var(--color-border)] rounded-xl shadow-xl w-96 max-w-[90vw] p-6">
        <h3 class="font-semibold text-sm text-[var(--color-text)] mb-2">有未保存的更改</h3>
        <p class="text-xs text-[var(--color-text-secondary)] mb-5">关闭前是否保存当前的工作内容？如果不保存，更改将丢失。</p>
        <div class="flex justify-end gap-2">
          <button @click="handleCloseCancel" class="px-4 py-2 text-sm rounded-lg border border-[var(--color-border)] text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] transition-colors">取消</button>
          <button @click="handleCloseDiscard" class="px-4 py-2 text-sm rounded-lg border border-red-400 text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors">不保存</button>
          <button @click="handleCloseSave" class="px-4 py-2 text-sm rounded-lg bg-[var(--color-primary)] text-white hover:opacity-90 transition-opacity">保存并退出</button>
        </div>
      </div>
    </div>
  </div>
</template>