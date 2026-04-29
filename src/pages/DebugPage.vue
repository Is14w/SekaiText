<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft } from 'lucide-vue-next'
import { useDebugLog } from '../composables/useDebugLog'

const router = useRouter()
const { logs, clear: clearFrontend } = useDebugLog()

interface ServerLog {
  timestamp: string
  message: string
}

const serverLogs = ref<ServerLog[]>([])
let pollTimer: ReturnType<typeof setInterval> | null = null
const logContainer = ref<HTMLElement | null>(null)

interface MergedEntry {
  id: string
  ts: string
  msg: string
  type: 'info' | 'warn' | 'error' | 'server'
}

const mergedLogs = ref<MergedEntry[]>([])

function mergeLogs() {
  const all: MergedEntry[] = []

  for (const entry of logs.value) {
    all.push({ id: `f-${all.length}`, ts: entry.ts, msg: entry.msg, type: entry.type })
  }

  for (const entry of serverLogs.value) {
    all.push({ id: `s-${all.length}`, ts: entry.timestamp, msg: entry.message, type: 'server' })
  }

  // Sort by timestamp (approximate - just use position to keep order roughly)
  mergedLogs.value = all
}

async function fetchServerLogs() {
  try {
    const res = await fetch('/api/v1/debug/logs')
    if (res.ok) {
      serverLogs.value = await res.json()
      mergeLogs()
      scrollToBottom()
    }
  } catch {
    // server not available
  }
}

function scrollToBottom() {
  nextTick(() => {
    if (logContainer.value) {
      logContainer.value.scrollTop = logContainer.value.scrollHeight
    }
  })
}

watch(logs, () => {
  mergeLogs()
  scrollToBottom()
}, { deep: true })

onMounted(() => {
  mergeLogs()
  fetchServerLogs()
  pollTimer = setInterval(fetchServerLogs, 2000)
})

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer)
})

async function saveLogs() {
  const lines = mergedLogs.value.map(e => {
    const tag = e.type === 'server' ? 'server' : 'front'
    return `[${e.ts}] [${tag}] ${e.msg}`
  }).join('\n')
  try {
    const res = await fetch('/api/v1/debug/save', { method: 'POST' })
    if (res.ok) {
      const data = await res.json()
      addLogLine(`日志已保存 (${data.lines} 行 → debug.log)`, 'info')
    }
  } catch {
    // Fallback: download as blob
    const blob = new Blob([lines], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `sekai-debug-${Date.now()}.log`
    a.click()
    URL.revokeObjectURL(url)
  }
}

let lineId = 0
function addLogLine(msg: string, type: 'info' | 'warn' | 'error') {
  const ts = new Date().toLocaleTimeString()
  mergedLogs.value.push({ id: `log-${lineId++}`, ts, msg, type })
  scrollToBottom()
}

function clearAll() {
  clearFrontend()
  serverLogs.value = []
  mergedLogs.value = []
}
</script>

<template>
  <div class="h-screen bg-[#0c0c0c] flex flex-col">
    <!-- Header -->
    <header class="flex items-center justify-between px-4 py-2 bg-[#1a1a1a] border-b border-[#2a2a2a] flex-shrink-0 select-none">
      <div class="flex items-center gap-3">
        <button
          @click="router.push('/')"
          class="text-xs text-[#888] hover:text-[#fff] transition-colors flex items-center gap-1"
        >
          <ArrowLeft :size="14" />
          返回
        </button>
        <span class="text-xs font-medium text-[#666]">debug.log</span>
      </div>
      <div class="flex items-center gap-2">
        <span class="text-xs text-[#555]">{{ mergedLogs.length }} lines</span>
        <button
          @click="saveLogs"
          class="text-xs text-[#666] hover:text-[#fff] transition-colors px-2 py-0.5 rounded hover:bg-[#2a2a2a]"
        >
          保存日志
        </button>
        <button
          @click="clearAll"
          class="text-xs text-[#666] hover:text-[#fff] transition-colors px-2 py-0.5 rounded hover:bg-[#2a2a2a]"
        >
          清空
        </button>
      </div>
    </header>

    <!-- Terminal Output -->
    <main
      ref="logContainer"
      class="flex-1 overflow-y-auto p-3 font-mono text-xs leading-relaxed"
    >
      <div v-if="mergedLogs.length === 0" class="text-[#444] py-4">
        <div class="mb-1"><span class="text-[#555]">$</span> <span class="text-[#888]">waiting for logs...</span></div>
        <div class="text-[#444]">[sekai:debug] console capture active</div>
        <div class="text-[#444]">[sekai:debug] polling server logs at /api/v1/debug/logs</div>
      </div>

      <div v-else class="space-y-0">
        <div
          v-for="entry in mergedLogs"
          :key="entry.id"
          class="flex gap-2"
        >
          <span class="text-[#444] flex-shrink-0 w-16 select-none">{{ entry.ts }}</span>
          <span v-if="entry.type === 'server'" class="text-[#569cd6] flex-shrink-0">[server]</span>
          <span v-else class="text-[#444] flex-shrink-0">[front]</span>
          <span :class="{
            'text-[#d4d4d4]': entry.type === 'info',
            'text-[#dcdcaa]': entry.type === 'warn',
            'text-[#f44747]': entry.type === 'error',
            'text-[#6a9955]': entry.type === 'server',
          }" class="whitespace-pre-wrap break-all">{{ entry.msg }}</span>
        </div>
      </div>
    </main>
  </div>
</template>
