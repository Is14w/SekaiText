import { ref } from 'vue'

export interface LogEntry {
  ts: string
  msg: string
  type: 'info' | 'warn' | 'error'
}

const logs = ref<LogEntry[]>([])
const enabled = ref(false)

export function useDebugLog() {
  function log(msg: string, type: LogEntry['type'] = 'info') {
    const ts = new Date().toLocaleTimeString()
    logs.value.push({ ts, msg, type })
    console.log(`[${ts}] ${msg}`)
  }

  function clear() {
    logs.value = []
  }

  return { logs, enabled, log, clear }
}
