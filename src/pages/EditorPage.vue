<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useAppStore } from '../stores/app'
import { useEditorStore } from '../stores/editor'
import { useSettingsStore } from '../stores/settings'
import { api } from '../api/client'
import { useToast } from '../composables/useToast'
import { Minus, Square, X, Minimize, Pencil, Check, CircleDot, ChevronLeft, ChevronRight, Cog, Download, Bug } from 'lucide-vue-next'
import { getCurrentWindow } from '@tauri-apps/api/window'
import StoryNavigator from '../components/navigation/StoryNavigator.vue'
import EditorWorkspace from '../components/editor/EditorWorkspace.vue'
import SpeakerCountDialog from '../components/dialogs/SpeakerCountDialog.vue'
import SpeakerCheckDialog from '../components/dialogs/SpeakerCheckDialog.vue'


const app = useAppStore()
const editor = useEditorStore()
const settings = useSettingsStore()
const toast = useToast()

const isTauri = typeof window !== 'undefined' && !!(window as any).__TAURI_INTERNALS__

async function minimizeWin(e: MouseEvent) {
  e.stopPropagation()
  try {
    await getCurrentWindow().minimize()
  } catch (e: any) {
    tauriErr.value = `minimize: ${e.message || e}`
  }
}

const isMaximized = ref(false)

onMounted(async () => {
  if (!isTauri) return
  try {
    const win = getCurrentWindow()
    isMaximized.value = await win.isMaximized()
    await win.onResized(async () => {
      isMaximized.value = await win.isMaximized()
    })
  } catch (e: any) {
    tauriErr.value = `init: ${e.message || e}`
  }
})

async function toggleMaxWin(e: MouseEvent) {
  e.stopPropagation()
  try {
    const win = getCurrentWindow()
    if (await win.isMaximized()) {
      await win.unmaximize()
    } else {
      await win.maximize()
    }
  } catch (e: any) {
    tauriErr.value = `toggle: ${e.message || e}`
  }
}

async function closeWin(e: MouseEvent) {
  e.stopPropagation()
  try { await getCurrentWindow().close() } catch (e: any) {
    tauriErr.value = `close: ${e.message || e}`
  }
}

const showSpeakerCount = ref(false)
const tauriErr = ref('')
const showSpeakerCheck = ref(false)

const sidebarOpen = ref(true)

function setMode(key: number) {
  editor.switchMode(key as 0 | 1 | 2)
  app.setEditorMode(key as 0 | 1 | 2)
}

const modes = [
  { key: 0, label: '翻译' },
  { key: 1, label: '校对' },
  { key: 2, label: '合意' },
]

const modeIcons: Record<number, typeof Pencil> = {
  0: Pencil,
  1: Check,
  2: CircleDot,
}

async function handleOpen() {
  const path = window.prompt('输入翻译文件路径：')
  if (!path) return
  try {
    const loaded = await api.translationLoad(path)
    editor.setTalks(loaded, loaded, [])
    editor.currentFilePath = path
    editor.markSaved()
    toast.show('已打开: ' + path, 'success')
  } catch (e: any) {
    toast.show('打开失败: ' + (e.message || '未知错误'), 'error')
  }
}

async function handleSave() {
  if (editor.talks.length === 0) return
  let path: string | null = editor.currentFilePath
  if (!path) {
    path = window.prompt('输入保存路径：')
    if (!path) return
  }
  try {
    await api.translationSave(path, editor.dstTalks, app.saveN)
    editor.currentFilePath = path
    editor.markSaved()
    toast.show('已保存', 'success')
  } catch (e: any) {
    toast.show('保存失败: ' + (e.message || '未知错误'), 'error')
  }
}

function handleClear() {
  if (editor.hasUnsavedChanges) {
    if (!confirm('有未保存的更改，确定清空吗？')) return
  }
  editor.clearAll()
  toast.show('已清空', 'info')
}

async function handleFullCheck() {
  if (editor.talks.length === 0) return
  let hasIssues = false
  const msgs: string[] = []
  for (const talk of editor.talks) {
    if (!talk.checked && talk.save) {
      hasIssues = true
      msgs.push(`行 ${talk.idx}: ${talk.text.split('\n')[0]}`)
    }
  }
  if (hasIssues) {
    toast.show('发现 ' + msgs.length + ' 个问题', 'error')
  } else {
    toast.show('全文检查通过', 'success')
  }
}
</script>

<template>
  <div class="h-screen bg-[var(--color-bg)] flex flex-col">
    <!-- Custom Title Bar -->
    <div
      v-show="isTauri"
      class="h-10 flex items-center px-3 bg-[var(--color-surface)] border-b border-[var(--color-border)] flex-shrink-0 select-none"
    >
      <div class="flex items-center gap-2 flex-1 self-stretch" style="-webkit-app-region: drag">
        <img src="/app-icon.png" alt="" class="w-5 h-5" />
        <span class="font-bold" style="color: var(--color-primary); font-size: 15px">SekaiText</span>
      </div>
      <div class="flex items-center gap-1">
        <button @mousedown.stop @click="minimizeWin" class="w-10 h-8 flex items-center justify-center rounded hover:text-[var(--color-primary)] transition-colors text-[var(--color-text-secondary)]">
          <Minus :size="14" />
        </button>
        <button @mousedown.stop @click="toggleMaxWin" class="w-10 h-8 flex items-center justify-center rounded hover:text-[var(--color-primary)] text-[var(--color-text-secondary)] transition-colors">
          <Minimize v-if="isMaximized" :size="12" />
          <Square v-else :size="12" />
        </button>
        <button @mousedown.stop @click="closeWin" class="w-10 h-8 flex items-center justify-center rounded hover:bg-red-500 hover:text-white text-[var(--color-text-secondary)] transition-colors">
          <X :size="14" />
        </button>
      </div>
    </div>

    <!-- Body -->
    <div class="flex flex-1 min-h-0">
      <!-- Left Sidebar -->
    <aside
      class="flex flex-col border-r border-[var(--color-border)] bg-[var(--color-surface)] flex-shrink-0 transition-all duration-200 overflow-hidden"
      :class="sidebarOpen ? 'w-36' : 'w-12'"
    >
      <button
        @click="sidebarOpen = !sidebarOpen"
        class="flex items-center gap-2 h-10 px-3 text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] transition-colors flex-shrink-0"
      >
        <ChevronLeft v-if="sidebarOpen" :size="18" />
        <ChevronRight v-else :size="18" />
        <span v-if="sidebarOpen" class="text-xs font-medium">模式</span>
      </button>

      <div class="border-b border-[var(--color-border)]" />

      <div class="flex flex-col gap-0.5 p-1.5">
        <button
          v-for="m in modes"
          :key="m.key"
          @click="setMode(m.key)"
          class="flex items-center gap-2.5 h-9 px-2 rounded-lg transition-colors text-sm flex-shrink-0"
          :class="app.editorMode === m.key
            ? 'bg-[var(--color-primary)]/10 text-[var(--color-primary)] font-medium'
            : 'text-[var(--color-text-secondary)] hover:text-[var(--color-primary)]'"
        >
          <component :is="modeIcons[m.key]" :size="18" />
          <span v-if="sidebarOpen" class="whitespace-nowrap">{{ m.label }}</span>
        </button>
      </div>

      <div class="flex-1" />

      <div class="border-t border-[var(--color-border)] p-1.5 space-y-0.5">
        <router-link
          to="/download"
          class="flex items-center gap-2.5 h-9 w-full px-2 rounded-lg transition-colors text-sm text-[var(--color-text-secondary)] hover:text-[var(--color-primary)]"
        >
          <Download :size="18" />
          <span v-if="sidebarOpen" class="whitespace-nowrap">下载</span>
        </router-link>
        <router-link
          v-if="settings.settings.debugEnabled"
          to="/debug"
          class="flex items-center gap-2.5 h-9 w-full px-2 rounded-lg transition-colors text-sm text-[var(--color-text-secondary)] hover:text-[var(--color-primary)]"
        >
          <Bug :size="18" />
          <span v-if="sidebarOpen" class="whitespace-nowrap">调试</span>
        </router-link>
        <router-link
          to="/settings"
          class="flex items-center gap-2.5 h-9 w-full px-2 rounded-lg transition-colors text-sm text-[var(--color-text-secondary)] hover:text-[var(--color-primary)]"
        >
          <Cog :size="18" />
          <span v-if="sidebarOpen" class="whitespace-nowrap">设置</span>
        </router-link>
      </div>
    </aside>

    <!-- Main Area -->
    <div class="flex-1 flex flex-col min-w-0">
      <header class="border-b border-[var(--color-border)] bg-[var(--color-surface)] px-4 py-2">
        <StoryNavigator />
      </header>

      <div class="border-b border-[var(--color-border)] bg-[var(--color-surface)] px-4 py-1.5">
        <div class="flex items-center gap-2 flex-wrap text-sm">
          <button @click="handleOpen" class="px-2.5 py-1 rounded text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] transition-colors">打开</button>
          <button @click="handleSave" class="px-2.5 py-1 rounded text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] transition-colors">保存</button>
          <button @click="handleClear" class="px-2.5 py-1 rounded text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] transition-colors">清空</button>

          <div class="w-px h-4 bg-[var(--color-border)]" />

          <label class="flex items-center gap-1 cursor-pointer text-[var(--color-text-secondary)] hover:text-[var(--color-primary)]">
            <input v-model="app.showFlashback" type="checkbox" class="accent-[var(--color-primary)] w-3 h-3" />闪回
          </label>
          <label class="flex items-center gap-1 cursor-pointer text-[var(--color-text-secondary)] hover:text-[var(--color-primary)]">
            <input v-model="app.syncScroll" type="checkbox" class="accent-[var(--color-primary)] w-3 h-3" />同步
          </label>
          <label class="flex items-center gap-1 cursor-pointer text-[var(--color-text-secondary)] hover:text-[var(--color-primary)]">
            <input v-model="app.showDiff" type="checkbox" class="accent-[var(--color-primary)] w-3 h-3" />差异
          </label>

          <div class="w-px h-4 bg-[var(--color-border)]" />

          <button @click="showSpeakerCheck = true" class="px-2.5 py-1 rounded text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] transition-colors">说话人</button>
          <button @click="handleFullCheck" class="px-2.5 py-1 rounded text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] transition-colors">检查</button>
          <button @click="showSpeakerCount = true" class="px-2.5 py-1 rounded text-[var(--color-text-secondary)] hover:text-[var(--color-primary)] transition-colors">统计</button>
        </div>
      </div>

      <main class="flex-1 min-h-0">
        <EditorWorkspace />
      </main>
    </div>
    </div>

    <SpeakerCountDialog v-if="showSpeakerCount" @close="showSpeakerCount = false" />
    <SpeakerCheckDialog v-if="showSpeakerCheck" @close="showSpeakerCheck = false" />

  </div>
</template>

