<script setup lang="ts">
import { User } from '@gewis/grooster-backend-ts';
import { computed, onMounted, ref, watch } from 'vue';
import { useRosterStore } from '@/stores/roster.store';
import { useOrganStore } from '@/stores/organ.store';
import { Role, useAuthStore } from '@/stores/auth.store';
import ApiService from '@/services/ApiService';

const props = defineProps<{
    id: number;
}>();

const rosterStore = useRosterStore();
const organStore = useOrganStore();
const authStore = useAuthStore();

const organUsers = ref<User[]>([]);
const selectedUserId = ref<number | null>(null);

const responsibles = computed(() => rosterStore.getResponsibles(props.id));

const responsibleUsers = computed(() =>
    responsibles.value
        .map((responsible) => organUsers.value.find((user) => user.id === responsible.userId))
        .filter((user): user is User => !!user),
);

const availableUsers = computed(() =>
    organUsers.value.filter((user) => !responsibles.value.some((responsible) => responsible.userId === user.id)),
);

onMounted(async () => {
    try {
        const response = await ApiService.user.userGet(organStore.organ);
        organUsers.value = response.data;
    } catch (error) {
        console.error('Failed to fetch organ users:', error);
    }
});

watch(
    () => props.id,
    (rosterId) => {
        void rosterStore.fetchResponsibles(rosterId);
    },
    { immediate: true },
);

const addResponsible = async (userId: number) => {
    if (!userId) return;

    try {
        await rosterStore.addResponsible(props.id, userId);
    } catch (error) {
        console.error('Failed to add roster responsible:', error);
    }

    selectedUserId.value = null;
};

const removeResponsible = async (userId: number) => {
    try {
        await rosterStore.removeResponsible(props.id, userId);
    } catch (error) {
        console.error('Failed to remove roster responsible:', error);
    }
};
</script>

<template>
    <div class="space-y-3 border-b border-slate-200 pb-4">
        <div>
            <h2 class="text-xl font-bold text-slate-800">Responsibles</h2>
            <p class="text-sm text-slate-500">Users responsible for organising this activity.</p>
        </div>

        <div class="flex flex-wrap items-center gap-3">
            <div
                v-for="user in responsibleUsers"
                :key="user.id"
                class="group flex items-center bg-slate-50 border border-slate-200 rounded-lg pl-3 pr-1 py-1 transition-colors hover:bg-slate-100"
            >
                <span class="text-sm font-medium text-slate-700">{{ user.name }}</span>
                <Button
                    v-if="authStore.can([Role.Admin, Role.Owner])"
                    class="!p-0 !w-7 !h-7 !text-slate-400 group-hover:!text-red-500"
                    icon="pi pi-times"
                    rounded
                    text
                    @click="removeResponsible(user.id)"
                />
            </div>

            <p v-if="responsibleUsers.length === 0" class="text-sm text-slate-400">No responsibles assigned yet.</p>

            <div v-if="authStore.can([Role.Admin, Role.Owner])" class="relative min-w-[180px]">
                <Select
                    v-model="selectedUserId"
                    class="!w-full !rounded-lg !border-dashed !border-slate-300 !bg-transparent hover:!border-indigo-400 transition-all !text-sm"
                    option-label="name"
                    option-value="id"
                    :options="availableUsers"
                    placeholder="+ Add Responsible"
                    @change="(event) => addResponsible(event.value)"
                />
            </div>
        </div>
    </div>
</template>

<style scoped></style>
