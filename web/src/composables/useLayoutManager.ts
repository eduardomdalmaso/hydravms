import { ref, computed } from "vue"
import type { CustomLayout, GridLayout, WorkspaceSlot } from "../types/mosaic"

export function useLayoutManager(username: string = "operador") {
  const storageKey = `hydravms_layouts_${username}`
  let savedOperatorLayouts: CustomLayout[] = []
  try {
    const raw = typeof window !== "undefined" ? localStorage.getItem(storageKey) : null
    if (raw) savedOperatorLayouts = JSON.parse(raw)
  } catch {}

  const layouts = ref<CustomLayout[]>([...savedOperatorLayouts])
  const activeLayoutId = ref<string>(layouts.value[0]?.id || "")
  const activeLayout = computed(() => layouts.value.find(l => l.id === activeLayoutId.value) || null)

  const saveOperatorLayouts = () => {
    try {
      const operatorOnly = layouts.value.filter(l => !l.is_system)
      if (typeof window !== "undefined") localStorage.setItem(storageKey, JSON.stringify(operatorOnly))
    } catch {}
  }

  const createOperatorLayout = (grid: GridLayout = "2x2", initialSlots: WorkspaceSlot[] = []): CustomLayout => {
    const operatorCount = layouts.value.filter(l => !l.is_system).length + 1
    const newLayout: CustomLayout = {
      id: `lay_usr_${Date.now()}`,
      name: `Layout ${operatorCount}`,
      grid,
      is_system: false,
      created_by: username,
      slots: [...initialSlots]
    }
    layouts.value.unshift(newLayout)
    activeLayoutId.value = newLayout.id
    saveOperatorLayouts()
    return newLayout
  }

  const renameLayout = (id: string, newName: string) => {
    const target = layouts.value.find(l => l.id === id)
    if (target && !target.is_system) { target.name = newName; saveOperatorLayouts() }
  }

  const deleteLayout = (id: string) => {
    const idx = layouts.value.findIndex(l => l.id === id)
    if (idx >= 0 && !layouts.value[idx].is_system) {
      layouts.value.splice(idx, 1)
      activeLayoutId.value = layouts.value[0]?.id || ""
      saveOperatorLayouts()
    }
  }

  const duplicateLayout = (layout: CustomLayout) => { createOperatorLayout(layout.grid, layout.slots || []) }

  return {
    layouts, activeLayoutId, activeLayout, createOperatorLayout, renameLayout, deleteLayout, duplicateLayout, saveOperatorLayouts
  }
}
