<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { OiChevronRight, OiPlus } from '@openvue/openicons-vue';
import type { ShiftGroup } from '@gewis/grooster-backend-ts';
import ApiService from '@/services/ApiService';
import ShiftGroupAdd from '@/components/templates/dialogs/ShiftGroupAdd.vue';
import ShiftGroupPriority from '@/components/templates/ShiftGroupPriority.vue';
import { dialogPt } from '@/theme/pt';

const props = defineProps<{ organId: number; isAdmin: boolean }>();
const emit = defineEmits<{ 'update:groups': [groups: ShiftGroup[]] }>();

const groups = ref<ShiftGroup[]>([]);
const addOpen = ref(false);
const openGroup = ref<ShiftGroup | null>(null);

async function fetchGroups() {
    try {
        const response = await ApiService.shiftGroupApi.getShiftGroups(props.organId);
        groups.value = response.data ?? [];
        emit('update:groups', groups.value);
    } catch (error) {
        console.error('Failed to fetch groups:', error);
    }
}

onMounted(fetchGroups);
</script>

<template>
    <section aria-labelledby="shift-groups-title" class="flex flex-col gap-4 pt-2">
        <div class="flex flex-col gap-3 md:flex-row md:items-end md:gap-4">
            <div class="flex grow flex-col gap-1">
                <h2 id="shift-groups-title" class="font-display text-[22px] font-semibold">Shift groups</h2>
                <p class="text-sm text-surface-600 dark:text-surface-200">
                    Who can be scheduled on a shift. Link a group to a shift in a template.
                </p>
            </div>
            <Button
                v-if="isAdmin"
                class="h-10 self-start bg-surface-0 text-sm md:self-auto dark:bg-surface-900"
                label="New group"
                outlined
                severity="secondary"
                @click="addOpen = true"
            >
                <template #icon><OiPlus aria-hidden="true" :size="16" /></template>
            </Button>
        </div>

        <ul
            v-if="groups.length"
            class="overflow-hidden rounded-card border border-surface bg-surface-0 dark:bg-surface-900"
        >
            <li v-for="group in groups" :key="group.id" class="border-t border-surface first:border-t-0">
                <button
                    class="flex w-full items-center gap-4 px-5 py-4 text-left hover:bg-surface-50 focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-primary dark:hover:bg-surface-800"
                    :disabled="!isAdmin"
                    type="button"
                    @click="openGroup = group"
                >
                    <span class="grow text-[15px] font-semibold">{{ group.name }}</span>
                    <span v-if="isAdmin" class="text-[13px] text-surface-600 dark:text-surface-200">Priorities</span>
                    <OiChevronRight
                        v-if="isAdmin"
                        aria-hidden="true"
                        class="text-surface-600 dark:text-surface-200"
                        :size="18"
                    />
                </button>
            </li>
        </ul>
        <p
            v-else
            class="rounded-card border border-dashed border-surface-400 px-6 py-8 text-center text-sm text-surface-600 dark:border-surface-500 dark:text-surface-200"
        >
            No shift groups yet.
        </p>
    </section>

    <ShiftGroupAdd v-model:visible="addOpen" :organ-id="organId" @created="fetchGroups" />
    <Dialog
        :header="openGroup ? `${openGroup.name}: priorities` : ''"
        modal
        :pt="{ ...dialogPt, root: { class: `${dialogPt.root.class} md:w-[800px]` } }"
        :visible="openGroup !== null"
        @update:visible="(open: boolean) => !open && (openGroup = null)"
    >
        <ShiftGroupPriority v-if="openGroup" :group-id="openGroup.id" :organ-id="organId" />
    </Dialog>
</template>
