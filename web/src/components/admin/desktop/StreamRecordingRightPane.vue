<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import type { StreamItem } from '../../../types/streamTree'
import type { RecordingProfile } from '../../../types/recordingSchedule'
import { fetchRemoteRecordingProfiles, saveRemoteRecordingProfile, deleteRemoteRecordingProfile } from '../../../services/recordingApi'
import ScheduleRecordingModal from './ScheduleRecordingModal.vue'

const props = defineProps<{ stream: StreamItem }>()
const emit = defineEmits<{ (e: 'saved', msg: string): void }>()
const profiles = ref<RecordingProfile[]>([]), isModalOpen = ref(false), selectedProfile = ref<RecordingProfile | null>(null)

const syncRecordMode = () => { props.stream.recordMode = profiles.value.find(p => p.isActive)?.mode || 'disabled' }
const loadProfiles = async () => { if (!props.stream?.id) return; profiles.value = await fetchRemoteRecordingProfiles(props.stream.id); syncRecordMode() }
onMounted(loadProfiles); watch(() => props.stream.id, loadProfiles)

const handleToggleActive = async (p: RecordingProfile) => {
  p.isActive = !p.isActive; await saveRemoteRecordingProfile(props.stream.id, p); syncRecordMode()
  emit('saved', p.isActive ? `[STATUS] Perfil "${p.name}" ATIVADO.` : `[STATUS] Perfil "${p.name}" DESATIVADO.`)
}
const handleOpenCreate = () => { selectedProfile.value = null; isModalOpen.value = true }
const handleOpenEdit = (p: RecordingProfile) => { selectedProfile.value = p; isModalOpen.value = true }
const handleDelete = async (id: string) => {
  await deleteRemoteRecordingProfile(props.stream.id, id); profiles.value = profiles.value.filter(p => p.id !== id); syncRecordMode()
  emit('saved', `Perfil excluido com sucesso.`)
}
const handleSaveProfile = async (profile: RecordingProfile) => {
  await saveRemoteRecordingProfile(props.stream.id, profile); await loadProfiles()
  emit('saved', `Perfil "${profile.name}" salvo com sucesso.`)
}
</script>

<template>
  <div class="vms-split-pane">
    <div class="vms-split-header">
      <div class="vms-flex-col" style="gap: 2px;">
        <span class="vms-font-bold" style="color: var(--vms-neu-accent-orange); font-size: 13px;">DIRETRIZES DE GRAVACAO</span>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">// CADASTRO DE PERFIS, MODOS E HORARIOS</span>
      </div>
      <button class="vms-btn vms-btn-primary" style="padding: 0.3rem 0.7rem; font-size: 16px; font-weight: 700; line-height: 1; display: flex; align-items: center; justify-content: center;" title="Novo Perfil de Gravacao" @click="handleOpenCreate">
        <span>+</span>
      </button>
    </div>

    <!-- Recording Profiles Table -->
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
        <tbody v-if="profiles.length > 0">
          <tr v-for="p in profiles" :key="p.id">
            <td class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange); text-align: left;">{{ p.id }}</td>
            <td style="text-align: left;">
              <span class="vms-font-semibold" style="color: #fff;">{{ p.name }}</span>
            </td>
            <td style="text-align: center;">
              <div style="display: flex; align-items: center; justify-content: center;">
                <input type="checkbox" class="vms-checkbox" :checked="p.isActive" @change="handleToggleActive(p)" />
              </div>
            </td>
            <td style="text-align: center;">
              <div style="display: flex; align-items: center; justify-content: center; gap: 6px;">
                <button class="vms-table-action-btn" title="Editar Grade de Horarios" @click="handleOpenEdit(p)">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/><path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/></svg>
                </button>
                <button class="vms-table-action-btn" title="Excluir Perfil" @click="handleDelete(p.id)">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18"/><path d="M8 6V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6"/></svg>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
        <tbody v-else>
          <tr>
            <td colspan="4" class="vms-text-mono vms-text-2xs vms-text-dim" style="text-align: center; padding: 2.5rem 1rem;">
              // NENHUM PERFIL DE GRAVAÇÃO CADASTRADO. CLIQUE EM [+] PARA CRIAR.
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Modal for Schedule & Advanced Settings -->
    <ScheduleRecordingModal :is-open="isModalOpen" :profile="selectedProfile" :next-id="`REC_0${profiles.length + 1}`" @close="isModalOpen = false" @save="handleSaveProfile" />
  </div>
</template>

