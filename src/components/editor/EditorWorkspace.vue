<script setup lang="ts">
import { computed, ref } from 'vue'
import { useAppStore } from '../../stores/app'
import { useStoryStore } from '../../stores/story'
import { useEditorStore } from '../../stores/editor'
import { api } from '../../api/client'
import { useToast } from '../../composables/useToast'
import VoicePlayButton from './VoicePlayButton.vue'
import type { DstTalk } from '../../types/translation'

const iconErrors = ref<Set<number>>(new Set())

const app = useAppStore()
const story = useStoryStore()
const editor = useEditorStore()
const toast = useToast()

// ---- Flashback (from SourcePanel) ----
const talksWithFlashback = computed(() => {
  if (!app.showFlashback) {
    return story.sourceTalks.map(t => ({ ...t, isFlashback: false }))
  }

  const clueCounts = new Map<string, number>()
  for (const talk of story.sourceTalks) {
    if (talk.clues) {
      for (const clue of talk.clues) {
        clueCounts.set(clue, (clueCounts.get(clue) || 0) + 1)
      }
    }
  }

  let majorClue: string | null = null
  let maxCount = 0
  for (const [clue, count] of clueCounts) {
    if (count > maxCount) {
      maxCount = count
      majorClue = clue
    }
  }

  return story.sourceTalks.map(talk => {
    if (!talk.clues || talk.clues.length === 0) {
      return { ...talk, isFlashback: false }
    }
    const isFlashback = talk.clues.some(c => c !== majorClue)
    return { ...talk, isFlashback }
  })
})

// ---- Helpers ----
function srcIdx(talk: DstTalk): number {
  return talk.idx - 1
}

function srcTalk(talk: DstTalk) {
  return story.sourceTalks[srcIdx(talk)]
}

function flashbackItem(talk: DstTalk) {
  return talksWithFlashback.value[srcIdx(talk)]
}

function srcTalkCharIndex(talk: DstTalk) {
  return srcTalk(talk)?.charIndex ?? -1
}

function flashbackClues(talk: DstTalk) {
  const fb = flashbackItem(talk)
  return fb?.isFlashback && fb?.clues ? fb.clues.filter((c: string) => c).join(', ') : undefined
}

// Group consecutive dest lines sharing the same source idx
const talkGroups = computed(() => {
  const groups: { srcIdx: number; items: { talk: DstTalk; globalIdx: number }[] }[] = []
  for (let i = 0; i < editor.talks.length; i++) {
    const talk = editor.talks[i]
    const last = groups[groups.length - 1]
    if (last && last.srcIdx === talk.idx) {
      last.items.push({ talk, globalIdx: i })
    } else {
      groups.push({ srcIdx: talk.idx, items: [{ talk, globalIdx: i }] })
    }
  }
  return groups
})

// ---- Editing (from DestPanel) ----
let editTimeout: ReturnType<typeof setTimeout> | null = null

const MAX_LINES_PER_SRC = 10

function getRowClass(talk: DstTalk): Record<string, boolean> {
  return {
    'opacity-40': !talk.save && !app.showDiff,
    'hidden': talk.proofread === false && !app.showDiff,
  }
}

function getDestBorder(talk: DstTalk): string {
  if (talk.proofread === true && app.showDiff) return 'border-l-green-400'
  if (talk.proofread === false) return 'border-l-yellow-400'
  if (talk.proofread === true && talk.checked && app.editorMode === 2) return 'border-l-blue-400'
  if (!talk.checked && talk.save) return 'border-l-red-400'
  return ''
}

async function handleTextChange(row: number, newText: string) {
  if (editTimeout) clearTimeout(editTimeout)
  editTimeout = setTimeout(async () => {
    try {
      const result = await api.changeText({
        row,
        text: newText,
        editorMode: app.editorMode,
        talks: editor.talks,
        dstTalks: editor.dstTalks,
        referTalks: editor.referTalks,
        sourceTalks: editor.sourceTalks,
      })
      editor.setTalks(result.talks, result.dstTalks, editor.referTalks)
      editor.markUnsaved()
    } catch {
      toast.show('文本保存失败', 'error')
    }
  }, 300)
}

function onBlur(e: Event, idx: number) {
  handleTextChange(idx, (e.target as HTMLElement).innerText)
}

async function handleAddLine(row: number) {
  const currentIdx = editor.talks[row]?.idx
  if (currentIdx && editor.talks.filter(t => t.idx === currentIdx).length >= MAX_LINES_PER_SRC) {
    toast.show(`每个原文行最多添加 ${MAX_LINES_PER_SRC} 行`, 'warn')
    return
  }
  try {
    const result = await api.addLine({
      row,
      talks: editor.talks,
      dstTalks: editor.dstTalks,
      isProofreading: app.editorMode !== 0,
      sourceTalks: editor.sourceTalks,
    })
    editor.setTalks(result.talks, result.dstTalks, editor.referTalks)
    editor.markUnsaved()
  } catch (e: any) {
    toast.show('添加行失败：' + e.message, 'error')
  }
}

async function handleRemoveLine(row: number) {
  try {
    const result = await api.removeLine({
      row,
      talks: editor.talks,
      dstTalks: editor.dstTalks,
    })
    editor.setTalks(result.talks, result.dstTalks, editor.referTalks)
    editor.markUnsaved()
  } catch (e: any) {
    toast.show('删除行失败：' + e.message, 'error')
  }
}

function handleBracketsReplace(row: number, brackets: string) {
  api.replaceBrackets({ row, brackets, talks: editor.talks }).then(newTalks => {
    editor.talks = newTalks
    editor.markUnsaved()
  })
}

function handleContextMenu(e: MouseEvent, row: number) {
  e.preventDefault()
  const brackets = window.prompt("选择括号类型：1=「」 2=『』 3=（） 4=\"\" 5=''")
  if (!brackets) return
  const map: Record<string, string> = {
    '1': '「」', '2': '『』', '3': '（）', '4': '""', '5': "''",
  }
  const b = map[brackets.trim()]
  if (b) handleBracketsReplace(row, b)
}
</script>

<template>
  <div class="flex h-full">
    <div
      class="flex-1 overflow-y-auto border border-[var(--color-border)] rounded-lg bg-[var(--color-surface)]"
    >
      <!-- Column headers -->
      <div class="grid grid-cols-2 border-b border-[var(--color-border)] bg-[var(--color-surface)] sticky top-0 z-10">
        <div class="flex items-center justify-between px-3 py-2">
          <span class="font-semibold text-sm text-[var(--color-text-secondary)]">原文</span>
          <span v-if="story.scenarioId" class="text-xs text-[var(--color-text-secondary)]">{{ story.scenarioId }}</span>
        </div>
        <div class="flex items-center px-3 py-2 border-l border-[var(--color-border)]">
          <span class="font-semibold text-sm text-[var(--color-text-secondary)]">译文</span>
          <input
            v-model="editor.currentFilePath"
            type="text"
            placeholder="标题/路径..."
            class="ml-2 flex-1 text-sm px-2 py-0.5 rounded border border-[var(--color-border)] bg-[var(--color-surface)]"
          />
        </div>
      </div>

      <template v-if="story.sourceTalks.length === 0">
        <div class="p-8 text-center text-[var(--color-text-secondary)] text-sm">
          选择故事并载入以查看原文
        </div>
      </template>

      <template v-else>
        <div class="flex flex-col gap-1.5 px-2 py-1">
          <template v-for="(group, gi) in talkGroups" :key="gi">
            <div class="grid grid-cols-2 gap-2">
              <!-- ===== Source Side (merged for group) ===== -->
              <div
                class="flex flex-col justify-center p-3 rounded-lg border border-[var(--color-border)] transition-colors"
                :class="{ 'bg-[var(--color-flashback)]': flashbackItem(group.items[0].talk)?.isFlashback }"
                :title="flashbackClues(group.items[0].talk) ? '闪回线索: ' + flashbackClues(group.items[0].talk) : undefined"
              >
                <div class="flex items-center gap-3">
                  <div
                    class="w-8 h-8 rounded-full flex-shrink-0 overflow-hidden bg-[var(--color-surface)] border border-[var(--color-border)]"
                  >
                    <img
                      v-if="srcTalkCharIndex(group.items[0].talk) >= 0 && !iconErrors.has(srcTalkCharIndex(group.items[0].talk)) && !['场景', '左上场景', '选项', ''].includes(srcTalk(group.items[0].talk)?.speaker)"
                      :src="api.characterIconUrl(srcTalkCharIndex(group.items[0].talk) + 1)"
                      :alt="srcTalk(group.items[0].talk)?.speaker"
                      class="w-full h-full object-cover"
                      @error="iconErrors.add(srcTalkCharIndex(group.items[0].talk))"
                    />
                    <div
                      v-else
                      class="w-full h-full flex items-center justify-center text-white text-xs font-medium select-none"
                      style="background-color: #9ca3af"
                    >
                      {{ srcTalk(group.items[0].talk)?.speaker?.charAt(0) || '' }}
                    </div>
                  </div>

                  <div class="flex-1 min-w-0">
                    <div class="text-xs font-medium text-[var(--color-text-secondary)] mb-0.5">
                      {{ srcTalk(group.items[0].talk)?.speaker }}
                    </div>
                    <div class="leading-relaxed whitespace-pre-wrap break-words" style="font-size: var(--editor-font-size)">
                      {{ srcTalk(group.items[0].talk)?.text }}
                    </div>
                  </div>

                  <VoicePlayButton
                    v-if="srcTalk(group.items[0].talk)?.voices && srcTalk(group.items[0].talk)?.voices.length > 0"
                    :scenario-id="story.scenarioId"
                    :voice-ids="srcTalk(group.items[0].talk).voices"
                    :volume="srcTalk(group.items[0].talk).volume"
                    :source="story.selectedSource"
                  />
                </div>
              </div>

              <!-- ===== Dest Side (stacked per sub-line) ===== -->
              <div class="flex flex-col gap-1 h-full">
                <div
                  v-for="item in group.items"
                  :key="item.globalIdx"
                  :class="['p-2 rounded-lg border border-[var(--color-border)] transition-colors hover:bg-[var(--color-primary)]/[0.04]', group.items.length === 1 ? 'flex-1 flex flex-col justify-center' : '', getRowClass(item.talk), getDestBorder(item.talk) ? `border-l-4 ${getDestBorder(item.talk)}` : '']"
                >
                  <div class="flex items-start gap-2">
                    <div class="w-8 flex-shrink-0 text-xs text-[var(--color-text-secondary)] pt-1">
                      <span v-if="item.talk.start" class="font-mono">{{ item.talk.idx }}</span>
                    </div>

                    <div v-if="item.talk.start" class="w-16 flex-shrink-0 text-xs text-[var(--color-text-secondary)] pt-1 truncate">
                      {{ item.talk.speaker }}
                    </div>
                    <div v-else class="w-16 flex-shrink-0" />

                    <div
                      class="flex-1 min-w-0"
                      @contextmenu="handleContextMenu($event, item.globalIdx)"
                    >
                      <div
                        :contenteditable="item.talk.save && !['场景', '左上场景', '选项', ''].includes(item.talk.speaker)"
                        class="leading-relaxed outline-none rounded px-1 -mx-1"
                        style="font-size: var(--editor-font-size)"
                        :class="{ 'cursor-text': item.talk.save && !['场景', '左上场景', '选项', ''].includes(item.talk.speaker) }"
                        @blur="onBlur($event, item.globalIdx)"
                      >{{ item.talk.text }}</div>
                      <div v-if="item.talk.message" class="text-xs text-red-400 mt-0.5">
                        {{ item.talk.message }}
                      </div>
                    </div>

                    <div class="flex items-center gap-1 flex-shrink-0">
                      <span v-if="!item.talk.end && item.talk.save" class="text-xs text-[var(--color-text-secondary)] font-mono">\N</span>
                      <button
                        v-if="item.talk.end && !['场景', '左上场景', '选项', ''].includes(item.talk.speaker) && item.talk.save"
                        class="w-6 h-6 rounded border border-[var(--color-border)] text-xs hover:text-[var(--color-primary)]"
                        title="添加行"
                        @click="handleAddLine(item.globalIdx)"
                      >+</button>
                      <button
                        v-if="!item.talk.start"
                        class="w-6 h-6 rounded border border-[var(--color-border)] text-xs hover:bg-red-50 dark:hover:bg-red-900/30"
                        title="删除行"
                        @click="handleRemoveLine(item.globalIdx)"
                      >−</button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </template>
        </div>
      </template>
    </div>
  </div>
</template>
