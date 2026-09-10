<script setup lang="ts">
import { ref } from 'vue'
import type { StorageSourceType, StorageRole, NewStoragePayload, UnallocatedDiskDevice } from '../../../types/storagePool'

const props = defineProps<{ unallocatedDisks: UnallocatedDiskDevice[]; hasBufferDisk?: boolean }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'save', payload: NewStoragePayload): void }>()

const sourceType = ref<StorageSourceType>('LOCAL_DISK')
const role = ref<StorageRole>('HOT_BUFFER')
const name = ref(''), serverHost = ref(''), sharePath = ref(''), s3Bucket = ref('')
const selectedDisk = ref(props.unallocatedDisks[0]?.devicePath || '')

const handleSubmit = () => {
  const isLocal = sourceType.value === 'LOCAL_DISK', isNas = sourceType.value === 'NETWORK_NAS'
  const disk = props.unallocatedDisks.find(d => d.devicePath === selectedDisk.value)
  emit('save', {
    name: name.value.trim() || `Pool ${sourceType.value}`, sourceType: sourceType.value, role: role.value,
    nodeOrServer: isLocal ? 'MAQUINA LOCAL // NO 01' : `${isNas ? 'SERVIDOR' : 'CLUSTER S3'} // ${serverHost.value}`,
    pathOrEndpoint: isLocal ? selectedDisk.value : (isNas ? `nfs://${serverHost.value}${sharePath.value}` : `s3://${serverHost.value}/${s3Bucket.value}`),
    filesystem: isLocal ? (disk?.filesystem || 'XFS') : (isNas ? 'NFSv4' : 'S3_API'),
    totalGb: isLocal ? (disk?.sizeGb || 1000) : (isNas ? 32000 : 50000)
  })
}
</script>

<template>
  <div class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-dialog" style="max-width: 480px;">
      <div class="vms-modal-header">
        <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange);">NOVO STORAGE // DISCO OU SERVIDOR</h3>
        <button class="vms-btn vms-btn-ghost vms-btn-sm" style="padding: 4px;" title="Fechar" @click="emit('close')">✕</button>
      </div>
      <div class="vms-modal-body vms-flex-col" style="gap: 10px;">
        <div class="vms-form-group">
          <label class="vms-label">TIPO DE ORIGEM</label>
          <select v-model="sourceType" class="vms-auth-input">
            <option value="LOCAL_DISK">[DISCO FISICO LOCAL // MAQUINA]</option>
            <option value="NETWORK_NAS">[SERVIDOR DE REDE // NAS NFS/SMB]</option>
            <option value="OBJECT_S3">[OBJETOS S3 // CLUSTER MINIO]</option>
          </select>
        </div>
        <div class="vms-form-group">
          <label class="vms-label">NOME DE IDENTIFICACAO</label>
          <input v-model="name" class="vms-auth-input" placeholder="Ex: SSD Local Buffer Quente" />
        </div>
        <div class="vms-form-group">
          <label class="vms-label">FUNCAO / TIER OPERACIONAL</label>
          <select v-model="role" class="vms-auth-input">
            <option value="HOT_BUFFER">[HOT_BUFFER] DISCO DE BUFFER NVMe (ESCRITA RAPIDA)</option>
            <option value="WARM_ARCHIVE">[WARM_ARCHIVE] ARMAZENAMENTO DE GRAVACOES (LONGO PRAZO)</option>
            <option value="SNAPSHOTS">[SNAPSHOTS] ARMAZENAMENTO DE FOTOS/ALARMES</option>
            <option value="DATABASE">[DATABASE] BANCO DE DADOS & METADADOS</option>
          </select>
        </div>
        <div v-if="sourceType === 'LOCAL_DISK'" class="vms-form-group">
          <label class="vms-label">DISPOSITIVO FISICO DETECTADO</label>
          <select v-model="selectedDisk" class="vms-auth-input">
            <option v-for="d in unallocatedDisks" :key="d.devicePath" :value="d.devicePath">
              [{{ d.busType }}] {{ d.devicePath }} // {{ d.model }} ({{ d.sizeGb }} GB)
            </option>
          </select>
        </div>
        <template v-else>
          <div class="vms-form-group">
            <label class="vms-label">HOST / IP DO SERVIDOR</label>
            <input v-model="serverHost" class="vms-auth-input" placeholder="Ex: 192.168.1.100 ou minio.corp.lan" />
          </div>
          <div v-if="sourceType === 'NETWORK_NAS'" class="vms-form-group">
            <label class="vms-label">CAMINHO DO COMPARTILHAMENTO (NFS/SMB)</label>
            <input v-model="sharePath" class="vms-auth-input" placeholder="/volume1/cftv_archive" />
          </div>
          <div v-else class="vms-form-group">
            <label class="vms-label">BUCKET S3 / MINIO</label>
            <input v-model="s3Bucket" class="vms-auth-input" placeholder="hydravms-recordings" />
          </div>
        </template>
      </div>
      <div class="vms-modal-footer">
        <button class="vms-btn vms-btn-secondary" @click="emit('close')">CANCELAR</button>
        <button class="vms-btn vms-btn-primary" @click="handleSubmit">SALVAR</button>
      </div>
    </div>
  </div>
</template>
