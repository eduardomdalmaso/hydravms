<script setup lang="ts">
import type { EnterpriseLayoutItem } from '../../../types/layoutTree'
import LayoutDetailLeftPane from './LayoutDetailLeftPane.vue'
import LayoutSlotsRightPane from './LayoutSlotsRightPane.vue'

defineProps<{ layout: EnterpriseLayoutItem }>()
const emit = defineEmits<{ (e: 'saved', msg: string): void }>()
</script>

<template>
  <div class="vms-desktop-split-layout">
    <LayoutDetailLeftPane :layout="layout" @saved="(msg) => emit('saved', msg)" />
    <LayoutSlotsRightPane :layout="layout" @saved="(msg) => emit('saved', msg)" />
  </div>
</template>
