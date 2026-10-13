<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useToast } from 'openvue/usetoast';
import { OiArrowCircleDown } from '@openvue/openicons-vue';
import {
    GEWISRoosterInternalModelsGroupPriority,
    type ShiftGroupPriority,
    type User,
} from '@gewis/grooster-backend-ts';
import ApiService from '@/services/ApiService';
import { fieldSelectPt } from '@/theme/pt';

export type GroupPriority =
    (typeof GEWISRoosterInternalModelsGroupPriority)[keyof typeof GEWISRoosterInternalModelsGroupPriority];

const priorityOptions = Object.entries(GEWISRoosterInternalModelsGroupPriority).map(([label, value]) => ({
    label,
    value,
}));

const props = defineProps<{ groupId: number; organId: number }>();

const toast = useToast();
const users = ref<User[]>([]);
const priorities = ref<ShiftGroupPriority[]>([]);
const selected = ref<User | null>(null);
const pushing = ref(false);

const priorityOf = (userId: number): GroupPriority =>
    priorities.value.find((p) => p.userId === userId)?.priority ?? GEWISRoosterInternalModelsGroupPriority.Default;
const priorityLabel = (userId: number) => priorityOptions.find((o) => o.value === priorityOf(userId))?.label;
const selectedPriority = computed(() => (selected.value ? priorityOf(selected.value.id) : null));

async function fetchUsers() {
    try {
        users.value = (await ApiService.user.userGet(props.organId)).data;
    } catch (error) {
        console.error(error);
    }
}

async function fetchPriorities() {
    try {
        priorities.value = (await ApiService.shiftGroupApi.getShiftGroupPriorities(props.groupId)).data ?? [];
    } catch (error) {
        console.error(error);
    }
}

async function setPriority(userId: number, priority: GroupPriority) {
    try {
        await ApiService.shiftGroupApi.updateShiftGroupPriority(props.groupId, { userId, priority });
        const existing = priorities.value.find((p) => p.userId === userId);
        if (existing) existing.priority = priority;
        else priorities.value.push({ userId, priority });
    } catch (error) {
        console.error(error);
        toast.add({ severity: 'error', summary: 'Could not change the priority', life: 5000 });
    }
}

async function pushToBottom(user: User) {
    pushing.value = true;
    try {
        await ApiService.shiftGroupApi.pushUserToBottom(props.groupId, { userId: user.id });
        await fetchPriorities();
        toast.add({
            severity: 'success',
            summary: 'Pushed to bottom',
            detail: `${user.name} was moved to the bottom of the group.`,
            life: 3000,
        });
    } catch (error) {
        console.error(error);
        toast.add({
            severity: 'error',
            summary: 'Failed to push user to bottom',
            detail: 'Something went wrong while updating the shift group order. Please try again.',
            life: 5000,
        });
    } finally {
        pushing.value = false;
    }
}

onMounted(() => Promise.all([fetchUsers(), fetchPriorities()]));
</script>

<template>
    <div class="flex flex-col overflow-hidden rounded-card border border-surface md:h-[420px] md:flex-row">
        <ul
            aria-label="Members"
            class="max-h-60 shrink-0 overflow-y-auto border-b border-surface md:max-h-none md:w-64 md:border-b-0 md:border-r"
        >
            <li v-for="user in users" :key="user.id">
                <button
                    :aria-current="selected?.id === user.id ? 'true' : undefined"
                    :class="[
                        'flex h-12 w-full items-center gap-2 border-l-[3px] px-4 text-left text-sm focus-visible:outline-2 focus-visible:-outline-offset-2 focus-visible:outline-primary',
                        selected?.id === user.id
                            ? 'border-l-primary bg-highlight font-semibold'
                            : 'border-l-transparent hover:bg-surface-50 dark:hover:bg-surface-800',
                    ]"
                    type="button"
                    @click="selected = user"
                >
                    <span class="min-w-0 grow truncate">{{ user.name }}</span>
                    <span class="text-xs text-surface-600 dark:text-surface-200">{{ priorityLabel(user.id) }}</span>
                </button>
            </li>
        </ul>

        <div class="flex grow flex-col gap-5 p-6">
            <template v-if="selected">
                <h3 class="font-display text-xl font-semibold">{{ selected.name }}</h3>
                <div class="flex max-w-xs flex-col gap-1.5">
                    <label
                        class="text-[13px] font-semibold text-surface-700 dark:text-surface-100"
                        for="group-priority"
                    >
                        Priority
                    </label>
                    <Select
                        class="h-[42px] w-full items-center"
                        input-id="group-priority"
                        :model-value="selectedPriority"
                        option-label="label"
                        option-value="value"
                        :options="priorityOptions"
                        :pt="fieldSelectPt"
                        @update:model-value="setPriority(selected.id, $event)"
                    />
                </div>
                <Button
                    class="h-10 self-start bg-surface-0 text-sm dark:bg-surface-900"
                    label="Push to bottom"
                    :loading="pushing"
                    outlined
                    severity="secondary"
                    @click="pushToBottom(selected)"
                >
                    <template #icon><OiArrowCircleDown aria-hidden="true" :size="16" /></template>
                </Button>
            </template>
            <p v-else class="m-auto text-sm text-surface-600 dark:text-surface-200">
                Pick a member to set their priority.
            </p>
        </div>
    </div>
</template>
