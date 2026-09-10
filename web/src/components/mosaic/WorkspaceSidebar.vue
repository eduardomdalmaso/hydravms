<script setup lang="ts">
import { ref, onMounted } from "vue"
import type { MapResource, CarouselConfig, CameraStreamInfo, CustomLayout, AlarmItemInfo } from "../../types/mosaic"
import { fetchCameras } from "../../services/api"
import SidebarCamerasSection from "./SidebarCamerasSection.vue"
import SidebarAlarmsSection from "./SidebarAlarmsSection.vue"
import SidebarLayoutsSection from "./SidebarLayoutsSection.vue"
import SidebarMapsSection from "./SidebarMapsSection.vue"
import SidebarCarouselsSection from "./SidebarCarouselsSection.vue"

defineProps<{ layouts: CustomLayout[]; activeLayoutId: string }>()
const emit = defineEmits<{
  (e: "selectCamera", camera: CameraStreamInfo): void
  (e: "selectAlarm", alarm: AlarmItemInfo): void
  (e: "selectMap", map: MapResource): void
  (e: "selectCarousel", carousel: CarouselConfig): void
  (e: "openCarouselModal"): void
  (e: "selectLayout", layoutId: string): void
  (e: "layoutContextMenu", payload: { event: MouseEvent; layout?: CustomLayout; isHeader?: boolean }): void
}>()

const isCollapsed = ref(false), cameras = ref<CameraStreamInfo[]>([])
const openSections = ref({ cameras: true, alarms: true, layouts: true, maps: false, carousel: false })
const toggleSection = (s: "cameras" | "alarms" | "layouts" | "maps" | "carousel") => { openSections.value[s] = !openSections.value[s] }

const liveAlarms = ref<AlarmItemInfo[]>([])
const liveMaps = ref<MapResource[]>([])
const liveCarousels = ref<CarouselConfig[]>([])

onMounted(async () => {
  const remote = await fetchCameras()
  if (remote.length > 0) {
    cameras.value = remote.map(c => ({
      id: c.id, name: c.name, location: c.location || 'Local', status: (c.status || 'online') as any,
      protocol: (c.protocol?.toUpperCase() as any) || 'RTSP', has_ptz: c.has_ptz, fps: c.fps || 30, resolution: c.resolution || '1080P'
    }))
  }
})
</script>

<template>
  <aside class="vms-asset-sidebar" :class="{ collapsed: isCollapsed }">
    <div class="vms-sidebar-toggle-line" @click="isCollapsed = !isCollapsed">
      <div class="vms-sidebar-toggle-pill">{{ isCollapsed ? "▶" : "◀" }}</div>
    </div>
    <div v-show="!isCollapsed" style="display: flex; flex-direction: column; height: 100%; width: 250px; overflow-y: auto;">
      <SidebarCamerasSection :cameras="cameras" :isOpen="openSections.cameras" @toggle="toggleSection('cameras')" @selectCamera="emit('selectCamera', $event)" />
      <SidebarAlarmsSection :alarms="liveAlarms" :isOpen="openSections.alarms" @toggle="toggleSection('alarms')" @selectAlarm="emit('selectAlarm', $event)" />
      <SidebarLayoutsSection :layouts="layouts" :activeLayoutId="activeLayoutId" :isOpen="openSections.layouts" @toggle="toggleSection('layouts')" @selectLayout="emit('selectLayout', $event)" @contextMenu="emit('layoutContextMenu', $event)" />
      <SidebarMapsSection :maps="liveMaps" :isOpen="openSections.maps" @toggle="toggleSection('maps')" @selectMap="emit('selectMap', $event)" />
      <SidebarCarouselsSection :carousels="liveCarousels" :isOpen="openSections.carousel" @toggle="toggleSection('carousel')" @selectCarousel="emit('selectCarousel', $event)" @openModal="emit('openCarouselModal')" />
    </div>
  </aside>
</template>
