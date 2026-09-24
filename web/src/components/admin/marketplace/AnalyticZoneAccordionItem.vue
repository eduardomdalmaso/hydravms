<script setup lang="ts">
import { ref } from 'vue'
import type { ZoneConfig, AnalyticMode, ScheduleInterval } from '../../../types/marketplace'

const props = defineProps<{
  zone: ZoneConfig
  color: string
  isOpen: boolean
  canDelete: boolean
}>()

const emit = defineEmits<{
  (e: 'toggle'): void
  (e: 'delete'): void
  (e: 'update:name', val: string): void
  (e: 'update:mode', val: AnalyticMode): void
  (e: 'toggle:class', clsId: string): void
  (e: 'add:schedule'): void
  (e: 'delete:schedule', idx: number): void
  (e: 'update:scheduleTime', payload: { idx: number; field: 'start_time' | 'end_time'; val: string }): void
  (e: 'toggle:scheduleDay', payload: { idx: number; day: number }): void
}>()

const isObjectsDropdownOpen = ref(false)

const availableClasses = [
  { id: 'person', label: 'Pessoa' }, { id: 'cell_phone', label: 'Celular' },
  { id: 'car', label: 'Carro' }, { id: 'motorcycle', label: 'Moto' },
  { id: 'bicycle', label: 'Bicicleta' }, { id: 'truck', label: 'Caminhão' },
  { id: 'bus', label: 'Ônibus' }, { id: 'bird', label: 'Pássaros' },
  { id: 'cat', label: 'Gato' }, { id: 'dog', label: 'Cachorro' }
]

const modes: { id: AnalyticMode; label: string }[] = [
  { id: 'intrusion', label: 'INTRUSÃO // POLÍGONO DE SEGURANÇA' },
  { id: 'crowd', label: 'MULTIDÃO // DENSIDADE E OCUPAÇÃO' },
  { id: 'counting', label: 'CONTAGEM // LINHA BIDIRECIONAL (A <-> B)' },
  { id: 'dwell_time', label: 'TEMPO EXCEDIDO // PERMANÊNCIA NA ÁREA' }
]

const days = ['Dom', 'Seg', 'Ter', 'Qua', 'Qui', 'Sex', 'Sáb']
</script>

<template>
  <div class="vms-zone-item" :style="{ borderColor: isOpen ? color : 'var(--vms-border)' }">
    <!-- Header da Zona -->
    <div class="vms-flex-between vms-zone-header" @click="emit('toggle')">
      <div class="vms-flex-row" style="gap: 0.5rem; align-items: center;">
        <span class="vms-zone-dot" :style="{ backgroundColor: color }"></span>
        <span class="vms-text-xs vms-font-semibold" style="color: #ffffff;">{{ zone.name.toUpperCase() }}</span>
        <span class="vms-badge vms-badge-secondary" style="font-size: 8.5px;">{{ zone.mode.toUpperCase() }}</span>
      </div>
      <div class="vms-flex-row" style="gap: 0.4rem; align-items: center;" @click.stop>
        <button v-if="canDelete" class="vms-btn-icon" style="color: #ff003c; padding: 2px;" title="Excluir Zona" @click="emit('delete')">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>
        </button>
        <button class="vms-btn-icon" style="padding: 2px;" @click="emit('toggle')">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" :style="{ transform: isOpen ? 'rotate(180deg)' : 'rotate(0deg)' }"><polyline points="6 9 12 15 18 9"/></svg>
        </button>
      </div>
    </div>

    <!-- Conteúdo da Zona Expandido -->
    <div v-if="isOpen" class="vms-zone-body">
      <!-- Nome e Modo em Dropdown -->
      <div class="vms-flex-col" style="gap: 0.4rem;">
        <div class="vms-flex-col" style="gap: 2px;">
          <label class="vms-text-mono vms-text-2xs vms-text-dim">IDENTIFICADOR:</label>
          <input :value="zone.name" class="vms-auth-input" style="padding: 4px 6px; font-size: 11px;" @input="emit('update:name', ($event.target as HTMLInputElement).value)" />
        </div>

        <div class="vms-flex-col" style="gap: 2px;">
          <label class="vms-text-mono vms-text-2xs vms-text-dim">MODO DA REGRA (LISTAGEM):</label>
          <select :value="zone.mode" class="vms-auth-input" style="padding: 4px 6px; font-size: 11px;" @change="emit('update:mode', ($event.target as HTMLSelectElement).value as AnalyticMode)">
            <option v-for="m in modes" :key="m.id" :value="m.id">[{{ m.label }}]</option>
          </select>
        </div>
      </div>

      <!-- Objetos Alvo em Listagem Selecionável -->
      <div class="vms-flex-col" style="gap: 2px; position: relative;">
        <label class="vms-text-mono vms-text-2xs vms-text-dim">OBJETOS ALVO (LISTAGEM):</label>
        <button class="vms-auth-input vms-flex-between" style="padding: 4px 8px; font-size: 10.5px; text-align: left; cursor: pointer;" @click="isObjectsDropdownOpen = !isObjectsDropdownOpen">
          <span style="overflow: hidden; text-overflow: ellipsis; white-space: nowrap; color: #ffffff;">
            {{ zone.target_classes.length > 0 ? zone.target_classes.map(c => availableClasses.find(a => a.id === c)?.label || c).join(', ') : '[NENHUM SELECIONADO]' }}
          </span>
          <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange);">▼</span>
        </button>

        <div v-if="isObjectsDropdownOpen" class="vms-card" style="position: absolute; top: calc(100% + 2px); left: 0; right: 0; z-index: 30; padding: 0.5rem; background: #0c0f17; border: 1px solid var(--vms-border); max-height: 140px; overflow-y: auto; display: flex; flex-direction: column; gap: 4px;">
          <label v-for="cls in availableClasses" :key="cls.id" class="vms-flex-row vms-text-xs" style="gap: 0.4rem; align-items: center; cursor: pointer; padding: 2px 4px;">
            <input type="checkbox" :checked="zone.target_classes.includes(cls.id)" style="accent-color: var(--vms-neu-accent-orange);" @change="emit('toggle:class', cls.id)" />
            <span :style="{ color: zone.target_classes.includes(cls.id) ? '#ff5e3a' : '#cbd5e1' }">{{ cls.label }}</span>
          </label>
        </div>
      </div>

      <!-- Múltiplos Horários de Ativação -->
      <div class="vms-flex-col" style="gap: 0.35rem; border-top: 1px solid rgba(255,255,255,0.06); padding-top: 0.45rem;">
        <div class="vms-flex-between" style="align-items: center;">
          <label class="vms-text-mono vms-text-2xs vms-text-dim">HORÁRIOS DE ATIVAÇÃO:</label>
          <button class="vms-btn vms-btn-ghost vms-btn-sm" style="font-size: 9px; padding: 1px 6px; color: var(--vms-neu-accent-orange);" @click="emit('add:schedule')">+ ADICIONAR HORÁRIO</button>
        </div>

        <div v-for="(sch, sIdx) in (zone.schedules || [])" :key="sch.id || sIdx" class="vms-schedule-row-inline">
          <div class="vms-flex-row" style="gap: 0.25rem; align-items: center; flex-shrink: 0;">
            <input :value="sch.start_time" type="time" class="vms-auth-input" style="padding: 2px 4px; font-size: 10px; width: 70px;" @input="emit('update:scheduleTime', { idx: sIdx, field: 'start_time', val: ($event.target as HTMLInputElement).value })" />
            <span class="vms-text-mono vms-text-2xs vms-text-dim">→</span>
            <input :value="sch.end_time" type="time" class="vms-auth-input" style="padding: 2px 4px; font-size: 10px; width: 70px;" @input="emit('update:scheduleTime', { idx: sIdx, field: 'end_time', val: ($event.target as HTMLInputElement).value })" />
          </div>

          <div class="vms-flex-row" style="gap: 2px; flex: 1; justify-content: center; flex-wrap: nowrap;">
            <span v-for="(dName, dIdx) in days" :key="dIdx" class="vms-chip" :class="{ active: sch.days.includes(dIdx) }" style="font-size: 8.5px; padding: 2px 4px;" @click="emit('toggle:scheduleDay', { idx: sIdx, day: dIdx })">{{ dName }}</span>
          </div>

          <button v-if="(zone.schedules || []).length > 1" class="vms-btn-icon" style="color: #ff003c; padding: 2px; flex-shrink: 0;" title="Remover Intervalo" @click="emit('delete:schedule', sIdx)">
            <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.vms-zone-item { background: rgba(0,0,0,0.3); border: 1px solid var(--vms-border); border-radius: 5px; overflow: hidden; transition: border-color 0.15s; }
.vms-zone-header { padding: 0.45rem 0.65rem; background: rgba(255,255,255,0.02); cursor: pointer; user-select: none; }
.vms-zone-dot { width: 8px; height: 8px; border-radius: 50%; display: inline-block; flex-shrink: 0; box-shadow: 0 0 6px currentColor; }
.vms-zone-body { padding: 0.6rem 0.7rem; display: flex; flex-direction: column; gap: 0.55rem; border-top: 1px solid rgba(255,255,255,0.06); }
.vms-schedule-row-inline { display: flex; align-items: center; justify-content: space-between; gap: 0.45rem; background: rgba(255,255,255,0.02); padding: 0.35rem 0.5rem; border-radius: 4px; border: 1px dashed rgba(255,255,255,0.08); width: 100%; box-sizing: border-box; }
.vms-chip { font-size: 8.5px; font-family: var(--vms-font-mono); padding: 2px 4px; border-radius: 3px; background: rgba(255,255,255,0.06); color: #8b94a0; cursor: pointer; border: 1px solid transparent; user-select: none; }
.vms-chip.active { background: rgba(255,94,58,0.25); color: #ff5e3a; border-color: rgba(255,94,58,0.4); font-weight: bold; }
</style>
