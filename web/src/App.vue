<script setup lang="ts">
import { ref, defineAsyncComponent } from 'vue'
import { useAuth } from './composables/useAuth'
import { useAnalyticsAlerts } from './composables/useAnalyticsAlerts'
import LoginView from './views/auth/LoginView.vue'
import AppTopHeader from './components/layout/AppTopHeader.vue'
import LiveMosaicView from './views/mosaic/LiveMosaicView.vue'
import AnalyticsAlertsDrawer from './components/analytics/AnalyticsAlertsDrawer.vue'

const AdminCenterView = defineAsyncComponent(() => import('./views/admin/AdminCenterView.vue'))

const { isAuthenticated, isLoading, errorMessage, username, isAdmin, handleLogin, handleLogout } = useAuth()
const { isDrawerOpen, alerts, unreadCount, toggleDrawer, acknowledgeAlert, clearAllAlerts } = useAnalyticsAlerts()

const currentMode = ref<'vms' | 'admin'>('vms')

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
