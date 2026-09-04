<script setup lang="ts">
import { ref } from "vue"
import type { MapResource, CarouselConfig, CameraStreamInfo, CustomLayout, AlarmItemInfo } from "../../types/mosaic"
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

const isCollapsed = ref(false)
const openSections = ref({ cameras: true, alarms: true, layouts: true, maps: false, carousel: false })

const toggleSection = (sec: "cameras" | "alarms" | "layouts" | "maps" | "carousel") => {
  openSections.value[sec] = !openSections.value[sec]
}

const sampleCameras: CameraStreamInfo[] = [
  { id: "cam_01", name: "Portaria Principal", location: "Entrada", status: "online", protocol: "RTSP", has_ptz: true, fps: 30, resolution: "1080P" },
  { id: "cam_02", name: "Estacionamento", location: "Patio", status: "online", protocol: "ONVIF", has_ptz: false, fps: 25, resolution: "1080P" },
  { id: "cam_03", name: "Corredor Galpao", location: "Docas", status: "online", protocol: "RTMP", has_ptz: false, fps: 30, resolution: "1080P" },
  { id: "cam_04", name: "Perimetro Fundos", location: "Muro", status: "online", protocol: "RTSP", has_ptz: true, fps: 30, resolution: "4K" },
  { id: "cam_05", name: "Acesso Docas B", location: "Docas", status: "offline", protocol: "ONVIF", has_ptz: false, fps: 30, resolution: "1080P" }
]

const sampleAlarms: AlarmItemInfo[] = [
  { id: "alm_01", name: "Sensor Perimetro Norte", zone: "Zona 01", status: "online", type: "IVS" },
  { id: "alm_02", name: "Barreira Infravermelha Docas", zone: "Zona 02", status: "alert", type: "PIR" },
  { id: "alm_03", name: "Porta Sala Servidores", zone: "Zona 03", status: "online", type: "MAG" },
  { id: "alm_04", name: "Detector Fumaca Bloco B", zone: "Zona 04", status: "offline", type: "SMK" }
]

const sampleMaps: MapResource[] = [
  { id: "map_01", name: "Planta Baixa - Galpao 01", image_url: "/map1.png", cameras_count: 8 },
  { id: "map_02", name: "Perimetro & Acessos", image_url: "/map2.png", cameras_count: 14 }
]

const sampleCarousels: CarouselConfig[] = [
  { id: "car_01", name: "Ronda Entradas", camera_ids: ["cam_01", "cam_02"], interval_seconds: 10 }
]
</script>

<template>
  <aside class="vms-asset-sidebar" :class="{ collapsed: isCollapsed }">
    <div class="vms-sidebar-toggle-line" @click="isCollapsed = !isCollapsed">
      <div class="vms-sidebar-toggle-pill">{{ isCollapsed ? "▶" : "◀" }}</div>
    </div>

    <div v-show="!isCollapsed" style="display: flex; flex-direction: column; height: 100%; width: 250px; overflow-y: auto;">
      <SidebarCamerasSection
        :cameras="sampleCameras"
        :isOpen="openSections.cameras"
        @toggle="toggleSection('cameras')"
        @selectCamera="emit('selectCamera', $event)"
      />
      <SidebarAlarmsSection
        :alarms="sampleAlarms"
        :isOpen="openSections.alarms"
        @toggle="toggleSection('alarms')"
        @selectAlarm="emit('selectAlarm', $event)"
      />
      <SidebarLayoutsSection
        :layouts="layouts"
        :activeLayoutId="activeLayoutId"
        :isOpen="openSections.layouts"
        @toggle="toggleSection('layouts')"
        @selectLayout="emit('selectLayout', $event)"
        @contextMenu="emit('layoutContextMenu', $event)"
      />
      <SidebarMapsSection
        :maps="sampleMaps"
        :isOpen="openSections.maps"
        @toggle="toggleSection('maps')"
        @selectMap="emit('selectMap', $event)"
      />
      <SidebarCarouselsSection
        :carousels="sampleCarousels"
        :isOpen="openSections.carousel"
        @toggle="toggleSection('carousel')"
        @selectCarousel="emit('selectCarousel', $event)"
        @openModal="emit('openCarouselModal')"
      />
    </div>
  </aside>
</template>
