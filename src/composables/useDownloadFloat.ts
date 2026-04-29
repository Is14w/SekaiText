import { ref } from 'vue'
import type { DownloadTask } from '../components/DownloadFloat.vue'

const tasks = ref<DownloadTask[]>([])
let nextId = 0

export function useDownloadFloat() {
  function add(name: string): number {
    const id = nextId++
    tasks.value.push({ id, name, status: 'pending' })
    return id
  }

  function start(id: number) {
    const t = tasks.value.find(t => t.id === id)
    if (t) t.status = 'downloading'
  }

  function progress(id: number, read: number, total: number, pct: number) {
    const t = tasks.value.find(t => t.id === id)
    if (t) {
      t.read = read
      t.total = total
      t.percent = pct
    }
  }

  function done(id: number, result: string) {
    const t = tasks.value.find(t => t.id === id)
    if (t) {
      t.status = 'done'
      t.result = result
    }
    setTimeout(() => remove(id), 5000)
  }

  function fail(id: number, error: string) {
    const t = tasks.value.find(t => t.id === id)
    if (t) {
      t.status = 'error'
      t.error = error
    }
    setTimeout(() => remove(id), 6000)
  }

  function remove(id: number) {
    tasks.value = tasks.value.filter(t => t.id !== id)
  }

  return { tasks, add, start, progress, done, fail }
}
