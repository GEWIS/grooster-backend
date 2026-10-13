<script setup lang="ts">
import Tag from 'openvue/tag';
import { OiLock } from '@openvue/openicons-vue';

// A roster is either open (collecting answers) or locked (`Roster.saved`).
defineProps<{ saved: boolean }>();
</script>

<template>
    <Tag
        :pt="{
            root: {
                class: [
                    'inline-flex items-center gap-1.5 rounded-full px-2.5 py-[3px] text-xs font-semibold',
                    saved
                        ? 'bg-primary-100 text-primary-800 dark:bg-primary-900 dark:text-primary-100'
                        : 'bg-status-warn text-status-warn-fg',
                ],
            },
            label: { class: 'font-semibold' },
        }"
        :value="saved ? 'Locked' : 'Open'"
    >
        <template #icon>
            <OiLock v-if="saved" aria-hidden="true" :size="12" />
            <span v-else aria-hidden="true" class="size-[7px] rounded-full bg-status-warn-dot" />
        </template>
    </Tag>
</template>
