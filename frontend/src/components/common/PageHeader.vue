<script setup lang="ts">
defineProps<{
    title: string;
    /** Smaller line under the title (used when there is no meta row) */
    description?: string;
    /** Shown next to the title in a lighter weight, e.g. the number of members */
    count?: number;
}>();
</script>

<template>
    <header class="flex flex-col gap-4 md:flex-row md:items-end md:gap-3">
        <div class="flex min-w-0 grow flex-col gap-1.5">
            <nav
                v-if="$slots.breadcrumb"
                aria-label="Breadcrumb"
                class="text-[13px] text-surface-600 dark:text-surface-200"
            >
                <slot name="breadcrumb" />
            </nav>
            <div class="flex flex-wrap items-center gap-x-3.5 gap-y-1">
                <h1 class="min-w-0 font-display text-[22px] font-semibold tracking-[-0.02em] md:text-[30px]">
                    {{ title }}
                    <span v-if="count !== undefined" class="font-medium text-surface-600 dark:text-surface-200">
                        {{ count }}
                    </span>
                </h1>
                <slot name="status" />
            </div>
            <p v-if="description" class="text-[15px] text-surface-600 dark:text-surface-200">{{ description }}</p>
            <div
                v-if="$slots.meta"
                class="flex flex-wrap items-center gap-x-[18px] gap-y-1 text-sm text-surface-600 dark:text-surface-200"
            >
                <slot name="meta" />
            </div>
        </div>
        <div v-if="$slots.actions" class="flex shrink-0 flex-wrap items-center gap-3">
            <slot name="actions" />
        </div>
    </header>
</template>
