<script setup>
import { useRoute } from 'vue-router'

const route = useRoute()

const isActive = (path) => route.path === path;
</script>

<template>
    <nav class="bg-green-800 shadow-md p-4 flex justify-between items-center">
        <div class="space-x-4">
            <RouterLink
                to="/"
                class="px-4 py-2 rounded-lg hover:bg-gray-400"
                :class="{ 'bg-gray-400': isActive('/') }"
            >
                Home
            </RouterLink>
        </div>
    </nav>
</template>

<style scoped>

</style>