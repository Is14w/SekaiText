import { ref } from 'vue'
import { useEditorStore } from '../stores/editor'
import { useAppStore } from '../stores/app'
import { api } from '../api/client'

/**
 * Periodically saves editor state to a recovery file (autosave).
 * Also saves to the real file if currentFilePath is set.
 */
export function useAutoSave(intervalMs = 30000) {
  const editor = useEditorStore()
  const app = useAppStore()
  const lastSaved = ref(Date.now())
  let timer: ReturnType<typeof setInterval> | null = null

  function start() {
    if (timer) return
    timer = setInterval(async () => {
      if (!editor.hasUnsavedChanges || editor.talks.length === 0) return

      // Always save recovery file
      try {
        await api.recoverySave({
          talks: editor.dstTalks,
          saveN: app.saveN,
          filePath: editor.currentFilePath,
          editorMode: app.editorMode,
        })
      } catch {
        // Silent fail on recovery save
      }

      // In proofread/合意 mode, never overwrite the opened source files
      if (app.editorMode !== 0) return

      // Also save to real file if a path is known
      if (editor.currentFilePath) {
        try {
          await api.translationSave(editor.currentFilePath, editor.dstTalks, app.saveN)
          editor.markSaved()
          lastSaved.value = Date.now()
        } catch {
          // Silent fail on auto-save
        }
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
