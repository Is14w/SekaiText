import { useDebugLog } from '../composables/useDebugLog'

const BASE_URL = '/api/v1'

const { log: debugLog } = useDebugLog()

class ApiError extends Error {
  status: number
  constructor(status: number, message: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const method = (options.method || 'GET').toUpperCase()
  const url = `${BASE_URL}${path}`
  const start = Date.now()
  debugLog(`→ ${method} ${path}`, 'info')

  try {
    const res = await fetch(url, {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    })

    const elapsed = Date.now() - start

    if (!res.ok) {
      const body = await res.json().catch(() => ({ error: res.statusText }))
      debugLog(`← ${method} ${path} ${res.status} ${elapsed}ms`, 'error')
      throw new ApiError(res.status, body.error || res.statusText)
    }

    const data = await res.json()
    debugLog(`← ${method} ${path} ${res.status} ${elapsed}ms`, 'info')
    return data
  } catch (e) {
    if (!(e instanceof ApiError)) {
      debugLog(`← ${method} ${path} FAILED`, 'error')
    }
    throw e
  }
}

export const api = {
  // Story navigation
  storyTypes: () => request<string[]>('/story/types'),

  storySorts: (type: string) =>
    request<{ label: string; value: string }[]>(`/story/sorts?type=${encodeURIComponent(type)}`),

  storyIndex: (type: string, sort: string) =>
    request<{ label: string; value: string; chapters?: number[] }[]>(
      `/story/index?type=${encodeURIComponent(type)}&sort=${encodeURIComponent(sort)}`,
    ),

  storyChapter: (type: string, sort: string, index: string) =>
    request<{ number: number; label: string }[]>(
      `/story/chapter?type=${encodeURIComponent(type)}&sort=${encodeURIComponent(sort)}&index=${encodeURIComponent(index)}`,
    ),

  jsonPath: (type: string, sort: string, index: string, chapter: number, source: string) =>
    request<{ url: string; fileName: string }>(
      `/story/json-path?type=${encodeURIComponent(type)}&sort=${encodeURIComponent(sort)}&index=${encodeURIComponent(index)}&chapter=${chapter}&source=${encodeURIComponent(source)}`,
    ),

  storyLoad: (data: {
    storyType: string
    sort: string
    index: string
    chapter: number
    source: string
  }) =>
    request<{ scenarioId: string; sourceTalks: import('../types/translation').SourceTalk[]; saveTitle: string; chapterTitle: string }>(
      '/story/load',
      { method: 'POST', body: JSON.stringify(data) },
    ),

  storyLoadLocal: (content: string) =>
    request<{ scenarioId: string; sourceTalks: import('../types/translation').SourceTalk[] }>(
      '/story/load-local',
      { method: 'POST', body: JSON.stringify({ content }) },
    ),

  // Translation
  translationLoadContent: (content: string) =>
    request<{
      talks: import('../types/translation').DstTalk[]
      meta: import('../types/api').SaveMetadata | null
    }>('/translation/load-content', {
      method: 'POST',
      body: JSON.stringify({ content }),
    }),

  translationSerialize: (data: {
    talks: import('../types/translation').DstTalk[]
    saveN: boolean
    meta?: import('../types/api').SaveMetadata
  }) =>
    request<{ content: string }>('/translation/serialize', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  translationCreate: (data: {
    sourceTalks: import('../types/translation').SourceTalk[]
    jp: boolean
  }) => request<import('../types/translation').DstTalk[]>('/translation/create', {
    method: 'POST',
    body: JSON.stringify(data),
  }),

  translationLoad: (filePath: string) =>
    request<{
      talks: import('../types/translation').DstTalk[]
      meta: import('../types/api').SaveMetadata | null
    }>('/translation/load', {
      method: 'POST',
      body: JSON.stringify({ filePath }),
    }),

  translationSave: (filePath: string, talks: import('../types/translation').DstTalk[], saveN: boolean, meta?: import('../types/api').SaveMetadata) =>
    request<{ status: string }>('/translation/save', {
      method: 'POST',
      body: JSON.stringify({ filePath, talks, saveN, meta }),
    }),

  checkLines: (data: {
    sourceTalks: import('../types/translation').SourceTalk[]
    loadedTalks: import('../types/translation').DstTalk[]
  }) => request<import('../types/translation').DstTalk[]>('/translation/check-lines', {
    method: 'POST',
    body: JSON.stringify(data),
  }),

  // Editor
  changeText: (data: {
    row: number
    text: string
    editorMode: number
    talks: import('../types/translation').DstTalk[]
    dstTalks: import('../types/translation').DstTalk[]
    referTalks: import('../types/translation').DstTalk[]
    sourceTalks: import('../types/translation').SourceTalk[]
  }) =>
    request<{ talks: import('../types/translation').DstTalk[]; dstTalks: import('../types/translation').DstTalk[] }>(
      '/editor/change-text',
      { method: 'POST', body: JSON.stringify(data) },
    ),

  addLine: (data: {
    row: number
    talks: import('../types/translation').DstTalk[]
    dstTalks: import('../types/translation').DstTalk[]
    isProofreading: boolean
    sourceTalks: import('../types/translation').SourceTalk[]
  }) =>
    request<{ talks: import('../types/translation').DstTalk[]; dstTalks: import('../types/translation').DstTalk[] }>(
      '/editor/add-line',
      { method: 'POST', body: JSON.stringify(data) },
    ),

  removeLine: (data: {
    row: number
    talks: import('../types/translation').DstTalk[]
    dstTalks: import('../types/translation').DstTalk[]
  }) =>
    request<{ talks: import('../types/translation').DstTalk[]; dstTalks: import('../types/translation').DstTalk[] }>(
      '/editor/remove-line',
      { method: 'POST', body: JSON.stringify(data) },
    ),

  compare: (data: {
    referTalks: import('../types/translation').DstTalk[]
    checkTalks: import('../types/translation').DstTalk[]
    editorMode: number
  }) => request<{ talks: import('../types/translation').DstTalk[]; dstTalks: import('../types/translation').DstTalk[] }>('/editor/compare', {
    method: 'POST',
    body: JSON.stringify(data),
  }),

  replaceBrackets: (data: {
    row: number
    brackets: string
    talks: import('../types/translation').DstTalk[]
  }) => request<import('../types/translation').DstTalk[]>('/editor/replace-brackets', {
    method: 'POST',
    body: JSON.stringify(data),
  }),

  // Text check
  checkText: (data: { speaker: string; text: string }) =>
    request<{ text: string; checked: boolean; message?: string }>('/check/text', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  // Flashback
  flashbackAnalyze: (sourceTalks: import('../types/translation').SourceTalk[]) =>
    request<{ majorClue: string; sourceTalks: import('../types/translation').SourceTalk[] }>(
      '/flashback/analyze',
      { method: 'POST', body: JSON.stringify({ sourceTalks }) },
    ),

  clueHints: (clue: string, lang = 'zh-cn') =>
    request<Record<string, string>>(`/flashback/clue-hints?clue=${encodeURIComponent(clue)}&lang=${encodeURIComponent(lang)}`),

  voiceClues: () => request<Record<string, unknown>>('/flashback/voice-clues'),

  // Voice
  voiceUrl: (scenarioId: string, voiceId: string, source: string) =>
    request<{ url: string }>(
      `/voice/url?scenarioId=${encodeURIComponent(scenarioId)}&voiceId=${encodeURIComponent(voiceId)}&source=${encodeURIComponent(source)}`,
    ),

  // Speaker
  speakerCount: (data: {
    talks: import('../types/translation').DstTalk[]
    sourceTalks: import('../types/translation').SourceTalk[]
  }) =>
    request<{ speakers: { japanese: string; chinese: string; count: number }[] }>('/speaker/count', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  // Recovery (autosave)
  recoverySave: (data: {
    talks: import('../types/translation').DstTalk[]
    saveN: boolean
    filePath: string
    editorMode: number
  }) =>
    request<{ status: string }>('/recovery/save', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  recoveryLoad: () =>
    request<{
      exists: boolean
      content?: string
      filePath?: string
      editorMode?: number
      savedAt?: string
    }>('/recovery/load'),

  recoveryClear: () =>
    request<{ status: string }>('/recovery/clear', { method: 'DELETE' }),

  // Settings
  getSettings: () => request<import('../types/api').Settings>('/settings'),
  putSettings: (settings: import('../types/api').Settings) =>
    request<import('../types/api').Settings>('/settings', {
      method: 'PUT',
      body: JSON.stringify(settings),
    }),

  // Update (CDN refresh)
  update: () => request<{ status: string }>('/update', { method: 'POST' }),
  updateProgress: () =>
    request<{ current: number; total: number; message?: string; done: boolean }>('/update/progress'),

  // JSON Download
  downloadJson: (data: {
    storyType: string
    sort: string
    index: string
    chapter: number
    source: string
    outputDir: string
  }) =>
    request<{ taskId: string }>('/story/download-json', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  downloadProgress: (taskId: string) =>
    request<{ taskId: string; status: string; read: number; total: number; filePath?: string; error?: string }>(
      '/story/download-progress?task=' + encodeURIComponent(taskId),
    ),

  // Assets
  characters: () =>
    request<import('../types/dictionary').CharacterInfo[]>('/assets/characters'),
  units: () => request<import('../types/dictionary').UnitInfo[]>('/assets/units'),
  areas: () => request<string[]>('/assets/areas'),
  characterIconUrl: (index: number) => `${BASE_URL}/assets/character-icon/${index}`,
}
