<script setup>
import { useRoute } from 'vue-router';
import { useOrganStore } from '@/stores/organ.store.js';

const route = useRoute();
const store = useOrganStore();

const isActive = (path) => route.path === path;
</script>

<template>
  <nav class="bg-green-800 shadow-md p-4 flex justify-between items-center">
    <div class="space-x-4">
      <RouterLink to="/" class="px-4 py-2 rounded-lg hover:bg-gray-400" :class="{ 'bg-gray-400': isActive('/') }">
        Home
      </RouterLink>
      <RouterLink
        v-if="route.path !== '/' && store.organ"
        :to="`/rosters/${store.organ}`"
        class="px-4 py-2 rounded-lg hover:bg-gray-400"
        :class="{ 'bg-gray-400': isActive('/rosters/' + store.organ) }"
      >
        Rosters
      </RouterLink>
      <RouterLink
        v-if="route.path !== '/' && store.organ"
        :to="`/templates/${store.organ}`"
        class="px-4 py-2 rounded-lg hover:bg-gray-400"
        :class="{ 'bg-gray-400': isActive('/templates/' + store.organ) }"
      >
        Templates
      </RouterLink>
    </div>
  </nav>
</template>

<style scoped></style>
