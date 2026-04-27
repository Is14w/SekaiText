import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { DstTalk, SourceTalk } from '../types/translation'

export const useEditorStore = defineStore('editor', () => {
  const talks = ref<DstTalk[]>([])
  const dstTalks = ref<DstTalk[]>([])
  const referTalks = ref<DstTalk[]>([])
  const sourceTalks = ref<SourceTalk[]>([])

  const currentFilePath = ref('')
  const hasUnsavedChanges = ref(false)
  const majorClue = ref<string | null>(null)

  function setSourceTalks(talks: SourceTalk[]) {
    sourceTalks.value = talks
  }

  function setTalks(newTalks: DstTalk[], newDstTalks: DstTalk[], newReferTalks: DstTalk[]) {
    talks.value = newTalks
    dstTalks.value = newDstTalks
    referTalks.value = newReferTalks
  }

  function markUnsaved() {
    hasUnsavedChanges.value = true
  }

  function markSaved() {
    hasUnsavedChanges.value = false
  }

  function clearAll() {
    talks.value = []
    dstTalks.value = []
    referTalks.value = []
    sourceTalks.value = []
    currentFilePath.value = ''
    hasUnsavedChanges.value = false
    majorClue.value = null
  }

  return {
    talks, dstTalks, referTalks, sourceTalks,
    currentFilePath, hasUnsavedChanges, majorClue,
    setSourceTalks, setTalks, markUnsaved, markSaved, clearAll,
  }
})
