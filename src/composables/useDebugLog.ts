import { ref, watch } from 'vue'

export interface LogEntry {
  ts: string
  msg: string
  type: 'info' | 'warn' | 'error'
}

const logs = ref<LogEntry[]>([])
const enabled = ref(localStorage.getItem('debug-enabled') === 'true')
const MAX_LOGS = 500
let captured = false

watch(enabled, (v) => {
  localStorage.setItem('debug-enabled', String(v))
})

export function useDebugLog() {
  function log(msg: string, type: LogEntry['type'] = 'info') {
    const ts = new Date().toLocaleTimeString()
    if (logs.value.length >= MAX_LOGS) {
      logs.value.splice(0, 100)
    }
    logs.value.push({ ts, msg, type })
  }

  function clear() {
    logs.value = []
  }

  function initConsoleCapture() {
    if (captured) return
    captured = true

    const origLog = console.log
    const origWarn = console.warn
    const origError = console.error
    const origDebug = console.debug

    console.log = (...args: any[]) => {
      log(args.map(a => typeof a === 'object' ? safeStringify(a) : String(a)).join(' '), 'info')
      origLog.apply(console, args)
    }
    console.warn = (...args: any[]) => {
      log(args.map(a => typeof a === 'object' ? safeStringify(a) : String(a)).join(' '), 'warn')
      origWarn.apply(console, args)
    }
    console.error = (...args: any[]) => {
      log(args.map(a => typeof a === 'object' ? safeStringify(a) : String(a)).join(' '), 'error')
      origError.apply(console, args)
    }
    console.debug = (...args: any[]) => {
      log(args.map(a => typeof a === 'object' ? safeStringify(a) : String(a)).join(' '), 'info')
      origDebug.apply(console, args)
    }

    log('Console capture initialized', 'info')
  }

  return { logs, enabled, log, clear, initConsoleCapture }
}

function safeStringify(obj: any): string {
  try {
    return JSON.stringify(obj)
  } catch {
    return String(obj)
  }
}
