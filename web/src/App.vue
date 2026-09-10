<script setup lang="ts">
import { ref, watch, onMounted, defineAsyncComponent } from 'vue'
import { useAuth } from './composables/useAuth'
import { useAnalyticsAlerts } from './composables/useAnalyticsAlerts'
import LoginView from './views/auth/LoginView.vue'
import AppTopHeader from './components/layout/AppTopHeader.vue'
import LiveMosaicView from './views/mosaic/LiveMosaicView.vue'
import AnalyticsAlertsDrawer from './components/analytics/AnalyticsAlertsDrawer.vue'

const AdminCenterView = defineAsyncComponent(() => import('./views/admin/AdminCenterView.vue'))

const { isAuthenticated, isLoading, errorMessage, username, isAdmin, handleLogin, handleLogout } = useAuth()
const { isDrawerOpen, alerts, unreadCount, toggleDrawer, acknowledgeAlert, clearAllAlerts } = useAnalyticsAlerts()

const getInitialMode = (): 'vms' | 'admin' => {
  if (typeof window !== 'undefined' && window.location.hash === '#admin') return 'admin'
  return 'vms'
}

const currentMode = ref<'vms' | 'admin'>(getInitialMode())

const syncFromHash = () => {
  const isHashAdmin = window.location.hash === '#admin'
  if (isHashAdmin) {
    currentMode.value = 'admin'
    window.name = 'hydravms_admin_center'
    document.title = 'HydraVMS - ADMIN CENTER'
  } else {
    currentMode.value = 'vms'
    window.name = 'hydravms_main_vms'
    document.title = 'HydraVMS'
  }
}

watch(currentMode, (mode) => {
  document.title = mode === 'admin' ? 'HydraVMS - ADMIN CENTER' : 'HydraVMS'
  const targetHash = mode === 'admin' ? '#admin' : '#vms'
  window.name = mode === 'admin' ? 'hydravms_admin_center' : 'hydravms_main_vms'
  if (window.location.hash !== targetHash) {
    window.history.replaceState(null, '', targetHash)
  }
}, { immediate: true })

onMounted(() => {
  syncFromHash()
  window.addEventListener('hashchange', syncFromHash)
})

const handleSwitchMode = (mode: 'vms' | 'admin') => {
  if (mode === 'admin' && !isAdmin.value) return
  currentMode.value = mode
}
</script>

<template>
  <div
    style="height: 100vh; width: 100vw; background: var(--vms-neu-bg); overflow: hidden; display: flex; flex-direction: column;"
    @contextmenu.prevent
  >
    <!-- Neumorphic Login -->
    <LoginView
      v-if="!isAuthenticated"
      :isLoading="isLoading"
      :errorMessage="errorMessage"
      @login="handleLogin"
    />

    <!-- Master Top Bar + Active View Mode -->
    <template v-else>
      <AppTopHeader
        :username="username"
        :isAdmin="isAdmin"
        :currentMode="currentMode"
        :unreadAlertsCount="unreadCount"
        @toggleAlerts="toggleDrawer"
        @switchMode="handleSwitchMode"
        @logout="handleLogout"
      />

      <div style="flex: 1; display: flex; overflow: hidden; position: relative;">
        <LiveMosaicView v-if="currentMode === 'vms' || !isAdmin" />
        <AdminCenterView v-else-if="currentMode === 'admin' && isAdmin" />
      </div>

      <Transition name="vms-drawer">
        <AnalyticsAlertsDrawer
          v-if="isDrawerOpen"
          :alerts="alerts"
          @close="isDrawerOpen = false"
          @acknowledge="acknowledgeAlert"
          @clearAll="clearAllAlerts"
        />
      </Transition>
    </template>
  </div>
</template>
