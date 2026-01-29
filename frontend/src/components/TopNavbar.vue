<script setup lang="ts">
import { useRoute } from 'vue-router';
import { useOrganStore } from '@/stores/organ.store.js';

const route = useRoute();
const store = useOrganStore();

const isActive = (path) => route.path === path;
</script>

<template>
  <nav
    class="sticky top-0 z-50 flex gap-2 items-center border-b border-emerald-800 bg-emerald-950 px-6 py-2.5 shadow-md"
  >
    <div class="flex items-center gap-3">
      <RouterLink to="/">
        <div
          class="flex h-10 w-10 items-center justify-center overflow-hidden rounded-lg bg-emerald-50 p-1.5 shadow-sm"
        >
          <img alt="GRooster Logo" class="h-full w-full object-contain" src="@/assets/grooster-logo.png" />
        </div>
      </RouterLink>
    </div>

    <div class="flex items-center gap-1">
      <template v-if="route.path !== '/' && store.organ">
        <RouterLink
          :class="[
            'flex items-center gap-2 rounded-md px-4 py-2 text-sm font-medium transition-all duration-200',
            isActive('/rosters/' + store.organ)
              ? 'bg-emerald-800 text-white shadow-sm'
              : 'text-emerald-100 hover:bg-emerald-900 hover:text-white',
          ]"
          :to="`/rosters/${store.organ}`"
        >
          <i class="pi pi-table text-xs" />
          <span>Rosters</span>
        </RouterLink>

        <RouterLink
          :class="[
            'flex items-center gap-2 rounded-md px-4 py-2 text-sm font-medium transition-all duration-200',
            isActive('/templates/' + store.organ)
              ? 'bg-emerald-800 text-white shadow-sm'
              : 'text-emerald-100 hover:bg-emerald-900 hover:text-white',
          ]"
          :to="`/templates/${store.organ}`"
        >
          <i class="pi pi-clone text-xs" />
          <span>Templates</span>
        </RouterLink>
      </template>
    </div>
  </nav>
</template>
