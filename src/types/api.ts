import type { SourceTalk, DstTalk } from './translation'

export interface Settings {
  fontSize: number
  downloadSource: string
  saveN: boolean
  saveVoice: boolean
  disableSSL: boolean
  indexOrder: 'asc' | 'desc'
  voiceOutputDir?: string
  lastStoryType?: string
  lastStorySort?: string
  lastStoryIndex?: string
  lastChapter?: number
  lastDataSource?: string
}

export interface UpdateProgress {
  current: number
  total: number
  message?: string
  done: boolean
}

export interface LoadRequest {
  storyType: string
  sort: string
  index: string
  chapter: number
  source: string
}

export interface LoadResponse {
  scenarioId: string
  sourceTalks: SourceTalk[]
}

export interface TranslationCreateRequest {
  sourceTalks: SourceTalk[]
  jp: boolean
}

export interface TranslationLoadRequest {
  filePath: string
}

export interface TranslationSaveRequest {
  filePath: string
  talks: DstTalk[]
  saveN: boolean
}

export interface CheckLinesRequest {
  sourceTalks: SourceTalk[]
  loadedTalks: DstTalk[]
}

export interface CheckTextRequest {
  speaker: string
  text: string
}

export interface CheckTextResponse {
  text: string
  checked: boolean
  message?: string
}

export interface SpeakerCountRequest {
  talks: DstTalk[]
  sourceTalks: SourceTalk[]
}

export interface SpeakerCountResponse {
  speakers: SpeakerEntry[]
}

export interface SpeakerEntry {
  japanese: string
  chinese: string
  count: number
}

export interface VoiceURLRequest {
  scenarioId: string
  voiceId: string
  source: string
}

export interface VoiceURLResponse {
  url: string
}
