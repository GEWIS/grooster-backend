<script setup lang="ts">
import { computed } from 'vue';
import PreferenceTag from '@/components/common/PreferenceTag.vue';
import { NOT_SET_LABEL } from '@/constants/preferences';

const props = withDefaults(
    defineProps<{
        modelValue?: string | null;
        /** The roster's answer values (`Roster.values`) */
        options: string[];
        size?: 'cell' | 'pill' | 'row';
        assigned?: boolean;
        /** Accessible name of the select, e.g. "{member}, {shift} (assigned): preference" */
        label: string;
    }>(),
    { modelValue: null, size: 'cell', assigned: false },
);

const emit = defineEmits<{ 'update:modelValue': [value: string] }>();

// Keep an answer visible even if the roster no longer offers that value
const choices = computed(() =>
    props.modelValue && !props.options.includes(props.modelValue)
        ? [props.modelValue, ...props.options]
        : props.options,
);

function onChange(event: Event) {
    emit('update:modelValue', (event.target as HTMLSelectElement).value);
}
</script>

<template>
    <PreferenceTag :assigned="assigned" chevron :size="size" :value="modelValue">
        <select
            :aria-label="label"
            :class="[
                'absolute inset-x-0 w-full cursor-pointer appearance-none text-base opacity-0',
                // The pill is 36px tall; extend the tap area to 44px on phones
                size === 'pill' ? '-top-1 h-[calc(100%+8px)]' : 'top-0 h-full',
            ]"
            :value="modelValue ?? ''"
            @change="onChange"
        >
            <option v-if="!modelValue" disabled value="">{{ NOT_SET_LABEL }}</option>
            <option v-for="option in choices" :key="option" :value="option">{{ option }}</option>
        </select>
    </PreferenceTag>
</template>
