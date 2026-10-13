<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useRoute } from 'vue-router';
import { useToast } from 'openvue/usetoast';
import Menu from 'openvue/menu';
import Drawer from 'openvue/drawer';
import SelectButton from 'openvue/selectbutton';
import { OiClone, OiDownload, OiEllipsisH, OiPencil, OiPlus, OiTrash } from '@openvue/openicons-vue';
import type { Roster } from '@gewis/grooster-backend-ts';
import PageHeader from '@/components/common/PageHeader.vue';
import StatusTag from '@/components/common/StatusTag.vue';
import RosterGrid from '@/components/roster/RosterGrid.vue';
import RosterShiftList from '@/components/roster/RosterShiftList.vue';
import RosterAssignPanel from '@/components/roster/RosterAssignPanel.vue';
import RosterResponsibles from '@/components/roster/RosterResponsibles.vue';
import MyPreferences from '@/components/roster/MyPreferences.vue';
import NewRosterDialog from '@/components/roster/dialogs/NewRosterDialog.vue';
import EditDialog from '@/components/roster/dialogs/EditDialog.vue';
import AddShiftDialog from '@/components/roster/dialogs/AddShiftDialog.vue';
import { useRosterStore } from '@/stores/roster.store';
import { useTemplateStore } from '@/stores/template.store';
import { useRosterGrid } from '@/composables/useRosterGrid';
import { useMediaQuery } from '@/composables/useMediaQuery';
import { useConfirmDelete } from '@/composables/useConfirmDelete';
import { exportRoster } from '@/helpers/exportRoster';
import { formatRosterDate, isoWeek, parseRosterDate } from '@/helpers/date';

const route = useRoute();
const toast = useToast();
const rosterStore = useRosterStore();
const templateStore = useTemplateStore();
const confirmDelete = useConfirmDelete();
const isLg = useMediaQuery('(min-width: 1024px)');

const organId = computed(() => parseInt(route.params.id as string));

// TODO Bring this filtering to backend instead of frontend
const rosters = computed(() =>
    Object.values(rosterStore.rosters)
        .filter((r) => r.organId === organId.value)
        .sort((a, b) => (a.date ?? '').localeCompare(b.date ?? '') || (a.name ?? '').localeCompare(b.name ?? '')),
);
const roster = computed<Roster | undefined>(() =>
    rosterStore.selectedRosterId ? rosterStore.getRoster(rosterStore.selectedRosterId) : undefined,
);
const rosterId = computed(() => roster.value?.id ?? 0);
const loaded = ref(false);

const grid = useRosterGrid(rosterId, organId);

async function loadRosters() {
    loaded.value = false;
    try {
        await rosterStore.fetchRosters(organId.value);
    } catch (e) {
        console.error(e);
        toast.add({
            severity: 'error',
            summary: 'Failed to load rosters',
            detail: 'Could not load the rosters for this organ. Please reload the page to try again.',
            life: 5000,
        });
    }
    if (!roster.value && rosters.value.length > 0) {
        const today = new Date();
        today.setHours(0, 0, 0, 0);
        const upcoming = rosters.value.find((r) => (parseRosterDate(r.date) ?? today) >= today);
        rosterStore.setSelectedRoster((upcoming ?? rosters.value[0]).id);
    }
    loaded.value = true;
    void templateStore.fetchTemplates(organId.value);
}
watch(organId, loadRosters, { immediate: true });

const date = computed(() => parseRosterDate(roster.value?.date));
const template = computed(() =>
    roster.value?.templateId ? templateStore.getTemplate(roster.value.templateId) : undefined,
);
const rosterOptionLabel = (r: Roster) => {
    const d = parseRosterDate(r.date);
    return d ? `${r.name} · ${formatRosterDate(d)}` : r.name;
};

type ViewMode = 'mine' | 'all';
const chosenView = ref<ViewMode | null>(null);
const view = computed<ViewMode>(() => {
    if (grid.isLocked.value || !grid.me.value) return 'all';
    return chosenView.value ?? (grid.isAdmin.value ? 'all' : 'mine');
});
const viewOptions = [
    { label: 'My answers', value: 'mine' },
    { label: 'Everyone', value: 'all' },
];

const selectedShiftId = ref<number | null>(null);
watch([rosterId, () => grid.isLocked.value], () => (selectedShiftId.value = null));
const drawerVisible = computed({
    get: () => !isLg.value && selectedShiftId.value !== null,
    set: (open) => {
        if (!open) selectedShiftId.value = null;
    },
});

const responsibleNames = computed(() =>
    rosterStore
        .getResponsibles(rosterId.value)
        .map((r) => grid.members.value.find((m) => m.id === r.userId)?.displayName)
        .filter((name): name is string => !!name),
);

type DialogType = 'new-roster' | 'edit' | 'add-shift';
const activeDialog = ref<DialogType | null>(null);
const closeDialog = () => (activeDialog.value = null);
const actionsMenu = ref<InstanceType<typeof Menu>>();

const actionItems = computed(() => {
    const items: {
        label: string;
        iconComponent: unknown;
        danger?: boolean;
        phoneOnly?: boolean;
        command: () => void;
    }[] = [];
    if (grid.isAdmin.value) {
        items.push({
            label: 'New roster',
            iconComponent: OiPlus,
            phoneOnly: true,
            command: () => (activeDialog.value = 'new-roster'),
        });
        if (!grid.isLocked.value) {
            items.push({
                label: 'Add shift',
                iconComponent: OiPlus,
                phoneOnly: true,
                command: () => (activeDialog.value = 'add-shift'),
            });
        }
        items.push({ label: 'Edit details', iconComponent: OiPencil, command: () => (activeDialog.value = 'edit') });
    }
    if (grid.isLocked.value) {
        items.push({
            label: 'Export as image',
            iconComponent: OiDownload,
            command: () =>
                void exportRoster(rosterId.value).catch((error) => {
                    console.error('Export failed', error);
                    toast.add({ severity: 'error', summary: 'Export failed', life: 5000 });
                }),
        });
    }
    if (grid.isAdmin.value) {
        items.push({
            label: 'Delete roster',
            iconComponent: OiTrash,
            danger: true,
            command: () =>
                confirmDelete({
                    header: 'Delete roster?',
                    message: `"${roster.value?.name}" and all its answers are deleted. This cannot be undone.`,
                    acceptLabel: 'Delete roster',
                    onAccept: () => rosterStore.deleteRoster(rosterId.value),
                }),
        });
    }
    return items;
});

const desktopActionItems = computed(() => actionItems.value.filter((item) => !item.phoneOnly));
const menuItems = ref(actionItems.value);

function openActions(event: Event, phone: boolean) {
    menuItems.value = phone ? actionItems.value : desktopActionItems.value;
    actionsMenu.value?.toggle(event);
}
</script>

<template>
    <div class="flex flex-1 flex-col md:gap-5 md:px-10 md:pb-8 md:pt-7">
        <template v-if="roster">
            <div
                class="flex flex-col gap-3 border-b border-surface bg-surface-0 px-4 pb-3 pt-4 md:gap-4 md:border-0 md:bg-transparent md:p-0 dark:bg-surface-900 md:dark:bg-transparent"
            >
                <PageHeader :title="roster.name">
                    <template #status><StatusTag :saved="grid.isLocked.value" /></template>
                    <template #meta>
                        <span v-if="date" class="inline-flex items-center gap-1.5">
                            {{ formatRosterDate(date) }} · Week {{ isoWeek(date) }}
                        </span>
                        <span v-if="template" class="inline-flex items-center gap-1.5">
                            <OiClone aria-hidden="true" :size="16" />
                            From template
                            <RouterLink class="text-primary underline underline-offset-2" :to="`/templates/${organId}`">
                                {{ template.name }}
                            </RouterLink>
                        </span>
                    </template>
                    <template #actions>
                        <Select
                            aria-label="Switch roster"
                            class="w-full md:w-auto"
                            :model-value="rosterId"
                            :option-label="rosterOptionLabel"
                            option-value="id"
                            :options="rosters"
                            :pt="{
                                root: {
                                    class: 'h-10 items-center border-surface-300 bg-surface-0 dark:border-surface-600 dark:bg-surface-900',
                                },
                                label: { class: 'py-0 text-sm font-semibold text-color' },
                            }"
                            @update:model-value="rosterStore.setSelectedRoster($event)"
                        >
                            <template #value>Switch roster</template>
                        </Select>
                        <Button
                            v-if="grid.isAdmin.value"
                            class="hidden h-10 bg-surface-0 text-sm md:inline-flex dark:bg-surface-900"
                            label="New roster"
                            outlined
                            severity="secondary"
                            @click="activeDialog = 'new-roster'"
                        >
                            <template #icon><OiPlus aria-hidden="true" :size="16" /></template>
                        </Button>
                        <button
                            v-if="desktopActionItems.length"
                            aria-controls="roster-actions"
                            aria-haspopup="menu"
                            aria-label="More actions"
                            class="hidden size-10 items-center justify-center rounded-control border border-surface-300 bg-surface-0 text-surface-700 hover:bg-surface-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary md:flex dark:border-surface-600 dark:bg-surface-900 dark:text-surface-100 dark:hover:bg-surface-800"
                            type="button"
                            @click="openActions($event, false)"
                        >
                            <OiEllipsisH aria-hidden="true" :size="18" />
                        </button>
                    </template>
                </PageHeader>

                <SelectButton
                    v-if="!grid.isLocked.value && grid.me.value"
                    aria-label="Show"
                    :model-value="view"
                    option-label="label"
                    option-value="value"
                    :options="viewOptions"
                    :pt="{ root: { class: 'self-start' }, pcToggleButton: { root: { class: 'h-10 px-4 text-sm' } } }"
                    @update:model-value="(value: ViewMode | null) => value && (chosenView = value)"
                />
            </div>

            <Teleport to="#topbar-actions">
                <button
                    v-if="actionItems.length"
                    aria-controls="roster-actions"
                    aria-haspopup="menu"
                    aria-label="Roster actions"
                    class="flex size-11 items-center justify-center rounded-lg text-white hover:bg-topbar-hover focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white"
                    type="button"
                    @click="openActions($event, true)"
                >
                    <OiEllipsisH aria-hidden="true" :size="20" />
                </button>
            </Teleport>

            <MyPreferences
                v-if="view === 'mine'"
                class="px-4 py-4 md:p-0"
                :grid="grid"
                :responsible-names="responsibleNames"
            />
            <template v-else>
                <RosterShiftList class="flex-1 md:hidden" :grid="grid" />
                <div class="hidden items-start gap-5 md:flex">
                    <RosterGrid
                        v-model:selected-shift-id="selectedShiftId"
                        class="min-w-0 grow"
                        :grid="grid"
                        @add-shift="activeDialog = 'add-shift'"
                    />
                    <aside
                        v-if="isLg && selectedShiftId !== null"
                        aria-label="Assign people"
                        class="sticky top-5 max-h-[calc(100vh-2.5rem)] w-80 shrink-0 overflow-y-auto rounded-card border border-surface bg-surface-0 dark:bg-surface-900"
                    >
                        <RosterAssignPanel :grid="grid" :shift-id="selectedShiftId" @close="selectedShiftId = null" />
                    </aside>
                </div>
            </template>

            <RosterResponsibles
                class="mx-4 mb-4 md:m-0"
                :is-admin="grid.isAdmin.value"
                :members="grid.members.value"
                :roster-id="rosterId"
            />

            <Drawer
                v-model:visible="drawerVisible"
                position="right"
                :pt="{
                    root: { class: 'w-[min(22rem,100vw)]' },
                    header: { class: 'hidden' },
                    content: { class: 'p-0' },
                }"
            >
                <RosterAssignPanel
                    v-if="selectedShiftId !== null"
                    :grid="grid"
                    :shift-id="selectedShiftId"
                    @close="selectedShiftId = null"
                />
            </Drawer>
            <EditDialog :roster="roster" :visible="activeDialog === 'edit'" @update:visible="closeDialog" />
            <AddShiftDialog
                :visible="activeDialog === 'add-shift'"
                @submit="grid.addShift"
                @update:visible="closeDialog"
            />
        </template>

        <div v-else-if="loaded" class="flex flex-col gap-5 px-4 py-6 md:p-0">
            <PageHeader description="Plan who works which shift." title="Rosters" />
            <div
                class="flex flex-col items-center gap-3 rounded-card border border-dashed border-surface-400 px-6 py-12 text-center dark:border-surface-500"
            >
                <p class="text-[15px] font-semibold">No rosters yet</p>
                <p class="text-sm text-surface-600 dark:text-surface-200">
                    {{
                        grid.isAdmin.value
                            ? 'Create the first roster to start collecting answers.'
                            : 'Rosters will show up here once an organiser creates one.'
                    }}
                </p>
                <Button
                    v-if="grid.isAdmin.value"
                    class="h-10 text-sm"
                    label="New roster"
                    @click="activeDialog = 'new-roster'"
                >
                    <template #icon><OiPlus aria-hidden="true" :size="16" /></template>
                </Button>
            </div>
        </div>

        <Menu id="roster-actions" ref="actionsMenu" :model="menuItems" popup>
            <template #item="{ item, props: itemProps }">
                <a v-bind="itemProps.action" :class="['gap-2', item.danger ? 'text-red-700 dark:text-red-300' : '']">
                    <component :is="item.iconComponent" aria-hidden="true" class="shrink-0" :size="16" />
                    <span>{{ item.label }}</span>
                </a>
            </template>
        </Menu>
        <NewRosterDialog :organ-id="organId" :visible="activeDialog === 'new-roster'" @update:visible="closeDialog" />
    </div>
</template>
