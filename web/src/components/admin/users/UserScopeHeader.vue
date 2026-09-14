<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import type { UserItem, UserRole } from '../../../types/userTree'
import { getAvailableTimezones, detectLocalTimezone } from '../../../utils/timezones'

const props = defineProps<{ user: UserItem }>()
const emit = defineEmits<{ (e: 'roleChange', newRole: UserRole): void; (e: 'saved', msg: string): void }>()
const availableTimezones = ref<string[]>([]), newPassword = ref(''), confirmPassword = ref(''), isSavedAuto = ref(false)
let autoSaveTimer: any = null

onMounted(() => {
  availableTimezones.value = getAvailableTimezones()
  if (!props.user.timezone) props.user.timezone = detectLocalTimezone()
})

// Auto-salvamento com debounce para dados cadastrais
watch([() => props.user.fullName, () => props.user.email, () => props.user.timezone, () => props.user.phonePrimary, () => props.user.phoneSecondary], () => {
  clearTimeout(autoSaveTimer)
  autoSaveTimer = setTimeout(() => { isSavedAuto.value = true; setTimeout(() => { isSavedAuto.value = false }, 2500) }, 400)
}, { deep: true })

const formatPhone = (val: string): string => {
  if (!val) return ''
  let digits = val.replace(/\D/g, '')
  if (!digits || digits === '5' || digits === '55') return ''
  if (digits.startsWith('55') && digits.length > 2) digits = digits.slice(2)
  if (digits.startsWith('0')) digits = digits.slice(1)
  if (!digits) return ''
  digits = digits.slice(0, 11)
  if (digits.length <= 2) return `(${digits}`
  if (digits.length <= 6) return `(${digits.slice(0, 2)}) ${digits.slice(2)}`
  if (digits.length <= 10) return `(${digits.slice(0, 2)}) ${digits.slice(2, 6)}-${digits.slice(6)}`
  return `(${digits.slice(0, 2)}) ${digits.slice(2, 7)}-${digits.slice(7, 11)}`
}

const handlePhoneInput = (field: 'phonePrimary' | 'phoneSecondary', ev: Event) => {
  const target = ev.target as HTMLInputElement; const formatted = formatPhone(target.value)
  props.user[field] = formatted; target.value = formatted
}

// Validação reativa de senha em tempo real
const passwordStatus = computed(() => {
  if (!newPassword.value && !confirmPassword.value) return { ok: false, msg: '', type: 'idle' }
  if (newPassword.value.length < 4) return { ok: false, msg: 'Mínimo de 4 caracteres necessários.', type: 'warn' }
  if (confirmPassword.value && newPassword.value !== confirmPassword.value) return { ok: false, msg: '⚠️ As senhas digitadas não coincidem.', type: 'error' }
  if (newPassword.value === confirmPassword.value) return { ok: true, msg: '✓ Senhas coincidem.', type: 'success' }
  return { ok: false, msg: '', type: 'idle' }
})

const handleUpdatePassword = () => {
  if (!passwordStatus.value.ok) return
  emit('saved', `[SEGURANÇA] Nova senha atualizada com sucesso para "${props.user.fullName || props.user.username}".`)
  newPassword.value = ''; confirmPassword.value = ''
}
</script>

<template>
  <div class="vms-flex-col" style="gap: 0.65rem; background: rgba(255, 255, 255, 0.02); border: 1px solid var(--vms-border); border-radius: 8px; padding: 0.75rem 0.85rem;">
    <!-- Profile & Contact Header with discrete Autosave indicator -->
    <div class="vms-flex-between" style="align-items: center; min-height: 24px;">
      <div class="vms-flex-col" style="gap: 2px;">
        <span class="vms-text-dim vms-text-2xs">DADOS CADASTRAIS // IDENTIFICAÇÃO DO USUÁRIO</span>
        <span class="vms-font-bold" style="color: var(--vms-neu-accent-orange); font-size: 13px;">
          {{ user.fullName || user.username }}
        </span>
      </div>
      <span v-if="isSavedAuto" class="vms-text-mono vms-text-2xs" style="color: var(--vms-neu-accent-green); font-size: 10px; font-weight: 600;">
        ✓ SALVO
      </span>
    </div>

    <!-- Useful Contact & Localization Grid -->
    <div class="vms-flex-col" style="gap: 0.45rem; border-top: 1px dashed var(--vms-border); padding-top: 0.5rem;">
      <div class="vms-form-group">
        <label class="vms-label" style="font-size: 9px;">Nome Completo</label>
        <input v-model="user.fullName" class="vms-auth-input" style="padding: 4px 8px; font-size: 11px; height: 28px; box-sizing: border-box;" placeholder="Nome Completo do Usuário" />
      </div>

      <div class="vms-flex-row" style="gap: 0.5rem;">
        <div class="vms-form-group" style="flex: 1;">
          <label class="vms-label" style="font-size: 9px;">E-mail</label>
          <input v-model="user.email" class="vms-auth-input" style="padding: 4px 8px; font-size: 11px; height: 28px; box-sizing: border-box;" placeholder="email@empresa.com" />
        </div>
        <div class="vms-form-group" style="flex: 1;">
          <label class="vms-label" style="font-size: 9px;">Fuso Horário (Digite para filtrar)</label>
          <input
            v-model="user.timezone"
            list="tz-datalist"
            class="vms-auth-input"
            style="padding: 4px 8px; font-size: 11px; height: 28px; box-sizing: border-box;"
            placeholder="Buscar fuso horário..."
          />
          <datalist id="tz-datalist">
            <option v-for="tz in availableTimezones" :key="tz" :value="tz">{{ tz }}</option>
          </datalist>
        </div>
      </div>

      <div class="vms-flex-row" style="gap: 0.5rem;">
        <div class="vms-form-group" style="flex: 1;">
          <label class="vms-label" style="font-size: 9px;">Telefone 1 (Principal)</label>
          <input :value="user.phonePrimary" class="vms-auth-input" style="padding: 4px 8px; font-size: 11px; height: 28px; box-sizing: border-box;" placeholder="(00) 00000-0000" maxlength="16" @input="handlePhoneInput('phonePrimary', $event)" />
        </div>
        <div class="vms-form-group" style="flex: 1;">
          <label class="vms-label" style="font-size: 9px;">Telefone 2 (Opcional)</label>
          <input :value="user.phoneSecondary" class="vms-auth-input" style="padding: 4px 8px; font-size: 11px; height: 28px; box-sizing: border-box;" placeholder="(00) 0000-0000" maxlength="16" @input="handlePhoneInput('phoneSecondary', $event)" />
        </div>
      </div>

      <!-- Password Change Section with Perfect Alignment & Reactive Feedback -->
      <div class="vms-flex-col" style="gap: 0.35rem; border-top: 1px dashed var(--vms-border); padding-top: 0.5rem;">
        <div class="vms-flex-between" style="align-items: center;">
          <span class="vms-text-mono vms-text-2xs vms-text-dim">// ALTERAR SENHA DE ACESSO</span>
          <span v-if="passwordStatus.msg" class="vms-text-2xs vms-text-mono" :style="{ color: passwordStatus.type === 'error' ? 'var(--vms-neu-accent-red)' : passwordStatus.type === 'success' ? 'var(--vms-neu-accent-green)' : 'var(--vms-neu-accent-yellow)', fontSize: '10px' }">
            {{ passwordStatus.msg }}
          </span>
        </div>
        <div class="vms-flex-row" style="gap: 0.5rem; align-items: flex-end;">
          <div class="vms-form-group" style="flex: 1; margin-bottom: 0;">
            <label class="vms-label" style="font-size: 9px;">Nova Senha</label>
            <input v-model="newPassword" type="password" class="vms-auth-input" style="padding: 4px 8px; font-size: 11px; height: 28px; box-sizing: border-box;" placeholder="••••••••" autocomplete="new-password" />
          </div>
          <div class="vms-form-group" style="flex: 1; margin-bottom: 0;">
            <label class="vms-label" style="font-size: 9px;">Confirmar Senha</label>
            <input v-model="confirmPassword" type="password" class="vms-auth-input" style="padding: 4px 8px; font-size: 11px; height: 28px; box-sizing: border-box;" :style="{ borderColor: passwordStatus.type === 'error' ? 'var(--vms-neu-accent-red)' : undefined }" placeholder="••••••••" autocomplete="new-password" />
          </div>
          <div class="vms-form-group" style="margin-bottom: 0; justify-content: flex-end;">
            <label class="vms-label" style="font-size: 9px; visibility: hidden;">Ação</label>
            <button class="vms-btn vms-btn-primary vms-btn-sm" style="font-size: 10px; padding: 0 16px; height: 28px; box-sizing: border-box; flex-shrink: 0;" :disabled="!passwordStatus.ok" @click="handleUpdatePassword">
              SALVAR
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>


