import { ref, watch } from 'vue'

export interface BrandingConfig {
  systemName: string
  adminTitle: string
  slogan: string
  headerLogo: string
  faviconUrl: string
  loginLogo: string
  brandColor: string
  companyName: string
  supportInfo: string
}

export const defaultBranding: BrandingConfig = {
  systemName: 'HYDRA VMS',
  adminTitle: 'ADMIN CENTER',
  slogan: 'Video Management & AI Studio',
  headerLogo: '/hydra.svg',
  faviconUrl: '/favicon.svg',
  loginLogo: '/hydra.svg',
  brandColor: '#ff5e3a',
  companyName: 'Hydra Core Enterprise',
  supportInfo: 'suporte@hydravms.io'
}

const STORAGE_KEY = 'hydravms_branding_v1'

const loadStoredBranding = (): BrandingConfig => {
  if (typeof window === 'undefined') return { ...defaultBranding }
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (raw) return { ...defaultBranding, ...JSON.parse(raw) }
  } catch {}
  return { ...defaultBranding }
}

const branding = ref<BrandingConfig>(loadStoredBranding())

const updateFaviconInDom = (url: string) => {
  if (typeof document === 'undefined') return
  let link = document.querySelector("link[rel~='icon']") as HTMLLinkElement
  if (!link) {
    link = document.createElement('link')
    link.rel = 'icon'
    document.head.appendChild(link)
  }
  link.href = url
}

const applyBrandingEffects = (cfg: BrandingConfig) => {
  if (typeof document === 'undefined') return
  updateFaviconInDom(cfg.faviconUrl || defaultBranding.faviconUrl)
  if (cfg.brandColor) {
    document.documentElement.style.setProperty('--vms-neu-accent-orange', cfg.brandColor)
  }
}

applyBrandingEffects(branding.value)

watch(branding, (cfg) => {
  if (typeof window !== 'undefined') {
    try { localStorage.setItem(STORAGE_KEY, JSON.stringify(cfg)) } catch {}
  }
  applyBrandingEffects(cfg)
}, { deep: true })

export function useBranding() {
  const saveBranding = (newConfig: Partial<BrandingConfig>) => {
    branding.value = { ...branding.value, ...newConfig }
  }

  const resetDefaults = () => {
    branding.value = { ...defaultBranding }
  }

  const readFileAsDataUrl = (file: File): Promise<string> => {
    return new Promise((resolve, reject) => {
      const reader = new FileReader()
      reader.onload = () => resolve(reader.result as string)
      reader.onerror = reject
      reader.readAsDataURL(file)
    })
  }

  return { branding, saveBranding, resetDefaults, readFileAsDataUrl }
}
