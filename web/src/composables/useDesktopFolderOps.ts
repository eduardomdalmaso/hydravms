import { ref } from 'vue'

export interface FolderItemLike {
  id: string
  name: string
}

export interface FolderNodeLike<T extends FolderItemLike> {
  id: string
  name: string
  items: T[]
}

export function useFolderModalState() {
  const folderToDelete = ref<{ id: string; name: string; itemCount: number } | null>(null)
  const isConfirmDeleteOpen = ref(false)

  const openDeletePrompt = (id: string, name: string, itemCount: number) => {
    folderToDelete.value = { id, name, itemCount }
    isConfirmDeleteOpen.value = true
  }

  const closeDeletePrompt = () => {
    isConfirmDeleteOpen.value = false
    folderToDelete.value = null
  }

  return { folderToDelete, isConfirmDeleteOpen, openDeletePrompt, closeDeletePrompt }
}
