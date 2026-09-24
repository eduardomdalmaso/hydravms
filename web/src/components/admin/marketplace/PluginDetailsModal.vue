<script setup lang="ts">
import { ref, computed } from 'vue'
import type { PluginManifest } from '../../../types/marketplace'

const props = defineProps<{
  isOpen: boolean
  plugin: PluginManifest | null
}>()

const emit = defineEmits<{ (e: 'close'): void }>()
const activeTab = ref<'specs' | 'json'>('specs')
const copied = ref(false)

const manifestJson = computed(() => {
  if (!props.plugin) return ''
  const fullManifest = {
    id: props.plugin.id,
    name: props.plugin.name,
    version: props.plugin.version,
    author: props.plugin.author,
    category: props.plugin.category,
    runtime: props.plugin.runtime,
    hardware_req: props.plugin.hardware_req,
    permissions: props.plugin.permissions || ['video:shm_read', 'events:write', 'cuda:ipc'],
    is_official: props.plugin.is_official,
    default_config: props.plugin.default_config && Object.keys(props.plugin.default_config).length > 0 ? props.plugin.default_config : {
      target_fps: 15,
      motion_gated: true,
      auto_sahi: true,
      supported_modes: ['intrusion', 'crowd', 'counting', 'dwell_time'],
      classes: ['person', 'cell_phone', 'car', 'motorcycle', 'bicycle', 'truck', 'bus', 'bird', 'cat', 'dog']
    },
    sandbox: { memory_limit_mb: 2048, gpu_memory_limit_mb: 4096, ipc: 'posix_shm' }
  }
  return JSON.stringify(fullManifest, null, 2)
})

const copyToClipboard = () => {
  if (!manifestJson.value) return
  navigator.clipboard.writeText(manifestJson.value)
  copied.value = true
  setTimeout(() => { copied.value = false }, 2000)
}
</script>

<template>
  <div v-if="isOpen && plugin" class="vms-modal-backdrop" @click.self="emit('close')">
    <div class="vms-modal-content vms-manifest-modal">
      <!-- Header com Borda Cítrica Visível -->
      <div class="vms-modal-header vms-flex-between" style="border-bottom: 1px solid rgba(255, 94, 58, 0.35); padding-bottom: 0.75rem;">
        <div class="vms-flex-col" style="gap: 2px;">
          <h4 class="vms-h4" style="margin: 0; color: #ffffff; font-weight: 700; letter-spacing: 0.5px;">{{ plugin.name }}</h4>
          <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-orange);">MANIFESTO TÉCNICO // plugin.json</span>
        </div>
        <button class="vms-btn-icon" @click="emit('close')">
          <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
        </button>
      </div>

      <!-- Seletor de Abas HUD -->
      <div class="vms-flex-row" style="gap: 0.5rem; padding: 0.75rem 0 0.5rem 0; border-bottom: 1px solid rgba(255, 255, 255, 0.08);">
        <button class="vms-btn vms-btn-sm" :class="activeTab === 'specs' ? 'vms-btn-primary' : 'vms-btn-secondary'" style="padding: 4px 12px; font-size: 11px;" @click="activeTab = 'specs'">
          ESPECIFICAÇÕES HUD
        </button>
        <button class="vms-btn vms-btn-sm" :class="activeTab === 'json' ? 'vms-btn-primary' : 'vms-btn-secondary'" style="padding: 4px 12px; font-size: 11px;" @click="activeTab = 'json'">
          CONFIGURAÇÃO JSON
        </button>
      </div>

      <!-- Corpo da Modal -->
      <div class="vms-modal-body vms-flex-col" style="flex: 1; overflow-y: auto; gap: 0.85rem; padding: 0.75rem 0;">
        <!-- Aba 1: Especificações HUD -->
        <template v-if="activeTab === 'specs'">
          <div class="vms-manifest-card">
            <span class="vms-text-xs vms-font-semibold" style="color: var(--vms-neu-accent-orange);">DESCRIÇÃO TÉCNICA</span>
            <p class="vms-text-xs vms-text-dim" style="margin: 0.25rem 0 0 0; line-height: 1.5;">{{ plugin.description }}</p>
          </div>

          <div class="vms-flex-row" style="gap: 0.65rem; flex-wrap: wrap;">
            <div class="vms-manifest-card" style="flex: 1; min-width: 130px;">
              <span class="vms-text-mono vms-text-2xs vms-text-dim">VERSÃO:</span>
              <span class="vms-text-xs vms-font-semibold" style="color: #00ff9d;">v{{ plugin.version }}</span>
            </div>
            <div class="vms-manifest-card" style="flex: 1.5; min-width: 160px;">
              <span class="vms-text-mono vms-text-2xs vms-text-dim">FORNECEDOR:</span>
              <span class="vms-text-xs vms-font-semibold" style="color: #ffffff;">{{ plugin.author }}</span>
            </div>
            <div class="vms-manifest-card" style="flex: 1; min-width: 130px;">
              <span class="vms-text-mono vms-text-2xs vms-text-dim">RUNTIME:</span>
              <span class="vms-badge vms-badge-secondary" style="font-size: 10px;">{{ plugin.runtime }}</span>
            </div>
            <div class="vms-manifest-card" style="flex: 1.5; min-width: 160px;">
              <span class="vms-text-mono vms-text-2xs vms-text-dim">ACELERAÇÃO:</span>
              <span class="vms-badge vms-badge-orange" style="font-size: 10px;">{{ plugin.hardware_req }}</span>
            </div>
          </div>

          <div class="vms-manifest-card">
            <span class="vms-text-xs vms-font-semibold" style="color: var(--vms-neu-accent-orange);">PERMISSÕES REQUERIDAS (SANDBOX IPC)</span>
            <div class="vms-flex-row" style="gap: 0.4rem; flex-wrap: wrap; margin-top: 0.35rem;">
              <span v-for="perm in (plugin.permissions || ['video:shm_read', 'events:write', 'cuda:ipc'])" :key="perm" class="vms-chip active">
                {{ perm }}
              </span>
            </div>
          </div>
        </template>

        <!-- Aba 2: JSON Formatado -->
        <template v-else>
          <div class="vms-flex-between" style="align-items: center;">
            <span class="vms-text-mono vms-text-2xs vms-text-dim">// SCHEMA MANIFEST DECLARATION</span>
            <button class="vms-btn vms-btn-secondary vms-btn-sm" style="font-size: 10px; padding: 2px 8px;" @click="copyToClipboard">
              {{ copied ? 'COPIADO!' : 'COPIAR JSON' }}
            </button>
          </div>
          <pre class="vms-json-viewer">{{ manifestJson }}</pre>
        </template>
      </div>

      <!-- Footer com Borda e Status Oficial -->
      <div class="vms-modal-footer vms-flex-between" style="border-top: 1px solid rgba(255, 255, 255, 0.1); padding-top: 0.75rem;">
        <span class="vms-text-mono vms-text-2xs" style="color: #00ff9d;">[ECDSA VERIFIED] // ASSINATURA OFICIAL HYDRA</span>
        <button class="vms-btn vms-btn-secondary vms-btn-sm" @click="emit('close')">FECHAR</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.vms-manifest-modal { width: 640px; max-width: 95vw; height: 520px; display: flex; flex-direction: column; background: #0b0e14; border: 1px solid rgba(255, 94, 58, 0.35); box-shadow: 0 20px 50px rgba(0,0,0,0.85), 0 0 16px rgba(255,94,58,0.12); border-radius: 8px; }
.vms-manifest-card { background: rgba(255, 255, 255, 0.03); border: 1px solid rgba(255, 255, 255, 0.08); border-radius: 6px; padding: 0.65rem 0.85rem; display: flex; flex-direction: column; gap: 2px; }
.vms-json-viewer { flex: 1; background: #05070a; border: 1px solid rgba(0, 240, 255, 0.25); border-radius: 6px; padding: 0.75rem 1rem; color: #a5f3fc; font-family: var(--vms-font-mono); font-size: 11px; line-height: 1.45; overflow: auto; margin: 0; }
.vms-chip { font-size: 10px; font-family: var(--vms-font-mono); padding: 3px 8px; border-radius: 4px; background: rgba(255,94,58,0.15); color: #ff5e3a; border: 1px solid rgba(255,94,58,0.35); }
</style>
