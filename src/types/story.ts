export type StoryType = 'event' | 'mainstory' | 'card' | 'festival' | 'areatalk' | 'greet' | 'special'

export interface StorySort {
  label: string
  value: string
}

export interface StoryIndex {
  label: string
  value: string
  chapters?: number[]
}

export interface StoryChapter {
  number: number
  label: string
}

export interface JsonPathResult {
  url: string
  fileName: string
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
