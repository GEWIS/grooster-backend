<script setup lang="ts">
import MemberUsername from '@/components/organ-settings/MemberUsername.vue';
import MemberRoles from '@/components/organ-settings/MemberRoles.vue';
import { Role, useAuthStore } from '@/stores/auth.store';

const authStore = useAuthStore();
</script>

<template>
  <div class="max-w-2xl mx-auto mt-10 px-4">
    <MemberUsername />

    <MemberRoles v-if="authStore.can([Role.Admin, Role.Owner])" />
  </div>
</template>
