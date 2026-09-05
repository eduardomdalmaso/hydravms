<script setup lang="ts">
import type { UserItem } from '../../../types/userTree'
import UserDetailLeftPane from './UserDetailLeftPane.vue'
import UserPermissionsRightPane from './UserPermissionsRightPane.vue'

defineProps<{ user: UserItem }>()
const emit = defineEmits<{ (e: 'saved', msg: string): void }>()
</script>

<template>
  <div class="vms-desktop-split-layout">
    <UserDetailLeftPane :user="user" />
    <UserPermissionsRightPane :user="user" @saved="(msg) => emit('saved', msg)" />
  </div>
</template>
