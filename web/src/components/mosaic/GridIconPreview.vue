<script setup lang="ts">
import type { GridLayout } from '../../types/mosaic'

defineProps<{ grid: GridLayout }>()
</script>

<template>
  <div class="vms-grid-mini-icon">
    <!-- 1x1 / Auto -->
    <div v-if="grid === '1x1' || grid === 'auto'" class="vms-mini-grid vms-mini-1x1">
      <span class="vms-mini-cell"></span>
    </div>
    <!-- 1x2 -->
    <div v-else-if="grid === '1x2'" class="vms-mini-grid vms-mini-1x2">
      <span class="vms-mini-cell"></span><span class="vms-mini-cell"></span>
    </div>
    <!-- 2x2 -->
    <div v-else-if="grid === '2x2'" class="vms-mini-grid vms-mini-2x2">
      <span v-for="i in 4" :key="i" class="vms-mini-cell"></span>
    </div>
    <!-- 3x3 -->
    <div v-else-if="grid === '3x3'" class="vms-mini-grid vms-mini-3x3">
      <span v-for="i in 9" :key="i" class="vms-mini-cell"></span>
    </div>
    <!-- 4x4 / 8x8 / 10x10 -->
    <div v-else-if="grid === '4x4' || grid === '8x8' || grid === '10x10'" class="vms-mini-grid vms-mini-4x4">
      <span v-for="i in 16" :key="i" class="vms-mini-cell"></span>
    </div>
    <!-- 1+5 -->
    <div v-else-if="grid === '1+5'" class="vms-mini-grid vms-mini-1-5">
      <span class="vms-mini-cell vms-mini-hero"></span>
      <span v-for="i in 5" :key="i" class="vms-mini-cell"></span>
    </div>
    <!-- 1+7 -->
    <div v-else-if="grid === '1+7'" class="vms-mini-grid vms-mini-1-7">
      <span class="vms-mini-cell vms-mini-hero-lg"></span>
      <span v-for="i in 7" :key="i" class="vms-mini-cell"></span>
    </div>
    <!-- 1+12 -->
    <div v-else-if="grid === '1+12'" class="vms-mini-grid vms-mini-1-12">
      <span class="vms-mini-cell vms-mini-hero-xl"></span>
      <span v-for="i in 12" :key="i" class="vms-mini-cell"></span>
    </div>
  </div>
</template>
