<script setup lang="ts">
import { ref } from "vue"
import type { AlarmItemInfo } from "../../types/mosaic"

const alarms = ref<AlarmItemInfo[]>([
  { id: "alm_01", name: "Sensor Perimetro Norte", zone: "Zona 01 // Muro Norte", status: "online", type: "IVS" },
  { id: "alm_02", name: "Barreira Infravermelha Docas", zone: "Zona 02 // Patio Cargas", status: "alert", type: "PIR" },
  { id: "alm_03", name: "Porta Sala Servidores", zone: "Zona 03 // CPD Central", status: "online", type: "MAG" },
  { id: "alm_04", name: "Detector Fumaca Bloco B", zone: "Zona 04 // Galpao 02", status: "offline", type: "SMK" }
])

const isModalOpen = ref(false)
const newAlarmName = ref("")
const newAlarmZone = ref("")
const newAlarmType = ref("IVS")

const handleAddAlarm = () => {
  if (!newAlarmName.value) return
  alarms.value.push({
    id: `alm_0${alarms.value.length + 1}`,
    name: newAlarmName.value,
    zone: newAlarmZone.value || "Zona Geral",
    status: "online",
    type: newAlarmType.value
  })
  isModalOpen.value = false; newAlarmName.value = ""; newAlarmZone.value = ""
}
</script>

<template>
  <div class="vms-flex-col" style="gap: 1.25rem;">
    <div class="vms-flex-between">
      <div class="vms-flex-col" style="gap: 2px;">
        <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange);">CONFIGURACAO DE ALARMES & SENSORES</h3>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">GESTAO DE ZONAS, DISPARADORES E SENSORES IVS/PIR/MAG</span>
      </div>
      <button class="vms-btn vms-btn-primary" @click="isModalOpen = true">[+ NOVO ALARME]</button>
    </div>

    <!-- Alarms Table -->
    <div class="vms-admin-table-wrapper">
      <table class="vms-table">
        <thead><tr><th>ID</th><th>NOME DO ALARME</th><th>ZONA / LOCAL</th><th>TIPO</th><th>STATUS</th><th>ACOES</th></tr></thead>
        <tbody>
          <tr v-for="a in alarms" :key="a.id">
            <td class="vms-text-mono">{{ a.id }}</td>
            <td class="vms-font-semibold" style="color: #fff;">{{ a.name }}</td>
            <td class="vms-text-mono vms-text-xs vms-text-dim">{{ a.zone }}</td>
            <td><span class="vms-badge" style="background: rgba(255,255,255,0.06); color: #fff;">{{ a.type }}</span></td>
            <td>
              <div class="vms-flex-row" style="gap: 0.35rem; align-items: center;">
                <span class="vms-status-led" :class="a.status"></span>
                <span class="vms-text-mono vms-text-2xs" :style="{ color: a.status === 'online' ? '#00ff9d' : a.status === 'alert' ? '#fcee0a' : '#ff003c' }">
                  [{{ a.status.toUpperCase() }}]
                </span>
              </div>
            </td>
            <td><button class="vms-btn vms-btn-ghost vms-btn-sm" style="font-size: 11px;">[EDITAR]</button></td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Modal Novo Alarme -->
    <div v-if="isModalOpen" class="vms-modal-backdrop" @click.self="isModalOpen = false">
      <div class="vms-modal-dialog">
        <div class="vms-modal-header">
          <h3 class="vms-h3">CADASTRAR NOVO ALARME</h3>
          <button class="vms-btn vms-btn-ghost vms-btn-sm" @click="isModalOpen = false">[X]</button>
        </div>
        <div class="vms-modal-body">
          <div class="vms-form-group"><label class="vms-label">Nome do Alarme</label><input v-model="newAlarmName" class="vms-auth-input" placeholder="Ex: Sensor Perimetro Sul" /></div>
          <div class="vms-form-group"><label class="vms-label">Zona / Localizacao</label><input v-model="newAlarmZone" class="vms-auth-input" placeholder="Ex: Zona 05 // Bloco C" /></div>
          <div class="vms-form-group">
            <label class="vms-label">Tipo de Sensor</label>
            <div class="vms-flex-row" style="gap: 0.5rem;">
              <button v-for="t in ['IVS', 'PIR', 'MAG', 'SMK']" :key="t" class="vms-btn vms-btn-sm" :class="newAlarmType === t ? 'vms-btn-primary' : 'vms-btn-secondary'" @click="newAlarmType = t">{{ t }}</button>
            </div>
          </div>
        </div>
        <div class="vms-modal-footer">
          <button class="vms-btn vms-btn-secondary" @click="isModalOpen = false">CANCELAR</button>
          <button class="vms-btn vms-btn-primary" @click="handleAddAlarm">SALVAR ALARME</button>
        </div>
      </div>
    </div>
  </div>
</template>
