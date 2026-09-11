<script setup lang="ts">
import { ref, computed } from "vue"
import { useI18n } from "../../composables/useI18n"
import type { AlarmItemInfo, AlarmStatus } from "../../types/mosaic"

const props = defineProps<{ alarms: AlarmItemInfo[]; isOpen: boolean }>()
const emit = defineEmits<{
  (e: "toggle"): void
  (e: "selectAlarm", alarm: AlarmItemInfo): void
}>()
const { t } = useI18n()

const isSearchOpen = ref(false)
const searchQuery = ref("")
const selectedStatus = ref<AlarmStatus | "ALL">("ALL")

const filteredAlarms = computed(() => {
  return props.alarms.filter(a => {
    const matchesName = a.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      a.zone.toLowerCase().includes(searchQuery.value.toLowerCase())
    const matchesStatus = selectedStatus.value === "ALL" || a.status === selectedStatus.value
    return matchesName && matchesStatus
  })
})

const getStatusTitle = (status: AlarmStatus) => {
  if (status === "online") return "Armado / Online"
  if (status === "alert") return "Desarmado"
  return "Offline / Falha"
}
</script>

<template>
  <div>
    <!-- Clickable Header Row with Orange Counter -->
    <div class="vms-accordion-header" :class="{ active: isOpen }" title="Clique na linha para expandir ou recuar" @click="emit('toggle')">
      <div class="vms-flex-row" style="gap: 0.35rem; align-items: center;">
        <span class="vms-text-xs vms-font-semibold" style="color: var(--vms-neu-accent-orange);">[{{ t('alarms').toUpperCase() }}]</span>
        <span class="vms-badge" style="background: rgba(255, 94, 58, 0.18); color: var(--vms-neu-accent-orange); border: 1px solid rgba(255, 94, 58, 0.35); font-size: 8.5px; font-weight: 700; padding: 1px 5px;">
          {{ filteredAlarms.length }}
        </span>
      </div>
      <div class="vms-flex-row" style="gap: 0.4rem; align-items: center;">
        <button class="vms-sidebar-icon-btn" :title="t('search_alm')" @click.stop="isSearchOpen = !isSearchOpen">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#ff5e3a" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="11" cy="11" r="8"></circle><line x1="21" y1="21" x2="16.65" y2="16.65"></line>
          </svg>
        </button>
      </div>
    </div>

    <!-- Search & Filter Bar -->
    <div v-if="isOpen && isSearchOpen" style="padding: 0.4rem 0.5rem; background: #11141b; border-bottom: 1px solid var(--vms-border); display: flex; flex-direction: column; gap: 0.3rem;">
      <input v-model="searchQuery" class="vms-auth-input" style="padding: 0.25rem 0.5rem; font-size: 11px; height: 26px;" :placeholder="t('search_alm')" autofocus />
      <div class="vms-flex-row" style="gap: 0.25rem;">
        <button v-for="s in (['ALL', 'online', 'alert', 'offline'] as const)" :key="s" class="vms-btn vms-btn-sm" :class="selectedStatus === s ? 'vms-btn-primary' : 'vms-btn-ghost'" style="font-size: 9px; padding: 1px 4px;" @click="selectedStatus = s">
          {{ s === 'ALL' ? t('all') : (s === 'online' ? t('armed') : (s === 'alert' ? t('disarmed') : 'OFFLINE')) }}
        </button>
      </div>
    </div>

    <!-- Alarms List with LED Status (Green Online, Yellow Alert, Red Offline) -->
    <div v-if="isOpen" style="padding: 0.4rem; display: flex; flex-direction: column; gap: 0.3rem; background: #0c0e14; max-height: 200px; overflow-y: auto;">
      <div v-for="alm in filteredAlarms" :key="alm.id" class="vms-asset-item" style="padding: 0.4rem 0.55rem;" @click="emit('selectAlarm', alm)">
        <span class="vms-text-xs" style="color: #ffffff; font-family: var(--vms-font-roboto); font-weight: 500; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; max-width: 120px;">
          {{ alm.name }}
        </span>
        <div class="vms-flex-row" style="gap: 0.45rem; align-items: center;">
          <span class="vms-badge" style="background: rgba(255, 255, 255, 0.08); color: #ffffff; border: 1px solid rgba(255, 255, 255, 0.15); font-family: var(--vms-font-roboto); font-weight: 500; font-size: 8px; padding: 1px 4px;">
            {{ alm.type || 'ALM' }}
          </span>
          <!-- LED Indicator (Green: Online, Yellow: Alert, Red: Offline) -->
          <span class="vms-status-led" :class="alm.status" :title="getStatusTitle(alm.status)"></span>
        </div>
      </div>
    </div>
  </div>
</template>
