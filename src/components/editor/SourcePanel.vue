<script setup lang="ts">
import { computed, ref } from 'vue'
import { useStoryStore } from '../../stores/story'
import { useAppStore } from '../../stores/app'
import { api } from '../../api/client'
import VoicePlayButton from './VoicePlayButton.vue'

const iconErrors = ref<Set<number>>(new Set())

const story = useStoryStore()
const app = useAppStore()

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
</script>

<template>
  <div class="flex items-center justify-between mb-2 px-1">
    <span class="font-semibold text-sm text-[var(--color-text-secondary)]">原文</span>
    <span v-if="story.scenarioId" class="text-xs text-[var(--color-text-secondary)]">{{ story.scenarioId }}</span>
  </div>

  <div v-if="story.sourceTalks.length === 0" class="flex-1 p-8 text-center text-[var(--color-text-secondary)] text-sm">
    选择故事并载入以查看原文
  </div>

  <div v-else class="divide-y divide-[var(--color-border)]">
    <div
      v-for="(talk, idx) in talksWithFlashback"
      :key="idx"
      class="p-3 transition-colors"
      :class="{ 'bg-[var(--color-flashback)]': talk.isFlashback }"
      :title="talk.isFlashback && talk.clues ? '闪回线索: ' + talk.clues.filter(c => c).join(', ') : undefined"
    >
      <div class="flex items-start gap-3">
        <div
          class="w-8 h-8 rounded-full flex-shrink-0 overflow-hidden bg-[var(--color-surface)] border border-[var(--color-border)]"
        >
          <img
            v-if="talk.charIndex >= 0 && !iconErrors.has(talk.charIndex) && !['场景', '左上场景', '选项', ''].includes(talk.speaker)"
            :src="api.characterIconUrl(talk.charIndex + 1)"
            :alt="talk.speaker"
            class="w-full h-full object-cover"
            @error="iconErrors.add(talk.charIndex)"
          />
          <div
            v-else
            class="w-full h-full flex items-center justify-center text-white text-xs font-medium"
            style="background-color: #9ca3af"
          >
            {{ talk.speaker.charAt(0) }}
          </div>
        </div>

        <div class="flex-1 min-w-0">
          <div class="text-xs font-medium text-[var(--color-text-secondary)] mb-0.5">
            {{ talk.speaker }}
          </div>
          <div class="text-sm leading-relaxed whitespace-pre-wrap break-words">{{ talk.text }}</div>
        </div>

        <VoicePlayButton
          v-if="talk.voices && talk.voices.length > 0"
          :scenario-id="story.scenarioId"
          :voice-ids="talk.voices"
          :volume="talk.volume"
          :source="story.selectedSource"
        />
      </div>
    </div>
  </div>
</template>
