<script setup lang="ts">
import { ref } from 'vue'
import { useBranding, type BrandingConfig } from '../../composables/useBranding'
import { useI18n } from '../../composables/useI18n'
import BrandingFormPane from '../../components/admin/branding/BrandingFormPane.vue'
import BrandingPreviewPane from '../../components/admin/branding/BrandingPreviewPane.vue'

const { branding, saveBranding, resetDefaults, readFileAsDataUrl } = useBranding()
const { t } = useI18n()
const formData = ref<BrandingConfig>({ ...branding.value })
const notification = ref<string | null>(null)

const showNotification = (msg: string) => {
  notification.value = msg
  setTimeout(() => { notification.value = null }, 3500)
}

const handleUpload = async (type: 'headerLogo' | 'faviconUrl' | 'loginLogo', file: File) => {
  try {
    const dataUrl = await readFileAsDataUrl(file)
    formData.value[type] = dataUrl
    showNotification(`[UPLOAD] Imagem carregada para ${type}.`)
  } catch (err: any) {
    showNotification(`[ERRO] Falha ao processar imagem: ${err?.message || ''}`)
  }
}

const handleSave = () => {
  saveBranding(formData.value)
  showNotification('[LOGOMARCA] Configurações de marca salvas com sucesso!')
}

const handleReset = () => {
  resetDefaults()
  formData.value = { ...branding.value }
  showNotification('[LOGOMARCA] Configurações restauradas para o padrão Hydra.')
}
</script>

<template>
  <div class="vms-flex-col" style="gap: 1rem;">
    <!-- Top Bar -->
    <div class="vms-flex-between">
      <div class="vms-flex-col" style="gap: 2px;">
        <h3 class="vms-h3" style="color: var(--vms-neu-accent-orange);">{{ t('branding_title') }}</h3>
        <span class="vms-text-mono vms-text-2xs vms-text-dim">{{ t('branding_sub') }}</span>
      </div>
      <div class="vms-flex-row" style="gap: 0.5rem; align-items: center;">
        <button class="vms-btn vms-btn-secondary" style="height: 32px; padding: 0 12px; font-size: 11px;" @click="handleReset">
          RESTAURAR PADRÕES
        </button>
        <button class="vms-btn vms-btn-primary" style="height: 32px; padding: 0 16px; font-size: 11px;" @click="handleSave">
          SALVAR
        </button>
      </div>
    </div>

    <!-- Toast Notification -->
    <Transition name="vms-toast">
      <div v-if="notification" class="vms-toast-notification" @click="notification = null">
        <span>{{ notification }}</span>
      </div>
    </Transition>

    <!-- Split Grid: Form vs Live Preview -->
    <div class="vms-desktop-container" style="padding: 1rem; overflow-y: auto;">
      <div style="display: flex; flex-direction: row; gap: 1.5rem; flex-wrap: wrap; width: 100%;">
        <BrandingFormPane :model-value="formData" @upload="handleUpload" />
        <BrandingPreviewPane :config="formData" />
      </div>
    </div>
  </div>
</template>
