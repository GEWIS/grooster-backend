<script setup lang="ts">
import TopNavbar from '@/components/TopNavbar.vue';
</script>

<template>
  <div class="flex flex-col min-h-screen">
    <TopNavbar />
    <div class="flex-1 p-4">
      <RouterView />
    </div>
  </div>
</template>

<style scoped></style>
