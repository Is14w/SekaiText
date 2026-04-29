<script setup lang="ts">
import { useEditorStore } from '../../stores/editor'
import { useAppStore } from '../../stores/app'
import { api } from '../../api/client'
import { useToast } from '../../composables/useToast'
import type { DstTalk } from '../../types/translation'

const editor = useEditorStore()
const app = useAppStore()
const toast = useToast()

let editTimeout: ReturnType<typeof setTimeout> | null = null

function onBlur(e: Event, idx: number) {
  handleTextChange(idx, (e.target as HTMLElement).innerText)
}

function getRowClass(talk: DstTalk): Record<string, boolean> {
  return {
    'border-l-4': true,
    'border-l-green-400': talk.proofread === true && app.showDiff,
    'border-l-yellow-400': talk.proofread === false,
    'border-l-blue-400': talk.proofread === true && talk.checked && app.editorMode === 2,
    'border-l-red-400': !talk.checked && talk.save,
    'opacity-40': !talk.save && !app.showDiff,
    'hidden': talk.proofread === false && !app.showDiff,
  }
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

async function handleAddLine(row: number) {
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
  <div class="flex items-center justify-between mb-2 px-1">
    <span class="font-semibold text-sm text-[var(--color-text-secondary)]">译文</span>
    <div class="flex items-center gap-2">
      <input
        v-model="editor.currentFilePath"
        type="text"
        placeholder="标题/路径..."
        class="text-sm px-2 py-1 rounded border border-[var(--color-border)] bg-[var(--color-surface)] w-56"
      />
    </div>
  </div>

  <div v-if="editor.talks.length === 0" class="flex-1 p-8 text-center text-[var(--color-text-secondary)] text-sm">
    载入故事后点击"载入"以创建翻译模板，或使用"打开"加载已有翻译文件
  </div>

  <div v-else class="divide-y divide-[var(--color-border)]">
    <div
      v-for="(talk, idx) in editor.talks"
      :key="idx"
      :class="getRowClass(talk)"
      class="p-2 hover:bg-[var(--color-primary)]/[0.04] transition-colors"
    >
      <div class="flex items-start gap-2">
        <div class="w-8 flex-shrink-0 text-xs text-[var(--color-text-secondary)] pt-1">
          <span v-if="talk.start" class="font-mono">{{ talk.idx }}</span>
        </div>

        <div v-if="talk.start" class="w-16 flex-shrink-0 text-xs text-[var(--color-text-secondary)] pt-1 truncate">
          {{ talk.speaker }}
        </div>
        <div v-else class="w-16 flex-shrink-0" />

        <div
          class="flex-1 min-w-0"
          @contextmenu="handleContextMenu($event, idx)"
        >
          <div
            :contenteditable="talk.save && !['场景', '左上场景', '选项', ''].includes(talk.speaker) && talk.start"
            class="text-sm leading-relaxed outline-none rounded px-1 -mx-1"
            :class="{ 'cursor-text': talk.start }"
            @blur="onBlur($event, idx)"
          >{{ talk.text }}</div>
          <div v-if="!talk.checked && talk.save && talk.text.includes('【')" class="text-xs text-red-400 mt-0.5">
            {{ talk.text.split('【')[1]?.replace('】', '') }}
          </div>
        </div>

        <div class="flex items-center gap-1 flex-shrink-0">
          <span v-if="!talk.end && talk.save" class="text-xs text-[var(--color-text-secondary)] font-mono">\N</span>
          <button
            v-if="talk.end && !['场景', '左上场景', '选项', ''].includes(talk.speaker) && talk.save"
            class="w-6 h-6 rounded border border-[var(--color-border)] text-xs hover:text-[var(--color-primary)]"
            title="添加行"
            @click="handleAddLine(idx)"
          >+</button>
          <button
            v-if="!talk.start"
            class="w-6 h-6 rounded border border-[var(--color-border)] text-xs hover:bg-red-50 dark:hover:bg-red-900/30"
            title="删除行"
            @click="handleRemoveLine(idx)"
          >−</button>
        </div>
      </div>
    </div>
  </div>
</template>
