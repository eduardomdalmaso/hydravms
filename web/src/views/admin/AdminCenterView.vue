<script setup lang="ts">
import { ref, computed } from "vue"
import AdminSidebar, { type AdminPageId } from "./AdminSidebar.vue"
import AdminCamerasView from "./AdminCamerasView.vue"
import AdminAlarmsView from "./AdminAlarmsView.vue"
import AdminUsersView from "./AdminUsersView.vue"
import AdminLayoutsView from "./AdminLayoutsView.vue"
import AdminCarouselsView from "./AdminCarouselsView.vue"
import AdminMapsView from "./AdminMapsView.vue"
import AdminWorkflowsView from "./AdminWorkflowsView.vue"
import AdminMarketplaceView from "./AdminMarketplaceView.vue"
import AnalyticPluginCrudView from "./AnalyticPluginCrudView.vue"
import AnalyticPluginEventsView from "./AnalyticPluginEventsView.vue"
import StorageDisksView from "../storage/StorageDisksView.vue"
import AdminPerformanceView from "./AdminPerformanceView.vue"
import AdminLogsView from "./AdminLogsView.vue"

const activePage = ref<AdminPageId>("video_streams")

const parsedPlugin = computed(() => {
  if (!activePage.value.startsWith('plugin_')) return null
  const isInst = activePage.value.endsWith('_instances')
  const isEvt = activePage.value.endsWith('_events')
  if (!isInst && !isEvt) return null
  const pluginId = activePage.value.replace('plugin_', '').replace('_instances', '').replace('_events', '')
  return { pluginId, type: isInst ? 'instances' : 'events' }
})
</script>

<template>
  <div class="vms-admin-container" style="display: flex; flex-direction: row; flex: 1; height: 100%; overflow: hidden; background: var(--vms-neu-bg);">
    <AdminSidebar :activePage="activePage" @selectPage="activePage = $event" />

    <main class="vms-admin-content" style="flex: 1; min-width: 0; padding: 1.5rem 2rem; overflow-y: auto;">
      <AdminCamerasView v-if="activePage === 'video_streams'" />
      <AdminMarketplaceView v-else-if="activePage === 'marketplace'" />
      <template v-else-if="parsedPlugin">
        <AnalyticPluginCrudView v-if="parsedPlugin.type === 'instances'" :plugin-id="parsedPlugin.pluginId" />
        <AnalyticPluginEventsView v-else-if="parsedPlugin.type === 'events'" :plugin-id="parsedPlugin.pluginId" />
      </template>
      <AdminAlarmsView v-else-if="activePage === 'alarms'" />
      <AdminUsersView v-else-if="activePage === 'users'" />
      <AdminLayoutsView v-else-if="activePage === 'layouts'" />
      <AdminCarouselsView v-else-if="activePage === 'carousels'" />
      <AdminMapsView v-else-if="activePage === 'maps'" />
      <AdminWorkflowsView v-else-if="activePage === 'workflows'" />
      <StorageDisksView v-else-if="activePage === 'storage'" />
      <AdminPerformanceView v-else-if="activePage === 'performance'" />
      <AdminLogsView v-else-if="activePage === 'logs'" />
    </main>
  </div>
</template>
