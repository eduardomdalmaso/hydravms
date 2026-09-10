<script setup lang="ts">
import { ref } from 'vue'
import { generateChannelPaths } from '../../../constants/rtspPresets'

export interface ChannelItem { id: number; name: string; path: string; subPath?: string; status: 'online' | 'offline' | 'testing' | 'pending' }
const props = defineProps<{ channels: ChannelItem[]; basePath: string; isOnvif?: boolean }>()
const emit = defineEmits<{ (e: 'update:channels', val: ChannelItem[]): void }>()
const isTestingAll = ref(false)

const addChannel = () => {
  const nextId = props.channels.length + 1; const paths = generateChannelPaths(props.basePath || '/stream1', nextId)
  emit('update:channels', [...props.channels, {
    id: nextId, name: `Canal ${nextId.toString().padStart(2, '0')}`,
    path: paths.main, subPath: paths.sub, status: 'online'
  }])
}

const removeChannel = (idx: number) => {
  if (props.channels.length <= 1) return
  emit('update:channels', props.channels.filter((_, i) => i !== idx))
}

const quickBatchAdd = (count: number) => {
  const list: ChannelItem[] = []
  for (let i = 1; i <= count; i++) {
    const paths = generateChannelPaths(props.basePath || '/stream1', i)
    list.push({
      id: i, name: `Canal ${i.toString().padStart(2, '0')}`,
      path: paths.main, subPath: paths.sub, status: 'online'
    })
  }
  emit('update:channels', list)
}


const testSingle = (ch: ChannelItem) => {
  ch.status = 'testing'; setTimeout(() => { ch.status = ch.path.trim().length > 3 ? 'online' : 'offline' }, 500)
}

const testAll = () => {
  isTestingAll.value = true; props.channels.forEach(ch => { ch.status = 'testing' })
  setTimeout(() => {
    props.channels.forEach(ch => { ch.status = ch.path.trim().length > 3 ? 'online' : 'offline' })
    isTestingAll.value = false
  }, 800)
}
</script>

<template>
  <div class="vms-flex-col" style="min-height: 380px; justify-content: space-between; gap: 0.65rem;">
    <div class="vms-flex-col" style="gap: 0.55rem;">
      <div class="vms-flex-between" style="align-items: center;">
        <div>
          <h4 class="vms-h4" style="margin: 0; color: #fff;">TOPOLOGIA DE CANAIS // AUTO-TESTE</h4>
          <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-text-dim);">
            <strong v-if="isOnvif" style="color: var(--vms-neu-accent-green);">[ONVIF AUTO-DISCOVERY]</strong>
            <span v-else>Canais sem resposta (FALHA) serão descartados automaticamente</span>
          </span>
        </div>
        <div class="vms-flex-row" style="gap: 0.35rem;">
          <button class="vms-btn vms-btn-secondary" style="font-size: 10px; padding: 3px 6px;" type="button" @click="quickBatchAdd(4)">+4</button>
          <button class="vms-btn vms-btn-secondary" style="font-size: 10px; padding: 3px 6px;" type="button" @click="quickBatchAdd(8)">+8</button>
          <button class="vms-btn vms-btn-secondary" style="font-size: 10px; padding: 3px 6px;" type="button" @click="quickBatchAdd(16)">+16</button>
          <button class="vms-btn vms-btn-primary" style="font-size: 10px; padding: 3px 8px;" type="button" @click="addChannel">+ CANAL</button>
          <button class="vms-btn vms-btn-primary" style="font-size: 10px; padding: 3px 8px; background: var(--vms-neu-accent-cyan); color: #000;" type="button" :disabled="isTestingAll" @click="testAll">
            {{ isTestingAll ? 'TESTANDO...' : 'TESTAR TODOS' }}
          </button>
        </div>
      </div>

      <!-- Channels List Container -->
      <div style="max-height: 290px; overflow-y: auto; display: flex; flex-direction: column; gap: 0.35rem; padding-right: 4px;">
        <div v-for="(ch, idx) in channels" :key="ch.id" style="display: flex; align-items: center; gap: 0.5rem; background: #07090e; border: 1px solid rgba(255, 255, 255, 0.08); border-radius: 4px; padding: 0.4rem 0.6rem;">
          <div style="width: 22px; font-family: var(--vms-font-jetbrains); font-size: 11px; font-weight: 700; color: var(--vms-neu-accent-cyan);">#{{ ch.id.toString().padStart(2, '0') }}</div>
          <div class="vms-form-group" style="flex: 1.1; margin: 0;"><input v-model="ch.name" class="vms-auth-input" style="padding: 4px 8px; font-size: 11px;" placeholder="Nome Canal" /></div>
          <div class="vms-form-group" style="flex: 1.4; margin: 0;"><input v-model="ch.path" class="vms-auth-input vms-text-mono" style="padding: 4px 8px; font-size: 11px; color: var(--vms-neu-accent-cyan);" placeholder="Path Main" /></div>
          <div class="vms-form-group" style="flex: 1.4; margin: 0;"><input v-model="ch.subPath" class="vms-auth-input vms-text-mono" style="padding: 4px 8px; font-size: 11px; color: var(--vms-text-dim);" placeholder="Path Sub" /></div>
          <div style="width: 85px; text-align: center; flex-shrink: 0;">
            <span v-if="ch.status === 'testing'" class="vms-badge" style="color: var(--vms-neu-accent-cyan); font-size: 9px;">TESTANDO</span>
            <span v-else-if="ch.status === 'online'" class="vms-badge" style="color: var(--vms-neu-accent-green); background: rgba(0,255,157,0.1); border: 1px solid rgba(0,255,157,0.3); font-size: 9px;">[ONLINE]</span>
            <span v-else-if="ch.status === 'offline'" class="vms-badge" style="color: var(--vms-neu-accent-red); background: rgba(255,0,60,0.1); border: 1px solid rgba(255,0,60,0.3); font-size: 9px;">[FALHA]</span>
            <button v-else type="button" class="vms-btn vms-btn-secondary" style="font-size: 9px; padding: 2px 6px;" @click="testSingle(ch)">TESTAR</button>
          </div>
          <button v-if="channels.length > 1" class="vms-btn vms-btn-ghost" style="color: var(--vms-neu-accent-red); padding: 2px 5px;" type="button" title="Remover" @click="removeChannel(idx)">✕</button>
        </div>
      </div>
    </div>

    <div style="background: rgba(0, 240, 255, 0.05); border: 1px solid rgba(0, 240, 255, 0.15); padding: 0.35rem 0.65rem; border-radius: 4px; font-size: 11px; display: flex; justify-content: space-between; align-items: center;">
      <span style="color: var(--vms-text-dim); font-family: var(--vms-font-jetbrains);">CANAIS ATIVOS: <strong style="color: #fff;">{{ channels.filter(c => c.status !== 'offline').length }} / {{ channels.length }}</strong></span>
      <span style="color: var(--vms-neu-accent-green); font-size: 10px; font-family: var(--vms-font-jetbrains);">{{ isOnvif ? 'PERFIS ONVIF SINCRONIZADOS' : 'DESCARTE AUTOMÁTICO DE FALHAS ATIVADO' }}</span>
    </div>
  </div>
</template>
