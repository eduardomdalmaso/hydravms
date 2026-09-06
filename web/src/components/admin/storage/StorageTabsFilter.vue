<script setup lang="ts">
import type { StorageRole } from '../../../types/storagePool'

defineProps<{
  activeFilter: StorageRole | 'ALL'
  searchQuery: string
  counts: { all: number; buffer: number; recordings: number; snapshots: number; database: number }
}>()

const emit = defineEmits<{
  (e: 'update:filter', filter: StorageRole | 'ALL'): void
  (e: 'update:search', query: string): void
}>()
</script>

<template>
  <div class="vms-flex-between" style="gap: 1rem; flex-wrap: wrap; align-items: center;">
    <div class="vms-flex-row" style="gap: 6px; flex-wrap: wrap;">
      <button
        class="vms-btn vms-btn-sm"
        :class="activeFilter === 'ALL' ? 'vms-btn-primary' : 'vms-btn-secondary'"
        style="font-weight: 700; font-size: 11px;"
        @click="emit('update:filter', 'ALL')"
      >
        TODOS ({{ counts.all }})
      </button>
      <button
        class="vms-btn vms-btn-sm"
        :class="activeFilter === 'HOT_BUFFER' ? 'vms-btn-primary' : 'vms-btn-secondary'"
        style="font-weight: 700; font-size: 11px;"
        @click="emit('update:filter', 'HOT_BUFFER')"
      >
        DISCO DE BUFFER ({{ counts.buffer }})
      </button>
      <button
        class="vms-btn vms-btn-sm"
        :class="activeFilter === 'WARM_ARCHIVE' ? 'vms-btn-primary' : 'vms-btn-secondary'"
        style="font-weight: 700; font-size: 11px;"
        @click="emit('update:filter', 'WARM_ARCHIVE')"
      >
        ARMAZENAMENTO DE GRAVACOES ({{ counts.recordings }})
      </button>
      <button
        class="vms-btn vms-btn-sm"
        :class="activeFilter === 'SNAPSHOTS' ? 'vms-btn-primary' : 'vms-btn-secondary'"
        style="font-weight: 700; font-size: 11px;"
        @click="emit('update:filter', 'SNAPSHOTS')"
      >
        ARMAZENAMENTO DE SNAPSHOTS ({{ counts.snapshots }})
      </button>
      <button
        class="vms-btn vms-btn-sm"
        :class="activeFilter === 'DATABASE' ? 'vms-btn-primary' : 'vms-btn-secondary'"
        style="font-weight: 700; font-size: 11px;"
        @click="emit('update:filter', 'DATABASE')"
      >
        BANCO DE DADOS ({{ counts.database }})
      </button>
    </div>

    <input
      :value="searchQuery"
      class="vms-auth-input"
      style="width: 240px; font-size: 12px; padding: 4px 10px;"
      placeholder="Filtrar por nome, nó ou caminho..."
      @input="emit('update:search', ($event.target as HTMLInputElement).value)"
    />
  </div>
</template>
