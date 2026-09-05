<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue"
import UserMenuDropdown from "./UserMenuDropdown.vue"

const props = withDefaults(
  defineProps<{ username?: string; isAdmin?: boolean; currentMode?: "vms" | "admin"; unreadAlertsCount?: number }>(),
  { username: "admin", isAdmin: true, currentMode: "vms", unreadAlertsCount: 0 }
)

const emit = defineEmits<{ (e: "toggleAlerts"): void; (e: "switchMode", mode: "vms" | "admin"): void; (e: "logout"): void }>()

const isUserMenuOpen = ref(false)
const isCalendarOpen = ref(false)
const currentTimeStr = ref("")

const updateTime = () => {
  const now = new Date()
  const days = ["Dom", "Seg", "Ter", "Qua", "Qui", "Sex", "Sab"]
  const months = ["Jan", "Fev", "Mar", "Abr", "Mai", "Jun", "Jul", "Ago", "Set", "Out", "Nov", "Dez"]
  currentTimeStr.value = `${days[now.getDay()]}, ${String(now.getDate()).padStart(2, "0")} ${months[now.getMonth()]} • ${now.toTimeString().slice(0, 5)}`
}

let timer: any = null
onMounted(() => { updateTime(); timer = setInterval(updateTime, 1000) })
onUnmounted(() => { if (timer) clearInterval(timer) })

const handleSelectMode = (mode: "vms" | "admin") => {
  if (mode === "admin" && !props.isAdmin) return
  emit("switchMode", mode)
  isUserMenuOpen.value = false
}
</script>

<template>
  <header class="vms-header" style="background: #15181d; border-bottom: 1px solid var(--vms-border); height: 44px; padding: 0 1rem;">
    <!-- Left: Brand Title -->
    <div class="vms-flex-row" style="align-items: center;">
      <span style="color: #ffffff; font-family: var(--vms-font-roboto); font-size: 15px; font-weight: 800; letter-spacing: 0.8px;">
        HYDRA VMS<template v-if="currentMode === 'admin'"> - ADMIN CENTER</template>
      </span>
    </div>

    <!-- Center: GNOME Style Calendar / Clock Trigger in Roboto -->
    <div style="position: relative;">
      <button class="vms-btn vms-btn-ghost vms-btn-sm" style="color: #ffffff; font-family: var(--vms-font-roboto); font-size: 13.5px; font-weight: 600; letter-spacing: 0.3px; padding: 0.25rem 0.85rem;" @click="isCalendarOpen = !isCalendarOpen">
        {{ currentTimeStr || "Qua, 03 Set • 20:55" }}
      </button>

      <!-- GNOME Style Calendar Popup -->
      <div v-if="isCalendarOpen" class="vms-gnome-calendar-popup">
        <div class="vms-flex-between" style="border-bottom: 1px solid var(--vms-border); padding-bottom: 0.5rem;">
          <span class="vms-font-semibold vms-text-xs" style="color: #fff; font-family: var(--vms-font-roboto);">SETEMBRO 2026</span>
          <button class="vms-ubuntu-close-btn" style="position: relative; top: 0; right: 0; width: 18px; height: 18px;" @click="isCalendarOpen = false">
            <svg width="8" height="8" viewBox="0 0 12 12" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round"><path d="M2 2L10 10M10 2L2 10" /></svg>
          </button>
        </div>
        <div class="vms-gnome-calendar-grid">
          <span v-for="d in ['D','S','T','Q','Q','S','S']" :key="d" class="vms-text-2xs vms-text-dim">{{ d }}</span>
          <span v-for="n in 30" :key="n" class="vms-gnome-cal-day" :class="{ active: n === 3 }">{{ n }}</span>
        </div>
      </div>
    </div>

    <!-- Right: Bell Notification + User Dropdown -->
    <div class="vms-flex-row" style="gap: 0.75rem;">
      <button class="vms-btn vms-btn-ghost vms-btn-sm" style="position: relative; padding: 0.35rem 0.65rem;" title="Alertas IA" @click="emit('toggleAlerts')">
        <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="color: var(--vms-neu-accent-orange);">
          <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"></path><path d="M13.73 21a2 2 0 0 1-3.46 0"></path>
        </svg>
        <span v-if="unreadAlertsCount > 0" class="vms-status-led alert" style="position: absolute; top: 4px; right: 4px;"></span>
      </button>

      <!-- User Dropdown Menu with RBAC protection -->
      <div class="vms-user-dropdown-container">
        <button class="vms-btn vms-btn-secondary vms-btn-sm" style="display: flex; align-items: center; gap: 0.5rem; padding: 0.3rem 0.65rem;" @click="isUserMenuOpen = !isUserMenuOpen">
          <span class="vms-text-xs" style="color: #fff; font-family: var(--vms-font-roboto); font-weight: 500;">{{ username }}</span>
          <span style="font-size: 9px; color: var(--vms-neu-accent-orange);">▼</span>
        </button>

        <UserMenuDropdown
          :username="username"
          :isAdmin="isAdmin"
          :currentMode="currentMode"
          :isOpen="isUserMenuOpen"
          @switchMode="handleSelectMode"
          @logout="emit('logout')"
          @close="isUserMenuOpen = false"
        />
      </div>
    </div>
  </header>
</template>
