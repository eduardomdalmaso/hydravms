export interface ContextMenuTarget {
  type: 'canvas' | 'folder' | 'stream'
  id?: string
  name?: string
  currentFolderId?: string | null
}

export interface SimpleFolderOption {
  id: string
  name: string
}
