<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useToast } from 'openvue/usetoast';
import { OiTimes } from '@openvue/openicons-vue';
import { useRosterStore } from '@/stores/roster.store';
import type { RosterMember } from '@/composables/useRosterGrid';

const props = defineProps<{
    rosterId: number;
    members: RosterMember[];
    isAdmin: boolean;
}>();

const rosterStore = useRosterStore();
const toast = useToast();
const selectedUserId = ref<number | null>(null);

const responsibles = computed(() => rosterStore.getResponsibles(props.rosterId));
const responsibleMembers = computed(() =>
    responsibles.value
        .map((responsible) => props.members.find((member) => member.id === responsible.userId))
        .filter((member): member is RosterMember => !!member),
);
const availableMembers = computed(() =>
    props.members.filter((member) => !responsibles.value.some((responsible) => responsible.userId === member.id)),
);

watch(
    () => props.rosterId,
    (rosterId) => void rosterStore.fetchResponsibles(rosterId),
    { immediate: true },
);

async function addResponsible(userId: number | null) {
    if (!userId) return;
    try {
        await rosterStore.addResponsible(props.rosterId, userId);
    } catch (error) {
        console.error('Failed to add roster responsible:', error);
        toast.add({ severity: 'error', summary: 'Could not add the responsible', life: 5000 });
    }
    selectedUserId.value = null;
}

async function removeResponsible(userId: number) {
    try {
        await rosterStore.removeResponsible(props.rosterId, userId);
    } catch (error) {
        console.error('Failed to remove roster responsible:', error);
        toast.add({ severity: 'error', summary: 'Could not remove the responsible', life: 5000 });
    }
}
</script>

<template>
    <section
        aria-labelledby="responsibles-title"
        class="flex flex-col gap-3 rounded-card border border-surface bg-surface-0 px-5 py-[18px] md:flex-row md:items-center md:gap-5 dark:bg-surface-900"
    >
        <div class="flex flex-col gap-0.5 md:w-[300px] md:shrink-0">
            <h2 id="responsibles-title" class="text-base font-semibold">Responsibles</h2>
            <p class="text-[13px] text-surface-600 dark:text-surface-200">
                Members responsible for organising this activity.
            </p>
        </div>
        <ul class="flex grow flex-wrap items-center gap-2">
            <li
                v-for="member in responsibleMembers"
                :key="member.id"
                class="inline-flex h-10 items-center gap-1 rounded-full border border-surface pl-3.5 pr-1.5 text-sm font-medium"
            >
                <span :title="member.name">{{ member.displayName }}</span>
                <button
                    v-if="isAdmin"
                    :aria-label="`Remove ${member.displayName} as responsible`"
                    class="flex size-7 items-center justify-center rounded-full text-surface-700 hover:bg-surface-100 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary dark:text-surface-100 dark:hover:bg-surface-800"
                    type="button"
                    @click="removeResponsible(member.id)"
                >
                    <OiTimes aria-hidden="true" :size="16" />
                </button>
            </li>
            <li v-if="!responsibleMembers.length && !isAdmin" class="text-sm text-surface-600 dark:text-surface-200">
                No responsibles yet.
            </li>
            <li v-if="isAdmin">
                <Select
                    v-model="selectedUserId"
                    aria-label="Add responsible"
                    filter
                    filter-placeholder="Search members"
                    option-label="displayName"
                    option-value="id"
                    :options="availableMembers"
                    placeholder="+ Add responsible"
                    :pt="{
                        root: {
                            class: 'h-10 items-center rounded-control border-dashed border-surface-400 bg-transparent dark:border-surface-500',
                        },
                        label: { class: 'py-0 text-sm font-semibold text-primary' },
                        dropdown: { class: 'hidden' },
                    }"
                    @update:model-value="addResponsible"
                />
            </li>
        </ul>
    </section>
</template>
