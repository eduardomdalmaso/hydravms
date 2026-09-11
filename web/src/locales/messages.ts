import { pt } from './pt'
import { en } from './en'
import { es } from './es'

export type Lang = 'PT' | 'EN' | 'ES'

export const translations: Record<Lang, Record<string, string>> = {
  PT: pt,
  EN: en,
  ES: es
}
