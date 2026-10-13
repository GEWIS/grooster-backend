<script setup lang="ts">
import { computed } from 'vue';
import { OiCheck, OiChevronDown } from '@openvue/openicons-vue';
import { isPreferenceSet, preferenceLabel, preferenceToneClass } from '@/constants/preferences';

export type PreferenceSize = 'tag' | 'cell' | 'pill' | 'row';

const props = withDefaults(
    defineProps<{
        value?: string | null;
        size?: PreferenceSize;
        assigned?: boolean;
        /** Overrides the visible text (e.g. "Assigned" in the legend) */
        label?: string;
        /** Shows a chevron; used when a native select sits on top (PreferenceSelect) */
        chevron?: boolean;
    }>(),
    { value: null, size: 'tag', assigned: false, label: undefined, chevron: false },
);

const sizeClass: Record<PreferenceSize, string> = {
    tag: 'gap-1.5 rounded-full px-2.5 py-[3px] text-xs font-semibold',
    cell: 'h-[52px] w-full min-w-0 gap-[5px] pl-2 pr-1.5 text-[13px] font-medium focus-within:z-1 focus-within:outline-2 focus-within:-outline-offset-2 focus-within:outline-primary',
    pill: 'h-9 gap-1.5 rounded-full pl-3 pr-2.5 text-[13px] font-semibold focus-within:outline-2 focus-within:outline-offset-2 focus-within:outline-primary',
    row: 'h-11 w-[150px] gap-2 rounded-[10px] pl-3.5 pr-3 text-[15px] font-semibold focus-within:outline-2 focus-within:outline-offset-2 focus-within:outline-primary',
};

const ringClass = computed(() => {
    if (props.assigned) return 'ring-2 ring-inset ring-surface-900 dark:ring-surface-0';
    // Empty pills/rows need a visible edge on white
    if (!isPreferenceSet(props.value) && (props.size === 'pill' || props.size === 'row')) {
        return 'ring-1 ring-inset ring-surface-400 dark:ring-surface-500';
    }
    return '';
});
</script>

<template>
    <span
        :class="[
            'group relative inline-flex shrink-0 items-center whitespace-nowrap',
            sizeClass[size],
            // An assigned cell without an answer still reads as body text, not muted
            assigned && !isPreferenceSet(value) ? 'bg-pref-none text-color' : preferenceToneClass(value),
            ringClass,
            assigned && (size === 'tag' || size === 'pill') ? 'pr-[26px]' : '',
        ]"
    >
        <span class="min-w-0 truncate">{{ label ?? preferenceLabel(value) }}</span>
        <OiChevronDown
            v-if="chevron"
            aria-hidden="true"
            class="ml-auto shrink-0 opacity-50 group-focus-within:opacity-100 group-hover:opacity-100"
            :size="size === 'row' ? 16 : 12"
        />
        <span
            v-if="assigned"
            aria-hidden="true"
            class="pointer-events-none absolute right-1 top-1 flex size-3.5 items-center justify-center rounded-full bg-surface-900 text-white dark:bg-surface-0 dark:text-surface-900"
        >
            <OiCheck :size="11" />
        </span>
        <slot />
    </span>
</template>
