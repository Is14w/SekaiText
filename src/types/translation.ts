export interface SourceTalk {
  speaker: string
  text: string
  voices?: string[]
  volume?: number[]
  charIndex: number
  clues?: string[]
}

export interface DstTalk {
  idx: number
  speaker: string
  text: string
  start: boolean
  end: boolean
  checked: boolean
  save: boolean
  message?: string
  dstidx: number
  referid?: number
  proofread?: boolean | null
  checkmode?: boolean
}

export type EditorMode = 0 | 1 | 2

export const EditorModeLabel: Record<EditorMode, string> = {
  0: '翻译',
  1: '校对',
  2: '合意',
}

export type TalkColor = 'white' | 'red' | 'yellow' | 'green' | 'blue'
