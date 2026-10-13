<script setup lang="ts">
import ConfirmDialog from 'openvue/confirmdialog';
import AppTopBar from '@/components/AppTopBar.vue';
import { confirmDialogPt } from '@/theme/pt';
</script>

<template>
    <div class="flex min-h-screen flex-col">
        <Toast />
        <ConfirmDialog :pt="confirmDialogPt" />
        <AppTopBar />
        <main class="flex flex-1 flex-col">
            <RouterView />
        </main>
    </div>
</template>
