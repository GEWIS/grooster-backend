<script setup lang="ts">
import { computed } from 'vue';
import { OiTimes } from '@openvue/openicons-vue';
import PreferenceTag from '@/components/common/PreferenceTag.vue';
import type { RosterGrid } from '@/composables/useRosterGrid';

const props = defineProps<{
    grid: RosterGrid;
    shiftId: number;
}>();

const emit = defineEmits<{ close: [] }>();

const shift = computed(() => props.grid.shifts.value.find((s) => s.id === props.shiftId));
const assigned = computed(() => props.grid.assignedMembers(props.shiftId));
const others = computed(() => props.grid.candidates(props.shiftId));
</script>

<template>
    <div v-if="shift" class="flex flex-col">
        <div class="flex items-center gap-2 border-b border-surface pb-3.5 pl-5 pr-3 pt-4">
            <div class="flex min-w-0 grow flex-col gap-0.5">
                <span class="text-xs font-semibold uppercase tracking-[0.06em] text-surface-600 dark:text-surface-200">
                    Assign
                </span>
                <h2 class="truncate font-display text-2xl font-semibold tracking-[-0.01em]">{{ shift.name }}</h2>
            </div>
            <button
                aria-label="Close assign panel"
                class="flex size-10 shrink-0 items-center justify-center rounded-lg text-surface-700 hover:bg-surface-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary dark:text-surface-100 dark:hover:bg-surface-800"
                type="button"
                @click="emit('close')"
            >
                <OiTimes aria-hidden="true" :size="20" />
            </button>
        </div>

        <section aria-labelledby="assign-assigned" class="flex flex-col gap-0.5 pb-2 pl-5 pr-4 pt-3.5">
            <h3 id="assign-assigned" class="mb-1 text-[13px] font-semibold text-surface-700 dark:text-surface-100">
                Assigned
            </h3>
            <ul v-if="assigned.length">
                <li v-for="member in assigned" :key="member.id" class="flex h-11 items-center gap-2.5">
                    <span class="min-w-0 grow truncate text-sm font-medium" :title="member.name">
                        {{ member.displayName }}
                    </span>
                    <PreferenceTag :value="grid.answerValue(member.id, shiftId)" />
                    <button
                        :aria-label="`Remove ${member.displayName} from ${shift.name}`"
                        class="flex size-8 shrink-0 items-center justify-center rounded-lg text-surface-700 hover:bg-surface-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary disabled:opacity-50 dark:text-surface-100 dark:hover:bg-surface-800"
                        :disabled="grid.isBusy(member.id, shiftId)"
                        type="button"
                        @click="grid.setAssigned(member.id, shiftId, false)"
                    >
                        <OiTimes aria-hidden="true" :size="18" />
                    </button>
                </li>
            </ul>
            <p v-else class="py-2 text-[13px] text-surface-600 dark:text-surface-200">
                Nobody yet. Pick someone below.
            </p>
        </section>

        <section
            aria-labelledby="assign-others"
            class="flex flex-col gap-0.5 border-t border-surface-100 pb-4 pl-5 pr-4 pt-3 dark:border-surface-800"
        >
            <h3 id="assign-others" class="mb-1 text-[13px] font-semibold text-surface-700 dark:text-surface-100">
                Everyone else
                <span class="font-normal text-surface-600 dark:text-surface-200">· best match first</span>
            </h3>
            <ul>
                <li v-for="member in others" :key="member.id" class="flex h-11 items-center gap-2.5">
                    <span class="min-w-0 grow truncate text-sm font-medium" :title="member.name">
                        {{ member.displayName }}
                    </span>
                    <PreferenceTag :value="grid.answerValue(member.id, shiftId)" />
                    <Button
                        :aria-label="`Assign ${member.displayName} to ${shift.name}`"
                        class="h-[34px] w-[72px] shrink-0 bg-surface-0 px-3 text-[13px] dark:bg-surface-900"
                        :disabled="grid.isBusy(member.id, shiftId)"
                        label="Assign"
                        outlined
                        severity="secondary"
                        @click="grid.setAssigned(member.id, shiftId, true)"
                    />
                </li>
            </ul>
            <p v-if="!others.length" class="py-2 text-[13px] text-surface-600 dark:text-surface-200">
                Everyone is assigned.
            </p>
        </section>
    </div>
</template>
