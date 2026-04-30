import { ref } from 'vue'
import { api } from '../api/client'

const hintCache = new Map<string, string[]>()

export function useFlashbackTooltip() {
  const visible = ref(false)
  const tooltipStyle = ref<Record<string, string>>({})
  const clueGroups = ref<{ clue: string; hints: string[] }[]>([])

  let showTimer: ReturnType<typeof setTimeout> | null = null

  async function show(e: MouseEvent, clues: string[]) {
    if (!clues || clues.length === 0) return

    const newHints: Record<string, string[]> = {}
    for (const clue of clues) {
      if (!clue) continue
      if (hintCache.has(clue)) {
        newHints[clue] = hintCache.get(clue)!
        continue
      }
      try {
        const res = await api.clueHints(clue)
        const hints = res.hints?.length ? res.hints : [clue]
        hintCache.set(clue, hints)
        newHints[clue] = hints
      } catch {
        newHints[clue] = [clue]
      }
    }

    clueGroups.value = Object.entries(newHints).map(([clue, hints]) => ({ clue, hints }))

    const target = e.currentTarget as HTMLElement
    const rect = target.getBoundingClientRect()
    tooltipStyle.value = {
      position: 'fixed',
      left: rect.left + 'px',
      top: (rect.bottom + 4) + 'px',
      minWidth: Math.min(rect.width, 320) + 'px',
      maxWidth: '360px',
      zIndex: '9999',
    }

    if (showTimer) clearTimeout(showTimer)
    showTimer = setTimeout(() => {
      visible.value = true
    }, 300)
  }

  function hide() {
    if (showTimer) {
      clearTimeout(showTimer)
      showTimer = null
    }
    visible.value = false
  }

  return { visible, tooltipStyle, clueGroups, show, hide }
}
