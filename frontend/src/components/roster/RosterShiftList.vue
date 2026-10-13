<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue';
import { OiCheck, OiComment, OiComments, OiLock, OiLockOpen, OiPlus } from '@openvue/openicons-vue';
import PreferenceSelect from '@/components/common/PreferenceSelect.vue';
import PreferenceTag from '@/components/common/PreferenceTag.vue';
import RosterCommentPopover from '@/components/roster/RosterCommentPopover.vue';
import type { RosterGrid } from '@/composables/useRosterGrid';

const props = defineProps<{ grid: RosterGrid }>();

const comments = ref<InstanceType<typeof RosterCommentPopover>>();
const activeId = ref<number | null>(null);
const active = computed(
    () => props.grid.shifts.value.find((shift) => shift.id === activeId.value) ?? props.grid.shifts.value[0],
);
watch(
    () => props.grid.shifts.value,
    (shifts) => {
        if (!shifts.some((shift) => shift.id === activeId.value)) activeId.value = shifts[0]?.id ?? null;
    },
    { immediate: true },
);

const canAssign = computed(() => props.grid.isLocked.value && props.grid.isAdmin.value);

// Roving tabindex for the tablist: arrow keys move between shifts
const tabs = ref<HTMLButtonElement[]>([]);

async function onTabKey(event: KeyboardEvent, index: number) {
    const count = props.grid.shifts.value.length;
    const next = { ArrowRight: index + 1, ArrowLeft: index - 1, Home: 0, End: count - 1 }[event.key];

    if (next === undefined) return;

    event.preventDefault();
    const target = props.grid.shifts.value[(next + count) % count];

    if (!target) return;

    activeId.value = target.id;

    await nextTick();
    tabs.value[(next + count) % count]?.focus();
}
</script>

<template>
    <div class="flex flex-col">
        <div
            v-if="grid.shifts.value.length"
            aria-label="Shifts"
            class="flex shrink-0 gap-2 overflow-x-auto border-b border-surface bg-surface-0 px-4 py-3 dark:bg-surface-900"
            role="tablist"
        >
            <button
                v-for="(shift, index) in grid.shifts.value"
                :id="`shift-tab-${shift.id}`"
                :key="shift.id"
                ref="tabs"
                :aria-controls="`shift-panel-${shift.id}`"
                :aria-selected="active?.id === shift.id"
                :class="[
                    'h-10 shrink-0 rounded-full border px-4 text-sm font-semibold focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary',
                    active?.id === shift.id
                        ? 'border-topbar bg-topbar text-white'
                        : 'border-surface-300 bg-surface-0 text-surface-700 dark:border-surface-600 dark:bg-surface-900 dark:text-surface-100',
                ]"
                role="tab"
                :tabindex="active?.id === shift.id ? 0 : -1"
                type="button"
                @click="activeId = shift.id"
                @keydown="onTabKey($event, index)"
            >
                {{ shift.name }}
            </button>
        </div>

        <div
            v-if="active"
            :id="`shift-panel-${active.id}`"
            :aria-labelledby="`shift-tab-${active.id}`"
            class="flex flex-col px-4 pb-4 pt-3"
            role="tabpanel"
        >
            <div class="flex items-baseline px-1 pb-2.5 pt-1">
                <h2 class="grow text-[15px] font-semibold">{{ active.name }}</h2>
                <span
                    v-if="!grid.isLocked.value && grid.me.value"
                    class="text-[13px] text-surface-600 dark:text-surface-200"
                >
                    Tap your answer to change it
                </span>
            </div>
            <ul class="overflow-hidden rounded-card border border-surface bg-surface-0 dark:bg-surface-900">
                <li
                    v-for="member in grid.members.value"
                    :key="member.id"
                    class="flex min-h-16 items-center gap-2 border-b border-surface-100 py-2 pl-3 pr-2.5 last:border-b-0 dark:border-surface-800"
                >
                    <div class="flex min-w-0 grow flex-col items-start gap-1">
                        <span class="max-w-full truncate text-[15px] font-semibold" :title="member.name">
                            {{ member.displayName }}
                            <span v-if="member.isMe" class="font-normal text-surface-600 dark:text-surface-200"
                                >(you)</span
                            >
                        </span>
                        <PreferenceSelect
                            v-if="grid.canEdit(member)"
                            :assigned="grid.isAssigned(member.id, active.id)"
                            :label="`${member.displayName}, ${active.name}: answer`"
                            :model-value="grid.answerValue(member.id, active.id)"
                            :options="grid.values.value"
                            size="pill"
                            @update:model-value="grid.setAnswer(member.id, active.id, $event)"
                        />
                        <PreferenceTag
                            v-else
                            :assigned="grid.isAssigned(member.id, active.id)"
                            size="pill"
                            :value="grid.answerValue(member.id, active.id)"
                        />
                    </div>
                    <button
                        :aria-label="member.isMe ? 'Your comment' : `Comment by ${member.displayName}`"
                        :class="[
                            'flex size-11 shrink-0 items-center justify-center rounded-full hover:bg-surface-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary dark:hover:bg-surface-800',
                            grid.commentFor(member.id)
                                ? 'text-primary-700 dark:text-primary-300'
                                : 'text-surface-600 dark:text-surface-200',
                        ]"
                        type="button"
                        @click="comments?.open($event, member)"
                    >
                        <OiComments v-if="grid.commentFor(member.id)" aria-hidden="true" :size="20" />
                        <OiComment v-else aria-hidden="true" :size="20" />
                    </button>
                    <button
                        v-if="canAssign"
                        :aria-label="`Assign ${member.displayName} to ${active.name}`"
                        :aria-pressed="grid.isAssigned(member.id, active.id)"
                        :class="[
                            'flex size-11 shrink-0 items-center justify-center rounded-full border-[1.5px] focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary disabled:opacity-50',
                            grid.isAssigned(member.id, active.id)
                                ? 'border-primary bg-primary text-primary-contrast'
                                : 'border-surface-400 bg-surface-0 text-primary dark:border-surface-500 dark:bg-surface-900',
                        ]"
                        :disabled="grid.isBusy(member.id, active.id)"
                        type="button"
                        @click="grid.setAssigned(member.id, active.id, !grid.isAssigned(member.id, active.id))"
                    >
                        <OiCheck v-if="grid.isAssigned(member.id, active.id)" aria-hidden="true" :size="20" />
                        <OiPlus v-else aria-hidden="true" :size="20" />
                    </button>
                </li>
            </ul>
        </div>
        <p v-else class="px-5 py-6 text-sm text-surface-600 dark:text-surface-200">This roster has no shifts yet.</p>

        <div
            v-if="grid.isAdmin.value"
            class="sticky bottom-0 mt-auto flex gap-2.5 border-t border-surface bg-surface-0 px-4 pb-[max(12px,env(safe-area-inset-bottom))] pt-3 dark:bg-surface-900"
        >
            <Button
                class="h-12 grow-0 basis-[40%] text-sm"
                label="Auto-fill"
                outlined
                severity="secondary"
                @click="grid.autoFill()"
            />
            <Button
                v-if="grid.isLocked.value"
                class="h-12 grow text-sm"
                label="Unlock roster"
                outlined
                severity="secondary"
                @click="grid.unlock()"
            >
                <template #icon><OiLockOpen aria-hidden="true" :size="16" /></template>
            </Button>
            <Button v-else class="h-12 grow text-sm" label="Lock roster" @click="grid.lock()">
                <template #icon><OiLock aria-hidden="true" :size="16" /></template>
            </Button>
        </div>
        <RosterCommentPopover ref="comments" :grid="grid" />
    </div>
</template>
