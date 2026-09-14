<script setup lang="ts">
import { ref } from 'vue'
import type { BrandingConfig } from '../../../composables/useBranding'

const props = defineProps<{ modelValue: BrandingConfig }>()
const emit = defineEmits<{ (e: 'update:modelValue', val: BrandingConfig): void; (e: 'upload', type: 'headerLogo' | 'faviconUrl' | 'loginLogo', file: File): void }>()

const colorPresets = ['#ff5e3a', '#3b82f6', '#00f0ff', '#00ff9d', '#eab308', '#a855f7', '#ec4899', '#ffffff']

const triggerUpload = (type: 'headerLogo' | 'faviconUrl' | 'loginLogo') => {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = 'image/svg+xml, image/png, image/jpeg, image/x-icon, image/webp'
  input.onchange = (e) => {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (file) emit('upload', type, file)
  }
  input.click()
}
</script>

<template>
  <div class="vms-flex-col" style="gap: 1rem; flex: 1; min-width: 320px;">
    <!-- Textos & Nomes da Marca -->
    <div class="vms-telemetry-card" style="padding: 1rem; gap: 0.75rem;">
      <div class="vms-flex-between">
        <span class="vms-font-bold vms-text-xs" style="color: var(--vms-neu-accent-orange);">// NOMES & IDENTIFICAÇÃO DO SISTEMA</span>
      </div>
      <div class="vms-flex-col" style="gap: 0.6rem;">
        <div class="vms-form-group">
          <label class="vms-label" style="font-size: 10px;">Nome do Sistema (VMS Operacional)</label>
          <input v-model="modelValue.systemName" class="vms-auth-input" style="height: 32px;" placeholder="Ex: HYDRA VMS / VISION CORE" />
        </div>
        <div class="vms-form-group">
          <label class="vms-label" style="font-size: 10px;">Nome do Painel Administrativo</label>
          <input v-model="modelValue.adminTitle" class="vms-auth-input" style="height: 32px;" placeholder="Ex: ADMIN CENTER / PAINEL GESTOR" />
        </div>
        <div class="vms-form-group">
          <label class="vms-label" style="font-size: 10px;">Slogan / Subtítulo</label>
          <input v-model="modelValue.slogan" class="vms-auth-input" style="height: 32px;" placeholder="Ex: Video Management & AI Studio" />
        </div>
      </div>
    </div>

    <!-- Upload de Imagens e Logotipos -->
    <div class="vms-telemetry-card" style="padding: 1rem; gap: 0.75rem;">
      <div class="vms-flex-between">
        <span class="vms-font-bold vms-text-xs" style="color: var(--vms-neu-accent-orange);">// ÍCONES & ARQUIVOS DE LOGOMARCA</span>
      </div>
      <div class="vms-flex-col" style="gap: 0.6rem;">
        <div class="vms-flex-between" style="align-items: center; background: rgba(0,0,0,0.25); padding: 0.5rem 0.75rem; border-radius: 6px; border: 1px solid var(--vms-border);">
          <div class="vms-flex-row" style="gap: 0.6rem; align-items: center;">
            <img :src="modelValue.headerLogo" alt="Header Logo" style="width: 26px; height: 26px; object-fit: contain; background: rgba(255,255,255,0.05); padding: 2px; border-radius: 4px;" />
            <div class="vms-flex-col"><span class="vms-text-xs vms-font-semibold" style="color: #fff;">Ícone da Barra Superior</span><span class="vms-text-mono vms-text-2xs vms-text-dim">SVG ou PNG (24x24)</span></div>
          </div>
          <button class="vms-btn vms-btn-secondary vms-btn-sm" style="font-size: 10px;" @click="triggerUpload('headerLogo')">ENVIAR ARQUIVO</button>
        </div>

        <div class="vms-flex-between" style="align-items: center; background: rgba(0,0,0,0.25); padding: 0.5rem 0.75rem; border-radius: 6px; border: 1px solid var(--vms-border);">
          <div class="vms-flex-row" style="gap: 0.6rem; align-items: center;">
            <img :src="modelValue.faviconUrl" alt="Favicon" style="width: 26px; height: 26px; object-fit: contain; background: rgba(255,255,255,0.05); padding: 2px; border-radius: 4px;" />
            <div class="vms-flex-col"><span class="vms-text-xs vms-font-semibold" style="color: #fff;">Favicon da Aba do Navegador</span><span class="vms-text-mono vms-text-2xs vms-text-dim">ICO, SVG ou PNG (16x16 / 32x32)</span></div>
          </div>
          <button class="vms-btn vms-btn-secondary vms-btn-sm" style="font-size: 10px;" @click="triggerUpload('faviconUrl')">ENVIAR ARQUIVO</button>
        </div>

        <div class="vms-flex-between" style="align-items: center; background: rgba(0,0,0,0.25); padding: 0.5rem 0.75rem; border-radius: 6px; border: 1px solid var(--vms-border);">
          <div class="vms-flex-row" style="gap: 0.6rem; align-items: center;">
            <img :src="modelValue.loginLogo" alt="Login Logo" style="width: 26px; height: 26px; object-fit: contain; background: rgba(255,255,255,0.05); padding: 2px; border-radius: 4px;" />
            <div class="vms-flex-col"><span class="vms-text-xs vms-font-semibold" style="color: #fff;">Logo da Tela de Login</span><span class="vms-text-mono vms-text-2xs vms-text-dim">SVG ou PNG em alta resolução</span></div>
          </div>
          <button class="vms-btn vms-btn-secondary vms-btn-sm" style="font-size: 10px;" @click="triggerUpload('loginLogo')">ENVIAR ARQUIVO</button>
        </div>
      </div>
    </div>

    <!-- Cor de Destaque & Empresa -->
    <div class="vms-telemetry-card" style="padding: 1rem; gap: 0.75rem;">
      <div class="vms-flex-between">
        <span class="vms-font-bold vms-text-xs" style="color: var(--vms-neu-accent-orange);">// COR DE DESTAQUE & EMPRESA</span>
      </div>
      <div class="vms-flex-col" style="gap: 0.6rem;">
        <div class="vms-form-group">
          <label class="vms-label" style="font-size: 10px;">Cor Primária / Destaque HUD</label>
          <div class="vms-flex-row" style="gap: 0.4rem; align-items: center; flex-wrap: wrap;">
            <button v-for="c in colorPresets" :key="c" type="button" class="vms-btn vms-btn-sm" :style="{ background: c, width: '24px', height: '24px', padding: 0, borderRadius: '4px', border: modelValue.brandColor === c ? '2px solid #fff' : '1px solid rgba(0,0,0,0.4)' }" @click="modelValue.brandColor = c" />
            <input v-model="modelValue.brandColor" type="color" style="width: 32px; height: 28px; padding: 0; background: none; border: none; cursor: pointer;" />
            <span class="vms-text-mono vms-text-xs" style="color: #fff; margin-left: 0.3rem;">{{ modelValue.brandColor }}</span>
          </div>
        </div>
        <div class="vms-flex-row" style="gap: 0.5rem;">
          <div class="vms-form-group" style="flex: 1;">
            <label class="vms-label" style="font-size: 10px;">Empresa Integradora</label>
            <input v-model="modelValue.companyName" class="vms-auth-input" style="height: 32px;" placeholder="Ex: Minha Empresa de Segurança" />
          </div>
          <div class="vms-form-group" style="flex: 1;">
            <label class="vms-label" style="font-size: 10px;">Contato de Suporte</label>
            <input v-model="modelValue.supportInfo" class="vms-auth-input" style="height: 32px;" placeholder="Ex: suporte@empresa.com" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
