<script setup lang="ts">
import PreferenceTag from '@/components/common/PreferenceTag.vue';

withDefaults(
    defineProps<{
        /** The roster's answer values (`Roster.values`) */
        values: string[];
        showAssigned?: boolean;
    }>(),
    { showAssigned: true },
);
</script>

<template>
    <ul aria-label="Legend" class="flex flex-wrap items-center gap-3.5">
        <li v-for="value in values" :key="value">
            <PreferenceTag :value="value" />
        </li>
        <li v-if="showAssigned">
            <PreferenceTag assigned label="Assigned" />
        </li>
    </ul>
</template>
