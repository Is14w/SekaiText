import { ref, watch } from 'vue'
import { useEditorStore } from '../stores/editor'
import { useAppStore } from '../stores/app'
import { api } from '../api/client'

/**
 * Auto-saves the translation file every N seconds when there are unsaved changes.
 */
export function useAutoSave(intervalMs = 30000) {
  const editor = useEditorStore()
  const app = useAppStore()
  const lastSaved = ref(Date.now())
  let timer: ReturnType<typeof setInterval> | null = null

  function start() {
    if (timer) return
    timer = setInterval(async () => {
      if (!editor.hasUnsavedChanges || !editor.currentFilePath || editor.talks.length === 0) {
        return
      }
      try {
        await api.translationSave(editor.currentFilePath, editor.dstTalks, app.saveN)
        editor.markSaved()
        lastSaved.value = Date.now()
      } catch {
        // Silent fail on auto-save
      }
    }, intervalMs)
  }

  function stop() {
    if (timer) {
      clearInterval(timer)
      timer = null
    }
  }

  return {
    lastSaved,
    start,
    stop,
  }
}
