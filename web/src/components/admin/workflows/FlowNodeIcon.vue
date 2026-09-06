<script setup lang="ts">
defineProps<{ kind: string; type: string }>()
</script>

<template>
  <div class="vms-icon-wrapper">
    <!-- 1. GATILHOS (TRIGGER / PLAY) -->
    <!-- CAMERA / FLUXO DE VIDEO (Icone Oficial de Camera) -->
    <svg v-if="type === 'TRIGGER' && kind === 'AI_DETECT'" width="26" height="26" viewBox="0 0 576 512" fill="#ff5e3a">
      <path d="M0 128C0 92.7 28.7 64 64 64H320c35.3 0 64 28.7 64 64V384c0 35.3-28.7 64-64 64H64c-35.3 0-64-28.7-64-64V128zM559.1 99.8c10.4 5.6 16.9 16.4 16.9 28.2V384c0 11.8-6.5 22.6-16.9 28.2s-23 5-32.9-1.6l-112-74.7c-9.8-6.5-16.1-17.4-16.1-29.9V205.1c0-12.5 6.3-23.4 16.1-29.9l112-74.7c9.9-6.6 22.5-7.3 32.9-1.6z"/>
    </svg>
    <svg v-else-if="type === 'TRIGGER' && kind === 'ALARM_SENSOR'" width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2">
      <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 0 1-3.46 0"/>
    </svg>
    <svg v-else-if="type === 'TRIGGER'" width="26" height="26" viewBox="0 0 24 24" fill="#ff5e3a" stroke="#ff5e3a" stroke-width="1.5">
      <polygon points="5 3 19 12 5 21 5 3"/>
    </svg>

    <!-- 2. CONDICOES (CONDITION) -->
    <!-- JUNCAO EM Y (OU / MERGE) -->
    <svg v-else-if="kind === 'BRANCH_OR'" width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="#ffffff" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
      <path d="M4 6L12 12M4 18L12 12M12 12H20"/>
    </svg>
    <!-- PORTA LOGICA E (AND) -->
    <svg v-else-if="kind === 'LOGIC_AND'" width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="#ffffff" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <path d="M5 4H11C15 4 19 7.6 19 12C19 16.4 15 20 11 20H5V4Z"/>
    </svg>
    <!-- DIVISAO EM Y (SPLIT / BIFURCACAO) -->
    <svg v-else-if="kind === 'SPLIT_Y'" width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="#ffffff" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
      <path d="M4 12H12M12 12L20 6M12 12L20 18"/>
    </svg>
    <svg v-else-if="kind === 'SCHEDULE'" width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="#ffffff" stroke-width="2">
      <circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/>
    </svg>
    <svg v-else-if="kind === 'ANTI_SPAM'" width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="#ffffff" stroke-width="2">
      <circle cx="12" cy="12" r="10"/><line x1="4.93" y1="4.93" x2="19.07" y2="19.07"/>
    </svg>

    <!-- 3. ACOES (ACTION) -->
    <svg v-else-if="kind === 'TELEGRAM'" width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2">
      <line x1="22" y1="2" x2="11" y2="13"/><polygon points="22 2 15 22 11 13 2 9 22 2"/>
    </svg>
    <svg v-else-if="kind === 'WEBHOOK'" width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2">
      <polyline points="16 18 22 12 16 6"/><polyline points="8 6 2 12 8 18"/>
    </svg>
    <svg v-else-if="kind === 'RELAY_SIREN'" width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2">
      <polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5"/><path d="M19.07 4.93a10 10 0 0 1 0 14.14M15.54 8.46a5 5 0 0 1 0 7.07"/>
    </svg>

    <!-- Bloco Terminal FIM / PARE -->
    <svg v-else-if="kind === 'TERMINAL_STOP'" width="24" height="24" viewBox="0 0 24 24" fill="#ff5e3a" stroke="#ff5e3a">
      <rect x="5" y="5" width="14" height="14" rx="2"/>
    </svg>

    <svg v-else width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#ffffff" stroke-width="2">
      <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"/>
    </svg>
  </div>
</template>

<style scoped>
.vms-icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
