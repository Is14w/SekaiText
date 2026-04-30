import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { DstTalk, SourceTalk, EditorMode } from '../types/translation'

interface ModeState {
  talks: DstTalk[]
  dstTalks: DstTalk[]
  referTalks: DstTalk[]
  currentFilePath: string
  hasUnsavedChanges: boolean
  majorClue: string | null
}

function emptyModeState(): ModeState {
  return {
    talks: [],
    dstTalks: [],
    referTalks: [],
    currentFilePath: '',
    hasUnsavedChanges: false,
    majorClue: null,
  }
}

export const useEditorStore = defineStore('editor', () => {
  const talks = ref<DstTalk[]>([])
  const dstTalks = ref<DstTalk[]>([])
  const referTalks = ref<DstTalk[]>([])
  const sourceTalks = ref<SourceTalk[]>([])

  const currentFilePath = ref('')
  const hasUnsavedChanges = ref(false)
  const majorClue = ref<string | null>(null)

  const currentMode = ref<EditorMode>(0)
  const modeCache = new Map<EditorMode, ModeState>()

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
    modeCache.clear()
  }

  function saveModeState(mode: EditorMode) {
    modeCache.set(mode, {
      talks: talks.value,
      dstTalks: dstTalks.value,
      referTalks: referTalks.value,
      currentFilePath: currentFilePath.value,
      hasUnsavedChanges: hasUnsavedChanges.value,
      majorClue: majorClue.value,
    })
  }

  function loadModeState(mode: EditorMode) {
    const state = modeCache.get(mode) || emptyModeState()
    talks.value = state.talks
    dstTalks.value = state.dstTalks
    referTalks.value = state.referTalks
    currentFilePath.value = state.currentFilePath
    hasUnsavedChanges.value = state.hasUnsavedChanges
    majorClue.value = state.majorClue
  }

  function switchMode(newMode: EditorMode) {
    saveModeState(currentMode.value)
    currentMode.value = newMode
    loadModeState(newMode)
  }

  return {
    talks, dstTalks, referTalks, sourceTalks,
    currentFilePath, hasUnsavedChanges, majorClue,
    currentMode,
    setSourceTalks, setTalks, markUnsaved, markSaved, clearAll,
    saveModeState, loadModeState, switchMode,
  }
})
