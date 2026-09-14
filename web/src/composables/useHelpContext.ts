import { ref } from 'vue'

const isHelpOpen = ref(false)
const currentHelpPageId = ref('video_streams')

export function useHelpContext() {
  const setHelpPageId = (id: string) => {
    currentHelpPageId.value = id
  }

  const toggleHelp = () => {
    isHelpOpen.value = !isHelpOpen.value
  }

  const closeHelp = () => {
    isHelpOpen.value = false
  }

  const openHelpForPage = (id?: string) => {
    if (id) currentHelpPageId.value = id
    isHelpOpen.value = true
  }

  return {
    isHelpOpen,
    currentHelpPageId,
    setHelpPageId,
    toggleHelp,
    closeHelp,
    openHelpForPage
  }
}
