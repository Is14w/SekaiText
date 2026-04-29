<script setup lang="ts">
import { ref } from 'vue'
import { api } from '../../api/client'

const props = defineProps<{
  scenarioId: string
  voiceIds: string[]
  volume?: number[]
  source?: string
}>()

const playing = ref(false)
const audioRef = ref<HTMLAudioElement | null>(null)

async function play() {
  if (playing.value && audioRef.value) {
    audioRef.value.pause()
    audioRef.value = null
    playing.value = false
    return
  }

  try {
    const result = await api.voiceUrl(props.scenarioId, props.voiceIds[0], props.source || 'sekai.best')
    if (result.url) {
      const audio = new Audio(result.url)
      audio.volume = props.volume?.[0] ? props.volume[0] / 100 : 1
      audio.onended = () => { playing.value = false }
      audioRef.value = audio
      audio.play()
      playing.value = true
    }
  } catch {
    // Silent fail
  }
}
</script>

<template>
  <button
    @click="play"
    class="w-8 h-8 rounded-full border border-[var(--color-border)] flex items-center justify-center hover:text-[var(--color-primary)] transition-colors text-xs"
    :class="{ 'bg-[var(--color-primary)] text-white border-[var(--color-primary)]': playing }"
    :title="playing ? '停止' : '播放语音'"
  >
    {{ playing ? '⏹' : '▶' }}
  </button>
</template>
