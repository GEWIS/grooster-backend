<script setup lang="ts">
import { useRoute } from 'vue-router';
import { ref } from 'vue';
import { useOrganStore } from '@/stores/organ.store.js';

const route = useRoute();
const store = useOrganStore();

const isMobileMenuOpen = ref(false);
const isActive = (path) => route.path === path;
</script>

<template>
  <nav class="relative bg-emerald-950">
    <div class="flex items-center justify-between px-4 py-2">
      <div class="flex items-center gap-6">
        <RouterLink to="/">
          <div
            class="flex h-10 w-10 items-center justify-center overflow-hidden rounded-lg bg-emerald-50 p-1.5 shadow-sm"
          >
            <img alt="GRooster Logo" class="h-full w-full object-contain" src="@/assets/grooster-logo.png" />
          </div>
        </RouterLink>

        <div class="hidden md:flex items-center gap-1">
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
      </div>

      <div class="flex items-center gap-2">
        <div class="hidden md:flex items-center">
          <template v-if="route.path !== '/' && store.organ">
            <RouterLink
              v-if="store.organ"
              :class="[
                'group flex items-center gap-3 rounded-lg pl-3 pr-1 py-1 transition-all duration-200 border border-transparent',
                isActive(`/organ/${store.organ}/profile`)
                  ? 'bg-emerald-800/50 border-emerald-700 text-white'
                  : 'text-emerald-100 hover:bg-emerald-900',
              ]"
              :to="`/organ/${store.organ}/profile`"
            >
              <span class="text-sm font-medium">{{ store.organ }}</span>
              <div
                class="flex h-8 w-8 items-center justify-center rounded-full bg-emerald-100 text-emerald-900 shadow-sm transition-transform group-hover:scale-105"
              >
                <i class="pi pi-user text-sm" />
              </div>
            </RouterLink>
          </template>
        </div>

        <button
          class="md:hidden flex items-center justify-center h-10 w-10 rounded-lg text-emerald-100 hover:bg-emerald-900"
          @click="isMobileMenuOpen = !isMobileMenuOpen"
        >
          <i :class="['pi', isMobileMenuOpen ? 'pi-times' : 'pi-bars', 'text-xl']"></i>
        </button>
      </div>
    </div>

    <div
      v-if="isMobileMenuOpen"
      class="absolute top-full left-0 right-0 z-50 bg-emerald-950 border-t border-emerald-800 p-4 md:hidden flex flex-col gap-2 shadow-lg"
    >
      <template v-if="route.path !== '/' && store.organ">
        <RouterLink class="text-emerald-100 py-2" to="/" @click="isMobileMenuOpen = false">Home</RouterLink>
        <RouterLink class="text-emerald-100 py-2" :to="`/rosters/${store.organ}`" @click="isMobileMenuOpen = false"
          >Rosters</RouterLink
        >
        <RouterLink class="text-emerald-100 py-2" :to="`/templates/${store.organ}`" @click="isMobileMenuOpen = false"
          >Templates</RouterLink
        >
        <RouterLink
          class="text-emerald-100 py-2"
          :to="`/organ/${store.organ}/profile`"
          @click="isMobileMenuOpen = false"
          >Profile</RouterLink
        >
      </template>
    </div>
  </nav>
</template>
