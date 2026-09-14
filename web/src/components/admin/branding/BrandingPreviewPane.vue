<script setup lang="ts">
import type { BrandingConfig } from '../../../composables/useBranding'

defineProps<{ config: BrandingConfig }>()
</script>

<template>
  <div class="vms-flex-col" style="gap: 1rem; flex: 1; min-width: 320px;">
    <!-- SIMULADOR 1: ABA DO NAVEGADOR -->
    <div class="vms-telemetry-card" style="padding: 1rem; gap: 0.75rem;">
      <span class="vms-font-bold vms-text-xs" style="color: var(--vms-neu-accent-orange);">// PRÉ-VISUALIZAÇÃO // ABA DO NAVEGADOR</span>
      <div style="background: #1e222a; border-radius: 8px 8px 0 0; padding: 0.5rem 0.75rem 0; border: 1px solid var(--vms-border);">
        <div class="vms-flex-row" style="gap: 0.4rem; align-items: center; background: #15181d; width: fit-content; padding: 0.35rem 0.85rem; border-radius: 6px 6px 0 0; border-top: 2px solid var(--vms-neu-accent-orange);">
          <img :src="config.faviconUrl" alt="Favicon" style="width: 14px; height: 14px; object-fit: contain;" />
          <span class="vms-text-xs vms-font-semibold" style="color: #fff; max-width: 180px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;">
            {{ config.adminTitle || config.systemName }}
          </span>
          <span style="color: #666; font-size: 11px; margin-left: 0.5rem;">×</span>
        </div>
      </div>
    </div>

    <!-- SIMULADOR 2: CABEÇALHO SUPERIOR -->
    <div class="vms-telemetry-card" style="padding: 1rem; gap: 0.75rem;">
      <span class="vms-font-bold vms-text-xs" style="color: var(--vms-neu-accent-orange);">// PRÉ-VISUALIZAÇÃO // BARRA SUPERIOR (HEADER)</span>
      <div style="background: #15181d; border: 1px solid var(--vms-border); border-radius: 6px; height: 40px; padding: 0 0.85rem; display: flex; justify-content: space-between; align-items: center;">
        <div class="vms-flex-row" style="align-items: center; gap: 0.55rem;">
          <img :src="config.headerLogo" alt="Logo" style="width: 20px; height: 20px; object-fit: contain;" />
          <span style="color: #ffffff; font-size: 13px; font-weight: 800; letter-spacing: 0.8px;">
            {{ config.adminTitle }}
          </span>
        </div>
        <span class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-green);">● ONLINE</span>
      </div>
    </div>

    <!-- SIMULADOR 3: CARD DE LOGIN -->
    <div class="vms-telemetry-card" style="padding: 1rem; gap: 0.75rem;">
      <span class="vms-font-bold vms-text-xs" style="color: var(--vms-neu-accent-orange);">// PRÉ-VISUALIZAÇÃO // TELA DE LOGIN</span>
      <div style="background: #0f1218; border: 1px solid var(--vms-border); border-radius: 8px; padding: 1.25rem; display: flex; flex-direction: column; align-items: center; gap: 0.65rem;">
        <div :style="{ background: 'rgba(255, 94, 58, 0.12)', border: `1px solid ${config.brandColor || '#ff5e3a'}` }" style="display: flex; align-items: center; justify-content: center; width: 48px; height: 48px; border-radius: 50%;">
          <img :src="config.loginLogo" alt="Login Logo" style="width: 28px; height: 28px; object-fit: contain;" />
        </div>
        <div class="vms-flex-col" style="text-align: center; gap: 2px;">
          <span class="vms-font-bold" style="color: #fff; font-size: 15px;">{{ config.systemName }}</span>
          <span class="vms-text-2xs vms-text-dim">{{ config.slogan }}</span>
        </div>
        <div style="width: 100%; max-width: 200px; height: 26px; background: #1a1e27; border-radius: 4px; border: 1px solid var(--vms-border); padding: 4px 8px; display: flex; align-items: center;">
          <span class="vms-text-mono vms-text-2xs vms-text-dim">admin@empresa.com</span>
        </div>
        <div style="width: 100%; max-width: 200px; height: 26px; border-radius: 4px; display: flex; align-items: center; justify-content: center;" :style="{ background: config.brandColor || '#ff5e3a' }">
          <span class="vms-font-bold vms-text-2xs" style="color: #fff;">ENTRAR</span>
        </div>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">{{ config.companyName }}</span>
      </div>
    </div>
  </div>
</template>
