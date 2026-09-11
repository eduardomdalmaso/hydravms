import { ref } from 'vue'
import { translations, type Lang } from '../locales/messages'

const currentLanguage = ref<Lang>((localStorage.getItem('hydravms_lang') as Lang) || 'PT')

export function useI18n() {
  const setLanguage = (lang: Lang) => {
    currentLanguage.value = lang
    localStorage.setItem('hydravms_lang', lang)
  }

  const t = (key: string): string => {
    return translations[currentLanguage.value]?.[key] || translations['PT']?.[key] || key
  }

  return { currentLanguage, setLanguage, t }
}
