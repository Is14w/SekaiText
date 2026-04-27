import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from '../api/client'
import type { SourceTalk } from '../types/translation'

export const useStoryStore = defineStore('story', () => {
  const storyTypes = ref<string[]>([])
  const sorts = ref<{ label: string; value: string }[]>([])
  const indices = ref<{ label: string; value: string; chapters?: number[] }[]>([])
  const chapters = ref<{ number: number; label: string }[]>([])

  const selectedType = ref('')
  const selectedSort = ref('')
  const selectedIndex = ref('')
  const selectedChapter = ref(0)
  const selectedSource = ref('haruki')

  const scenarioId = ref('')
  const sourceTalks = ref<SourceTalk[]>([])
  const loading = ref(false)

  async function fetchTypes() {
    storyTypes.value = await api.storyTypes()
  }

  async function fetchSorts(type: string) {
    sorts.value = await api.storySorts(type)
  }

  async function fetchIndex(type: string, sort: string) {
    indices.value = await api.storyIndex(type, sort)
  }

  async function fetchChapters(type: string, sort: string, index: string) {
    chapters.value = await api.storyChapter(type, sort, index)
  }

  async function loadStory() {
    loading.value = true
    try {
      const result = await api.storyLoad({
        storyType: selectedType.value,
        sort: selectedSort.value,
        index: selectedIndex.value,
        chapter: selectedChapter.value,
        source: selectedSource.value,
      })
      scenarioId.value = result.scenarioId
      sourceTalks.value = result.sourceTalks
    } finally {
      loading.value = false
    }
  }

  async function loadStoryLocal(content: string) {
    loading.value = true
    try {
      const result = await api.storyLoadLocal(content)
      scenarioId.value = result.scenarioId
      sourceTalks.value = result.sourceTalks
      selectedType.value = ''
      selectedSort.value = ''
      selectedIndex.value = ''
      selectedChapter.value = 0
    } finally {
      loading.value = false
    }
  }

  return {
    storyTypes, sorts, indices, chapters,
    selectedType, selectedSort, selectedIndex, selectedChapter, selectedSource,
    scenarioId, sourceTalks, loading,
    fetchTypes, fetchSorts, fetchIndex, fetchChapters, loadStory, loadStoryLocal,
  }
})
