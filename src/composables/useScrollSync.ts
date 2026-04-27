import { ref, type Ref } from 'vue'

/**
 * Syncs scroll position between source and dest panels.
 */
export function useScrollSync() {
  const sourceRef: Ref<HTMLElement | null> = ref(null)
  const destRef: Ref<HTMLElement | null> = ref(null)

  let syncing = false

  function onSourceScroll() {
    if (syncing || !sourceRef.value || !destRef.value) return
    syncing = true
    const ratio = sourceRef.value.scrollTop / (sourceRef.value.scrollHeight - sourceRef.value.clientHeight)
    destRef.value.scrollTop = ratio * (destRef.value.scrollHeight - destRef.value.clientHeight)
    requestAnimationFrame(() => { syncing = false })
  }

  function onDestScroll() {
    if (syncing || !sourceRef.value || !destRef.value) return
    syncing = true
    const ratio = destRef.value.scrollTop / (destRef.value.scrollHeight - destRef.value.clientHeight)
    sourceRef.value.scrollTop = ratio * (sourceRef.value.scrollHeight - sourceRef.value.clientHeight)
    requestAnimationFrame(() => { syncing = false })
  }

  return {
    sourceRef,
    destRef,
    onSourceScroll,
    onDestScroll,
  }
}
