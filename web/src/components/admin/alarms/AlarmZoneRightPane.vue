<script setup lang="ts">
import { ref } from 'vue'
import type { AlarmItem, AlarmZoneRule } from '../../../types/alarmTree'
import AlarmZoneRuleModal from './AlarmZoneRuleModal.vue'

const props = defineProps<{ alarm: AlarmItem }>()
const emit = defineEmits<{ (e: 'saved', msg: string): void }>()

const rules = ref<AlarmZoneRule[]>([
  { id: 'ZON_01', name: 'Notificacao Push + Snapshot de Camera', actionType: 'notification', isActive: true, delaySeconds: 0, targetOutput: 'APP_MOBILE' },
  { id: 'ZON_02', name: 'Sirene Local + Disparo de Relay PTZ', actionType: 'relay', isActive: false, delaySeconds: 0, targetOutput: 'RELAY_01' },
  { id: 'ZON_03', name: 'Trava Eletromagnetica de Portao', actionType: 'relay', isActive: true, delaySeconds: 5, targetOutput: 'RELAY_02' }
])

const isModalOpen = ref(false)
const selectedRule = ref<AlarmZoneRule | null>(null)

const handleToggleActive = (r: AlarmZoneRule) => {
  r.isActive = !r.isActive
  emit('saved', r.isActive ? `[STATUS] Regra "${r.name}" ATIVADA com sucesso.` : `[STATUS] Regra "${r.name}" DESATIVADA.`)
}

const handleOpenCreate = () => { selectedRule.value = null; isModalOpen.value = true }
const handleOpenEdit = (r: AlarmZoneRule) => { selectedRule.value = r; isModalOpen.value = true }
const handleDelete = (id: string) => {
  rules.value = rules.value.filter(r => r.id !== id)
  emit('saved', `Regra ${id} excluida com sucesso.`)
}

const handleSaveRule = (rule: AlarmZoneRule) => {
  const idx = rules.value.findIndex(r => r.id === rule.id)
  if (idx >= 0) rules.value[idx] = rule
  else rules.value.push(rule)
  emit('saved', `Regra ${rule.name} cadastrada com sucesso.`)
}
</script>

<template>
  <div class="vms-split-pane">
    <div class="vms-split-header">
      <div class="vms-flex-col" style="gap: 2px;">
        <span class="vms-font-bold" style="color: var(--vms-neu-accent-orange); font-size: 13px;">DIRETRIZES DE ZONA & DISPAROS</span>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">// CADASTRO DE ACOES, REGRAS E DISPAROS</span>
      </div>
      <button class="vms-btn vms-btn-primary" style="padding: 0.3rem 0.7rem; font-size: 16px; font-weight: 700; line-height: 1; display: flex; align-items: center; justify-content: center;" title="Novo Cadastro de Regra de Zona" @click="handleOpenCreate">
        <span>+</span>
      </button>
    </div>

    <!-- Zone Action Rules Table -->
    <div style="flex: 1; overflow-y: auto; border: 1px solid var(--vms-border); border-radius: 6px; background: #080a0e;">
      <table class="vms-table" style="font-size: 11px; width: 100%;">
        <thead>
          <tr>
            <th style="width: 75px; text-align: left;">ID</th>
            <th style="text-align: left;">NOME DO CADASTRO</th>
            <th style="width: 80px; text-align: center;">ATIVAR</th>
            <th style="width: 90px; text-align: center;">ACOES</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="r in rules" :key="r.id">
            <td class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange); text-align: left;">{{ r.id }}</td>
            <td style="text-align: left;">
              <span class="vms-font-semibold" style="color: #fff;">{{ r.name }}</span>
            </td>
            <td style="text-align: center;">
              <div style="display: flex; align-items: center; justify-content: center;">
                <input type="checkbox" class="vms-checkbox" :checked="r.isActive" @change="handleToggleActive(r)" />
              </div>
            </td>
            <td style="text-align: center;">
              <div style="display: flex; align-items: center; justify-content: center; gap: 6px;">
                <button class="vms-table-action-btn" title="Editar Regra" @click="handleOpenEdit(r)">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
                </button>
                <button class="vms-table-action-btn" title="Excluir Regra" @click="handleDelete(r.id)">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18"/><path d="M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6"/></svg>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Modal for Rule Cadastro -->
    <AlarmZoneRuleModal :is-open="isModalOpen" :rule="selectedRule" @close="isModalOpen = false" @save="handleSaveRule" />
  </div>
</template>
