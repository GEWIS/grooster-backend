<script setup lang="ts">
import { computed, ref } from 'vue';
import { OiComment, OiComments } from '@openvue/openicons-vue';
import PreferenceSelect from '@/components/common/PreferenceSelect.vue';
import RosterCommentPopover from '@/components/roster/RosterCommentPopover.vue';
import type { RosterGrid } from '@/composables/useRosterGrid';

const props = defineProps<{
    grid: RosterGrid;
    /** Names of the roster's responsibles, for the intro line */
    responsibleNames: string[];
}>();

const comments = ref<InstanceType<typeof RosterCommentPopover>>();
const me = computed(() => props.grid.me.value);
</script>

<template>
    <div class="mx-auto flex w-full max-w-[560px] flex-col gap-2.5">
        <div v-if="me" class="flex items-center gap-2 px-1">
            <button
                :aria-label="grid.commentFor(me.id) ? 'Edit your comment' : 'Add a comment'"
                class="flex h-9 shrink-0 items-center gap-1.5 whitespace-nowrap rounded-lg border border-surface-300 bg-surface-0 px-3 text-[13px] font-semibold hover:bg-surface-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary dark:border-surface-600 dark:bg-surface-900 dark:hover:bg-surface-800"
                type="button"
                @click="comments?.open($event, me)"
            >
                <OiComments v-if="grid.commentFor(me.id)" aria-hidden="true" :size="18" />
                <OiComment v-else aria-hidden="true" :size="18" />
                {{ grid.commentFor(me.id) ? 'Edit comment' : 'Add comment' }}
            </button>
        </div>

        <ul v-if="me" class="overflow-hidden rounded-card border border-surface bg-surface-0 dark:bg-surface-900">
            <li
                v-for="shift in grid.shifts.value"
                :key="shift.id"
                class="flex min-h-[60px] items-center gap-3 border-b border-surface-100 py-2 pl-4 pr-2 last:border-b-0 dark:border-surface-800"
            >
                <span class="min-w-0 grow truncate text-base font-semibold">{{ shift.name }}</span>
                <PreferenceSelect
                    :assigned="grid.isAssigned(me.id, shift.id)"
                    :label="`${shift.name}: your answer`"
                    :model-value="grid.answerValue(me.id, shift.id)"
                    :options="grid.values.value"
                    size="row"
                    @update:model-value="grid.setAnswer(me.id, shift.id, $event)"
                />
            </li>
            <li v-if="!grid.shifts.value.length" class="px-4 py-5 text-sm text-surface-600 dark:text-surface-200">
                This roster has no shifts yet.
            </li>
        </ul>
        <p v-else class="px-1 text-sm text-surface-600 dark:text-surface-200">You are not a member of this roster.</p>
        <RosterCommentPopover ref="comments" :grid="grid" />
    </div>
</template>
