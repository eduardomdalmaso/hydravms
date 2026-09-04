<script setup lang="ts">
import { ref } from "vue"

defineProps<{
  username: string
  isAdmin: boolean
  currentMode: "vms" | "admin"
  isOpen: boolean
}>()

const emit = defineEmits<{
  (e: "switchMode", mode: "vms" | "admin"): void
  (e: "logout"): void
  (e: "close"): void
}>()

const currentLang = ref<"PT" | "EN" | "ES">("PT")
</script>

<template>
  <div v-if="isOpen" class="vms-user-dropdown-menu">
    <!-- Compact User Profile Header -->
    <div style="padding: 0.4rem 0.6rem; border-bottom: 1px solid rgba(255,255,255,0.06); display: flex; align-items: center; gap: 0.55rem;">
      <div style="width: 24px; height: 24px; border-radius: 50%; background: #1a1e26; border: 1px solid var(--vms-neu-accent-orange); display: flex; align-items: center; justify-content: center; position: relative;">
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path><circle cx="12" cy="7" r="4"></circle></svg>
        <span class="vms-status-led online" style="position: absolute; bottom: -1px; right: -1px; width: 5px; height: 5px;"></span>
      </div>
      <span style="color: #ffffff; font-family: var(--vms-font-roboto); font-size: 12px; font-weight: 700; text-transform: uppercase;">{{ username }}</span>
    </div>

    <!-- Navigation Options -->
    <div style="display: flex; flex-direction: column; gap: 0.2rem; padding: 0.25rem 0;">
      <button class="vms-user-dropdown-item" :class="{ active: currentMode === 'vms' }" @click="emit('switchMode', 'vms')">
        <div class="vms-flex-row" style="gap: 0.45rem; align-items: center;">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="7" height="7"></rect><rect x="14" y="3" width="7" height="7"></rect><rect x="14" y="14" width="7" height="7"></rect><rect x="3" y="14" width="7" height="7"></rect></svg>
          <span style="font-family: var(--vms-font-roboto); font-weight: 500;">[PAINEL VMS]</span>
        </div>
        <span v-if="currentMode === 'vms'" class="vms-badge vms-badge-online" style="font-size: 8px; padding: 1px 4px;">[ATIVO]</span>
      </button>

      <button v-if="isAdmin" class="vms-user-dropdown-item" :class="{ active: currentMode === 'admin' }" @click="emit('switchMode', 'admin')">
        <div class="vms-flex-row" style="gap: 0.45rem; align-items: center;">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"></circle><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1-2.83 2.83l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 1 1-4 0v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 1 1-2.83-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 1 1 0-4h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 1 1 2.83-2.83l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 1 1 4 0v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 1 1 2.83 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 1 1 0 4h-.09a1.65 1.65 0 0 0-1.51 1z"></path></svg>
          <span style="font-family: var(--vms-font-roboto); font-weight: 500;">[ADMIN CENTER]</span>
        </div>
        <span v-if="currentMode === 'admin'" class="vms-badge vms-badge-info" style="font-size: 8px; padding: 1px 4px;">[ATIVO]</span>
      </button>
    </div>

    <!-- Bottom Action Row: Language Selectors + Logout Icon -->
    <div class="vms-flex-between" style="border-top: 1px solid rgba(255,255,255,0.06); padding: 0.35rem 0.25rem 0.1rem; align-items: center;">
      <div class="vms-flex-row" style="gap: 0.25rem; align-items: center;">
        <button
          v-for="l in (['PT', 'EN', 'ES'] as const)"
          :key="l"
          class="vms-btn vms-btn-sm"
          :style="{
            padding: '2px 5px',
            fontSize: '9px',
            fontWeight: currentLang === l ? '700' : '500',
            color: currentLang === l ? 'var(--vms-neu-accent-orange)' : '#94a3b8',
            background: currentLang === l ? 'rgba(255,94,58,0.15)' : 'transparent',
            border: currentLang === l ? '1px solid rgba(255,94,58,0.35)' : '1px solid transparent',
            borderRadius: '4px'
          }"
          :title="l === 'PT' ? 'Português (Brasil)' : l === 'EN' ? 'English (US)' : 'Español'"
          @click="currentLang = l"
        >
          {{ l === 'PT' ? 'PT-BR' : l === 'EN' ? 'US' : 'ES' }}
        </button>
      </div>

      <!-- Logout Icon Button -->
      <button
        class="vms-sidebar-icon-btn"
        style="padding: 4px; color: #ff003c; border-radius: 4px;"
        title="Desconectar / Encerrar Sessão"
        @click="emit('logout')"
      >
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
          <polyline points="16 17 21 12 16 7"></polyline>
          <line x1="21" y1="12" x2="9" y2="12"></line>
        </svg>
      </button>
    </div>
  </div>
</template>
