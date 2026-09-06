<script setup lang="ts">
import type { EnterpriseRondaItem } from '../../../types/rondaTree'
import RondaDetailLeftPane from './RondaDetailLeftPane.vue'
import RondaStreamsRightPane from './RondaStreamsRightPane.vue'

defineProps<{ ronda: EnterpriseRondaItem }>()
const emit = defineEmits<{ (e: 'saved', msg: string): void }>()
</script>

<template>
  <div class="vms-desktop-split-layout">
    <RondaDetailLeftPane :ronda="ronda" @saved="(msg) => emit('saved', msg)" />
    <RondaStreamsRightPane :ronda="ronda" @saved="(msg) => emit('saved', msg)" />
  </div>
</template>
