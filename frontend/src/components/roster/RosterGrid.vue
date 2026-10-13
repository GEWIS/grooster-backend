<script setup lang="ts">
import { computed, ref } from 'vue';
import Menu from 'openvue/menu';
import {
    OiArrowLeft,
    OiArrowRight,
    OiBolt,
    OiComment,
    OiComments,
    OiLock,
    OiLockOpen,
    OiPlus,
    OiTrash,
} from '@openvue/openicons-vue';
import type { RosterShift } from '@gewis/grooster-backend-ts';
import PreferenceSelect from '@/components/common/PreferenceSelect.vue';
import PreferenceTag from '@/components/common/PreferenceTag.vue';
import RosterCommentPopover from '@/components/roster/RosterCommentPopover.vue';
import type { RosterGrid } from '@/composables/useRosterGrid';
import { useConfirmDelete } from '@/composables/useConfirmDelete';

const props = defineProps<{
    grid: RosterGrid;
    selectedShiftId: number | null;
}>();

const emit = defineEmits<{
    'update:selectedShiftId': [shiftId: number | null];
    'add-shift': [];
}>();

const confirmDelete = useConfirmDelete();
const comments = ref<InstanceType<typeof RosterCommentPopover>>();

const canAssign = computed(() => props.grid.isLocked.value && props.grid.isAdmin.value);
const canManageShifts = computed(() => !props.grid.isLocked.value && props.grid.isAdmin.value);

function toggleShift(shiftId: number) {
    emit('update:selectedShiftId', props.selectedShiftId === shiftId ? null : shiftId);
}

const shiftMenu = ref<InstanceType<typeof Menu>>();
const menuShift = ref<RosterShift | null>(null);

const shiftMenuItems = computed(() => {
    const shift = menuShift.value;
    if (!shift) return [];
    const index = props.grid.shifts.value.findIndex((s) => s.id === shift.id);
    return [
        {
            label: 'Move left',
            iconComponent: OiArrowLeft,
            disabled: index === 0,
            command: () => props.grid.moveShift(shift.id, 'left'),
        },
        {
            label: 'Move right',
            iconComponent: OiArrowRight,
            disabled: index === props.grid.shifts.value.length - 1,
            command: () => props.grid.moveShift(shift.id, 'right'),
        },
        { separator: true },
        {
            label: 'Delete shift',
            iconComponent: OiTrash,
            danger: true,
            command: () =>
                confirmDelete({
                    header: `Delete "${shift.name}"?`,
                    message: 'All answers for this shift are deleted too. This cannot be undone.',
                    acceptLabel: 'Delete shift',
                    onAccept: () => props.grid.removeShift(shift.id),
                }),
        },
    ];
});

function openShiftMenu(event: Event, shift: RosterShift) {
    menuShift.value = shift;
    shiftMenu.value?.toggle(event);
}

const cellLabel = (memberName: string, shift: RosterShift, assigned: boolean) =>
    `${memberName}, ${shift.name}${assigned ? ' (assigned)' : ''}: answer`;

const headerBase =
    'flex h-[52px] min-w-0 items-center justify-center border-b-[3px] border-l border-l-surface-200 px-1.5 text-sm font-semibold dark:border-l-surface-700';
const focusRing = 'focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-primary';
</script>

<template>
    <section
        aria-label="Answers and assignments"
        class="overflow-hidden rounded-card border border-surface bg-surface-0 dark:bg-surface-900"
    >
        <div class="overflow-x-auto">
            <div
                class="grid grid-cols-[196px_repeat(var(--shifts),minmax(84px,1fr))]"
                role="table"
                :style="{ '--shifts': Math.max(grid.shifts.value.length, 1) }"
            >
                <div class="contents" role="row">
                    <div
                        class="sticky left-0 z-[2] flex h-[52px] items-center bg-surface-50 px-5 text-xs font-semibold uppercase tracking-[0.06em] text-surface-600 dark:bg-surface-800 dark:text-surface-200"
                        role="columnheader"
                    >
                        Member
                    </div>
                    <div v-for="shift in grid.shifts.value" :key="shift.id" class="flex min-w-0" role="columnheader">
                        <button
                            v-if="canAssign"
                            :aria-label="`${shift.name}: assign people`"
                            :aria-pressed="selectedShiftId === shift.id"
                            :class="[
                                headerBase,
                                'w-full',
                                focusRing,
                                selectedShiftId === shift.id
                                    ? 'border-b-primary bg-highlight'
                                    : 'border-b-transparent bg-surface-50 hover:bg-surface-100 dark:bg-surface-800 dark:hover:bg-surface-700',
                            ]"
                            type="button"
                            @click="toggleShift(shift.id)"
                        >
                            <span class="truncate">{{ shift.name }}</span>
                        </button>
                        <button
                            v-else-if="canManageShifts"
                            aria-controls="shift-menu"
                            aria-haspopup="menu"
                            :aria-label="`${shift.name}: shift actions`"
                            :class="[
                                headerBase,
                                'w-full border-b-transparent bg-surface-50 hover:bg-surface-100 dark:bg-surface-800 dark:hover:bg-surface-700',
                                focusRing,
                            ]"
                            type="button"
                            @click="openShiftMenu($event, shift)"
                        >
                            <span class="truncate">{{ shift.name }}</span>
                        </button>
                        <span
                            v-else
                            :class="[headerBase, 'w-full border-b-transparent bg-surface-50 dark:bg-surface-800']"
                        >
                            <span class="truncate">{{ shift.name }}</span>
                        </span>
                    </div>
                    <div
                        v-if="grid.shifts.value.length === 0"
                        class="flex h-[52px] items-center border-l border-surface bg-surface-50 px-4 text-sm text-surface-600 dark:bg-surface-800 dark:text-surface-200"
                        role="columnheader"
                    >
                        No shifts yet
                    </div>
                </div>

                <div v-for="member in grid.members.value" :key="member.id" class="group/row contents" role="row">
                    <div
                        class="sticky left-0 z-[1] flex h-[52px] min-w-0 items-center gap-1.5 border-t border-surface bg-surface-0 pl-5 pr-2 group-hover/row:bg-surface-50 dark:bg-surface-900 dark:group-hover/row:bg-surface-800"
                        role="rowheader"
                    >
                        <span class="min-w-0 grow truncate text-sm font-semibold" :title="member.name">
                            {{ member.displayName }}
                            <span v-if="member.isMe" class="font-normal text-surface-600 dark:text-surface-200"
                                >(you)</span
                            >
                        </span>
                        <button
                            :aria-label="
                                member.isMe
                                    ? 'Your comment'
                                    : grid.commentFor(member.id)
                                      ? `Comment by ${member.displayName}`
                                      : `${member.displayName} has no comment`
                            "
                            class="flex size-8 shrink-0 items-center justify-center rounded-lg text-surface-600 hover:bg-surface-100 focus-visible:outline-2 focus-visible:outline-primary dark:text-surface-200 dark:hover:bg-surface-700"
                            :class="{ 'text-primary-700 dark:text-primary-300': grid.commentFor(member.id) }"
                            type="button"
                            @click="comments?.open($event, member)"
                        >
                            <OiComments v-if="grid.commentFor(member.id)" aria-hidden="true" :size="18" />
                            <OiComment v-else aria-hidden="true" :size="18" />
                        </button>
                    </div>
                    <div
                        v-for="shift in grid.shifts.value"
                        :key="shift.id"
                        :class="[
                            'min-w-0 border-l border-t border-surface',
                            selectedShiftId === shift.id ? 'border-l-primary-200 dark:border-l-primary-800' : '',
                        ]"
                        role="cell"
                    >
                        <PreferenceSelect
                            v-if="grid.canEdit(member)"
                            :assigned="grid.isAssigned(member.id, shift.id)"
                            :label="cellLabel(member.displayName, shift, grid.isAssigned(member.id, shift.id))"
                            :model-value="grid.answerValue(member.id, shift.id)"
                            :options="grid.values.value"
                            @update:model-value="grid.setAnswer(member.id, shift.id, $event)"
                        />
                        <PreferenceTag
                            v-else
                            :assigned="grid.isAssigned(member.id, shift.id)"
                            size="cell"
                            :value="grid.answerValue(member.id, shift.id)"
                        />
                    </div>
                    <div v-if="grid.shifts.value.length === 0" class="border-l border-t border-surface" role="cell" />
                </div>
            </div>
        </div>

        <div
            v-if="grid.isAdmin.value"
            class="flex flex-wrap items-center gap-3 border-t border-surface bg-surface-50 px-4 py-3 dark:bg-surface-800"
        >
            <Button
                class="h-10 px-2.5 text-sm"
                :disabled="grid.isLocked.value"
                label="Add shift"
                text
                @click="emit('add-shift')"
            >
                <template #icon><OiPlus aria-hidden="true" :size="16" /></template>
            </Button>
            <span class="grow" />
            <Button
                class="h-10 bg-surface-0 text-sm dark:bg-surface-900"
                label="Auto-fill from preferences"
                outlined
                severity="secondary"
                @click="grid.autoFill()"
            >
                <template #icon><OiBolt aria-hidden="true" :size="16" /></template>
            </Button>
            <Button
                v-if="grid.isLocked.value"
                class="h-10 bg-surface-0 text-sm dark:bg-surface-900"
                label="Unlock roster"
                outlined
                severity="secondary"
                @click="grid.unlock()"
            >
                <template #icon><OiLockOpen aria-hidden="true" :size="16" /></template>
            </Button>
            <Button v-else class="h-10 text-sm" label="Lock roster" @click="grid.lock()">
                <template #icon><OiLock aria-hidden="true" :size="16" /></template>
            </Button>
        </div>

        <Menu id="shift-menu" ref="shiftMenu" :model="shiftMenuItems" popup>
            <template #item="{ item, props: itemProps }">
                <a v-bind="itemProps.action" :class="['gap-2', item.danger ? 'text-red-700 dark:text-red-300' : '']">
                    <component :is="item.iconComponent" aria-hidden="true" class="shrink-0" :size="16" />
                    <span>{{ item.label }}</span>
                </a>
            </template>
        </Menu>
        <RosterCommentPopover ref="comments" :grid="grid" />
    </section>
</template>
