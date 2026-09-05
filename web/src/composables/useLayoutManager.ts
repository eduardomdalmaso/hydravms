import { ref, computed } from "vue"
import type { CustomLayout, GridLayout, WorkspaceSlot } from "../types/mosaic"

const defaultSystemLayouts: CustomLayout[] = [
  { id: "lay_sys_01", name: "Grade Padrao (2x2)", grid: "2x2", is_system: true, created_by: "admin" },
  { id: "lay_alpha_01", name: "[ALPHA] Mural Operacional 2x2", grid: "2x2", is_system: true, created_by: "admin" },
  { id: "lay_vv_01", name: "[VISTA VERDE] Portaria 2x2", grid: "2x2", is_system: true, created_by: "admin" },
  { id: "lay_sys_02", name: "Destaque Hero (1+5)", grid: "1+5", is_system: true, created_by: "admin" },
  { id: "lay_sys_03", name: "Mural (3x3)", grid: "3x3", is_system: true, created_by: "admin" },
  { id: "lay_sys_04", name: "Central (4x4)", grid: "4x4", is_system: true, created_by: "admin" },
  { id: "lay_sys_05", name: "Hero Panoramico (1+7)", grid: "1+7", is_system: true, created_by: "admin" },
  { id: "lay_sys_06", name: "Matrix Master (1+12)", grid: "1+12", is_system: true, created_by: "admin" },
  { id: "lay_sys_07", name: "Super Wall (8x8)", grid: "8x8", is_system: true, created_by: "admin" },
  { id: "lay_sys_08", name: "Mega Wall (10x10)", grid: "10x10", is_system: true, created_by: "admin" }
]

export function useLayoutManager(username: string = "operador") {
  const storageKey = `hydravms_layouts_${username}`
  let savedOperatorLayouts: CustomLayout[] = []
  try {
    const raw = typeof window !== "undefined" ? localStorage.getItem(storageKey) : null
    if (raw) savedOperatorLayouts = JSON.parse(raw)
  } catch {}

  const layouts = ref<CustomLayout[]>([...savedOperatorLayouts, ...defaultSystemLayouts])
  const activeLayoutId = ref<string>(layouts.value[0]?.id || "lay_sys_01")
  const activeLayout = computed(() => layouts.value.find(l => l.id === activeLayoutId.value) || layouts.value[0])

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
      activeLayoutId.value = layouts.value[0]?.id || "lay_sys_01"
      saveOperatorLayouts()
    }
  }

  const duplicateLayout = (layout: CustomLayout) => { createOperatorLayout(layout.grid, layout.slots || []) }

  return {
    layouts, activeLayoutId, activeLayout, createOperatorLayout, renameLayout, deleteLayout, duplicateLayout, saveOperatorLayouts
  }
}
